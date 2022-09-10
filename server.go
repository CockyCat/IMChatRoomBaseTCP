package main

import (
	"fmt"
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
	user := NewUser(conn)

	s.mapLock.Lock()
	s.OnlineUsers[user.Name] = user
	s.mapLock.Unlock()

	s.Broadcast(user, "has been online!")

	select {}
}

func (s *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

//listen mesage and do brodcast to every body
func (s *Server) ListenMsg() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, user := range s.OnlineUsers {
			user.C <- msg
		}
		s.mapLock.Unlock()
	}
}
