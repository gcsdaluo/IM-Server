package main

import "net"

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

// DoMessage 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this, msg)
}

// ListenMessage 监听当前User channel的方法，一旦有消息，就直接发送给客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
