package node

import (
	"fmt"
	"sync"
	"mc/serializer"
	"mc/transport"
)

// The emitter node that sends out the merged messages

type MergeResultNode struct {
	Group string
	Name string
	Identity string
	Inputs []*MergeInputNode
	OutputChannel chan transport.WrappedMessage
	resultMux *sync.Mutex
}

// The received node

type MergeInputNode struct {
	Group string
	Name string
	Identity string
	result *MergeResultNode
	received map[string]map[string]interface{}
}


func MakeMergeNode(group, name string, merge_inputs []string, ch chan transport.WrappedMessage) []Node {
	identity := fmt.Sprintf("%s/merge_%s_output",group,name)
	output := &MergeResultNode{
		Group: group,
		Name: name,
		Identity: identity,
		OutputChannel: ch,
		resultMux: &sync.Mutex{}}
	ans := []Node{output}
	inputs := make([]*MergeInputNode, len(merge_inputs))
	for i, mi_name := range merge_inputs {
		input_identity := fmt.Sprintf("%s/merge_%s_input_%s",group,name,mi_name)
		inputs[i] = &MergeInputNode{
			Group: group,
			Name: mi_name,
			Identity: input_identity,
			result: output,
			received: make(map[string]map[string]interface{})}
		ans = append(ans, inputs[i])
	}
	output.Inputs = inputs
	return ans
}

func (n *MergeResultNode) GotAlive(identity string, transport transport.MangoTransport) {
	n.Identity = identity
}

func (n *MergeResultNode) SecsAgo() int {
	return 0
}

func (n *MergeResultNode) LastSeen() string {
	return "now"
}

func (n *MergeResultNode) GetId() string {
	return n.Identity
}

func (n *MergeResultNode) GetGroup() string {
	return n.Group
}

func (n *MergeResultNode) GetName() string {
	return n.Name
}

func (n *MergeResultNode) GetType() string {
	return "merge_output"
}

func (n *MergeResultNode) SendToNode(m serializer.Msg) error {
	return fmt.Errorf("ERROR: Merge output node should never directly receive messages")
}

func (n *MergeResultNode) ToString() string {
	return fmt.Sprintf("%s/%s", n.Group, n.Name)
}

func (n *MergeResultNode) CheckReady(mid string) {
	fmt.Println("[MC MERGE] CheckReady",&n)
	n.resultMux.Lock()
	defer n.resultMux.Unlock()
	all := true
	for _, input := range n.Inputs {
		fmt.Println("[MC MERGE] FROM",input.GetName(),input.received)
		if _, ok := input.received[mid]; !ok {
			//fmt.Println("[MC] MERGE Nothing yet from",input.GetName())
			all = false
		}
	}
	if !all {
		return
	}
	fmt.Println("[MC MERGE] MERGING!")
	ans := make(map[string]interface{})
	for _, input := range n.Inputs {
		ans[input.GetName()] = input.received[mid]
		delete(input.received, mid)
	}
	fmt.Println("[MC MERGE] MERGED", ans)
	wmsg := transport.WrappedMessage {
		Identity: n.ToString(),
		Transport: nil,
		Message: serializer.Msg {
			Sender: n.ToString(),
			MessageId: mid,
			Command: "merged",
			Cookie: n.Identity,
			Data: ans}}
	n.OutputChannel <- wmsg
}

// The receiver nodes

func (n *MergeInputNode) GetId() string {
	return n.Identity
}

func (n *MergeInputNode) GetGroup() string {
	return n.Group
}

func (n *MergeInputNode) GotAlive(identity string, transport transport.MangoTransport) {}


func (n *MergeInputNode) SecsAgo() int {
	return 0
}

func (n *MergeInputNode) LastSeen() string {
	return "now"
}

func (n *MergeInputNode) GetName() string {
	return n.Name
}

func (n *MergeInputNode) GetType() string {
	return "merge_input"
}

func (n *MergeInputNode) SendToNode(m serializer.Msg) error {
	n.received[m.MessageId] = m.Data
	n.result.CheckReady(m.MessageId)
	return nil
}

func (n *MergeInputNode) ToString() string {
	return fmt.Sprintf("%s/%s", n.Group, n.Name)
}
