package main

import (
	"fmt"
	"net"
	"time"
)

// 模拟客户端
func main() {
	fmt.Println("client is running...")
	time.Sleep((1 * time.Second))

	// 1. 直接连接远程服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start error, exit!")
		return
	}

	for {
		// 2. 写数据到服务器
		_, err := conn.Write([]byte("Hello W2-INX"))
		if err != nil {
			fmt.Println("write error:", err)
			return
		}

		// 3. 读取服务器回显的数据
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error:", err)
			return
		}

		fmt.Printf("Server echo: %s, length is: %d\n", buf[:cnt], cnt)
		time.Sleep(1 * time.Second) // 每隔1秒发送一次数据
	}
}
