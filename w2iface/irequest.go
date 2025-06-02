package w2iface

// IRequest 定义一个请求接口
// 实际上是把客户端请求的链接信息和请求数据包装到一起
type IRequest interface {
	// 得到当前链接
	GetConnection() IConnection
	// 得到请求的消息数据
	GetData() []byte
}
