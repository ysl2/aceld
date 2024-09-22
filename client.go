package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	key        int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		key:        999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.conn = conn
	return client
}

func (c *Client) DealResponse() {
	io.Copy(os.Stdout, c.conn)
	// io := make([]byte, 1024)
	// for {
	// 	n, err := c.conn.Read(io)
	// 	if err != nil {
	// 		fmt.Println("conn.Read err:", err)
	// 		return
	// 	}
	// 	fmt.Println(string(io[:n]))
	// }
}

func (c *Client) menu() bool {
	var key int
	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("4. 查看用户列表")
	fmt.Println("0. 退出")
	fmt.Scanln(&key)
	if key >= 0 && key <= 4 {
		c.key = key
		return true
	} else {
		fmt.Println("请输入正确的选项")
		return false
	}
}

func (c *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

func (c *Client) PrivateChat() {
	c.SelectUsers()
	fmt.Println("请选择聊天对象用户名:")
	var remoteName string
	fmt.Scanln(&remoteName)
	fmt.Println("请输入聊天内容, exit退出")
	for {
		var chatMsg string
		fmt.Scanln(&chatMsg)
		if chatMsg == "exit" {
			break
		}
		if len(chatMsg) > 0 {
			sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
			_, err := c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}
	}
}

func (c *Client) PublicChat() {
	var chatMsg string
	fmt.Println("请输入聊天内容, exit退出")
	for {
		fmt.Scanln(&chatMsg)
		if chatMsg == "exit" {
			break
		}
		if len(chatMsg) > 0 {
			sendMsg := chatMsg + "\n"
			_, err := c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}
	}
}

func (c *Client) UpdateName() bool {
	fmt.Println("请输入用户名:")
	fmt.Scanln(&c.Name)
	sendMsg := "rename|" + c.Name + "\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

func (c *Client) Run() {
	for c.key != 0 {
		for !c.menu() {
		}
		switch c.key {
		case 1:
			c.PublicChat()
		case 2:
			c.PrivateChat()
		case 3:
			c.UpdateName()
		case 4:
			c.SelectUsers()
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("连接服务器失败")
		return
	}
	go client.DealResponse()
	fmt.Println("连接服务器成功")
	client.Run()
}
