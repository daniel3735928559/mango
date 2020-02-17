package main

import (
	"fmt"
	"strings"
	"net"
	"time"
)

type SocketTransport struct {
	Port int
	Socket net.Listener
	CurrentId int
	SockidToSocket map[int]net.Conn
	SockidToNode map[int]string
	NodeToSockid map[string]int
	MessageInput chan MCMessage
}

func MakeSocketTransport(port int, msgs chan MCMessage) *SocketTransport {
	return &SocketTransport {
		Port: port,
		CurrentId: 0,
		SockidToSocket: make(map[int]net.Conn),
		SockidToNode: make(map[int]string),
		NodeToSockid: make(map[string]int),
		MessageInput: msgs}
}

func (t *SocketTransport) Tx(dest string, data []byte) {
	t.SockidToSocket[t.NodeToSockid[dest]].Write(data)
}

func (t *SocketTransport) HandleConnection(c net.Conn, sockid int, register func(*MCMessage, MCTransport) bool) {
        fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	buf_len := 65536
	buf_index := 0
	buf := make([]byte, buf_len)
	header_raw := make([]byte, buf_len)
	body_raw := make([]byte, buf_len)
	var header_str string
	var body_str string
	var msg *MCMessage
	state := 0 // 0 = header, 1 = body
        for {
		recv_len, err := c.Read(buf[buf_index:])
                if err != nil {
                        fmt.Println("ERR reading message",err)
                        return
                }
		if buf_index + recv_len >= buf_len {
			fmt.Println("oops overflowed")
			return
		}
		delim_index := strings.Index(string(buf), "\n")
		for delim_index >= 0 {
			if state == 0 {
				// Receive header
				fmt.Println("header",string(buf[:delim_index]))
				for i := 0; i < delim_index; i++ {
					header_raw[i] = buf[i]
				}
				header_str = string(header_raw[:delim_index])
				state = 1
			} else if state == 1 {
				// Receive body
				for i := 0; i < delim_index; i++ {
					body_raw[i] = buf[i]
				}
				body_str = string(body_raw[:delim_index])
				fmt.Println("body",body_str)
				msg, err = ParseMessage(string(header_str), string(body_str))
				fmt.Println("got",msg)
				if err != nil {
					fmt.Println(err)
				} else if msg.Command == "hellomango" {
					if register(msg, t) {
						t.SockidToSocket[sockid] = c
						t.SockidToNode[sockid] = msg.Sender
						t.NodeToSockid[msg.Sender] = sockid
					} else {
						fmt.Println("Registration failed")
					}
				} else if _, ok := t.NodeToSockid[msg.Sender]; !ok {
					fmt.Println("No such node",msg.Sender)
				} else if t.NodeToSockid[msg.Sender] != sockid || t.SockidToNode[sockid] != msg.Sender {
					fmt.Println("Bad message received",msg.Sender,sockid)
				} else {
					fmt.Println("Processing normally from",msg.Sender)
					fmt.Println(msg)
					t.MessageInput <- *msg
				}

				state = 0
			}
			
			for i := delim_index+1; i < buf_len; i++ {
				buf[i-delim_index-1] = buf[i]
			}
			delim_index = strings.Index(string(buf), "\n")
		}
        }
        c.Close()
}

func (t *SocketTransport) RunServer(register func(*MCMessage, MCTransport) bool) {
	var err error
	t.Socket, err = net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d",t.Port))
        if err != nil {
                fmt.Println(err)
                return
        }
        defer t.Socket.Close()

        for {
                c, err := t.Socket.Accept()
                if err != nil {
                        fmt.Println(err)
                        return
                }
		new_id := t.CurrentId
		t.CurrentId += 1
                go t.HandleConnection(c, new_id, register)
        }
}


func test_socket_client(port int) {
	time.Sleep(2*time.Second)
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d",port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	fmt.Println("sending")
	fmt.Fprintf(conn, `{"source":"exciter2","format":"json","command":"hellomango"}
{"group":"foo"}
`)
	
	fmt.Fprintf(conn,`{"source":"exciter2","format":"json","command":"excite"}
{"message":"asda"}
`)	
}
