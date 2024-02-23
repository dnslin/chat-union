package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:   ip,
		Port: port,
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
	fmt.Println("连接成功...:", conn.RemoteAddr().String())
}
