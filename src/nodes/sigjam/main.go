package main

import (
	"os"
	"fmt"
	"log"
	"net"
	"time"
	"bufio"
	"encoding/json"
	"libmango"
)

type Sigjam struct {
	SocketPath string
	s2w chan interface{}
	w2s chan interface{}
	account string
	node *libmango.Node
}

func (s *Sigjam) Listen(c chan interface{}) {
	socket, err := net.Dial("unix", s.SocketPath)
	if err != nil {
		log.Print("ERROR connecting: ",err)
		return
	}
	e := json.NewEncoder(socket)
	e.Encode(map[string]string{"type":"subscribe","username":s.account})
	var msg interface{}
	reader := bufio.NewReader(socket)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Print("ERROR reading: ",err)
		} else {
			err := json.Unmarshal([]byte(line), &msg)
			if err != nil {
				log.Print("ERROR deserialising: ", line, err)
			} else {
				fmt.Println("[SIGJAM] MSG",msg)
				c <- msg
			}
		}
	}
}

func (s *Sigjam) SendRequest(request interface{}) {
	b, err := json.Marshal(request)
	if err != nil {
		log.Print("[SIGJAM] ERROR serializing: ",err)
		return
	}
	log.Print("[SIGJAM] TX ", string(b))
	socket, err := net.Dial("unix", s.SocketPath)
	if err != nil {
		log.Print("ERROR: ",err)
		return
	}
	e := json.NewEncoder(socket)
	e.Encode(request)
}

func (s *Sigjam) Handle(msg interface{}) {
	fmt.Println("handling",msg)
	base := msg.(map[string]interface{})
	if base["type"].(string) == "message" {
		fmt.Println("message")
		data := base["data"].(map[string]interface{})
		fmt.Println("data")
		fmt.Println("D",data["dataMessage"])
		if data_msg, ok := data["dataMessage"].(map[string]interface{}); ok {
			fmt.Println("data msg")
			s.node.Send("recv",map[string]interface{}{
				"time":data_msg["timestamp"].(float64),
				"from":data["source"].(string),
				"msg":data_msg["message"].(string)})
		} else if sync_msg, ok := data["syncMessage"].(map[string]interface{}); ok {
			fmt.Println("sync msg")
			time.Sleep(1*time.Second)
			if sent, ok := sync_msg["sent"].(map[string]interface{}); ok {
				m := sent["message"].(map[string]interface{})
				s.node.Send("recv",map[string]interface{}{
					"time":m["timestamp"].(float64),
					"from":data["source"].(string),
					"msg":m["message"].(string)})
			}
		}
	}
}

func (s *Sigjam) Run() {
	for {
		select {
		case r := <- s.s2w:
			data, err := json.Marshal(r)
			if err != nil {
				fmt.Println("Problem serialising",r)
			} else {
				fmt.Println(string(data))
			}
			s.Handle(r)
			
		case r := <- s.w2s:
			s.SendRequest(r)
		}
	}

}

func (s *Sigjam) Send(args map[string]interface{}) (string, map[string]interface{}, error) {
	fmt.Println("[SIGJAM]",args)
	s.w2s <- map[string]interface{}{"type":"send","username":s.account,"recipientNumber":args["to"].(string),"messageBody":args["msg"].(string)}
	return "", nil, nil
}

func main() {
	fmt.Println("hello")
	
	signal2world := make(chan interface{}, 100)
	world2signal := make(chan interface{}, 100)
	
	s := &Sigjam{
		s2w:signal2world,
		w2s:world2signal,
		SocketPath: "/var/run/signald/signald.sock",
		account:os.Args[1]}

	n, err := libmango.NewNode("sigjam",map[string]libmango.MangoHandler{"send":s.Send})
	if err != nil {
		fmt.Println(err)
		return
	}
	s.node = n
	
	// Subscribe to account
	
	go s.Listen(signal2world)
	go s.Run()
	
	n.Start()
}
