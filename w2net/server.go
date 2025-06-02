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

// 写死当前客户端绑定的handle api
func CallbackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Printf("[Connection Handle] CallbackToClient: Received %d bytes from %s: %s\n", cnt, conn.RemoteAddr().String(), string(data[:cnt]))
	// Echo the data back to the client
	_, err := conn.Write(data[:cnt])
	if err != nil {
		fmt.Println("Write error:", err)
	}
	return err
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

		var cId uint32 = 0

		// 3. 阻塞等待客户端的链接，处理客户端链接业务(读写)
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error:", err)
				continue
			}

			// 将处理新链接的业务方法和conn进行绑定，得到我们的链接模块
			dealConn := NewConnection(conn, cId, CallbackToClient)
			cId++

			go dealConn.Start()
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
