package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	Conn net.Conn

	Server *Server
}

//NewUser 创建一个用户的API
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

//SendMsg 给当前User对应的客户端发送消息
func (user *User) SendMsg(msg string) {
	_, err := user.Conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("user SendMsg error:", err)
		return
	}
}

//DoMessage 用户处理消息的业务
func (user *User) DoMessage(msg string) {
	//查询当前所有在线在线用户信息
	if msg == "who" {
		user.Server.mapLock.Lock()
		for _, v := range user.Server.OnlineMap {
			onlineMsg := "[" + v.Addr + "]" + v.Name + ":在线...\n"
			user.SendMsg(onlineMsg)
		}
		user.Server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" { //修改用户名校验
		//消息格式 rename|张三
		reName := strings.Split(msg, "|")[1]
		//判断当前name是否存在
		_, ok := user.Server.OnlineMap[reName]
		if ok {
			msgName := "用户名:[" + reName + "]已被使用!!!"
			user.SendMsg(msgName)
		} else { //修改用户名操作
			//先要删除原来的数据 然后添加一个新数据  否则会出现两个数据
			user.Server.mapLock.Lock()
			delete(user.Server.OnlineMap, user.Name)
			user.Server.OnlineMap[reName] = user
			user.Server.mapLock.Unlock()

			user.Name = reName
			//告诉用户，用户名修改成功
			msgName := "用户名:[" + user.Name + "]修改成功!!!"
			user.SendMsg(msgName)
		}
	} else if len(msg) > 4 && msg[:3] == "to|" { //to| 协议  表示跟某个用户私聊
		//消息格式 to|张三|消息内容
		//1、获取用户姓名
		toName := strings.Split(msg, "|")[1]
		if toName == "" {
			user.SendMsg("消息格式不正确,请输入\"to|张三|消息内容\"格式!!!\n")
			return
		}
		if toName == user.Name {
			user.SendMsg("不能给自己发送消息,请重新输入!!!\n")
			return
		}
		//2、根据用户名获取用户的User对象
		nUser, ok := user.Server.OnlineMap[toName]
		if !ok {
			user.SendMsg("当前用户名不存在,请输入正确的用户名!!!")
			return
		}
		//3、获取消息内容，通过对方的User对象将消息	内容发送过去
		toContent := strings.Split(msg, "|")[2]
		if toContent == "" {
			user.SendMsg("不能发送为空的消息!!!\n")
			return
		}
		nUser.SendMsg(user.Name + "对您说:" + toContent)
	} else {
		//将得到的消息进行广播
		user.Server.BroadCast(user, msg)
	}
}

//ListenMessage 监听当前User channel的方法，一旦有消息，就直接发送给对端 客户端
func (user *User) ListenMessage() {
	for {
		msg, ok := <-user.C
		if !ok {
			fmt.Printf("用户%v连接已断开\n", user.Name)
			return
		}
		_, err := user.Conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("user channel error:", err)
			return
		}
	}
}
