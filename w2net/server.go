package w2net

import (
	"fmt"
	"net"

	"github.com/lnnlmario/w2-inx/w2iface"
)

type Server struct {
	// 服务器的而名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server is listening IP:%s, Port:%d, Version:%s\n", s.IP, s.Port, s.IPVersion)

	go func() {
		// 1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("ResolveTCPAddr error:", err)
			return
		}

		// 2. 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("ListenTCP error:", err)
			return
		}
		fmt.Println("start w2-inx server", s.Name, "success, now listening...")

		// 3. 阻塞等待客户端的链接，处理客户端链接业务(读写)
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error:", err)
				continue
			}

			// 与已经建立连接的客户端做一个基本的-最大512字节长度的回显业务
			fmt.Println("Client connected:", conn.RemoteAddr().String())
			go func() {
				for {
					buf := make([]byte, 512)
					// 从客户端读取数据到buf中
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("receive buffer error", err)
						continue
					}

					fmt.Println("Receive from client:", string(buf[:cnt]), "Length:", cnt)

					// 回显功能
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write buffer error", err)
						continue
					}
				}
			}()
		}
	}()
}

func (s *Server) Stop() {}

func (s *Server) Serve() {
	// 启动服务器
	s.Start()

	// 阻塞主线程
	select {}
}

/**
 * NewServer 创建一个服务器实例
 * @param name 服务器名称
 * @return IServer 服务器实例
 */
func NewServer(name string) w2iface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}

	return s
}
