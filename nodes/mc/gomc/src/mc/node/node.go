package node

import (
	"fmt"
	"os"
	"os/exec"
	"time"
	"math/rand"
	"mc/transport"
	"github.com/google/shlex"
)

type NodeStatus int

const (
	NODE_STATUS_READY NodeStatus = iota+1
	NODE_STATUS_RUNNING
	NODE_STATUS_UNRESPONSIVE
	NODE_STATUS_DEAD
)

type Node struct {
	Id string
	TransportId string
	Name string
	Group string
	NodeType string
	Transport transport.MangoTransport
	Command string
	Server string
	LastHeartbeat int64
	Status NodeStatus
	TransportQueue [][]byte
	Proc *exec.Cmd
}

func MakeNode(name, group, typename, command string, trans transport.MangoTransport) *Node {
	id := ""
	for i := 0; i < 4; i++ {
		id += fmt.Sprintf("%016x", rand.Uint64())
	}
	return &Node{
		Id: id,
		Name: name,
		Group: group,
		NodeType: typename,
		Command: command,
		Status: NODE_STATUS_READY,
		Transport: trans,
		TransportQueue: make([][]byte, 0),
		TransportId: ""}
}

// func (n *Node) Handler(src, mid, command string, args map[string]interface{}) {
// 	data := serializer.MakeMessage(src, mid, command, args).Serialize()
// 	n.Transport.Tx(dst, []byte(data))
// }

func (n *Node) RestartWorker() {
	n.Proc.Wait()
	fmt.Println("RESTARTING")
	n.Start(n.Server)
}

func (n *Node) Start(server string) {
	args, err := shlex.Split(n.Command)
	if err != nil {
		fmt.Println("ERROR starting command",n.Command,err)
	}
	n.Server = server
	n.Proc = exec.Command(args[0], args[1:]...)
	n.Proc.Stdin = os.Stdin
	n.Proc.Stdout = os.Stdout
	n.Proc.Env = os.Environ()
	n.Proc.Env = append(n.Proc.Env, fmt.Sprintf("MANGO_COOKIE=%s", n.Id))
	n.Proc.Env = append(n.Proc.Env, fmt.Sprintf("MANGO_SERVER=%s", n.Server))
	n.Status = NODE_STATUS_RUNNING
	n.Proc.Start()
}

func (n *Node) Heartbeat() {
	n.LastHeartbeat = time.Now().UnixNano()
}

func (n *Node) HeartbeatFailed() {
	n.Status = NODE_STATUS_UNRESPONSIVE
}

func (n *Node) Kill() {
	n.Status = NODE_STATUS_DEAD
}

func (n *Node) GotAlive(tid string, t transport.MangoTransport) {
	n.TransportId = tid
	n.Transport = t
	n.Status = NODE_STATUS_RUNNING
	// Send any enqueued things
	fmt.Println("GOT ALIVE--sending queue of length", len(n.TransportQueue))
	for _, bs := range n.TransportQueue {
		n.Transport.Tx(n.TransportId, bs)
	}
	n.TransportQueue = make([][]byte, 0)
}

func (n *Node) SendToNode(data []byte) {
	if n.Status == NODE_STATUS_RUNNING && len(n.TransportId) > 0 {
		fmt.Println("Sending")
		n.Transport.Tx(n.TransportId, data)
	} else {
		fmt.Println("Appending to queue")
		n.TransportQueue = append(n.TransportQueue, data)
	}
}

func (n *Node) ToString() string {
	return fmt.Sprintf("%s/%s", n.Group, n.Name)
}
