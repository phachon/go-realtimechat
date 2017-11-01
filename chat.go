package main

import (
	"github.com/gorilla/websocket"
	"sync"
	"fmt"
	"os"
)

// Group struct
type Group struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

// User struct
type User struct {
	Id string  `json:"id"`
	GroupId string `json:"group_id"`
	Name string `json:"name"`
	Ws *websocket.Conn `json:"ws"`
}

// Messages struct
type Messages struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
	UserId  string `json:"user_id"`
}

// Chat struct
type Chat struct {
	lock sync.Mutex
	groups []*Group
	users []*User
	MessageChan chan Messages
	signalChan  chan string
}

func NewChat() *Chat {
	return &Chat{
		groups: []*Group{},
		users: []*User{},
		MessageChan: make(chan Messages, 10),
		signalChan: make(chan string, 1),
	}
}

// add  group
func (chat *Chat) AddGroup(id string, name string) {
	chat.lock.Lock()
	defer chat.lock.Unlock()

	//check group is exists
	for _, chatGroup := range chat.groups {
		if chatGroup.Id == id {
			printError("add group error: id already exists!");
		}
		if chatGroup.Name == name {
			printError("add group error: name already exists!");
		}
	}
	group := &Group{
		Id:id,
		Name:name,
	}

	chat.groups = append(chat.groups, group)
}

// delete group and delete user under group
func (chat *Chat) DeleteGroup(id string) {
	chat.lock.Lock()
	defer chat.lock.Unlock()

	groups := []*Group{}
	for _, chatGroup := range chat.groups {
		if chatGroup.Id == id {
			continue
		}
		groups = append(groups, chatGroup)
	}
	chat.groups = groups

	users := []*User{}
	//clear group user and close ws
	for _, chatUser := range chat.users {
		if chatUser.GroupId == id {
			chatUser.Ws.Close()
			continue
		}
		users = append(users, chatUser)
	}

	chat.users = users
}

// add user
func (chat *Chat) AddUser(id string, groupId string, name string, ws *websocket.Conn) {
	chat.lock.Lock()
	defer chat.lock.Unlock()

	//check group is exists
	groupExist := false
	for _, chatGroup := range chat.groups {
		if chatGroup.Id == groupId {
			groupExist = true
			break
		}
	}
	if !groupExist {
		printError("add user error: group_id is not exists!");
	}

	//check user id name is exists
	for _, chatUser := range chat.users {
		if chatUser.Id == id {
			printError("add user error: id already exists!");
		}
		if chatUser.Name == name {
			printError("add user error: name already exists!");
		}
	}

	user := &User {
		Id:id,
		GroupId: groupId,
		Name: name,
		Ws: ws,
	}

	chat.users = append(chat.users, user)
}

// delete user and close user ws
func (chat *Chat) DeleteUser(id string)  {
	chat.lock.Lock()
	defer chat.lock.Unlock()

	users := []*User{}
	for _, chatUser := range chat.users {
		if chatUser.Id == id {
			//close ws
			chatUser.Ws.Close()
			continue
		}
		users = append(users, chatUser)
	}

	chat.users = users
}

// get user list
func (chat *Chat) Users() []*User {
	return chat.users
}

// get group list
func (chat *Chat) Groups() []*Group {
	return chat.groups
}

func printError(message string) {
	fmt.Println(message)
	os.Exit(0)
}
