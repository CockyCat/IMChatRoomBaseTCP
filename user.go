package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	Server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string), //Message Channel for broadcast everybody
		conn:   conn,
		Server: server,
	}
	//启动用户消息的监听，后台执行
	go user.ListenChannel()

	return user
}

func (user *User) Online() {
	//thread safety
	user.Server.mapLock.Lock()
	user.Server.OnlineUsers[user.Name] = user
	user.Server.mapLock.Unlock()

	//user online to brodcast
	user.Server.Broadcast(user, "has been online!")
}

func (user *User) Offline() {
	//thread safety
	user.Server.mapLock.Lock()

	//从map中删除此用户
	delete(user.Server.OnlineUsers, user.Name)

	user.Server.mapLock.Unlock()

	//user online to brodcast
	user.Server.Broadcast(user, "has been offline!")
}

func (user *User) DoMessage(msg string) {
	user.Server.Broadcast(user, msg)
}

//监听发送至user结构体的channel
func (u *User) ListenChannel() {
	for {
		//消费来自ListenMsg中生产的消息
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n")) //Response message to client
	}
}
