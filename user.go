package main

import (
	"net"
	"strings"
)

// User 后端表示用户的结构体封装
type User struct {
	Name string      //当前用户名字-默认地址
	Addr string      //当前用户名字-默认地址
	C    chan string //当前是否有数据回写给客户端
	conn net.Conn    //socket通信的链接

	server *Server // 和server类进行链接，可以使用广播等功能
}

// NewUser 创建一个用户API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	// 启动监听当前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

// Online 用户上线业务
func (this *User) Online() {

	//用户上线，将用户加入到onlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//广播当前用户上线消息
	this.server.BroadCast(this, "已上线")
}

// Offline 用户下线业务
func (this *User) Offline() {

	//用户下线，将用户从onlineMap中删除，delete是从map里面找到键Name进行对应删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	//广播当前用户上线消息
	this.server.BroadCast(this, "下线")
}

// SendMsg 给当前User对应的客户端发送消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// DoMessage 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前在线用户都有哪些？

		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 消息格式: rename|张三
		newName := strings.Split(msg, "|")[1]

		// 判断name是否存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("当前用户名被使用\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			// 将新的用户名赋给方法结构体
			this.Name = newName
			this.SendMsg("您已经更新用户名:" + this.Name + "\n")
		}
	} else {
		this.server.BroadCast(this, msg)
	}
}

// ListenMessage 监听当前User channel的方法，一旦有消息，就直接发送给客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
