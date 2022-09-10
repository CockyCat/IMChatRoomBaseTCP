package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
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
}
