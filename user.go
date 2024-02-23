package main

import (
	"net"
)

// User 用户类
type User struct {
	Name   string
	Addr   string
	C      chan string // 用户发送数据的管道
	conn   net.Conn    // 用户的连接
	server *Server     // 当前用户所在的server
}

// NewUser 创建一个用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	// 启动监听当前用户channel消息的goroutine
	go user.ListenMessage()

	return user
}

// ListenMessage 监听当前用户channel的方法
func (user *User) ListenMessage() {
	for {
		msg := <-user.C
		// 写回数据
		_, err := user.conn.Write([]byte(msg + "\n"))
		if err != nil {
			// todo 用户下线会报错
			continue
		}
	}
}

// Online 上线
func (user *User) Online() {
	// 用户上线
	user.server.mapLock.Lock()
	// 将用户添加到在线用户列表  map(名字,用户对象)
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()
	// 广播用户上线消息
	user.server.BroadCast(user, "已上线")
}

// Offline 下线
func (user *User) Offline() {
	// 将用户从在线用户列表中删除
	delete(user.server.OnlineMap, user.Name)
	// 广播用户下线消息
	user.server.BroadCast(user, "下线")
}

// SendMessage 发送消息
func (user *User) SendMessage(msg string) {
	user.server.BroadCast(user, msg)
}
