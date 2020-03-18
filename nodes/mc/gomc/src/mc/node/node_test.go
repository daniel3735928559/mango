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

type TestTransport struct {
	Sent []serializer.Msg
}

func (t *TestTransport) RunServer() {

}

func (t *TestTransport) Tx(identity string, m serializer.Msg) error {
	t.Sent = append(t.Sent, m)
	return nil
}

func TestExecNode(t *testing.T) {
	tr := &TestTransport{
		Sent: make([]serializer.Msg, 0)}
	
	execnode := MakeExecNode("test","exectest", "testtype", "echo hi", tr)
	
	execnode.SendToNode(serializer.Msg{
		Sender:"test",
		MessageId:"m1",
		Command:"test_cmd",
		Cookie:"cookie1",
		Data:map[string]interface{}{"arg1":"val1"}})
	execnode.SendToNode(serializer.Msg{
		Sender:"test",
		MessageId:"m2",
		Command:"test_cmd",
		Cookie:"cookie2",
		Data:map[string]interface{}{"arg2":"val2"}})
	execnode.SendToNode(serializer.Msg{
		Sender:"test",
		MessageId:"m3",
		Command:"test_cmd",
		Cookie:"cookie3",
		Data:map[string]interface{}{"arg3":"val3"}})
	if len(execnode.TransportQueue) != 3 {
		t.Errorf("Expected transport queue of length 3, got %d", len(execnode.TransportQueue))
	}
	if len(tr.Sent) != 0 {
		t.Errorf("Expected sent array of length 0, got %d", len(tr.Sent))
	}
	
	execnode.Start("testserver")
	
	if len(execnode.TransportQueue) != 3 {
		t.Errorf("Expected transport queue of length 3, got %d", len(execnode.TransportQueue))
	}
	if len(tr.Sent) != 0 {
		t.Errorf("Expected sent array of length 0, got %d", len(tr.Sent))
	}
	
	execnode.GotAlive("testid",tr)
	
	if len(execnode.TransportQueue) != 0 {
		t.Errorf("Expected transport queue of length 0, got %d", len(execnode.TransportQueue))
	}
	if len(tr.Sent) != 3 {
		t.Errorf("Expected sent array of length 3, got %d", len(tr.Sent))
	}
	
	execnode.SendToNode(serializer.Msg{
		Sender:"test",
		MessageId:"m4",
		Command:"test_cmd",
		Cookie:"cookie4",
		Data:map[string]interface{}{"arg4":"val4"}})

	
	if len(execnode.TransportQueue) != 0 {
		t.Errorf("Expected transport queue of length 0, got %d", len(execnode.TransportQueue))
	}
	if len(tr.Sent) != 4 {
		t.Errorf("Expected sent array of length 4, got %d", len(tr.Sent))
	}
}
