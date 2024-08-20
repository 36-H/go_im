package main

import (
	"net"
)

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

func (point *User) sendMsg(msg string){
	point.Channel <- msg
}

func (point *User) DoMessage(msg string){
	if msg == "/who"{
		onlineMsg := "=========在线用户==========\n";
		point.server.mapLock.RLock()
		for _,user := range point.server.OnlineMap {
			onlineMsg += "[" +user.Addr + "] " + user.Name + "\n" 
		}
		point.server.mapLock.RUnlock()
		onlineMsg += "==========================";
		point.sendMsg(onlineMsg)
	}else if len(msg) > 7 && msg[:8] == "/rename|" {
		newName := msg[8:]
		point.server.mapLock.Lock()
		_,ok := point.server.OnlineMap[newName]
		if ok {
			point.sendMsg("当前用户名已被使用\n")
			return
		} else{
			delete(point.server.OnlineMap,point.Name)
			point.Name = newName
			point.server.OnlineMap[point.Name] = point
			point.sendMsg("用户名修改成功,当前用户名为:" + point.Name + "\n")
		}
		defer point.server.mapLock.Unlock()
	}else{
		point.server.BroadCast(point, msg)
	}
}

