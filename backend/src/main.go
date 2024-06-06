package main

import (
	"fmt"
	"io"
	"net/http"
	"quizzy_game/handlers/sessionHandler"
	"quizzy_game/network"
)

func getFrontPage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is the quiz backend!\n")
}

func main() {
	http.HandleFunc("/", getFrontPage)
	http.HandleFunc("/questions", network.GetQuestionsWeb)
	http.HandleFunc("/categories", network.GetCategoriesWeb)
	http.HandleFunc("/quiz_options", network.GetQuizOptions)
	http.HandleFunc("/ws", sessionHandler.WsEndpoint)
	http.HandleFunc("/ws/", sessionHandler.WsEndpoint)

	http.ListenAndServe(":8000", nil)

}
