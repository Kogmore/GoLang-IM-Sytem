package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	Conn net.Conn

	Server *Server
}

//NewUser 创建一个用的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		Conn:   conn,
		Server: server,
	}

	//启动监听当前user  channel的goroutine
	go user.ListenMessage()

	return user
}

//Online 用户的上线业务
func (user *User) Online() {
	//用户上线 将用户添加到onlineMap里面
	user.Server.mapLock.Lock()
	user.Server.OnlineMap[user.Name] = user
	user.Server.mapLock.Unlock()

	//广播当前用户上线的消息
	user.Server.BroadCast(user, "已上线!!!")
}

//Offline 用户的下线业务
func (user *User) Offline() {
	//用户下线 将用户从onlineMap中删除
	user.Server.mapLock.Lock()
	delete(user.Server.OnlineMap, user.Name)
	user.Server.mapLock.Unlock()

	//广播当前用户下线的消息
	user.Server.BroadCast(user, "已下线!!!")
}

//DoMessage 用户处理消息的业务
func (user *User) DoMessage(msg string) {
	//将得到的消息进行广播
	user.Server.BroadCast(user, msg)
}

//ListenMessage 监听当前User channel的方法，一旦有消息，就直接发送给对端 客户端
func (user *User) ListenMessage() {
	for {
		msg := <-user.C
		_, err := user.Conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("user channel error", err)
			return
		}
	}
}
