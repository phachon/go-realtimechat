package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"fmt"
	"encoding/json"
)

var (
	// websocket upgrader
	upgrader = websocket.Upgrader{}
	// chat
	chat = NewChat()
)

//初始化配置
func init()  {
	initConfig();
}

func main() {

	// chat html
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	// user list
	http.HandleFunc("/users", userList)
	// group list
	http.HandleFunc("/groups", groupList)
	// handle websocket
	http.HandleFunc("/ws", handleConnections)

	// start read message
	go handleMessages()

	log.Println("http server started on : 8087")
	err := http.ListenAndServe(":8087", nil)
	if(err != nil) {
		log.Println("http server error: "+ err.Error())
	}
}

// all user
//return json
func userList(w http.ResponseWriter, req *http.Request)  {
	users := chat.Users()
	v, _ := json.Marshal(users)
	w.Write(v)
}

// all group
// return json
func groupList(w http.ResponseWriter, req *http.Request)  {
	groups := chat.Groups()
	v, _ := json.Marshal(groups)
	w.Write(v)
}

// handle ws conn
func handleConnections(w http.ResponseWriter, req *http.Request) {

	userId := req.FormValue("user_id")
	roomId := req.FormValue("room_id")
	if(userId == "") {
		return
	}
	if(roomId == "") {
		return
	}
	log.Println("room_id : "+ roomId +" user_id: "+ userId +" come in!")

	//将 get 请求升级为 websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if(err != nil) {
		fmt.Println("请求出错:" + err.Error())
		return
	}
	defer ws.Close()

	// add group
	chat.AddGroup(roomId, "default_"+roomId)
	// add user
	chat.AddUser(userId, roomId, "user_"+userId, ws)

	for {
		// read ws json data
		var msg Messages
		err := ws.ReadJSON(&msg)
		if(err != nil) {
			log.Println(err.Error())
			log.Println("room_id : "+ roomId +" user_id: "+ userId +" exit!")
			chat.DeleteUser(userId)
			break
		}

		// send message to messageChan
		chat.MessageChan <- msg
	}
}

// handle read message
func handleMessages()  {
	for {
		// read message from chat messageChan
		var msg Messages
		msg = <- chat.MessageChan

		// 暂时先广播给所有房间的所有人
		for _, user := range chat.Users() {
			ws := user.Ws
			err := ws.WriteJSON(msg)
			if(err != nil) {
				chat.DeleteUser(user.Id)
				log.Printf("error: %v", err)
			}
		}
	}
}

