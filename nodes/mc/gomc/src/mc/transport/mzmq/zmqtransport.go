package mzmq

import (
	"fmt"
	//"math/rand"
	"mc/transport"
	zmq "github.com/pebbe/zmq4"
)

type ZMQTransport struct {
	Port int
	Socket *zmq.Socket
	MessageInput chan transport.WrappedMessage
}

func MakeZMQTransport(port int, msgs chan transport.WrappedMessage) *ZMQTransport {
	return &ZMQTransport {
		Port: port,
		MessageInput: msgs}
}

func (t *ZMQTransport) GetServerAddr() string {
	return fmt.Sprintf("tcp://localhost:%d", t.Port)
}

func (t *ZMQTransport) Tx(identity string, data []byte) {
	t.Socket.Send(identity, zmq.SNDMORE)
	t.Socket.Send(string(data), 0)
}

func (t *ZMQTransport) RunServer() {
	t.Socket, _ = zmq.NewSocket(zmq.ROUTER)
	t.Socket.Bind(fmt.Sprintf("tcp://*:%d", t.Port))
	for {
		identity, _ := t.Socket.Recv(0)
		t.Socket.Recv(0)
		data, _ := t.Socket.Recv(0)
		fmt.Println("[zmqtransport.go] RX",identity,data)
		msg := transport.WrappedMessage {
			Transport: t,
			Identity: identity,
			Data: []byte(data)}
		t.MessageInput <- msg
	}
}

func TestZmqClient(port int, id string) {
	worker, _ := zmq.NewSocket(zmq.DEALER)
	defer worker.Close()
	set_id("test", worker) //  Set a printable identity
	worker.Connect(fmt.Sprintf("tcp://localhost:%d",port))

	worker.Send("",zmq.SNDMORE)
	worker.Send(fmt.Sprintf(`{"source":"system/test","format":"json","command":"echo","cookie":"%s"}
{"message":"foo"}`, id),0)
	
	worker.Send("",zmq.SNDMORE)
	worker.Send(fmt.Sprintf(`{"source":"system/test","format":"json","command":"excite","cookie":"%s"}
{"message":"bar"}`, id),0)
}

func set_id(identity string, soc *zmq.Socket) {
	//identity := fmt.Sprintf("%04X-%04X", rand.Intn(0x10000), rand.Intn(0x10000))
	soc.SetIdentity(identity)
}
