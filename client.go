package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

// net.Conn 是 Go 语言中用于表示通用网络连接的接口类型。
// 它提供了与网络连接相关的方法，例如读取和写入数据、关闭连接等。
// 通过 net.Conn 类型的对象，可以与远程服务器进行通信，发送和接收数据。

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn // 与服务器建立的网络连接
	flag       int      // 当前client的模式
}

func NewClient(serverIp string, serverPort int) *Client {
	// 1、创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
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

// Dealresponse 处理server回应的消息，直接显示到标准输出即可
func (client *Client) Dealresponse() {
	// 一旦client.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>请输入合法范围内的数字<<<<<<<<")
		return false
	}
}

// SearchUsers 查询在线用户
func (client *Client) SearchUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

// PrivateChat 私聊模式
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SearchUsers()
	fmt.Println(">>>>>请输入聊天对象[用户名],exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>请输入消息内容,exit退出")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			//发给服务器
			// 消息不为空则发送
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"

				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>请输入聊天内容,exit退出")
			fmt.Scanln(&chatMsg)
		}

		client.SearchUsers()
		fmt.Println(">>>>>请输入聊天对象[用户名],exit退出")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) PublicChat() {
	// 提示用户输入消息
	var chatMsg string

	fmt.Println(">>>>请输入聊天内容,exit退出.")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		//发给服务器

		// 消息不为空则发送
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>请输入聊天内容,exit退出")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) UpdateName() bool {

	fmt.Println(">>>>请输入用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}

		// 根据不同的模式处理不同的业务
		switch client.flag {
		case 1:
			//公聊模式
			client.PublicChat()
		case 2:
			//私聊模式
			client.PrivateChat()
		case 3:
			//更新用户名
			client.UpdateName()

		}
	}
}

// 两个全局形参绑定到flag包中
var serverIp string
var serverPort int

// 每个go文件都有一个init函数，这个函数会在main函数前执行
// 初始化命令行参数
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}

func main() {
	// 命令行解析，可以指定客户端
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>> 链接服务器失败...")
		return
	}

	// 单独开启一个goroutine去处理server的
	go client.Dealresponse()

	fmt.Println(">>>>> 链接服务器成功...")

	// 启动客户端的业务
	client.Run()
}
