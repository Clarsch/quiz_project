package requestcommand

type QuizCommands string

const (
	CreateQuiz     = "create"
	JoinQuiz       = "join"
	LeaveQuiz      = "leave"
	StartQuiz      = "start"
	StopQuiz       = "stop"
	ResetQuiz      = "reset"
	AnswerQuestion = "answer"
)
