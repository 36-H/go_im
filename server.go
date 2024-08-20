package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip string
	Port int
}

func NewServer(ip string,port int) *Server{
	return &Server{
		Ip: ip,
		Port: port,
	}
}

func (server_poit *Server) Handler(connect net.Conn){
	fmt.Println("connect accept")
}

func (server_poit *Server) Start(){
	listener,err  := net.Listen("tcp", fmt.Sprintf("%s:%d", server_poit.Ip, server_poit.Port))
	if err != nil {
		fmt.Printf("net listen err: %v\n", err)
		return
	}
	defer listener.Close()
	for {
		connect,err := listener.Accept();
		if err != nil {
			fmt.Printf("accept err: %v\n", err)
			continue
		}

		go server_poit.Handler(connect)
	}
}
