package main

import (
	"fmt"
	"time"
	"libmango"
)

func Cmd(args map[string]interface{}) (string, map[string]interface{}, error) {
	for k, v := range args {
		fmt.Printf("%s: %v\n", k, v)
	}
	return "cmd",args,nil
}

func main() {
	n, err := libmango.NewNode("test",map[string]libmango.MangoHandler{"cmd":Cmd})
	if err != nil {
		fmt.Println(err)
		return
	}
	n.Start()
	n.Send("cmd", map[string]interface{}{"foo":"bar","baz":1})
	n.Send("cmd", map[string]interface{}{"test":[]string{"one","two","three"}})
	for {
		time.Sleep(10*time.Second)
		fmt.Println("working...")
	}
}
