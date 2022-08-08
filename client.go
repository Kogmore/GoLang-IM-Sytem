package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
	Flag       int //当前Client的模式
}

//NewClient 创建客户端
func NewClient(serverIp string, serverPort int) *Client {
	//创建一个客户端对象 client
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		Flag:       0,
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
	fmt.Println("4、退出功能")

	var cFlag int
	//获取用户输入的模式 Scanln() 括号里面返回的是个指针
	_, err := fmt.Scanln(&cFlag)
	if err != nil {
		fmt.Println("fmt.Scanln() error:", err)
		return false
	}
	if cFlag > 0 && cFlag < 5 {
		client.Flag = cFlag
		return true
	} else {
		fmt.Println(">>>>请输入合法的菜单指令<<<<")
		return false
	}
}

//Run 菜单业务处理
func (client *Client) Run() {
	for client.Flag == 0 {
		for client.Menu() == false {

		}
		switch client.Flag {
		case 1: //1、公聊模式
			fmt.Println("公聊模式选择成功。。。")
			break
		case 2: //2、私聊模式
			fmt.Println("私聊模式选择成功。。。")
			break
		case 3: //3、更新用户名
			fmt.Println("更新用户名选择成功。。。")
			break
		case 4: //4、退出功能
			fmt.Println("退出成功。。。")
			break
		}
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

	//启动客户端的业务
	client.Run()
}
