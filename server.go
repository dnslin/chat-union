package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User // 在线用户列表
	Message   chan string      // 消息广播
	mapLock   sync.RWMutex     // 读写锁
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

func (sever *Server) Start() {
	// 监听 socket
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", sever.Ip, sever.Port))
	if err != nil {
		fmt.Println("网络socket监听异常...:", err)
	}
	// 关闭socket
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			fmt.Println("网络socket关闭异常...:", err)
		}
		fmt.Println("网络socket关闭...")
	}(listen)
	// 启动 监听 Message 的goroutine
	go sever.ListenMessage()
	// accept
	for {
		accept, err := listen.Accept()
		if err != nil {
			fmt.Println("网络socket接收异常...:", err)
			continue
		}
		// 处理业务逻辑
		go sever.Handler(accept)
	}
}

// Handler 处理业务
func (sever *Server) Handler(conn net.Conn) {
	// 创建一个用户
	user := NewUser(conn)
	// 用户上线
	sever.mapLock.Lock()
	// 将用户添加到在线用户列表  map(名字,用户对象)
	sever.OnlineMap[user.Name] = user
	sever.mapLock.Unlock()
	// 广播用户上线消息
	sever.BroadCast(user, "已上线")

	// 阻塞当前 不然 协程就死了
	select {}
}

// BroadCast 广播消息
func (sever *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	sever.Message <- sendMsg
}

// ListenMessage 监听广播消息 channel 的goroutine
func (sever *Server) ListenMessage() {
	for {
		msg := <-sever.Message
		// 将消息发送给在线用户
		sever.mapLock.Lock()
		for _, user := range sever.OnlineMap {
			user.C <- msg
		}
		sever.mapLock.Unlock()
	}
}
