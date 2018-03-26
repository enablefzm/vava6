package vaconn

import (
	"net"
)

// 连接到指定的服务器地址
// 	成功返回net.TCPConn，失败返回error
func LinkServer(sIp, port string, handle func()) (*net.TCPConn, error) {
	tcpIpAddr, err := net.ResolveTCPAddr("tcp4", sIp+":"+port)
	if err != nil {
		return nil, err
	}
	return net.DialTCP("tcp", nil, tcpIpAddr)
}

type SocketClient struct {
}
