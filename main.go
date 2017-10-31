package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"fmt"
)

//Room struct
type Room struct {
	Id      string
	Users   []User
}

//User struct
type User struct {
	Id      string
	Ws      *websocket.Conn
}

//Message struct
type Messages struct {
	Type    int    `json:"type"`
	Message    string `json:"message"`
	UserId    string `json:"user_id"`
}

var (
	users = []User{}
	rooms = map[string][]User{}
	//通道
	broadcast = make(chan Messages)
	//升级为websocket
	upgrader = websocket.Upgrader{}
)

//初始化配置
func init()  {
	initConfig();
}

func main() {

	// 静态文件服务
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	//处理 websocket 连接
	http.HandleFunc("/ws", handleConnections)

	//开始监听聊天信息
	go handleMessages()

	log.Println("http server started on : 8087")
	err := http.ListenAndServe(":8087", nil)
	if(err != nil) {
		log.Println("http server error: "+ err.Error())
	}
}

//处理请求
func handleConnections(w http.ResponseWriter, r *http.Request) {

	userId := r.FormValue("user_id")
	roomId := r.FormValue("room_id")
	if(userId == "") {
		return
	}
	if(roomId == "") {
		return
	}
	fmt.Println("room_id : "+ roomId +" user_id: "+ userId +" come in!")

	//将 get 请求升级为 websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if(err != nil) {
		fmt.Println("请求出错:" + err.Error())
		return
	}
	defer ws.Close()

	//开始注册
	user := &User{
		Id: userId,
		Ws: ws,
	}
	users = append(users, *user)
	rooms[roomId] = append(rooms[roomId], *user)

	//开始处理消息
	for {
		//读取 json 消息
		var msg Messages
		err := ws.ReadJSON(&msg)
		if(err != nil) {
			log.Println(err.Error())
		}
		//发送消息到channel
		broadcast <- msg
	}

}


//处理消息
func handleMessages()  {
	for {
		//从 channel 读取消息
		//读取 json 消息
		var msg Messages
		msg = <- broadcast

		//广播给所有房间的所有人
		for _, users := range rooms {
			//roomId := room.Id
			for _, user := range users {
				ws := user.Ws
				err := ws.WriteJSON(msg)
				if(err != nil) {
					log.Printf("error: %v", err)
					ws.Close()
				}
			}
		}
	}
}

