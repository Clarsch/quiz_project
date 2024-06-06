package sessionHandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"quizzy_game/handlers/quizHandler"
	"quizzy_game/internal/dataTypes"
	frontdto "quizzy_game/internal/dto/frontDTO"
	requestcommand "quizzy_game/internal/enums/requestCommands"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var users = make(map[string]*dataTypes.User)

func reader(conn *websocket.Conn, user *dataTypes.User) {
	var mutex sync.Mutex // Create a mutex to synchronize writes to the WebSocket connection

	go func() {
		for {
			data, more := <-user.MsgChannel
			if more {
				mutex.Lock()
				err := conn.WriteMessage(websocket.TextMessage, []byte(data))
				mutex.Unlock()
				if err != nil {
					fmt.Println("Error writing to WebSocket:", err)
					break
				}
			} else {
				break
			}
		}
	}()

	user.MsgChannel <- fmt.Sprintf("UserID: %s", user.Id)

	for {
		messageType, byteMsg, cErr := conn.ReadMessage()
		if cErr != nil {
			log.Println(cErr)
			close(user.MsgChannel)
			return
		}
		// Process the JSON request
		processJSONRequest(conn, user, messageType, byteMsg)
	}
}

func processJSONRequest(conn *websocket.Conn, user *dataTypes.User, messageType int, byteMsg []byte) {
	var request frontdto.Request
	mErr := json.Unmarshal(byteMsg, &request)
	if mErr != nil {
		fmt.Println("Error:", mErr)
		user.MsgChannel <- fmt.Sprintf("Invalid input! Error: %s\n", mErr)
		return
	}
	request.ReceivedTime = time.Now()

	if request.RequestType == requestcommand.SetUsername {
		jsonData, _ := json.Marshal(request.Data)
		var setUsernameData frontdto.SetUsernameRequestDTO
		err := json.Unmarshal(jsonData, &setUsernameData)
		if err != nil {
			fmt.Println("Error:", err)
			user.MsgChannel <- fmt.Sprintf("Invalid input! Error: %s\n", err)
			user.MsgChannel <- fmt.Sprintf("Example of correct json input: %s\n", setUsernameData.GetExample())
			return
		}
		user.Name = setUsernameData.Username
		user.MsgChannel <- "Your Username has been updated to: " + setUsernameData.Username
	} else {
		// Adding data to the channel should also be synchronized
		var mutex sync.Mutex
		mutex.Lock()
		go quizHandler.HandleQuizRequest(request, user)
		cErr := conn.WriteMessage(messageType, byteMsg)
		mutex.Unlock()

		if cErr != nil {
			log.Println(cErr)
			close(user.MsgChannel)
			return
		}
	}

}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	parts := strings.Split(r.URL.Path, "/")
	var id string
	if len(parts) > 2 {
		id = strings.TrimSpace(parts[2])
	}

	var user *dataTypes.User

	if id != "" {
		fmt.Println("User Connected with UserID: ", id)
		user = users[id]
		if user != nil {
			user.MsgChannel = make(chan string, 1)
		}
	}

	if user == nil {
		fmt.Println("User Connected without UserID or  with invalid UserID. Creating new User!")
		responseChan := make(chan string, 1)
		newUserId := uuid.NewString()
		newUser := dataTypes.User{
			Id:         newUserId,
			Name:       "user" + newUserId[0:4],
			MsgChannel: responseChan,
		}
		users[newUser.Id] = &newUser
		user = &newUser
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection

	reader(ws, user)
}
