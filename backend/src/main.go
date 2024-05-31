package main

import (
	"fmt"
	"io"
	"net/http"
	"quizzy_game/api"
	"quizzy_game/clientManagement"
)

func getFrontPage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}

func main() {
	http.HandleFunc("/", getFrontPage)
	http.HandleFunc("/questions", api.GetQuestionsWeb)
	http.HandleFunc("/categories", api.GetCategoriesWeb)
	http.HandleFunc("/ws", clientManagement.WsEndpoint)
	http.HandleFunc("/ws/", clientManagement.WsEndpoint)

	http.ListenAndServe(":8000", nil)

}
