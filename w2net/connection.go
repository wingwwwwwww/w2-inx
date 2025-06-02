package w2net

import (
	"fmt"
	"net"

	"github.com/lnnlmario/w2-inx/w2iface"
)

// 链接模块
type Connection struct {
	// 当前链接的socket TCP套接字
	Conn *net.TCPConn

	// 链接的ID
	ConnId uint32

	// 当前的链接状态
	isClosed bool

	// 当前链接所绑定的处理业务方法
	handleAPI w2iface.HandleFunc

	// 告知当前链接已经退出/停止 channel
	ExitChan chan bool
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connId uint32, callback w2iface.HandleFunc) *Connection {
	return &Connection{
		Conn:      conn,
		ConnId:    connId,
		isClosed:  false,
		handleAPI: callback,
		ExitChan:  make(chan bool, 1),
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Read goroutine is running..., connId=", c.ConnId)
	defer fmt.Println("Read connId=", c.ConnId, "remote addr=", c.RemoteAddr(), "is exit!")
	defer c.Stop()
	// 循环读取数据
	for {
		// 读取数据
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Read error:", err)
			continue
		}

		// 调用当前链接绑定的业务处理方法
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("Handle error:", err)
			break
		}
	}
}

// 启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Connection Start()..., connId=", c.ConnId)
	// 启动从当前链接读数据的业务
	go c.StartReader()
	// TODO: 启动从当前链接写数据的业务
}

// 停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Connection Stop()..., connId=", c.ConnId)

	// 如果当前链接已经关闭
	if c.isClosed {
		return
	}
	// 将当前链接状态设置为已关闭
	c.isClosed = true
	// 关闭socket链接
	if err := c.Conn.Close(); err != nil {
		fmt.Println("Connection Stop() Close() error:", err)
		return
	}
	// 通知当前链接已经退出
	c.ExitChan <- true
	// 关闭当前链接的ExitChan
	close(c.ExitChan)
}

// 获取当前链接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前链接模块的链接ID
func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

// 获取远程客户端的TCP状态地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据 将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil // TODO: 实现发送数据的逻辑
}
