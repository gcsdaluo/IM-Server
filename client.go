package main

import (
	"fmt"
	"net"
)

// net.Conn 是 Go 语言中用于表示通用网络连接的接口类型。
// 它提供了与网络连接相关的方法，例如读取和写入数据、关闭连接等。
// 通过 net.Conn 类型的对象，可以与远程服务器进行通信，发送和接收数据。

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn // 与服务器建立的网络连接
}

func NewClient(serverIp string, serverPort int) *Client {
	// 1、创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	// 2、链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn

	// 3、返回对象
	return client
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>> 链接服务器失败...")
		return
	}

	fmt.Println(">>>>> 链接服务器成功...")

	// 启动客户端的业务
	select {}
}
