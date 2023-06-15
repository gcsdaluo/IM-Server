package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 创建在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server

}

// ListenMessager 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		//将msg发送给全部的在线User
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// BroadCast 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

// Handler 处理函数
func (this *Server) Handler(conn net.Conn) {
	// .. 当前链接业务
	//fmt.Println("链接建立成功")

	user := NewUser(conn, this)

	user.Online()

	// 接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		// 持续捕获消息
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			// 提取用户的消息(去掉'\n')
			msg := string(buf[:n-1])

			//用户针对msg进行消息的处理
			user.DoMessage(msg)
		}
	}()

	// 当前handler阻塞
	select {}
}

// Start 启动服务器的接口
func (this *Server) Start() {
	// so
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}

	//close listen socket
	defer listener.Close()

	// 启动监听Message的goroutine
	go this.ListenMessager()

	//accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(fmt.Println("listen accept err:", err))
			continue
		}

		//do handler
		go this.Handler(conn)
	}
}
