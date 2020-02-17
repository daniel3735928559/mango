package router

import (
	"fmt"
	"testing"
)

func TestRouter(t *testing.T) {
	rt := MCRouter()	

	for i := 0; i < 3; i++ {
		rt.AddNode(fmt.Sprintf("node%d",i))

	}
	rt.AddRoute(rt.ParseRoute("node0 > node1"))
	rt.AddRoute(rt.ParseRoute("node0 > fe {} "))
}
