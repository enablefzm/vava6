package vaconn

import (
	"fmt"
	"io"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	vaConn := NewNaviConnect(CLIENT_SOCKET, "2016", "127.0.0.1", func(conn MConn) {
		go func() {
			i := 1
			for {
				fmt.Print("第", i, "次读取")
				result, err := conn.Read()
				fmt.Println("读到完成.")
				if err != nil {
					if err == io.EOF {
					}
					fmt.Println("断开连接了", err.Error())
					conn.Close()
					return
				}
				fmt.Println(result)
				i++
			}
		}()

		// conn.Send("login fff fff")
		go conn.Send("login fzm fzm1111")
		// time.Sleep(1 * time.Microsecond)
		go conn.Send("chat hi1")
		// time.Sleep(1 * time.Microsecond)
		go conn.Send("chat hi2")
		// time.Sleep(1 * time.Microsecond)
		go conn.Send("chat hi3")
	})
	vaConn.Listen()
	time.Sleep(6 * time.Second)
}
