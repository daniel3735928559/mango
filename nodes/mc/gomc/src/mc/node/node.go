package router

import (
	"fmt"
	serializer "mc/serializer"
)

type Node struct {
	Name string
	Group string
	Transport serializer.MCTransport
}

func (n *Node) Handler(src, mid, command string, args map[string]interface{}, dst string) {
	data := serializer.MakeMessage(src, mid, command, args).Serialize()
	n.Transport.Tx(dst, []byte(data))
}

func (n *Node) ToString() string {
	return fmt.Sprintf("%s/%s", n.Group, n.Name)
}

