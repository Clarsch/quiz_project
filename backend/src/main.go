package main

import (
	"fmt"
	"io"
	"net/http"
	"quizzy_game/api"
	"quizzy_game/handlers/sessionHandler"
)

func getFrontPage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is the quiz backend!\n")
}

func main() {
	http.HandleFunc("/", getFrontPage)
	http.HandleFunc("/questions", api.GetQuestionsWeb)
	http.HandleFunc("/categories", api.GetCategoriesWeb)
	http.HandleFunc("/quiz_options", api.GetQuizOptions)
	http.HandleFunc("/ws", sessionHandler.WsEndpoint)
	http.HandleFunc("/ws/", sessionHandler.WsEndpoint)

	http.ListenAndServe(":8000", nil)

}
