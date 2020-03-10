package main

import (
	"fmt"
	// "strings"
	// "strconv"
	"github.com/docopt/docopt-go"
	// "github.com/google/shlex"
	"time"
	"math/rand"
	"mc/route"
	"mc/node"
	"mc/nodetype"
	"mc/value"
	"mc/serializer"
	"mc/transport"
	"mc/transport/mzmq"
	"mc/registry"
	// "libmango/transport/msocket"
	// "encoding/json"
)

// Each server is responsible for registration and mapping node

type MangoHandler func(command string, data map[string]interface{})

type MangoCommander struct {
	zmqTransport *mzmq.ZMQTransport
	//socketTransport *msocket.SocketTransport
	MessageInput chan transport.WrappedMessage
	Commands map[string]MangoHandler
	Registry *registry.Registry
	Self *node.Node
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
		if src_type == nil {
			fmt.Println("ERROR: Could not find type: ",src.NodeType)
			continue
		}
		
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
				fmt.Println("ERROR: failed running route", rt.ToString(), result_err)
				continue
			}
			outval, err := dst_type.ValidateInput(result_cmd, result_val)
			if err != nil {
				fmt.Println("ERROR: route output failed validation", result_cmd, result_val.ToString(), err)
				continue
			}
			data, err := serializer.Serialize(src.ToString(), msg.MessageId, result_cmd, outval.ToObject())
			if err != nil {
				fmt.Println("ERROR: Message failed to serialize", err)
				continue
			}
			dst.Transport.Tx(dst.Id, data)
		}
	}
}

type MCLoopbackTransport struct {
	MC *MangoCommander
}

func (t *MCLoopbackTransport) Tx(id string, data []byte) {
	if id != t.MC.Self.Id {
		fmt.Println("ERROR: Message not actually intended for MC")
		return
	}
	msg, err := serializer.Deserialize(string(data))
	if err != nil {
		fmt.Println("ERROR: deserialization error")
		return
	}
	if handler, ok := t.MC.Commands[msg.Command]; ok {
		handler(msg.Command, msg.Data)
	} else {
		fmt.Printf("ERROR: Invalid command %s\n",msg.Command)
	}
}

func (t *MCLoopbackTransport) RunServer() {
	
}


func (mc *MangoCommander) RouteAdd(command string, args map[string]interface{}) {
	spec := args["spec"].(string)
	group := args["group"].(string)
	id := args["id"].(string)
	routes, err := route.Parse(spec)
	if err != nil {
		fmt.Println("ERROR: Routes failed to parse: ", err)
		return
	}
	for i, rt := range routes {
		rt.Id = fmt.Sprintf("%s_%d", id, i)
		rt.Group = group
		mc.Registry.AddRoute(rt)
	}
}

func (mc *MangoCommander) Echo(command string, args map[string]interface{}) {
	fmt.Println("ECHO", args)
}

func (mc *MangoCommander) Excite(command string, args map[string]interface{}) {
	fmt.Printf("%s!\n", args["message"].(string))
}

func main() {
	usage := `Usage: mc [-t]`
	args, err := docopt.ParseDoc(usage)
	if err != nil {
		fmt.Println("ERROR parsing args: ",err)
	}
	
	fmt.Println("hi")

	
	//go test_socket_client(1920)
	message_aggregator := make(chan transport.WrappedMessage, 100)
	
	MC := &MangoCommander{
		zmqTransport: mzmq.MakeZMQTransport(1919, message_aggregator),
		//socketTransport: MakeSocketTransport(1920, message_aggregator),
		MessageInput: message_aggregator,
		Registry: registry.MakeRegistry()}	
	
	MC.Commands = map[string]MangoHandler{
		"routeadd":MC.RouteAdd,
		"echo":MC.Echo}
	rand.Seed(time.Now().UnixNano())
	self_id := ""
	for i := 0; i < 2; i++ {
		self_id += fmt.Sprintf("%016x", rand.Uint64())
	}
	
	MC.Self = &node.Node{
		Id: self_id,
		Name: "mc",
		Group: "system",
		NodeType: "mc",
		Transport: &MCLoopbackTransport{MC: MC}}

	MC.Registry.AddNode(MC.Self)

	mc_type := `
[config]
name mc

[interface]
input echo {message:string}
`
	
	mcnodetype, err := nodetype.Parse(mc_type)
	if err != nil {
		fmt.Println("Failed parsing MC node type: ", err)
		return
	}
	MC.Registry.AddNodeType(mcnodetype)

	if args["-t"].(bool) {
		go mzmq.TestZmqClient(1919)

		testnodetype, err := nodetype.Parse(`
[config]
name test_node

[interface]
output echo {message:string}
output excite {message:string}
`)
		if err != nil {
			fmt.Println("Failed parsing test node type: ",err)
			return
		}
		MC.Registry.AddNodeType(testnodetype)
		
		TestNode := &node.Node{
			Id: "TEST",
			Name: "test",
			Group: "system",
			NodeType: "test_node",
			Transport: MC.zmqTransport}
		
		MC.Registry.AddNode(TestNode)

		testroutes, rerr := route.Parse(`system/test > ? excite % echo {message += "!";} > system/mc`)
		if rerr != nil {
			fmt.Println("Failed parsing test route: ",rerr)
			return
		}

		for i, rt := range testroutes {
			rt.Group = "system"
			rt.Id = fmt.Sprintf("test_%d", i)
			MC.Registry.AddRoute(rt)
		}
		
	}
	MC.Run()
}
