package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server

}

// Handler 处理函数
func (this *Server) Handler(conn net.Conn) {
	// .. 当前链接业务
	fmt.Println("链接建立成功")
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
