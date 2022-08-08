package server

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//在线用户的列表
	OnlineMap map[string]*User
	//因为这个map是全局的  所以需要加锁
	mapLock sync.RWMutex

	//消息广播的channel
	Message chan string
}

//NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

//ListenMessage 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
func (server *Server) ListenMessage() {
	for {
		msg := <-server.Message

		//将msg发送给全部的在线用户
		server.mapLock.Lock()
		for _, v := range server.OnlineMap {
			v.C <- msg
		}
		server.mapLock.Unlock()
	}
}

//BroadCast 广播用户上线信息的方法
func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	server.Message <- sendMsg
}

//Handler ...当前连接的业务
func (server *Server) Handler(conn net.Conn) { //fmt.Println("连接建立成功")
	//创建一个User
	user := NewUser(conn, server)
	//用户上线广播
	user.Online()

	//监听当前用户是否活跃的channel
	isLive := make(chan bool)

	//接收客户端发来的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				//用户下线广播
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read error:", err)
				return
			}
			//提取用户的消息(去除'\n')
			msg := string(buf)
			//用户针对msg进行消息处理
			user.DoMessage(msg)

			//用户的任意消息 代表用户是活跃的
			isLive <- true
		}
	}()

	//当前的Handler阻塞  监听用户是否超时 超时则回收资源  否则不断重置定时器
	for {
		select {
		case <-isLive:
			//当前用户是活跃的应该重置定时器
			//不做任何事情，为了激活select，更新下面的定时器
		case <-time.After(time.Second * 300): //如果channel中有数据可取 则进入
			//进入 则代表有用户已超时
			//将该用户强制关闭
			timeout := "用户:[" + user.Name + "]连接已超时,强制踢出!!"
			user.SendMsg(timeout)

			//销毁用户的资源
			close(user.C)
			//关闭连接
			conn.Close()
			//退出当前的 go Handler
			return //runtime.Goexit()
		}
	}
}

//Start 启动服务器的接口
func (server *Server) Start() {
	//socket listen		监听socket连接
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("new.Listen error:", err)
		return
	}
	//close listen socket
	defer listener.Close()

	//启动监听Message的goroutine
	go server.ListenMessage()

	for {
		//accept	建立连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept error:", err)
			return
		}
		//do handler
		go server.Handler(conn)
	}
}
