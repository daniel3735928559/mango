package main

import (
	"os"
	"fmt"
	"libmango"
)

var (
	logfile string
)

func Log(args map[string]interface{}) (string, map[string]interface{}, error) {
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	
	defer f.Close()

	msg := fmt.Sprintf("[%f %s] %s\n", args["timestamp"].(float64), args["source"].(string), args["message"].(string))
	
	if _, err = f.WriteString(msg); err != nil {
		panic(err)
	}
	return "", nil, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <logfile>\n", os.Args[0])
		os.Exit(1)
	}
	logfile = os.Args[1]
	n, err := libmango.NewNode("log",map[string]libmango.MangoHandler{"log":Log})
	if err != nil {
		fmt.Println(err)
		return
	}
	n.Start()
}
