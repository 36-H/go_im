package main

import "net"

type User struct {
	Name    string
	Addr    string
	Channel chan string
	connect net.Conn
	server	*Server
}

func NewUser(connect net.Conn,server *Server) *User {
	userAddr := connect.RemoteAddr().String()

	user :=  &User{
		Name:    userAddr,
		Addr:    userAddr,
		Channel: make(chan string),
		connect: connect,
		server: server,
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

func (point *User) Online(){
	point.server.mapLock.Lock()
	point.server.OnlineMap[point.Name] = point
	defer point.server.mapLock.Unlock()

	point.server.BroadCast(point, "已上线")
}

func (point *User) Offline(){
	point.server.mapLock.Lock()
	delete(point.server.OnlineMap,point.Name)
	defer point.server.mapLock.Unlock()

	point.server.BroadCast(point, "已下线")
}

func (point *User) DoMessage(msg string){
	point.server.BroadCast(point,msg)
}