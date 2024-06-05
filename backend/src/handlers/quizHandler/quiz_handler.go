package quizHandler

import (
	"fmt"
	"quizzy_game/api"
	"quizzy_game/internal/dataTypes"
	dbdto "quizzy_game/internal/dto/dbDTO"
	"quizzy_game/internal/enums"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var quizzes = make(map[string]*dataTypes.Quiz)
var categories = make(map[int]dbdto.CategoryIncomingDTO)
var answerTimeout = 10 * time.Second

func init() {
	var dbCategories []dbdto.CategoryIncomingDTO = api.GetCategories()
	for _, category := range dbCategories {
		categories[category.Id] = category
	}
	fmt.Println("Categories are fetched. Number of available categories: ", len(categories))

}

func HandleQuizUpdate(quizUpdate string, user *dataTypes.User) {

	responseChannel := user.MsgChannel

	fmt.Println("CLIENT REQUEST: ", quizUpdate)
	input := strings.Fields(quizUpdate)
	if len(input) < 1 {
		responseChannel <- "Not enough input parameters! Try: \n\tcreateQuiz \n\tjoinQuiz \n\tleaveQuiz \n\tstartQuiz \n\tresetQuiz \n\tprint"
		return
	}

	switch input[0] {
	case "createQuiz":
		if len(input) < 2 {
			responseChannel <- fmt.Sprintln("Not enough input parameters!\n",
				"Try:\n\t createQuiz <quizName> \n",
				"OR\n\t createQuiz <quizName> <categoryId(9-32)> <easy||medium||hard> <multiple||boolean> ")
			return
		}
		quizName := input[1]
		var newQuizId string
		if len(input) == 5 {
			categoryId, _ := strconv.Atoi(input[2])
			category := categories[categoryId]
			difficulty := enums.Difficulty(input[3])
			quizType := enums.QuestionType(input[4])
			newQuizId = createQuiz(quizName, category, difficulty, quizType)

		} else {
			newQuizId = createQuiz(quizName, categories[9], enums.Easy, enums.MultipleChoice)
		}

		newQuiz, ok := quizzes[newQuizId]
		if !ok {
			fmt.Println("Quiz not found.")
			responseChannel <- "Something went wrong. Error Creating Quiz!"
			return
		}
		joinQuiz(newQuizId, user)
		responseChannel <- fmt.Sprintf("QuizID: %s, QuizStatus: %s, participants: %s ", newQuizId, newQuiz.QuizStatus, newQuiz.ParticipantsAsString())
	case "joinQuiz":
		if len(input) < 2 {
			responseChannel <- "Not enough input parameters! Try: joinQuiz <quizID>"
			return
		}
		quizID := input[1]
		responseChannel <- joinQuiz(quizID, user)
	case "leaveQuiz":
		if len(input) < 2 {
			responseChannel <- "Not enough input parameters! Try: leaveQuiz <quizID>"
			return
		}
		quizID := input[1]
		responseChannel <- leaveQuiz(quizID, user)
	case "startQuiz":
		if len(input) < 2 {
			responseChannel <- "Not enough input parameters! Try: startQuiz <quizID>"
			return
		}
		quizID := input[1]
		responseChannel <- startQuiz(quizID)
	case "stopQuiz":
		if len(input) < 2 {
			responseChannel <- "Not enough input parameters! Try: stopQuiz <quizID>"
			return
		}
		quizID := input[1]
		responseChannel <- stopQuiz(quizID)
	case "resetQuiz":
		if len(input) < 2 {
			responseChannel <- "Not enough input parameters! Try: resetQuiz <quizID>"
			return
		}
		quizID := input[1]
		responseChannel <- resetQuiz(quizID)
	case "answerQuestion":
		if len(input) < 4 {
			responseChannel <- "Not enough input parameters! Try: answerQuestion <QuizID> <QuestionID> <Answer>"
			return
		}
		quizID := input[1]
		questionID := input[2]
		answer := strings.Join(input[3:], " ")
		timeReceived := time.Now()
		fmt.Printf("Received answer \"%s\" for questionID %s \n at time: %s\n", answer, questionID, timeReceived.String())

		handleAnswer(quizID, questionID, answer, timeReceived, user)
	case "print":
		quizPrintString := "Printing Quizzes: \n"
		for _, quiz := range quizzes {
			quizPrintString += quiz.String()
		}
		fmt.Println(quizPrintString)
		responseChannel <- quizPrintString
	default:
		responseChannel <- "Unknown command. Try: \n\tcreateQuiz \n\tjoinQuiz \n\tleaveQuiz \n\tstartQuiz \n\tstopQuiz \n\tresetQuiz \n\tprint"
	}

}

func createQuiz(name string, category dbdto.CategoryIncomingDTO, difficulty enums.Difficulty, quizType enums.QuestionType) string {

	questionIncomingDtos := api.GetQuestions(category.Id, difficulty, quizType)
	questions := make(map[string]*dataTypes.Question)
	for _, qIDto := range questionIncomingDtos {
		q := dataTypes.NewQuestion(qIDto)
		questions[q.Id] = &q
	}
	newQuiz := dataTypes.Quiz{
		Id:           uuid.NewString(),
		Name:         name,
		QuizStatus:   enums.StatusInitialized,
		Category:     category,
		Difficulty:   enums.Easy,
		Type:         enums.MultipleChoice,
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
	case quiz.QuizStatus != enums.StatusInitialized:
		return fmt.Sprintf("Quiz could not start. Expected status Initialized but got " + string(quiz.QuizStatus))
	case len(quiz.Participants) < 2:
		return fmt.Sprintf("Quiz could not start. Expected PARTICIPANTS to contain more than 2. got %d.", len(quiz.Participants))
	default:
		quiz.QuizStatus = enums.StatusStart
		fmt.Println("Quiz Status updated to: ", enums.StatusStart)

		var wg sync.WaitGroup
		statusChannel := make(chan enums.QuizStatus)
		quiz.StatusChannel = &statusChannel

		wg.Add(1)
		go timerRoutine(&wg, statusChannel)

		questionLoopRoutine(quizID, statusChannel)

		quiz.QuizStatus = enums.StatusEnded
		scoreBoardMsg := fmt.Sprintf("QuizID: %s, QuizStatus: %s, participants: %s\n", quiz.Id, quiz.QuizStatus, quiz.ScoreBoard())
		broadcastToParticipants(quizID, scoreBoardMsg)
		fmt.Println("Quiz Status updated to: ", enums.StatusEnded)
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

	quiz.QuizStatus = enums.StatusStopped
	*quiz.StatusChannel <- enums.StatusStopped
	msg := fmt.Sprintf("Quiz Status updated to: %s", enums.StatusStopped)
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

	quiz.QuizStatus = enums.StatusInitialized
	return fmt.Sprintf("Quiz Status updated to: %s", enums.StatusInitialized)

}

func questionLoopRoutine(quizID string, statusChannel chan enums.QuizStatus) {
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
		if quiz.QuizStatus == enums.StatusStopped {
			statusChannel <- enums.StatusStopped
			return
		}
		if question.IsAskedStatus {
			fmt.Printf("Skipping question. Has already been asked!")
			break
		}

		quiz.QuizStatus = enums.StatusQuizTime
		fmt.Println("Quiz Status updated to: ", enums.StatusQuizTime)

		quizMsg := fmt.Sprintf("QuestionID: %s, Question: %s, Options: %s, CorrectAnswer: %s",
			question.Id, question.GetQuestion(), question.GetOptions(), question.GetCorrectAnswer())
		statusChannel <- enums.StatusQuizTime
		question.IsAskedStatus = true
		question.LastAskedTime = time.Now()
		broadcastToParticipants(quizID, quizMsg)

		for status := range statusChannel {
			if status == enums.StatusQuizTimeEnded {
				fmt.Println("Quiz Time Ended!")
				quiz.QuizStatus = enums.StatusEvaluation
				fmt.Println("Quiz Status updated to: ", enums.StatusEvaluation)
				break
			} else if status == enums.StatusStopped {
				return

			}
		}
	}
}

func timerRoutine(wg *sync.WaitGroup, statusChannel chan enums.QuizStatus) {
	TAG := "TIMER_ROUTINE: "

	for status := range statusChannel {
		fmt.Println(TAG, "RECEIVED status: ", status)

		if status == enums.StatusQuizTime {
			answerTimer := time.NewTimer(answerTimeout)
			fmt.Println(TAG, "Answer timer started!")
			<-answerTimer.C
			statusChannel <- enums.StatusQuizTimeEnded
			fmt.Println(TAG, "Answer timer ENDED!")

		} else if status == enums.StatusEnded || status == enums.StatusStopped {
			fmt.Println(TAG, "Shutting down the timer GoRoutine!")
			defer wg.Done()
			return
		}
	}

}

func handleAnswer(quizID string, questionID string, answer string, timeReceived time.Time, user *dataTypes.User) {
	TAG := "ANSWER_HANDLER: "
	quiz, quizOk := quizzes[quizID]
	if !quizOk {
		fmt.Println(TAG, "Quiz not found.")
		return
	}
	question, questionOk := quiz.Questions[questionID]
	if !questionOk {
		fmt.Println(TAG, "Questions contains:", quiz.Questions)
		fmt.Println(TAG, "Question not found.")
		return
	}
	if !question.IsAskedStatus {
		fmt.Println(TAG, "Question has not been asked.")
		return
	}
	if answer != question.GetCorrectAnswer() {
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