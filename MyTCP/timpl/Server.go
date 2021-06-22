package timpl

import (
	"fmt"
	"net"
)

type Server struct {
	// 服务器名字
	Name string
	// 版本号
	IPVersion string
	// IP号
	IP string
	// 端口号
	Port string
}

func (s *Server) Start() {
	// 第一步，根据IP+Port获得ADDR对象
	tcpAddr, err := net.ResolveTCPAddr("tcp4", s.IP+":"+s.Port)
	if err != nil {
		fmt.Println("获取ADDR出错：", err)
		return
	}

	// 第二步，根据ADDR获取listener
	listenner, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		fmt.Println("获取listenner出错：", err)
		return
	}

	// 第三步，轮询监听并获取客户端连接的conn
	for {
		conn, err := listenner.AcceptTCP()
		if err != nil {
			fmt.Println("获取conn连接出错：", err)
			return
		}

		// 获取到监听端口后，即可进行业务，这里暂时写一个简单的回显函数
		go func() {
			for {
				buf := make([]byte, 512)
				// 返回一个字节长度
				len, err := conn.Read(buf)
				if err != nil {
					fmt.Println("conn读入出错：0", err)
					return
				}
				fmt.Println(string(buf[:len]))

				// 回显
				if _, err := conn.Write(buf[:len]); err != nil {
					fmt.Println("conn 写入出错:", err)
				}
			}
		}()
	}
}

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	s.Start()
	// 启动服务器后的业务

}

// 获得一个Server
func GetServer(name string, ipVersion string, ip string, port string) *Server {
	return &Server{
		Name:      name,
		IPVersion: ipVersion,
		IP:        ip,
		Port:      port,
	}
}
