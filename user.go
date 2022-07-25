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
}

//NewUser 创建一个用的API
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		Conn: conn,
	}

	//启动监听当前user  channel的goroutine
	go user.ListenMessage()

	return user
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
