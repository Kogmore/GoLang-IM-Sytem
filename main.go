package main

func main() {
	//发送地址给服务器
	server := NewServer("127.0.0.1", 8888)
	//启动服务器
	server.Start()
}
