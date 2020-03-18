package node

import (
	"fmt"
	"testing"
	"mc/serializer"
	"mc/transport"
)

func TestMergeNode(t *testing.T) {
	ch := make(chan transport.WrappedMessage, 100)
	mergenodes := MakeMerge("test","merger", []string{"in1","in2","in3"}, ch)
	in1 := mergenodes[1]
	in2 := mergenodes[2]
	in3 := mergenodes[3]
	in1.SendToNode(serializer.Msg{
		Sender:"test",
		MessageId:"m1",
		Command:"test_cmd",
		Cookie:"cookie1",
		Data:map[string]interface{}{"arg1":"val1"}})
	in2.SendToNode(serializer.Msg{
		Sender:"test",
		MessageId:"m1",
		Command:"test_cmd",
		Cookie:"cookie1",
		Data:map[string]interface{}{"arg2":"val2"}})
	in3.SendToNode(serializer.Msg{
		Sender:"test",
		MessageId:"m1",
		Command:"test_cmd",
		Cookie:"cookie1",
		Data:map[string]interface{}{"arg3":"val3"}})
	select {
	case wmsg := <-ch:
		outmsg := wmsg.Message.Data
		fmt.Println(outmsg)
		if arg, ok := outmsg["in1"].(map[string]interface{}); ok {
			if arg["arg1"].(string) != "val1" {
				t.Errorf("Expected in1.arg1 == val1 in output, but in1.arg1 = %v", arg["arg1"].(string))
			}
		} else {
			t.Errorf("Expected key `in1` in output, but did not find it")
		}
		if arg, ok := outmsg["in2"].(map[string]interface{}); ok {
			if arg["arg2"].(string) != "val2" {
				t.Errorf("Expected in2.arg2 == val2 in output, but in2.arg2 = %v", arg["arg2"].(string))
			}
		} else {
			t.Errorf("Expected key `in2` in output, but did not find it")
		}
		if arg, ok := outmsg["in3"].(map[string]interface{}); ok {
			if arg["arg3"].(string) != "val3" {
				t.Errorf("Expected in3.arg3 == val3 in output, but in3.arg3 = %v", arg["arg3"].(string))
			}
		} else {
			t.Errorf("Expected key `in3` in output, but did not find it")
		}
	}
}
