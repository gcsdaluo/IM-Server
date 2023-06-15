package main

import "net"

// User 后端表示用户的结构体封装
type User struct {
	Name string      //当前用户名字-默认地址
	Addr string      //当前用户名字-默认地址
	C    chan string //当前是否有数据回写给客户端
	conn net.Conn    //socket通信的链接
}

// NewUser 创建一个用户API
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// 启动监听当前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

// ListenMessage 监听当前User channel的方法，一旦有消息，就直接发送给客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
