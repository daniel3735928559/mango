package libmango

// import (
// 	"fmt"
// 	// "strings"
// 	// "strconv"
// 	// "github.com/docopt/docopt-go"
// 	// "github.com/google/shlex"
// 	mprotocol "libmango/protocol"
// 	// "libmango/transport/mzmq"
// 	// "libmango/transport/msocket"
// )

// // Each server is responsible for registration and mapping node

// type MHeader struct {
// 	Cmd string
// 	MessageId string
// 	Format string
// }

// type MMsg struct {
// 	header *MHeader
// 	args map[string]interface{}
// }

// type MInterface struct {
// 	Commands map[string]func(*MHeader,interface{})(string,interface{},error)
// }

// type MNode struct {
// 	Version string
// 	NodeId string
// 	server string
// 	m_interface *MInterface	
// }

// func (n *MNode) Dispatch(header *MHeader, args interface{})   {
// 	if handler, ok := n.m_interface.Commands[header.Cmd]; ok {
// 		result_cmd, result, err := handler()
// 		if err != nil {
// 			n.HandleError(err)
// 		} else if result != nil {
// 			n.MSend(result_cmd, result, header.MessageId)
// 		}
// 	} else {
// 		fmt.Println("No handler for",header.Cmd)
// 		n.HandleError(errors.New(fmt.Sprintf("No handler for %s",header.Cmd)))
// 	}
// }

// func (n *MNode) HandleError(e error) {
// 	fmt.Println(e)
// 	n.MSend("error",map[string]string{"source":n.Name,"message":e.ToString()})
// }

// func (n *MNode) Heartbeat() {
// 	n.MSend("alive",map[string]interface{}{},header.MessageId)
// }

// func (n *MNode) MakeHeader(cmd, mid string) *MHeader {
// 	return &MHeader{
// 		Cmd: cmd,
// 		MessageId: mid}
// }

// func (n *MNode) Exit() {
// 	os.Exit()
// }

// func (n *MNode) MSend(name, msg, mid, cmd string) {
// 	fmt.Println("sending",name,msg,mid)
// 	header := n.MakeHeader(name,mid,cmd)
	
// }

// type MCNode struct {
// 	NodeId string
// 	GroupId string
// 	Transport *mprotocol.MangoTransport
// 	CurrentMid int
// 	State int // alive, stalled, dead
// }

// func (n *MCNode) heartbeat_worker() {
// 	for {
// 		time.Sleep(5*time.Second)
// 		n.Transport.Tx("asda")
// 	}
// }


// func (n *MCNode) send_message(command string, args interface{}) {
// 	header := n.MakeHeader(command)
// 	body := JSON.Marsal(args)
// 	n.Transport.Send(header, body)
// }


func (n *MNode) Run() {
	
}
