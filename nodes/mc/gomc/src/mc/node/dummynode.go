package node

import (
	"fmt"
)

type DummyNode struct {
	Group string
	Name string
}

func (n *DummyNode) ToString() string {
	return fmt.Sprintf("%s/%s", n.Group, n.Name)
}
