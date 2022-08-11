package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
	Flag       string //当前Client的模式
}

//DealResponse 处理server响应的消息 直接显示到标准输出即可
func (client *Client) DealResponse() {
	//读取server端的响应数据并输出 第一个返回值是 读取的字节数
	//一旦client.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	_, err := io.Copy(os.Stdout, client.Conn)
	if err != nil {
		fmt.Println("DealResponse() | io.Copy() error:", err)
		return
	}
	//for {
	//	buf := make([]byte, 4096)
	//	_, err := client.Conn.Read(buf)
	//	if err != nil {
	//		fmt.Println("DealResponse() | client.Conn.Read() error:", err)
	//	}
	//	fmt.Println(buf)
	//}
}

//NewClient 创建客户端
func NewClient(serverIp string, serverPort int) *Client {
	//创建一个客户端对象 client
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		Flag:       "888", //判断用户选择的功能
	}
	//连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial() error:", err)
		return nil
	}
	client.Conn = conn
	//返回客户端对象
	return client
}

//Menu 绑定菜单的方法
func (client *Client) Menu() bool {
	fmt.Println("1、公聊模式")
	fmt.Println("2、私聊模式")
	fmt.Println("3、更新用户名")
	fmt.Println("0、退出功能")

	var cFlag string
	//获取用户输入的模式 Scan() 括号里面返回的是个指针
	n, err := fmt.Scan(&cFlag)
	if n == 0 {
		fmt.Println(">>>>请选择模式:")
		return false
	}
	if err != nil {
		fmt.Println("Menu() | fmt.Scan() error:", err)
		return false
	}
	if cFlag >= "0" && cFlag <= "3" {
		client.Flag = cFlag
		return true
	} else {
		fmt.Println(">>>>请输入合法的菜单指令<<<<")
		return false
	}
}

//Run 菜单业务处理
func (client *Client) Run() {
	for client.Flag != "0" {
		for client.Menu() != true {

		}
		switch client.Flag {
		case "1": //1、公聊模式
			client.PublicChat()
			break
		case "2": //2、私聊模式
			client.PrivateChat()
			break
		case "3": //3、更新用户名
			client.UpdateName()
			break
		case "0": //0、退出功能
			fmt.Println("退出成功。。。")
			break
		}
	}
}

//PublicChat 菜单功能之公聊模式
func (client *Client) PublicChat() {
	var chatMsg string
	//提示用户输入消息
	fmt.Println("PublicChat()1>>>>请输入聊天内容:")
	fmt.Println("PublicChat()1>>>>输入exit退出聊天.")
	//获取用户输入的消息
	n, err := fmt.Scan(&chatMsg)
	if n == 0 {
		fmt.Println("PublicChat()1>>>>不可输入内容为空的消息!!!")
		return
	}
	if err != nil {
		fmt.Println("PublicChat()1 | fmt.Scan() error:", err)
		return
	}
	for chatMsg != "exit" {
		//拼接server处理
		sendMsg := chatMsg + "\n"
		//把消息发送给服务器
		_, err = client.Conn.Write([]byte(sendMsg))
		if err != nil {
			fmt.Println("PublicChat() | client.Conn.Write() error:", err)
			break
		}
		//提示用户输入消息
		fmt.Println("PublicChat()2>>>>请输入聊天内容:")
		fmt.Println("PublicChat()2>>>>输入exit退出聊天.")
		//获取用户输入的消息
		n, err = fmt.Scan(&chatMsg)
		if n == 0 {
			fmt.Println("PublicChat()2>>>>不可输入内容为空的消息!!!")
			break
		}
		if err != nil {
			fmt.Println("PublicChat()2 | fmt.Scan() error:", err)
			break
		}
	}
}

//SelectUsers 菜单功能之私聊模式（查询所有在线用户）
func (client *Client) SelectUsers() {
	sendMsg := "all\n"
	//把消息发送给服务器
	_, err := client.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("SelectUsers() | client.Conn.Write() error:", err)
		return
	}
}

//PrivateChat 菜单功能之私聊模式
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	//查询所有在线用户
	client.SelectUsers()
	//提示用户输入消息
	fmt.Println("PrivateChat()1>>>>请输入聊天对象[用户名]")
	fmt.Println("PrivateChat()1>>>>输入exit退出聊天.")
	//获取用户输入的用户名
	n, err := fmt.Scan(&remoteName)
	if n == 0 {
		fmt.Println("PrivateChat()1>>>>不可输入为空的用户名!!!")
		return
	}
	if err != nil {
		fmt.Println("PrivateChat()1 | fmt.Scan() error:", err)
		return
	}
	for remoteName != "exit" {
		//提示用户输入消息
		fmt.Println("PrivateChat()1>>>>请输入消息内容")
		//获取用户输入的消息
		n, err = fmt.Scan(&chatMsg)
		if n == 0 {
			fmt.Println("PrivateChat()1>>>>不可输入内容为空的消息!!!")
			break
		}
		if err != nil {
			fmt.Println("PrivateChat()2 | fmt.Scan() error:", err)
			break
		}
		for chatMsg != "exit" {
			//拼接server处理
			sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
			//把消息发送给服务器
			_, err = client.Conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("PrivateChat() | client.Conn.Write() error:", err)
				return
			}

			chatMsg = ""
			//提示用户输入消息
			fmt.Println("PrivateChat()2>>>>请输入消息内容")
			//获取用户输入的消息
			n, err = fmt.Scan(&chatMsg)
			if n == 0 {
				fmt.Println("PrivateChat()2>>>>不可输入内容为空的消息!!!")
				break
			}
			if err != nil {
				fmt.Println("PrivateChat()3 | fmt.Scan() error:", err)
				return
			}
		}
		//查询所有在线用户
		client.SelectUsers()
		//提示用户输入消息
		fmt.Println("PrivateChat()2>>>>请输入聊天对象[用户名]")
		fmt.Println("PrivateChat()2>>>>输入exit退出聊天.")
		//获取用户输入的用户名
		n, err = fmt.Scan(&remoteName)
		if n == 0 {
			fmt.Println("PrivateChat2>>>>不可输入为空的用户名!!!")
			break
		}
		if err != nil {
			fmt.Println("PrivateChat()4 | fmt.Scan() error:", err)
			break
		}
	}
}

//UpdateName 菜单功能之更新用户名
func (client *Client) UpdateName() {
	fmt.Println(">>>>请输入用户名:")
	//获取用户输入的消息
	n, err := fmt.Scan(&client.Name)
	if n == 0 {
		fmt.Println(">>>>不可输入为空的用户名!!!")
		return
	}
	//if client.Name == "" {
	//	fmt.Println("不可输入为空的用户名!!!")
	//	return
	//}
	if err != nil {
		fmt.Println("UpdateName() | fmt.Scan() error:", err)
		return
	}

	//拼接server处理 更新用户名业务的指令
	sendMsg := "rename|" + client.Name + "\n"
	_, err = client.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("UpdateName() | client.Conn.Write() error:", err)
		return
	}
}

var serIp string
var serPort int

//./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serPort, "port", 8888, "设置服务器端口(默认是8888)")
}

func main() {
	//命令行解析
	flag.Parse()

	//创建客户端 连接服务器
	client := NewClient(serIp, serPort)
	if client == nil {
		fmt.Println("NewClient()》》》》》》》》》  连接服务器失败!!!")
		return
	}
	fmt.Println("NewClient()》》》》》》》》》  连接服务器成功!!!")

	//启动一个监听server端响应的goroutine
	go client.DealResponse()
	//启动客户端的业务
	client.Run()
}
