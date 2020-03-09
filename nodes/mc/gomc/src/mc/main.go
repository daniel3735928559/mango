package main

import (
	"fmt"
	// "strings"
	// "strconv"
	// "github.com/docopt/docopt-go"
	// "github.com/google/shlex"
	// "time"
	"mc/value"
	"mc/serializer"
	"mc/transport"
	"mc/transport/mzmq"
	"mc/registry"
	// "libmango/transport/msocket"
	// "encoding/json"
)

// Each server is responsible for registration and mapping node

type MangoHandler func(header map[string]string, args map[string]interface{})

type MangoCommander struct {
	zmqTransport *mzmq.ZMQTransport
	//socketTransport *msocket.SocketTransport
	MessageInput chan transport.WrappedMessage
	Commands map[string]MangoHandler
	Registry registry.Registry
}


// func (n *MCNode) heartbeat_worker() {
// 	for {
// 		time.Sleep(5*time.Second)
// 		n.Transport.Tx("heartbeat", []byte(""))
// 	}
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

// func (mc *MangoCommander) Register(msg *serializer.MCMessage, t serializer.MCTransport) bool {
// 	fmt.Println("Register")
// 	// Peel off registration commands directly
// 	var group string
// 	// if group, ok = msg.Data.(map[string]interface{})["group"].(string); !ok {
// 	// 	group = "root"
// 	// }
// 	group = "root"

// 	new_node := &node.Node{
// 		Name: msg.Sender,
// 		Group: group,
// 		Transport: t}
// 	mc.Router.AddNode(new_node)

// 	return true
// }

func (mc *MangoCommander) Run() {
	go mc.zmqTransport.RunServer()
	//go mc.socketTransport.RunServer(mc.Register)

	for wrapped_msg := range mc.MessageInput {
		msg, err := serializer.Deserialize(string(wrapped_msg.Data))
		if err != nil {
			fmt.Println("ERROR: deserialization error")
			continue
		}
		fmt.Println("FROM",msg.Sender)
		fmt.Println("DATA",msg.Data)
		src := mc.Registry.FindNodeById(msg.Cookie)
		if src == nil {
			fmt.Println("ERROR: node not found:", msg.Cookie)
			continue
		}
		// Since we got here, the node we're talking to is
		// using this transport/id combo, so set those
		// properties here (in case it has reconnected since
		// last time or something): 
		src.Id = wrapped_msg.Identity
		src.Transport = wrapped_msg.Transport
		
		// validate message is of acceptable format for node
		// output:
		cmd := msg.Command
		incoming_val, err := value.FromObject(msg.Data)
		if err != nil {
			fmt.Println("ERROR: failed to convert incoming value")
			continue
		}
		
		src_type := mc.Registry.FindNodeType(src.NodeType)
		
		validated_val, err := src_type.ValidateOutput(cmd, incoming_val)
		if err != nil {
			fmt.Println("ERROR: message failed to validate")
			continue
		}
		
		// If we made it here, the value is legitimate. Send
		// on all routes originating at src:
		routes := mc.Registry.FindRoutesBySrc(src.ToString())
		for _, rt := range routes {
			dst := mc.Registry.FindNodeByName(rt.Dest)
			dst_type := mc.Registry.FindNodeType(dst.NodeType)
			result_cmd, result_val, result_err := rt.Run(cmd, validated_val)
			if result_err != nil {
				fmt.Println("ERROR: failed running route", rt.ToString())
				continue
			}
			outval, err := dst_type.ValidateInput(result_cmd, result_val)
			if err != nil {
				fmt.Println("ERROR: route failed validation", rt.ToString())
				continue
			}
			data, err := serializer.Serialize(src.ToString(), msg.MessageId, result_cmd, outval.ToObject())
			if err != nil {
				fmt.Println("ERROR: Message failed to serialize")
				continue
			}
			dst.Transport.Tx(dst.Id, data)
		}
	}
}

// type MCLoopbackTransport struct {
// 	MC *MangoCommander
// }

// func (t *MCLoopbackTransport) Tx(dest string, data []byte) {
// 	if dest == "root/mc" {
// 		msg, err := serializer.Deserialize(string(data))
// 		if err == nil {
// 			if handler, ok := t.MC.Commands[msg.Command]; ok {
// 				handler(msg.RawHeader(), msg.Data)
// 			}
// 		}
// 	}
// }

// func (t *MCLoopbackTransport) RunServer(register func(*serializer.MCMessage, serializer.MCTransport) bool) {
// }


// func (mc *MangoCommander) RouteAdd(header map[string]string, args map[string]interface{}) {
// 	spec := args["spec"].(string)
// 	mc.Router.ParseAndAddRoutes(spec)
// }

func main() {
	fmt.Println("hi")
	go mzmq.TestZmqClient(1919)
	//go test_socket_client(1920)
	message_aggregator := make(chan transport.WrappedMessage, 100)
	
	MC := &MangoCommander{
		zmqTransport: mzmq.MakeZMQTransport(1919, message_aggregator),
		//socketTransport: MakeSocketTransport(1920, message_aggregator),
		MessageInput: message_aggregator}
	
	// MC.Commands = map[string]MangoHandler{
	// 	"routeadd":MC.RouteAdd}
	// MC.Router.AddNode(&router.Node{
	// 	Name: "mc",
	// 	Group: "root",
	// 	Transport: &MCLoopbackTransport{MC: MC}})
	MC.Run()
}
