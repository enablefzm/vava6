package vaconn

import (
	"net"
	"net/http"
	"vava6/valog"

	"golang.org/x/net/websocket"
)

const (
	TypeWEBSocket = "WEBSocket"
	TypeSocket    = "Socket"
)

func NewNaviConnect(nType, sPort, sIp string, onLinkFunc func(MConn)) *NaviConnect {
	ob := &NaviConnect{
		connType: nType,
		port:     sPort,
		ip:       sIp,
		onfunc:   onLinkFunc,
	}
	return ob
}

/**
 *	网络连接器
 **/
type NaviConnect struct {
	connType string
	port     string
	ip       string
	bufLen   int
	onfunc   func(MConn)   // 回调函数
	onCount  func() string // 返回总在线人数
	fnError  func(error)   // 错误时执行
}

// 通过WebSocket连接进来的事件处理
func (this *NaviConnect) onWebSocket(webConn *websocket.Conn) {
	conn := &WebSocketConn{Conn: webConn}
	// 将这个连接转交给游戏世界处理这个连接
	this.onfunc(conn)
}

// 通过Socket连接进来的事件处理
func (this *NaviConnect) onSocket(conn net.Conn) {
	sock := &SocketConn{
		Conn:   conn,
		bufLen: 128,
	}
	// 将这个连接转交给游戏世界处理这个连接
	this.onfunc(sock)
}

func (this *NaviConnect) OnError(err error) {
	if this.fnError != nil {
		this.fnError(err)
	}
}

func (this *NaviConnect) SetError(fn func(error)) {
	this.fnError = fn
}

func (this *NaviConnect) CountPlayer(rw http.ResponseWriter, req *http.Request) {
	if this.onCount != nil {
		rw.Write([]byte(this.onCount()))
	}
}

func (this *NaviConnect) SetOnCountFunc(handleFunc func() string) {
	this.onCount = handleFunc
}

func (this *NaviConnect) Listen() {
	// 如果是WEBSocket运行WEBSocket
	if this.connType == TypeWEBSocket {
		// 注册
		http.Handle("/", websocket.Handler(this.onWebSocket))
		http.HandleFunc("/count", this.CountPlayer)
		// 阻塞侦听
		valog.OBLog.LogMessage("开始WebSocket侦听 " + this.ip + ":" + this.port + " ...")
		err := http.ListenAndServe(this.ip+":"+this.port, nil)
		if err != nil {
			valog.OBLog.LogMessage("WebSocket服务器启动失败！" + err.Error())
			this.OnError(err)
			return
		}
	} else {
		// 通过普通SocketTCP连接
		addr, err := net.ResolveTCPAddr("tcp", this.ip+":"+this.port)
		if err != nil {
			valog.OBLog.LogMessage("Socket服务器启动 地址创建失败！")
			return
		}
		// socket连接
		server, err := net.ListenTCP("tcp", addr)
		if err != nil {
			valog.OBLog.LogMessage("Socket服务器启动失败！" + err.Error())
			this.OnError(err)
			return
		}
		// 侦听
		valog.OBLog.LogMessage("开始Socket侦听 " + this.ip + ":" + this.port + " ...")
		for {
			conn, err := server.Accept()
			if err != nil {
				// 客户端连接错误
				valog.OBLog.LogMessage("客户端Socket连接失败！" + err.Error())
			}
			// 移交给协程处理
			go this.onSocket(conn)
		}
	}
	valog.OBLog.LogMessage("侦听结束")
}

// 接入Socket类型统一接口
type MConn interface {
	Send(msg string) error // 发送
	Read() (string, error) // 接送
	Close()                // 关闭
	GetIPInfo() string     // 获取IP地址
}

/*******************************************************************************
 * WebSocketConn对象
 ******************************************************************************/
type WebSocketConn struct {
	Conn *websocket.Conn
}

func (this *WebSocketConn) Send(msg string) error {
	if len(msg) < 1 {
		return nil
	}
	err := websocket.Message.Send(this.Conn, msg)
	return err
}

func (this *WebSocketConn) Read() (string, error) {
	var err error
	var str string
	err = websocket.Message.Receive(this.Conn, &str)
	return str, err
}

func (this *WebSocketConn) Close() {
	this.Conn.Close()
}

func (this *WebSocketConn) GetIPInfo() string {
	r := this.Conn.Request()
	return this.Conn.RemoteAddr().Network() + "|" + r.RemoteAddr
}

// ============================================================================
// SocketConn对象
// ============================================================================
type SocketConn struct {
	Conn   net.Conn
	bufLen int
}

func (this *SocketConn) Send(msg string) error {
	if len(msg) < 1 {
		return nil
	}
	_, err := this.Conn.Write([]byte(msg + "\n\r"))
	if err != nil {
		valog.OBLog.LogMessage(err.Error())
	}
	return err
}

func (this *SocketConn) Read() (string, error) {
	data := make([]byte, 0)
	bufData := make([]byte, this.bufLen)
	var err error
	for {
		n, err := this.Conn.Read(bufData)
		// 发生错误退出
		if err != nil {
			/*
				if err == io.EOF {
					valog.OBLog.LogMessage("Socket客户端主动断开连接。")
				}
				// valog.OBLog.LogMessage("Socket错误事件" + err.Error())
			*/
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

func (this *SocketConn) Close() {
	this.Conn.Close()
}

func (this *SocketConn) GetIPInfo() string {
	ipAddr := this.Conn.RemoteAddr()
	return ipAddr.Network() + "|" + ipAddr.String()
}
