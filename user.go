package main

import "net"

type User struct {
	Name    string
	Addr    string
	Channel chan string
	connect net.Conn
}

func NewUser(connect net.Conn) *User {
	userAddr := connect.RemoteAddr().String()

	user :=  &User{
		Name:    userAddr,
		Addr:    userAddr,
		Channel: make(chan string),
		connect: connect,
	}

	go user.ListenMessage()

	return user
}

func (point *User) ListenMessage() {
	for {
		msg := <-point.Channel

		point.connect.Write([]byte(msg + "\n"))
	}
}
