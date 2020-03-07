package nodetype

import (
	"fmt"
	"testing"
	value "mc/value"
)

func TestInterfaceParser(t *testing.T) {
	
        examples := map[string][]string{
		"node0 >":[]string{"node0 > node1"},
		"node0 <>":[]string{"node0 > node1", "node1 > node0"},
		"node0 < {}":[]string{"node1 > node0"}}
        for s, ans := range examples {
		fmt.Println("PARSING",s)
		rs, _ := Parse(s)
		fmt.Println("RS",rs,ans)
		if len(rs) > 0 {
			t.Errorf("expected no routes, got: %d", len(rs))
		}
        }
}
