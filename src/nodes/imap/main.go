package main

import (
	"fmt"
	"log"
	"time"
	"io"
	"io/ioutil"
	"strconv"
	"libmango"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
	"github.com/docopt/docopt-go"
)

type MboxInfo struct {
	server string
	username string
	password string
	last_seqnum int
	last_checked time.Time
	done chan bool
}

type MailAttachment struct {
	Name string `json:"name"`
	Content []byte `json:"content"`
}

type MailMessage struct {
	SentAt time.Time `json:"senttime"`
	ReceivedAt time.Time `json:"recvtime"`
	Subject string `json:"subject"`
	Body []interface{} `json:"body"`
	Attachments []interface{} `json:"attachments"`
	From []interface{} `json:"from"`
	To []interface{} `json:"to"`
}

func (msg MailMessage) ToObject() map[string]interface{} {
	attachments := make([]map[string]interface{}, 0)
	for _, ai := range msg.Attachments{
		a := ai.(MailAttachment)
		attachments = append(attachments, map[string]interface{}{
			"name":a.Name,
			"content":string(a.Content)})
	}
	ans := map[string]interface{}{
		"senttime":fmt.Sprintf("%d",msg.SentAt.UnixNano()),
		"recvtime":fmt.Sprintf("%d",msg.ReceivedAt.UnixNano()),
		"subject":msg.Subject,
		"body":msg.Body,
		"from":msg.From,
		"to":msg.To,
		"attachments":attachments}
	return ans
}

var (
	current_mbox *MboxInfo
	node *libmango.Node
)

func InitMbox(server, username, password string) *MboxInfo{
	mi := &MboxInfo{
		server:server,
		username:username,
		password:password,
		last_checked:time.Now()}
	mi.last_seqnum, _ = mi.count_messages()
	return mi
}

func (mi *MboxInfo) Stop() {
	mi.done <- true
}
		
func (mi *MboxInfo) count_messages() (int, error) {
	c, err := client.DialTLS(mi.server, nil)
	if err != nil {
		return -1, err
	}
	defer c.Logout()
	if err := c.Login(mi.username, mi.password); err != nil {
		return -1, err
	}
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return -1, err
	}
	return int(mbox.Messages), nil
}

func (mi *MboxInfo) fetch() []MailMessage {
	fmt.Println("Connecting to server...",mi.server)
	c, err := client.DialTLS(mi.server, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected")
	defer c.Logout()
	if err := c.Login(mi.username, mi.password); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Logged in")

	// List mailboxes
	// mailboxes := make(chan *imap.MailboxInfo, 10)
	// done := make(chan error, 1)
	// go func () {
	// 	done <- c.List("", "*", mailboxes)
	// }()

	// fmt.Println("Mailboxes:")
	// for m := range mailboxes {
	// 	fmt.Println("* " + m.Name)
	// }

	// if err := <-done; err != nil {
	// 	log.Fatal(err)
	// }

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Flags for INBOX:", mbox.Flags)

	if mbox.Messages == 0 {
		fmt.Println("No messages")
		return nil
	}
	if int(mbox.Messages) == mi.last_seqnum {
		fmt.Println("No new messages")
		return nil
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(uint32(mi.last_seqnum+1), mbox.Messages)

	var section imap.BodySectionName
	section.Peek = true
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, section.FetchItem(), imap.FetchBody, imap.FetchBodyStructure, imap.FetchUid, imap.FetchInternalDate}, messages)
	}()

	fmt.Println("Last messages:")
	ans := make([]MailMessage, 0)
	for msg := range messages {
		m := MailMessage{
			ReceivedAt: msg.InternalDate,
			Body: make([]interface{}, 0),
			Attachments: make([]interface{}, 0),
			From: make([]interface{}, 0),
			To: make([]interface{}, 0)}
		
		fmt.Println("* " + fmt.Sprintf("%d %d %v ", int(msg.SeqNum), int(msg.Uid), msg.InternalDate) + msg.Envelope.Subject)
		
		r := msg.GetBody(&section)
		if r == nil {
			log.Fatal("body not found")
		}
		
		// Create a new mail reader
		mr, _ := mail.CreateReader(r)
		header := mr.Header
		if date, err := header.Date(); err == nil {
			fmt.Println("Date:", date)
			m.SentAt = date
		}
		if from, err := header.AddressList("From"); err == nil {
			fmt.Println("From:", from)
			for _, a := range from {
				m.From = append(m.From, a.String())
			}
		}
		if to, err := header.AddressList("To"); err == nil {
			fmt.Println("To:", to)
			for _, a := range to {
				m.To = append(m.To, a.String())
			}
		}
		if subject, err := header.Subject(); err == nil {
			fmt.Println("Subject:", subject)
			m.Subject = subject
		}
		
		// Process each message's part
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			switch h := p.Header.(type) {
			case *mail.InlineHeader:
				// This is the message's text (can be plain-text or HTML)
				b, _ := ioutil.ReadAll(p.Body)
				fmt.Printf("Got text: %v", string(b))
				m.Body = append(m.Body, string(b))
			case *mail.AttachmentHeader:
				// This is an attachment
				filename, _ := h.Filename()
				fmt.Printf("Got attachment: %v", filename)
				bs, _ := ioutil.ReadAll(p.Body)
				fmt.Printf("Content of attachment: %v", bs)
				m.Attachments = append(m.Attachments, MailAttachment{Name: filename, Content: bs})
			}
		}
		ans = append(ans, m)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")
	mi.last_seqnum = int(mbox.Messages)
	return ans
}

func Logout(args map[string]interface{}) (string, map[string]interface{}, error) {
	current_mbox.Stop()
	return "",nil,nil
}

func main() {
	usage := `Usage: imap <server> <username> <password> [<interval>]`
	args, err := docopt.ParseDoc(usage)
	srv := args["<server>"].(string)
	un := args["<username>"].(string)
	pw := args["<password>"].(string)
	ival := 10
	if v, ok := args["<interval>"]; ok && v != nil {
		ival, _ = strconv.Atoi(v.(string))
	}
	
	n, err := libmango.NewNode("checkemail",map[string]libmango.MangoHandler{"logout":Logout})
	if err != nil {
		fmt.Println(err)
		return
	}
	go n.Start()

	current_mbox = InitMbox(srv, un, pw)
	ticker := time.NewTicker(time.Duration(ival) * time.Second)
	for {
		select {
		case <-current_mbox.done:
			return
		case <-ticker.C:
			msgs := current_mbox.fetch()
			for _, m := range msgs {
				n.Send("recv",m.ToObject())
			}
		}
	}
}
