package main

import (
	"fmt"
	"net"
)

// User 用户类
type User struct {
	Name string
	Addr string
	C    chan string // 用户发送数据的管道
	conn net.Conn    // 用户的连接
}

// NewUser 创建一个用户
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
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
			fmt.Printf("用户%s消息发送异常...错误信息 %s:\n", user.Name, err)
			continue
		}
	}
}
