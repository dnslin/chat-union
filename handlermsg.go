package main

import (
	"fmt"
	"strings"
)

// Command 是所有命令的接口
type Command interface {
	Execute(args []string, user *User)
}

// RenameCommand 是一个具体的命令
type RenameCommand struct{}

func (c *RenameCommand) Execute(args []string, user *User) {
	fmt.Println("执行重命名操作，参数：", args)
	if len(args) == 0 {
		user.SendMessage("请指定用户名")
		return
	}
	// 判断 用户名是否存在
	if _, ok := user.server.OnlineMap[args[0]]; ok {
		user.SendMessage("当前用户名已被使用")
	}
	// 修改用户名
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.OnlineMap[args[0]] = user
	user.server.mapLock.Unlock()
	user.Name = args[0]
	// 发送更新成功消息
	user.SendMessage("用户名更新成功！")

}

// MessageCommand 是另一个具体的命令
type MessageCommand struct{}

func (c *MessageCommand) Execute(args []string, _ *User) {
	fmt.Println("执行消息发送操作，参数：", args)
}

// NormalMessageCommand 处理普通消息
type NormalMessageCommand struct{}

func (c *NormalMessageCommand) Execute(args []string, user *User) {
	user.SendMessage(strings.Join(args, " "))
}

// CommandParser 负责解析命令并执行
type CommandParser struct {
	commands map[string]Command
}

// WhoCommand 是另一个具体的命令
type WhoCommand struct {
}

func (c *WhoCommand) Execute(_ []string, user *User) {
	user.server.mapLock.Lock()
	for _, user := range user.server.OnlineMap {
		user.SendMessage("在线用户：" + user.Name + "\nIP地址: " + user.Addr)
	}
	user.server.mapLock.Unlock()
}

func NewCommandParser() *CommandParser {
	return &CommandParser{
		commands: map[string]Command{
			"rename":  &RenameCommand{},
			"message": &MessageCommand{},
			"who":     &WhoCommand{},
		},
	}
}

func (p *CommandParser) ParseAndExecute(message string, user *User) {
	if len(message) == 0 {
		fmt.Println("空消息")
		return
	}

	if prefix := message[:1]; prefix == "@" || prefix == "#" {
		content := strings.Fields(message[1:])
		if len(content) == 0 {
			user.SendMessage("指令格式错误")
			return
		}

		command, ok := p.commands[content[0]]
		if !ok {
			user.SendMessage("未知的指令:" + content[0])
			return
		}

		command.Execute(content[1:], user)
	} else {
		// 处理普通消息
		(&NormalMessageCommand{}).Execute(strings.Fields(message), user)
	}
}
