package emp

import (
	"fmt"
	"strings"
	"strconv"
)

type EMP struct {
	Group string
	NodeTypes map[string]string
	Nodes []string
	Routes []string
}

func (emp *EMP) Run() {
	for n := range Nodes {
		fmt.Println(n)
	}
}
