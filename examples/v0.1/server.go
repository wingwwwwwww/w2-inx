package main

import "github.com/lnnlmario/w2-inx/w2net"

func main() {
	// 1. 创建一个服务器实例
	s := w2net.NewServer("[W2-INX 0.1] Server")
	// 2. 启动server
	s.Serve()
}
