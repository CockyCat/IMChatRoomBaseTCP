package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string), //Message Channel for broadcast everybody
		conn: conn,
	}
	//启动用户消息的监听，后台执行
	go user.ListenChannel()

	return user
}

//监听发送至user结构体的channel
func (u *User) ListenChannel() {
	for {
		//消费来自ListenMsg中生产的消息
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n")) //Response message to client
	}
}
