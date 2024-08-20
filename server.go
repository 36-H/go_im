package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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
	user := NewUser(connect,server_poit)

	user.Online()

	isLive := make(chan bool)

	go func ()  {
		buf := make([]byte,4096)
		for{
			n, err := connect.Read(buf)
			if n == 0 {
				user.Offline()
				return 
			}
			if err != nil && err != io.EOF {
				fmt.Printf("Connect Read Err: %v\n", err)
				return
			}
			msg := string(buf[:n - 1])

			user.DoMessage(msg)

			isLive <- true
		}
	}()

	for{
		select{
		case <-isLive:
		case <-time.After(10*time.Second):
			user.sendMsg("超时下线")
			server_poit.mapLock.Lock()
			delete(server_poit.OnlineMap,user.Name)
			server_poit.mapLock.Unlock()
			close(user.Channel)
			connect.Close()
			return
		}
	}

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
