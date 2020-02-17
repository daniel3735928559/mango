package main

import (
	"fmt"
	"math/rand"
	zmq "github.com/pebbe/zmq4"
)

type ZMQTransport struct {
	Port int
	Socket *zmq.Socket
	IdentityToNode map[string]string
	NodeToIdentity map[string]string
	MessageInput chan MCMessage
}

func MakeZMQTransport(port int, msgs chan MCMessage) *ZMQTransport {
	return &ZMQTransport {
		Port: port,
		IdentityToNode: make(map[string]string),
		NodeToIdentity: make(map[string]string),
		MessageInput: msgs}
}

func (t *ZMQTransport) Tx(dest string, data []byte) {
	t.Socket.Send(t.NodeToIdentity[dest],zmq.SNDMORE)
	t.Socket.Send(string(data), 0)
}

func (t *ZMQTransport) RunServer(register func(*MCMessage, MCTransport) bool) {
	t.Socket, _ = zmq.NewSocket(zmq.ROUTER)
	t.Socket.Bind(fmt.Sprintf("tcp://*:%d", t.Port))
	for {
		identity, _ := t.Socket.Recv(0)
		t.Socket.Recv(0)
		header, _ := t.Socket.Recv(0)
		fmt.Println(header)
		body, _ := t.Socket.Recv(0)
		fmt.Println(body)
		msg, err := ParseMessage(header, body)

		if err != nil {
			fmt.Println(err)
		} else if msg.Command == "hellomango" {
			if register(msg, t) {
				t.NodeToIdentity[msg.Sender] = identity
				t.IdentityToNode[identity] = msg.Sender
			} else {
				fmt.Println("Registration failed")
			}
			continue
		} else if _, ok := t.NodeToIdentity[msg.Sender]; !ok {
			fmt.Println("No such node",msg.Sender)
			continue
		} else if t.NodeToIdentity[msg.Sender] != identity || t.IdentityToNode[identity] != msg.Sender {
			fmt.Println("Updating",msg.Sender,identity)
			prev_identity := t.NodeToIdentity[msg.Sender]
			delete(t.IdentityToNode, prev_identity)
			t.IdentityToNode[identity] = msg.Sender
		}
		fmt.Println("Processing normally from",msg.Sender)
		fmt.Println(msg)
		t.MessageInput <- *msg
	}
	
}

func test_zmq_client(port int) {
	worker, _ := zmq.NewSocket(zmq.DEALER)
	defer worker.Close()
	set_id(worker) //  Set a printable identity
	worker.Connect(fmt.Sprintf("tcp://localhost:%d",port))

	worker.Send("",zmq.SNDMORE)
	worker.Send(`{"source":"exciter","format":"json","command":"hellomango"}`,zmq.SNDMORE)
	worker.Send(`{"group":"foo"}`,0)
	
	worker.Send("",zmq.SNDMORE)
	worker.Send(`{"source":"exciter","format":"json","command":"excite"}`,zmq.SNDMORE)
	worker.Send(`{"message":"asda"}`,0)
}

func set_id(soc *zmq.Socket) {
	identity := fmt.Sprintf("%04X-%04X", rand.Intn(0x10000), rand.Intn(0x10000))
	soc.SetIdentity(identity)
}
