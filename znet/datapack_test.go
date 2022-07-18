package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		return
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}

			go func(conn net.Conn) {
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println(err)
						break
					}

					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println(err)
						return
					}

					if msgHead.GetDataLen() > 0 {
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println(err)
							return
						}

						fmt.Println("-->Recv MsgId:", msg.Id, "datalen:", msg.DataLen, "data:", string(msg.Data))
					}

				}
			}(conn)
		}
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println(err)
		return
	}

	dp := NewDataPack()
	//模拟粘包过程,封装两个msg一同发送
	//封装第一个
	m1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte("zinx"),
	}
	data1, err := dp.Pack(m1)
	if err != nil {
		fmt.Println(err)
		return
	}
	//封装第二个
	m2 := &Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte("nihao"),
	}
	data2, err := dp.Pack(m2)
	if err != nil {
		fmt.Println(err)
		return
	}
	//将两个包黏在一起
	data := append(data1, data2...)
	//一次性发送
	conn.Write(data)

	select {}
}
