package nodetype

import (
	"fmt"
	"strings"
	"mc/value"
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
	
func TestParser(t *testing.T) {
	inputs := []string{
		"3",
		"{key1:[4,5]}",
		`"hello"`,
		`{k1:["hello1"]}`,
		`{}`,
		`[["a","b","c"],["d"]]`}
	input_answers := [][]bool{
		[]bool{false,true,true,false,false,false},
		[]bool{true,false,false,true,true,true}}
	outputs := []string{
		`[1,2,3]`,
		`[[1,2],[3,4]]`,
		`40`,
		`{f:["a","b","c"]}`,
		`{f:[["a","b"],["c","d"]]}`,
		`true`}
	output_answers := [][]bool{
		[]bool{false,true,true,false,false,false},
		[]bool{true,false,true,true,true,true}}
        examples := []string{
		`[config]
name bar

[interface]
type foo [num]
input inp1 oneof({key1:foo},string)
output out1 oneof([foo],num)

[usage]
bar --arg1 <arg1>`,
	`[config]
name baz

[interface]
type foo oneof([string],[[string]])
input inp1 oneof(num,foo,{k1:[string]=["hi"],k2*:bool})
output out1 oneof(num,{f:foo},bool,[num],[string],[bool])

[usage]
baz -a --arg1 <arg1>...`}
        for i, ex := range examples {
		fmt.Println("PARSING",ex)
		nt, err := Parse(ex)
		if err != nil {
			t.Errorf("Failed to parse: %v", err)
		} else {
			for j, vs := range inputs {
				v, _ := value.Parse(vs)
				nv, _ := nt.ValidateInput("inp1", v)
				if nv == nil && input_answers[i][j] == true {
					t.Errorf("Interface %d: Input %s should have validated but did not", i, v.ToString())
				} else if nv != nil && input_answers[i][j] == false {
					t.Errorf("Interface %d: Input %s should not have validated but did", i, v.ToString())
				}
			}
			for j, vs := range outputs {
				v, _ := value.Parse(vs)
				nv, _ := nt.ValidateOutput("out1", v)
				if nv == nil && output_answers[i][j] == true {
					t.Errorf("Interface %d: Output %s should have validated but did not", i, v.ToString())
				} else if nv != nil && output_answers[i][j] == false {
					t.Errorf("Interface %d: Output %s should not have validated but did", i, v.ToString())
				}
			}
			
			v1 := value.MakeBoolValue(true)
			_, err := nt.ValidateInput("inp2", v1)
			if err == nil {
				t.Errorf("inp2 should not validate")
			}
			_, err = nt.ValidateOutput("out2", v1)
			if err == nil {
				t.Errorf("out2 should not validate")
			}
		}
        }
}


