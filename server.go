package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip             string
	Port           int
	OnlineMap      map[string]*User
	mapLock        sync.RWMutex
	MessageChannel chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:             ip,
		Port:           port,
		OnlineMap:      make(map[string]*User),
		MessageChannel: make(chan string),
	}
}

func (server_poit *Server) ListenMessage() {
	for {
		msg := <-server_poit.MessageChannel
		fmt.Printf("%v\n", msg)
		server_poit.mapLock.RLock()
		for _, user := range server_poit.OnlineMap {
			user.Channel <- msg
		}
		server_poit.mapLock.RUnlock()
	}
}

func (server_poit *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "] " + user.Name + " : " + msg
	server_poit.MessageChannel <- sendMsg
}

func (server_poit *Server) Handler(connect net.Conn) {
	User := NewUser(connect)

	server_poit.mapLock.Lock()
	server_poit.OnlineMap[User.Name] = User
	defer server_poit.mapLock.Unlock()

	server_poit.BroadCast(User, "已上线")

}

func (server_poit *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server_poit.Ip, server_poit.Port))
	if err != nil {
		fmt.Printf("net listen err: %v\n", err)
		return
	}
	defer listener.Close()

	go server_poit.ListenMessage()

	for {
		connect, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept err: %v\n", err)
			continue
		}

		go server_poit.Handler(connect)
	}
}
