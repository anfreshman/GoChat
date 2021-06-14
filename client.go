package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

// 定义Client类型
type Client struct {
	IP   string
	Port string
	Name string
	conn net.Conn

	// 新增模式选择功能
	flag int
}

func GetClient(IP, Port string) *Client {
	conn, err := net.Dial("tcp", IP+":"+Port)
	if err != nil {
		fmt.Println("客户端连接失败，错误信息为:", err)
		return nil
	} else {
		client := &Client{
			IP:   IP,
			Port: Port,
			conn: conn,
			flag: 999,
		}
		return client
	}
}

func (client *Client) menu() bool {
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更改用户名")
	fmt.Println("0.退出")

	var flag int
	fmt.Scanln(&flag)
	if flag >= 0 && flag < 4 {
		client.flag = flag
		return true
	} else {
		fmt.Println("您的输入不合法")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
			// 死循环，反复执行menu函数
		}
		// Go的switch不需要break语句
		switch client.flag {
		case 1:
			// 公聊模式代码块
			fmt.Println("您已进入公聊模式")
			client.PublicChat()
			break
		case 2:
			// 私聊模式代码块
			fmt.Println("您已进入私聊模式")
			client.PrivateChat()
			break
		case 3:
			fmt.Println("您选择了更改用户名")
			client.UpdateName()
			break
		}
	}
}

// 更新用户名的业务方法
func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>>>>请输入用户名>>>>>>>>>>")
	fmt.Scanln(&client.Name)
	_, err := client.conn.Write([]byte("rename|" + client.Name + "\n"))
	if err != nil {
		fmt.Println("更新出错:", err)
		return false
	}
	return true
}

// 公聊模式
func (client *Client) PublicChat() {
	fmt.Println(">>>>请输入您的消息: exit表示退出")
	var msg string
	fmt.Scanln(&msg)
	for msg != "exit" {
		if msg != "" {
			_, err := client.conn.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Println("消息发送失败")
				return
			}
			msg = ""
			fmt.Println(">>>>请输入您的消息: exit表示退出")
			fmt.Scanln(&msg)
		}

	}
}

// 私聊模式
func (client *Client) PrivateChat() {
	// 查询现在有哪些人在线，并由接收线程自动显示
	client.conn.Write([]byte("whos" + "\n"))
	fmt.Println(">>>>选择您要私聊的对象: exit表示退出")
	var remoteName string
	fmt.Scanln(&remoteName)
	if remoteName == "exit" || remoteName == "" {
		return
	}
	fmt.Println(">>>>请输入您的消息: exit表示退出")
	var msg string
	fmt.Scanln(&msg)
	for remoteName != "exit" {
		for msg != "exit" {
			if msg != "" {
				_, err := client.conn.Write([]byte("to|" + remoteName + "|" + msg + "\n"))
				if err != nil {
					fmt.Println("消息发送失败")
					return
				}
				msg = ""
				fmt.Println(">>>>请输入您的消息: exit表示退出")
				fmt.Scanln(&msg)
			}
		}
		remoteName = ""
		fmt.Println(">>>>选择您要私聊的对象: exit表示退出")
		fmt.Scanln(&remoteName)
	}

}

// 监听服务器返回消息的方法，单开go程
func (client *Client) DealResponse() {
	// 一种读入连接并显示到标准输出的简写方法
	fmt.Println("客户端显示方法执行")
	io.Copy(os.Stdout, client.conn)
}

var serverIP string
var serverPort string

func init() {
	// 输入 -h时，会显示最后的提示信息
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "设置服务器的IP地址")
	flag.StringVar(&serverPort, "port", "8888", "设置服务器的连接端口号")
}

func main() {
	// 解析命令行
	flag.Parse()
	client := GetClient("127.0.0.1", "8888")
	if client == nil {
		fmt.Println(">>>>>>>>服务器连接失败>>>>>>>>>")
		return
	}

	fmt.Println(">>>>>>>>>>服务器连接成功>>>>>>>>>>>")
	go client.DealResponse()
	client.Run()

	select {}

}
