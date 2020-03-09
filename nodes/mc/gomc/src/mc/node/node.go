package node

import (
	"fmt"
	"mc/transport"
)

type Node struct {
	Id string
	Name string
	Group string
	NodeType string
	Transport transport.MangoTransport
}

func MakeNode(name, group, typename string) *Node {
	return &Node{
		Name: name,
		Group: group,
		NodeType: typename}
}

// func (n *Node) Handler(src, mid, command string, args map[string]interface{}) {
// 	data := serializer.MakeMessage(src, mid, command, args).Serialize()
// 	n.Transport.Tx(dst, []byte(data))
// }

func (n *Node) ToString() string {
	return fmt.Sprintf("%s/%s", n.Group, n.Name)
}

