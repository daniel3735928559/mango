package nodetype

import (
	"fmt"
	"strings"
	"testing"
)

func TestInterfaceParser(t *testing.T) {
        examples := []string{
		`type foo num
input inp1 {key1:foo}
output out1 oneof([foo],num)`,
	`input inp1 num
input inp2 oneof(string,{k1:[string],k2:bool})
output out1 oneof(num,string,bool,[num],[string],[bool])`}
        for _, ans := range examples {
		fmt.Println("PARSING",ans)
		ni, err := ParseNodeInterface(ans)
		if err != nil {
			t.Errorf("Failed to parse: %v", err)
		} else {
			obs := strings.TrimSpace(ni.ToString())
			if obs != ans {
				t.Errorf("expected %s, got: %s", ans, obs)
			}
		}
        }
}
