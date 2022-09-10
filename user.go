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
	//Start listen channel's message
	go user.ListenChannel()

	return user
}

func (u *User) ListenChannel() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n")) //Response message to client
	}
}
