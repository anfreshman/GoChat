package timpl

import (
	"fmt"
	"net"
	"time"
)

func ServerTest() {

	fmt.Println("ServerTest Start")
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("客户端conn连接出错：", err)
		return
	}
	for {
		_, err := conn.Write([]byte("Hello,MyTCP"))
		if err != nil {
			fmt.Println("write error err", err)
			return
		}
		buf := make([]byte, 512)
		len, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read error err", err)
			return
		}
		fmt.Println("server call back :" + string(buf[:len]))
		time.Sleep(1 * time.Second)
	}
}

func Testmain() {
	s := GetServer("MyTCP01", "tcp4", "127.0.0.1", "8888")
	s.Start()
	go ServerTest()
}
