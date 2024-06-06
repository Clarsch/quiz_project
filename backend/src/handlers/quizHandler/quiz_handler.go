package quizHandler

import (
	"encoding/json"
	"fmt"
	"quizzy_game/internal/dataTypes"
	dbdto "quizzy_game/internal/dto/dbDTO"
	frontdto "quizzy_game/internal/dto/frontDTO"
	quizstatus "quizzy_game/internal/enums/quizStatus"
	requestcommand "quizzy_game/internal/enums/requestCommands"
	"quizzy_game/network"
	"sync"
	"time"

	"github.com/google/uuid"
)

var quizzes = make(map[string]*dataTypes.Quiz)
var categories = make(map[int]dbdto.CategoryIncomingDTO)
var answerTimeout = 10 * time.Second

func init() {
	var dbCategories []dbdto.CategoryIncomingDTO = network.GetCategories()
	for _, category := range dbCategories {
		categories[category.Id] = category
	}
	fmt.Println("Categories are fetched. Number of available categories: ", len(categories))

}

func HandleQuizRequest(request frontdto.Request, user *dataTypes.User) {

	responseChannel := user.MsgChannel

	jsonData, _ := json.Marshal(request.Data)

	if request.RequestType == requestcommand.CreateQuiz {
		newQuizId := createQuiz(jsonData, user)
		newQuiz, ok := quizzes[newQuizId]
		if !ok {
			fmt.Println("Quiz not found.")
			responseChannel <- "Something went wrong. Error Creating Quiz!"
			return
		}
		joinQuiz(newQuizId, user)
		responseChannel <- fmt.Sprintf("QuizID: %s, QuizStatus: %s, participants: %s ", newQuizId, newQuiz.QuizStatus, newQuiz.ParticipantsAsString())

	} else if request.RequestType == requestcommand.AnswerQuestion {
		handleAnswer(jsonData, request.ReceivedTime, user)

	} else {
		var updateData frontdto.QuizUpdateRequestDTO
		err := json.Unmarshal(jsonData, &updateData)
		if err != nil {
			fmt.Println("Error:", err)
			responseChannel <- fmt.Sprintf("Invalid input! Error: %s\n", err)
			responseChannel <- fmt.Sprintf("Example of correct json input: %s\n", updateData.GetExample())
			return
		}
		switch request.RequestType {
		case requestcommand.JoinQuiz:
			responseChannel <- joinQuiz(updateData.QuizId, user)
		case requestcommand.LeaveQuiz:
			responseChannel <- leaveQuiz(updateData.QuizId, user)
		case requestcommand.StartQuiz:
			responseChannel <- startQuiz(updateData.QuizId)
		case requestcommand.ResetQuiz:
			responseChannel <- resetQuiz(updateData.QuizId)
		default:
			responseChannel <- "Unknown command. Try: \n\tcreateQuiz \n\tjoinQuiz \n\tleaveQuiz \n\tstartQuiz \n\tstopQuiz \n\tresetQuiz \n\tprint"
		}
	}
}

func createQuiz(inputJson []byte, user *dataTypes.User) string {

	var createData frontdto.QuizCreateRequestDTO
	err := json.Unmarshal(inputJson, &createData)
	if err != nil {
		fmt.Println("Error:", err)
		user.MsgChannel <- fmt.Sprintf("Invalid input! Error: %s\n", err)
		user.MsgChannel <- fmt.Sprintf("Example of correct json input: %s\n", createData.GetExample())
		return ""
	}

	questionIncomingDtos := network.GetQuestions(
		createData.CategoryId,
		createData.Difficulty,
		createData.Type,
	)

	questions := make(map[string]*dataTypes.Question)
	for _, qIDto := range questionIncomingDtos {
		q := dataTypes.NewQuestion(qIDto)
		questions[q.Id] = &q
	}

	newQuiz := dataTypes.Quiz{
		Id:           uuid.NewString(),
		Name:         createData.Name,
		QuizStatus:   quizstatus.StatusInitialized,
		Category:     categories[createData.CategoryId],
		Difficulty:   createData.Difficulty,
		Type:         createData.Type,
		Questions:    questions,
		Participants: make(map[string]*dataTypes.Participant),
	}

	quizzes[newQuiz.Id] = &newQuiz
	fmt.Printf("Sucessfully created Quiz %s with ID: %s\n", newQuiz.Name, newQuiz.Id)
	return newQuiz.Id
}

func joinQuiz(quizID string, user *dataTypes.User) string {
	quiz, ok := quizzes[quizID]
	if !ok {
		fmt.Println("Quiz not found.")
		return "Error joining Quiz: " + quizID
	}
	if _, ok := quiz.Participants[user.Id]; ok {
		fmt.Println("Error Joining Quiz. UserName is already exist.")
		return fmt.Sprintf("Error joining Quiz: %s. User is already joined! ", quiz.Name)

	}

	quiz.Participants[user.Id] = dataTypes.NewParticipant(user)
	fmt.Printf("Added user %s to Quiz: %s\n", user.Name, quizID)

	msg := fmt.Sprintf("QuizID: %s, QuizStatus: %s, participants: %s\n", quiz.Id, quiz.QuizStatus, quiz.ParticipantsAsString())
	go broadcastToParticipants(quizID, msg)
	return "Sucessfully joined the quiz: " + quizID
}

func leaveQuiz(quizID string, user *dataTypes.User) string {
	quiz, ok := quizzes[quizID]
	if !ok {
		return "Error leaving Quiz. Quiz not found. ID: " + quizID
	}
	if _, ok := quiz.Participants[user.Id]; ok {
		delete(quiz.Participants, user.Id)
		fmt.Printf("Deleted user %s from Quiz %s: \n", user.Name, quizID)
	}
	msg := fmt.Sprintf("QuizID: %s, QuizStatus: %s, participants: %s\n", quiz.Id, quiz.QuizStatus, quiz.ParticipantsAsString())
	broadcastToParticipants(quizID, msg)
	return fmt.Sprintf("User: %s left Quiz QuizID: %s\n", user.Name, quiz.Id)
}

func startQuiz(quizID string) string {
	quiz, ok := quizzes[quizID]
	switch {
	case !ok:
		return "Error Starting. Quiz not found. ID: " + quizID
	case quiz.QuizStatus != quizstatus.StatusInitialized:
		return fmt.Sprintf("Quiz could not start. Expected status Initialized but got " + string(quiz.QuizStatus))
	case len(quiz.Participants) < 2:
		return fmt.Sprintf("Quiz could not start. Expected PARTICIPANTS to contain more than 2. got %d.", len(quiz.Participants))
	default:
		quiz.QuizStatus = quizstatus.StatusStart
		fmt.Println("Quiz Status updated to: ", quizstatus.StatusStart)

		var wg sync.WaitGroup
		statusChannel := make(chan quizstatus.QuizStatus)
		quiz.StatusChannel = &statusChannel

		wg.Add(1)
		go timerRoutine(&wg, statusChannel)

		questionLoopRoutine(quizID, statusChannel)

		quiz.QuizStatus = quizstatus.StatusEnded
		scoreBoardMsg := fmt.Sprintf("QuizID: %s, QuizStatus: %s, participants: %s\n", quiz.Id, quiz.QuizStatus, quiz.ScoreBoard())
		broadcastToParticipants(quizID, scoreBoardMsg)
		fmt.Println("Quiz Status updated to: ", quizstatus.StatusEnded)
		fmt.Println("Quiz Scoreboard: ", scoreBoardMsg)

		wg.Wait() // Wait for all goroutines to finish
		close(statusChannel)

	}
	return fmt.Sprintf("Sucessfully ran quiz with ID: %s\n", quiz.Id)
}

func stopQuiz(quizID string) string {
	quiz, ok := quizzes[quizID]
	if !ok {
		fmt.Println("Quiz not found.")
		return "Error stopping quiz with ID: " + quizID
	}
	fmt.Println("Stopping Quiz with ID: ", quiz.Id)

	quiz.QuizStatus = quizstatus.StatusStopped
	*quiz.StatusChannel <- quizstatus.StatusStopped
	msg := fmt.Sprintf("Quiz Status updated to: %s", quizstatus.StatusStopped)
	broadcastToParticipants(quizID, msg)
	return msg
}

func resetQuiz(quizID string) string {
	quiz, ok := quizzes[quizID]
	if !ok {
		fmt.Println("Quiz not found.")
		return "Error reset quiz with ID: " + quizID
	}
	fmt.Println("Resetting Quiz with ID: ", quiz.Id)

	for _, question := range quiz.Questions {
		question.IsAskedStatus = false
	}
	fmt.Println("Questions reset for Quiz with ID: ", quiz.Id)

	for _, participant := range quiz.Participants {
		participant.Score = 0
	}
	fmt.Println("Scores for users reset for Quiz with ID: ", quiz.Id)

	quiz.QuizStatus = quizstatus.StatusInitialized
	return fmt.Sprintf("Quiz Status updated to: %s", quizstatus.StatusInitialized)

}

func questionLoopRoutine(quizID string, statusChannel chan quizstatus.QuizStatus) {
	quiz, ok := quizzes[quizID]
	if !ok {
		fmt.Println("Question Loop: Quiz not found.")
		return
	}

	if quiz.RemainingQuestions() < 1 {
		fmt.Println("Question Loop: No questions remaining.")
		return
	}

	for _, question := range quiz.Questions {
		if quiz.QuizStatus == quizstatus.StatusStopped {
			statusChannel <- quizstatus.StatusStopped
			return
		}
		if question.IsAskedStatus {
			fmt.Printf("Skipping question. Has already been asked!")
			break
		}

		quiz.QuizStatus = quizstatus.StatusQuizTime
		fmt.Println("Quiz Status updated to: ", quizstatus.StatusQuizTime)

		quizMsg := fmt.Sprintf("QuestionID: %s, Question: %s, Options: %s, CorrectAnswer: %s",
			question.Id, question.GetQuestion(), question.Options, question.GetCorrectAnswer())
		statusChannel <- quizstatus.StatusQuizTime
		question.IsAskedStatus = true
		question.LastAskedTime = time.Now()
		broadcastToParticipants(quizID, quizMsg)

		for status := range statusChannel {
			if status == quizstatus.StatusQuizTimeEnded {
				fmt.Println("Quiz Time Ended!")
				quiz.QuizStatus = quizstatus.StatusEvaluation
				fmt.Println("Quiz Status updated to: ", quizstatus.StatusEvaluation)
				break
			} else if status == quizstatus.StatusStopped {
				return

			}
		}
	}
}

func timerRoutine(wg *sync.WaitGroup, statusChannel chan quizstatus.QuizStatus) {
	TAG := "TIMER_ROUTINE: "

	for status := range statusChannel {
		fmt.Println(TAG, "RECEIVED status: ", status)

		if status == quizstatus.StatusQuizTime {
			answerTimer := time.NewTimer(answerTimeout)
			fmt.Println(TAG, "Answer timer started!")
			<-answerTimer.C
			statusChannel <- quizstatus.StatusQuizTimeEnded
			fmt.Println(TAG, "Answer timer ENDED!")

		} else if status == quizstatus.StatusEnded || status == quizstatus.StatusStopped {
			fmt.Println(TAG, "Shutting down the timer GoRoutine!")
			defer wg.Done()
			return
		}
	}

}

func handleAnswer(inputJson []byte, timeReceived time.Time, user *dataTypes.User) {
	TAG := "ANSWER_HANDLER: "

	var answerData frontdto.QuestionAnswerDTO
	err := json.Unmarshal(inputJson, &answerData)
	if err != nil {
		fmt.Println("Error:", err)
		user.MsgChannel <- fmt.Sprintf("Invalid input! Error: %s\n", err)
		user.MsgChannel <- fmt.Sprintf("Example of correct json input: %s\n", answerData.GetExample())
		return
	}

	quiz, quizOk := quizzes[answerData.QuizId]
	if !quizOk {
		fmt.Println(TAG, "Quiz not found.")
		return
	}
	question, questionOk := quiz.Questions[answerData.QuestionId]
	if !questionOk {
		fmt.Println(TAG, "Questions contains:", quiz.Questions)
		fmt.Println(TAG, "Question not found.")
		return
	}
	if !question.IsAskedStatus {
		fmt.Println(TAG, "Question has not been asked.")
		return
	}
	if answerData.AnswerId != question.CorrectOptionId {
		fmt.Println(TAG, "Wrong answer. 0 points")
		return
	}
	timeSpent := timeReceived.Sub(question.LastAskedTime)
	fmt.Printf("%sTime spent on answer: %s\n", TAG, timeSpent)

	if timeSpent > answerTimeout {
		statusMsg := fmt.Sprintf("%sAnswer took too long. Spent time: %fm:%fs, allowed time: %f seconds.\n", TAG, timeSpent.Minutes(), timeSpent.Seconds(), answerTimeout.Seconds())
		user.MsgChannel <- statusMsg
		fmt.Println(statusMsg)
		return
	}
	factor := 100.0
	millSecRemain := float64(answerTimeout.Abs().Milliseconds()) - float64(timeSpent.Abs().Milliseconds())
	points := int((millSecRemain) / factor)
	participant := quiz.Participants[user.Id]
	participant.Score += int(points)
	statusMsg := fmt.Sprintf("%sUser: %s earned %d points giving a total score at %d in Quiz: %s", TAG, participant.Ref.Name, int(points), participant.Score, quiz.Name)
	user.MsgChannel <- statusMsg
	fmt.Println(statusMsg)
}

func broadcastToParticipants(quizID, msg string) {
	quiz, ok := quizzes[quizID]
	if !ok {
		fmt.Println("Quiz not found.")
		return
	}

	var wg sync.WaitGroup
	for _, participant := range quiz.Participants {
		wg.Add(1)
		go broadcastToParticipant(participant, msg, &wg)
	}

	wg.Wait()
}

func broadcastToParticipant(participant *dataTypes.Participant, msg string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	select {
	case participant.Ref.MsgChannel <- msg:
		fmt.Printf("BROADCAST: Message sent to Participant %s\n", participant.Ref.Id)
	default:
		fmt.Printf("BROADCAST: Error sending message to Participant %s\n", participant.Ref.Id)
	}
}
