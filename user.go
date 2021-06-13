package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
	// 用户对应的服务器
	server *Server
}

// 监听channel功能
func (this *User) ListenMessage() {
	for {
		mes := <-this.C
		// 传入mes，并将其转换为字节数组
		this.conn.Write([]byte(mes + "\n"))
	}
}

// 上线功能
func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, "上线")
}

// 下线功能
func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, "下线")
}

// 仅对当前客户端发送消息
func (this *User) SendMsgSimple(msg string) {
	this.conn.Read([]byte(msg))
}

// 消息处理
func (this *User) DoMessage(msg string) {

	// 规定:当接收到的msg为whos时，显示当前有哪些用户在线
	if msg == "whos" {
		for _, client := range this.server.OnlineMap {
			onlineMsg := "[" + client.Addr + "]" + client.Name + "在线"
			this.SendMsgSimple(onlineMsg)
		}
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 规定当接收到rename|新名称指令时，更改用户名
		newName := strings.Split(msg, "|")[1]
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsgSimple("当前用户名已被占用")
			return
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 规定私聊的信息为 to|用户名|消息
		// 获取用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsgSimple("您的消息格式有误，请按照\"to|张三|你好\"的格式进行私聊消息发送")
			return
		}
		// 根据用户名获取User对象
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsgSimple("对方用户不存在或已下线")
			return
		}
		// 调用对方用户的私发消息方法
		remotemsg := strings.Split(msg, "|")[2]
		if remotemsg == "" {
			this.SendMsgSimple("发送消息不能为空")
			return
		} else {
			remoteUser.SendMsgSimple(this.Name + "对您说:" + remotemsg)
		}
	} else {
		this.server.BroadCast(this, msg)
	}
}

// 当有新的TCP连接时，从该连接中获得User
func GetUser(conn net.Conn, server *Server) *User {
	// 从连接获得连接地址并转换为字符串
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	// 启动监听
	go user.ListenMessage()
	return user
}
