package main

import (
	"fmt"
	// "strings"
	// "strconv"
	// "github.com/docopt/docopt-go"
	// "github.com/google/shlex"
	// "time"
	serializer "mc/serializer"
	mzmq "mc/transport/mzmq"
	router "mc/router"
	// "libmango/transport/msocket"
	// "encoding/json"
)

// Each server is responsible for registration and mapping node

// type MCNode struct {
// 	NodeId string
// 	GroupId string
// 	Transport serializer.MCTransport
// 	State int // alive, stalled, dead
// }

type MangoHandler func(header map[string]string, args map[string]interface{})

type MangoCommander struct {
	zmqTransport *mzmq.ZMQTransport
	//socketTransport *msocket.SocketTransport
	MessageInput chan serializer.MCMessage
	Router *router.Router
	Commands map[string]MangoHandler
}


// func (n *MCNode) heartbeat_worker() {
// 	for {
// 		time.Sleep(5*time.Second)
// 		n.Transport.Tx("heartbeat", []byte(""))
// 	}
// }

// func (n *MCNode) sendMessage(command string, args map[string]interface{}) {
// 	header := n.MakeHeader(command)
// 	body, _ := json.Marshal(args)
// 	n.Transport.Tx(header, body)
// }

// func (mc *MangoCommander) reap(n *MCNode) {
	
// }

// func (mc *MangoCommander) reaperWorker() {	
// 	for {
// 		time.Sleep(5*time.Second)
// 		for _, n := range mc.Nodes {
// 			if n.State == 2 {
// 				mc.reap(n)
// 			}
// 		}
// 	}
// }

func (mc *MangoCommander) Register(msg *serializer.MCMessage, t serializer.MCTransport) bool {
	fmt.Println("Register")
	// Peel off registration commands directly
	var group string
	// if group, ok = msg.Data.(map[string]interface{})["group"].(string); !ok {
	// 	group = "root"
	// }
	group = "root"

	new_node := &router.Node{
		Name: msg.Sender,
		Group: group,
		Transport: t}
	mc.Router.AddNode(new_node)

	return true
}

func (mc *MangoCommander) Run() {
	go mc.zmqTransport.RunServer(mc.Register)
	//go mc.socketTransport.RunServer(mc.Register)

	for msg := range mc.MessageInput {
		fmt.Println("FROM",msg.Sender)
		fmt.Println("DATA",msg.Data)
		if src := mc.Router.FindNode(msg.Sender); src != nil {
			fmt.Println("Processing from",msg.Sender)
			mc.Router.Send(src, msg.MessageId, "cmd", msg.Data)
		} else {
			fmt.Println("No such node",msg.Sender)
		}
	}
}

type MCLoopbackTransport struct {
	MC *MangoCommander
}

func (t *MCLoopbackTransport) Tx(dest string, data []byte) {
	if dest == "root/mc" {
		msg, err := serializer.ParseMessage(string(data))
		if err == nil {
			if handler, ok := t.MC.Commands[msg.Command]; ok {
				handler(msg.RawHeader(), msg.Data)
			}
		}
	}
}

func (t *MCLoopbackTransport) RunServer(register func(*serializer.MCMessage, serializer.MCTransport) bool) {
}


func (mc *MangoCommander) RouteAdd(header map[string]string, args map[string]interface{}) {
	spec := args["spec"].(string)
	mc.Router.ParseAndAddRoutes(spec)
}

func main() {
	fmt.Println("hi")
	go mzmq.TestZmqClient(1919)
	//go test_socket_client(1920)
	message_aggregator := make(chan serializer.MCMessage, 100)
	MC := &MangoCommander{
		Router: router.MakeRouter(),
		zmqTransport: mzmq.MakeZMQTransport(1919, message_aggregator),
		//socketTransport: MakeSocketTransport(1920, message_aggregator),
		MessageInput: message_aggregator}
	MC.Commands = map[string]MangoHandler{
		"routeadd":MC.RouteAdd}
	MC.Router.AddNode(&router.Node{
		Name: "mc",
		Group: "root",
		Transport: &MCLoopbackTransport{MC: MC}})
	MC.Run()
}
