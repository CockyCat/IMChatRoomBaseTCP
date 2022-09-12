package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	OnlineUsers map[string]*User
	mapLock     sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:          ip,
		Port:        port,
		OnlineUsers: make(map[string]*User),
		Message:     make(chan string),
	}
	return server
}

func (s *Server) Start() {
	//socket listen
	lister, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		log.Println("net listen error:", err)
		return
	}
	defer lister.Close()

	go s.ListenMsg()

	//accept
	for {
		conn, err := lister.Accept()
		if err != nil {
			log.Println("net listen error:", err)
			continue
		}
		//do handler
		go s.Handler(conn)
	}

	//close listen
}

func (s *Server) Handler(conn net.Conn) {
	log.Println("connection establish success")

	//Add user to map
	user := NewUser(conn)

	//thread safety
	s.mapLock.Lock()
	s.OnlineUsers[user.Name] = user
	s.mapLock.Unlock()

	//user online to brodcast
	s.Broadcast(user, "has been online!")

	//Accept message from users(clients)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			//Read data is zero means this user have been offline.
			if n == 0 {
				s.Broadcast(user, "have been offline!")
				return
			}
			//EOF is a symbol that end of IO
			if err != nil && err != io.EOF {
				log.Println("read error:", err)
				return
			}
			//convert byte buf to string that message who sended
			msg := string(buf[:n-1])

			//broadcast message who client sended
			s.Broadcast(user, msg)

		}
	}()

	//blocking goroutine
	select {}
}

//生产者：生产消息
func (s *Server) Broadcast(fromUser *User, msg string) {
	sendMsg := "[" + fromUser.Addr + "]" + fromUser.Name + ":" + msg
	s.Message <- sendMsg
}

//listen mesage and do brodcast to every body
//消费者：监听消费消息
func (s *Server) ListenMsg() {
	for {
		msg := <-s.Message
		//加锁
		s.mapLock.Lock()
		for _, user := range s.OnlineUsers {
			//将消息消费并发送到User的channel中，成为user.C的消息的生产者
			user.C <- msg
		}
		s.mapLock.Unlock()
	}
}
