package node

import (
	"fmt"
	"os"
	"os/exec"
	"time"
	"math/rand"
	"mc/serializer"
	"mc/transport"
	"github.com/google/shlex"
)

type ExecNode struct {
	Id string
	TransportId string
	Name string
	Group string
	NodeType string
	Transport transport.MangoTransport
	Command string
	Environment map[string]string
	Server string
	LastHeartbeat int64
	Status NodeStatus
	TransportQueue []serializer.Msg
	Proc *exec.Cmd
}

func MakeExecNode(group, name, typename, command string, env map[string]string, trans transport.MangoTransport) *ExecNode {
	id := ""
	for i := 0; i < 4; i++ {
		id += fmt.Sprintf("%016x", rand.Uint64())
	}
	return &ExecNode{
		Id: id,
		Name: name,
		Group: group,
		NodeType: typename,
		Command: command,
		Environment: env,
		Status: NODE_STATUS_READY,
		Transport: trans,
		TransportQueue: make([]serializer.Msg, 0),
		TransportId: ""}
}

// func (n *Node) Handler(src, mid, command string, args map[string]interface{}) {
// 	data := serializer.MakeMessage(src, mid, command, args).Serialize()
// 	n.Transport.Tx(dst, []byte(data))
// }

func (n *ExecNode) GetName() string {
	return n.Name
}
func (n *ExecNode) GetId() string {
	return n.Id
}
func (n *ExecNode) GetGroup() string {
	return n.Group
}
func (n *ExecNode) GetType() string {
	return n.NodeType
}

func (n *ExecNode) RestartWorker() {
	n.Proc.Wait()
	fmt.Println("RESTARTING")
	n.Start(n.Server)
}

func (n *ExecNode) Start(server string) {
	args, err := shlex.Split(n.Command)
	if err != nil {
		fmt.Println("ERROR starting command",n.Command,err)
		return
	}
	fmt.Println("[execnode.go] STARTING",args)
	n.Server = server
	n.Proc = exec.Command(args[0], args[1:]...)
	n.Proc.Stdin = os.Stdin
	n.Proc.Stdout = os.Stdout
	n.Proc.Env = os.Environ()
	n.Proc.Env = append(n.Proc.Env, fmt.Sprintf("MANGO_COOKIE=%s", n.Id))
	n.Proc.Env = append(n.Proc.Env, fmt.Sprintf("MANGO_SERVER=%s", n.Server))
	for name, val := range n.Environment {
		n.Proc.Env = append(n.Proc.Env, fmt.Sprintf("%s=%s", name, val))
	}
	n.Status = NODE_STATUS_RUNNING
	n.Proc.Start()
	go n.watch_proc()
}

func (n *ExecNode) watch_proc() {
	n.Proc.Wait()
	n.Status = NODE_STATUS_DEAD
}

func (n *ExecNode) Kill() {
	n.Status = NODE_STATUS_DEAD
}

func (n *ExecNode) SecsAgo() int {
	if n.Status == NODE_STATUS_DEAD {
		return -1
	}
	return int((time.Now().UnixNano()-n.LastHeartbeat)/1000000000)
}

func (n *ExecNode) LastSeen() string {
	if n.Status == NODE_STATUS_DEAD {
		return "dead"
	}
	return fmt.Sprintf("%d secs ago", (time.Now().UnixNano()-n.LastHeartbeat)/1000000000)
}

func (n *ExecNode) GotAlive(tid string, t transport.MangoTransport) {
	n.LastHeartbeat = time.Now().UnixNano()
	n.TransportId = tid
	n.Transport = t
	n.Status = NODE_STATUS_RUNNING
	// Send any enqueued things
	//fmt.Println("GOT ALIVE--sending queue of length", len(n.TransportQueue))
	for _, bs := range n.TransportQueue {
		n.Transport.Tx(n.TransportId, bs)
	}
	n.TransportQueue = make([]serializer.Msg, 0)
}

func (n *ExecNode) SendToNode(m serializer.Msg) error {
	if n.Status == NODE_STATUS_RUNNING && len(n.TransportId) > 0 {
		//fmt.Println("Sending")
		return n.Transport.Tx(n.TransportId, m)
	}
	//fmt.Println("Appending to queue")
	n.TransportQueue = append(n.TransportQueue, m)
	return nil
}

func (n *ExecNode) ToString() string {
	return fmt.Sprintf("%s/%s", n.Group, n.Name)
}
