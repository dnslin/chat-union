package main

import (
	"flag"
	"fmt"
	"github.com/TwiN/go-color"
	"io"
	"net"
	"os"
)

// Client 构建 client
type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       9999,
	}
	// 连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIp, client.ServerPort))
	if err != nil {
		fmt.Println("链接服务器发生异常...", err)
		return nil
	}
	client.conn = conn
	return client
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}

// Run 定义业务处理
func (client *Client) Run() {
	// 如果 flag 不为0 则一直处理
	for client.flag != 0 {
		for client.menu() != true {

		}
		if client.flag >= 0 && client.flag <= 3 {
			switch client.flag {
			case 1:
				fmt.Println(color.Ize(color.Blue, "####进入公聊模式####"))
				break
			case 2:
				fmt.Println(color.Ize(color.Blue, "####进入私聊模式####"))
				break
			case 3:
				client.updateName()
				break
			case 0:
				fmt.Println(color.Ize(color.Red, "退出"))
				break
			}
		} else {
			fmt.Println(color.Ize(color.Red, "请输入合法的数字"))
		}
	}
}

// 更新用户名
func (client *Client) updateName() bool {
	fmt.Println(color.Ize(color.Blue, ">>>>>请输入用户名:"))
	_, err2 := fmt.Scanln(&client.Name)
	if err2 != nil {
		return false
	}
	sendMsg := "#rename " + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write 发生异常...", err)
		return false
	}
	return true

}

// DealResponse 处理服务端返回的消息
func (client *Client) DealResponse() {
	// 一旦client.conn有数据，就会执行下面的代码
	_, err := io.Copy(os.Stdout, client.conn)
	if err != nil {
		return
	}
}

// 定义菜单
func (client *Client) menu() bool {
	var key int
	fmt.Println(color.Ize(color.Cyan, "1. 公聊模式"))
	fmt.Println(color.Ize(color.Cyan, "2. 私聊模式"))
	fmt.Println(color.Ize(color.Cyan, "3. 更新用户名"))
	fmt.Println(color.Ize(color.Cyan, "0. 退出"))
	_, err := fmt.Scanln(&key)
	if err != nil {
		return false
	}
	client.flag = key
	return true

}

func main() {
	// 命令行解析
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("链接服务器失败...")
		return
	}
	fmt.Println("链接服务器成功...")
	// 处理 server 返回的消息
	go client.DealResponse()
	// 启动客户端业务
	client.Run()
}
