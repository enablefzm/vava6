package checkcode

import (
	"fmt"
	"net"
	"strings"
	"vava6/vatools"
)

type CheckCode struct {
	addres   string
	port     string
	conn     *net.TCPListener
	chClose  chan int
	blnCheck bool
	code     string
	key      string
	val      string
	tryCount int
}

func NewCheckCode(port, key string) (*CheckCode, error) {
	code := vatools.GetRndChar(10)
	addr, err := net.ResolveTCPAddr("tcp", ":"+port)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &CheckCode{
		port:    port,
		code:    code,
		key:     key,
		conn:    conn,
		val:     vatools.MD5(fmt.Sprintf("%s_%s", key, code)),
		chClose: make(chan int),
	}, nil
}

func (this *CheckCode) BlnCheck() bool {
	return this.blnCheck
}

func (this *CheckCode) ListenCheck() {
	for {
		select {
		case _, ok := <-this.chClose:
			if !ok {
				return
			}
		default:
			c, err := this.conn.Accept()
			if err != nil {
				break
			}
			// 封装对象
			sck := &SocketConn{
				Conn:   c,
				bufLen: 128,
			}
			go this.onLink(sck)
		}
	}
}

func (this *CheckCode) GetVal() string {
	if len(this.val) < 10 {
		return vatools.GetRndChar(10)
	}
	return this.val[8:20]
}

func (this *CheckCode) onLink(sck *SocketConn) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	sck.Send(fmt.Sprintf("请输入【%s】的激活码", this.code))
	for {
		str, err := sck.Read()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		str = strings.ReplaceAll(str, "\r", "")
		str = strings.ReplaceAll(str, "\n", "")
		if len(str) == 12 {
			if str == this.GetVal() {
				this.blnCheck = true
				sck.Send("激活成功")
				sck.Close()
				this.closeServer()
				return
			}
			this.tryCount++
			if this.tryCount >= 3 {
				sck.Send("尝试3次还未正确，激活失败!")
				sck.Close()
				this.closeServer()
				return
			}
		}
		sck.Send("激活不正确")
	}
}

func (this *CheckCode) closeServer() {
	close(this.chClose)
	this.conn.Close()
}

type SocketConn struct {
	Conn   net.Conn
	bufLen int
}

func (this *SocketConn) Send(msg string) error {
	if len(msg) < 1 {
		return nil
	}
	_, err := this.Conn.Write([]byte(msg + "\n\r"))
	return err
}

func (this *SocketConn) Read() (string, error) {
	data := make([]byte, 0, 1024)
	bufData := make([]byte, this.bufLen)
	var err error
	for {
		n, err := this.Conn.Read(bufData)
		// 发生错误退出
		if err != nil {
			return "", err
		}
		// 数组整合
		data = append(data, bufData[:n]...)
		// 判断是不是读取结束
		if n != this.bufLen {
			break
		}
	}
	return string(data), err
}

func (this *SocketConn) Close() error {
	return this.Conn.Close()
}

func (this *SocketConn) GetIPInfo() string {
	ipAddr := this.Conn.RemoteAddr()
	return ipAddr.Network() + "|" + ipAddr.String()
}
