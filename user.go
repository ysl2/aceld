package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}

func (u *User) Online() {
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
	u.server.BroadCast(u, "已上线")
}

func (u *User) Offline() {
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	delete(u.server.OnlineMap, u.Name)
	u.server.BroadCast(u, "下线")
}

func (u *User) DoMessage(msg string) {
	if msg == "who" {
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ": 在线...\n"
			u.SendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		if _, ok := u.server.OnlineMap[newName]; ok {
			u.SendMsg("当前用户名已被使用\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()
			u.Name = newName
			u.SendMsg("您已更新用户名: " + newName + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMsg("消息格式不正确，请使用\"to|用户名|消息内容\"格式\n")
			return
		}
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.SendMsg("该用户名不存在\n")
			return
		}
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMsg("消息内容不能为空\n")
			return
		}
		remoteUser.SendMsg(u.Name + "对您说: " + content)
	} else {
		u.server.BroadCast(u, msg)
	}
}

func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.SendMsg(msg + "\n")
	}
}
