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
	"nodes/mx/tui"
)

var (
	ctx *zmq.Context
)

func send(srv, data string) {
	if ctx == nil {
		fmt.Println("CONTEXT IS nil")
		return
	}
	socket, err := ctx.NewSocket(zmq.DEALER)
	if err != nil {
		fmt.Println("ERROR making socket:",err)
		return
	}
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
	socket, _ := ctx.NewSocket(zmq.DEALER)
	rand.Seed(time.Now().UnixNano())
	identity := fmt.Sprintf("%04X-%04X", rand.Intn(0x10000), rand.Intn(0x10000))
	socket.SetIdentity(identity)
	socket.Connect(srv)
	
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

func sendonerecvall(srv, data string, ch chan string) {
	s, _ := ctx.NewSocket(zmq.DEALER)
	rand.Seed(time.Now().UnixNano())
	identity := fmt.Sprintf("%04X-%04X", rand.Intn(0x10000), rand.Intn(0x10000))
	s.SetIdentity(identity)
	s.Connect(srv)
	
	s.Send("",zmq.SNDMORE)
	s.Send(data, 0)

	for {
		data, err := s.Recv(0)
		if err != nil {
			fmt.Println("ERROR: Failed receiving from agent",err)
		} else {
			ch <- data
		}
	}
	fmt.Println("LISTENING ENDED OOPS")
}

func main() {
	usage := `Usage: 
  mx send <command> [<args>...]
  mx start <nodetype> <nodegroup> <nodename> <nodeargs>...
  mx emp <group> <emp_file> [<args>...]
  mx nodes
  mx types
  mx routes
  mx connect <node>
  mx disconnect
  mx help [<command>]
  mx list
  mx clear
  mx peek
  mx pop
  mx get <id>
  mx listen <types>...
  mx shell <node>
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
		emp_args := make([]map[string]string, 0)
		arg_name := ""
		if len(args["<args>"].([]string)) % 2 != 0 {
			fmt.Println("ERROR: Odd number of args. Expected <args> of the form `arg val arg val ...`")
			os.Exit(1)
		}
		for i, arg := range args["<args>"].([]string) {
			if i%2 == 0 {
				arg_name = arg
			} else {
				emp_args = append(emp_args, map[string]string{"name":arg_name,"value":arg})
			}
		}
		data["args"] = emp_args
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
		bs, _ := json.Marshal(data)
		fmt.Println(sendrecv(agent_srv, string(bs)))
	} else if args["peek"].(bool) {
		data := make(map[string]interface{})
		data["operation"] = "peek"
		bs, _ := json.Marshal(data)
		fmt.Println(sendrecv(agent_srv, string(bs)))
	} else if args["clear"].(bool) {
		data := make(map[string]interface{})
		data["operation"] = "clear"
		bs, _ := json.Marshal(data)
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
	} else if args["listen"].(bool) {
		data := make(map[string]interface{})
		data["operation"] = "listen"
		data["types"] = args["<types>"].([]string)
		bs, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("ERROR: Failed to serialize data: %v\n",data)
			os.Exit(1)
		}
		c := make(chan string, 1000)
		go sendonerecvall(agent_srv, string(bs), c)
		for msg := range c {
			fmt.Println(msg)
		}
	} else if args["shell"].(bool) {
		bsc, _ := json.Marshal(map[string]interface{}{"operation":"connect","target":args["<node>"].(string)})
		sendrecv(agent_srv, string(bsc))
		bs, _ := json.Marshal(map[string]interface{}{"operation":"help"})
		iface := sendrecv(agent_srv, string(bs))
		if iface == "Not connected" {
			fmt.Println(iface)
			return
		}
		
		input_ch := make(chan string, 1000)
		cmd_ch := make(chan string, 1000)
		done_ch := make(chan bool)
		
		bs, _ = json.Marshal(map[string]interface{}{"operation":"listen","types":[]string{"any"}})
		go sendonerecvall(agent_srv, string(bs), input_ch)
		
		t := tui.MakeTui(iface, cmd_ch, done_ch)
		go t.Run()
		for {
			select{
			case msg := <- input_ch:
				d := map[string]interface{}{}
				json.Unmarshal([]byte(msg), &d)
				a, _ := json.Marshal(d["args"])
				t.GotOutput(d["command"].(string), string(a))
			case cmd := <- cmd_ch:
				send(agent_srv, cmd)
			case <- done_ch:
				fmt.Println("MX DONE")
				break
			}
		}
	}
}
