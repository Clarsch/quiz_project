package requestcommand

type QuizCommands string

const (
	CreateQuiz     = "create"
	JoinQuiz       = "join"
	LeaveQuiz      = "leave"
	StartQuiz      = "start"
	ResetQuiz      = "reset"
	AnswerQuestion = "answer"
)
