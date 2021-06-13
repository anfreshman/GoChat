package main

import (
	"flag"
	"fmt"
	"net"
)

// 定义Client类型
type Client struct {
	IP   string
	Port string
	Name string
	conn net.Conn
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
		}
		return client
	}
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
	} else {
		fmt.Println(">>>>>>>>>>服务器连接成功>>>>>>>>>>>")
	}
	select {}

}
