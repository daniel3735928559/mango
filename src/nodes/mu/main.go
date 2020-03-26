package main

import (
	"os"
	"fmt"
	"time"
	"math/rand"
	"mc/value"
	"encoding/json"
	"path/filepath"
	docopt "github.com/docopt/docopt-go"
	zmq "github.com/pebbe/zmq4"
)

var (
	socket *zmq.Socket
	ctx *zmq.Context
)

func send(srv, data string) {
	if ctx == nil {
		fmt.Println("CONTEXT IS nil")
		return
	}
	s, err := ctx.NewSocket(zmq.DEALER)
	if err != nil {
		fmt.Println("ERROR making socket:",err)
		return
	}
	socket = s
	rand.Seed(time.Now().UnixNano())
	identity := fmt.Sprintf("%04X-%04X", rand.Intn(0x10000), rand.Intn(0x10000))
	socket.SetIdentity(identity)
	//socket.SetLinger(-1)
	socket.Connect(srv)
	socket.Send("",zmq.SNDMORE)
	socket.Send(data, 0)
	socket.Recv(0)
	socket.Disconnect(srv)
}

func sendrecv(srv, data string) string {
	socket, _ = ctx.NewSocket(zmq.DEALER)
	rand.Seed(time.Now().UnixNano())
	identity := fmt.Sprintf("%04X-%04X", rand.Intn(0x10000), rand.Intn(0x10000))
	socket.SetIdentity(identity)
	socket.Connect(srv)
	socket.SetIdentity("aaaa")
	
	socket.Send("",zmq.SNDMORE)
	socket.Send(data, 0)
	
	data, err := socket.Recv(0)
	if err != nil {
		fmt.Println("ERROR: Failed receiving from agent",err)
		os.Exit(1)
	}
	socket.Disconnect(srv)
	return data
}

func main() {
	usage := `Usage: 
  mx send <command> <args>...
  mx start <nodetype> <nodegroup> <nodename> <nodeargs>...
  mx emp <group> <emp_file>
  mx nodes
  mx types
  mx routes
  mx connect <node>
  mx disconnect
  mx help [<command>]
  mx list
  mx pop
  mx get <id>
`
	args, err := docopt.ParseDoc(usage)
	if err != nil {
		fmt.Println("ERROR parsing args: ",err)
	}

	ctx, err = zmq.NewContext()
	if err != nil {
		fmt.Println("ERROR creating context:",err)
		return
	}
	if ctx == nil {
		fmt.Println("NO CONTEXT")
		return
	}
	
	agent_srv, ok := os.LookupEnv("MX_AGENT")
	if !ok || len(agent_srv) == 0 {
		fmt.Println("ERROR: MX_AGENT environment var not set")
		os.Exit(1)
	}

	if args["emp"].(bool) {
		emp_path, err := filepath.Abs(args["<emp_file>"].(string))
		if err != nil {
			fmt.Printf("ERROR: Problem accessing file %s\n", args["<emp_file>"].(string))
			os.Exit(1)
		}
		info, err := os.Stat(emp_path)
		if os.IsNotExist(err) || info.IsDir() {
			fmt.Printf("ERROR: File does not exist or is not regular file: %s\n", emp_path)
			os.Exit(1)
		}
		
		data := make(map[string]interface{})
		data["operation"] = "emp"
		data["empfile"] = emp_path
		data["group"] = args["<group>"].(string)
		bs, _ := json.Marshal(data)
		send(agent_srv, string(bs))
	} else if args["start"].(bool) {		
		data := make(map[string]interface{})
		data["operation"] = "start"
		data["type"] = args["<nodetype>"].(string)
		data["group"] = args["<nodegroup>"].(string)
		data["name"] = args["<nodename>"].(string)
		data["args"] = args["<nodeargs>"].([]string)
		bs, _ := json.Marshal(data)
		send(agent_srv, string(bs))
	} else if args["send"].(bool) {
		arg_val, _ := value.FromObject(make(map[string]interface{}))
		arg_name := ""
		if len(args["<args>"].([]string)) % 2 != 0 {
			fmt.Println("ERROR: Odd number of args. Expected <args> of the form `arg val arg val ...`")
			os.Exit(1)
		}
		for i, arg := range args["<args>"].([]string) {
			if i%2 == 0 {
				arg_name = arg
			} else {
				v, err := value.Parse(arg)
				if err != nil {
					v = value.MakeStringValue(arg)
				}
				arg_val.MapVal[arg_name] = v
			}
		}
		data := make(map[string]interface{})
		data["operation"] = "send"
		data["command"] = args["<command>"].(string)
		data["args"] = arg_val.ToObject()
		bs, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("ERROR: Failed to serialize data: %v\n",data)
			os.Exit(1)
		}
		//fmt.Println("Sending",string(bs))
		send(agent_srv, string(bs))
	} else if args["help"].(bool) {
		data := make(map[string]interface{})
		data["operation"] = "help"
		if v, ok := args["<command>"]; ok && v != nil {
			data["command"] = v.(string)
		}
		bs, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("ERROR: Failed to serialize data: %v\n",data)
			os.Exit(1)
		}
		fmt.Println(sendrecv(agent_srv, string(bs)))
	} else if args["nodes"].(bool) {
		bs, _ := json.Marshal(map[string]interface{}{"operation":"nodes"})
		fmt.Println(sendrecv(agent_srv, string(bs)))
	} else if args["types"].(bool) {
		bs, _ := json.Marshal(map[string]interface{}{"operation":"types"})
		fmt.Println(sendrecv(agent_srv, string(bs)))
	} else if args["routes"].(bool) {
		bs, _ := json.Marshal(map[string]interface{}{"operation":"routes"})
		fmt.Println(sendrecv(agent_srv, string(bs)))
	} else if args["connect"].(bool) {
		bs, _ := json.Marshal(map[string]interface{}{"operation":"connect","target":args["<node>"].(string)})
		send(agent_srv, string(bs))
	} else if args["disconnect"].(bool) {
		bs, _ := json.Marshal(map[string]interface{}{"operation":"disconnect"})
		send(agent_srv, string(bs))
	} else if args["list"].(bool) {
		data := make(map[string]interface{})
		data["operation"] = "list"
		bs, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("ERROR: Failed to serialize data: %v\n",data)
			os.Exit(1)
		}
		fmt.Println(sendrecv(agent_srv, string(bs)))
	} else if args["pop"].(bool) {
		data := make(map[string]interface{})
		data["operation"] = "pop"
		bs, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("ERROR: Failed to serialize data: %v\n",data)
			os.Exit(1)
		}
		fmt.Println(sendrecv(agent_srv, string(bs)))
	} else if args["get"].(bool) {
		data := make(map[string]interface{})
		data["operation"] = "get"
		data["id"] = args["<id>"].(string)
		bs, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("ERROR: Failed to serialize data: %v\n",data)
			os.Exit(1)
		}
		fmt.Println(sendrecv(agent_srv, string(bs)))
	}
}
