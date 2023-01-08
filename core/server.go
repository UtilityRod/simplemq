package core

import (
	"errors"
	"fmt"
	"net"
	"smq/packets"
	"sync"
)

const (
	connect    = 1
	publish    = 2
	subscribe  = 3
	disconnect = 4
)

type Topic struct {
	In      chan<- Message
	Clients []*SMQClient
}

type Message struct {
	Topic   string
	Message string
}

type SMQClient struct {
	Conn       net.Conn
	ClientName string
	Message    chan Message
	Quit       chan bool
}

type SMQServer struct {
	Addr            string
	Port            string
	Ln              net.Listener
	MandatoryTopics map[string]Topic
}

func NewSMQServer(addr, port string) (*SMQServer, error) {
	fmt.Printf("Starting SMQ Server on '%s:%s'\n", addr, port)
	ln, err := net.Listen("tcp", addr+":"+port)

	if err != nil {
		return nil, err
	}

	server := SMQServer{
		Addr: addr,
		Port: port,
		Ln:   ln,
	}

	server.MandatoryTopics = make(map[string]Topic)

	return &server, nil
}

func NewSMQClient(name string, conn net.Conn) *SMQClient {
	client := SMQClient{
		Conn:       conn,
		ClientName: name,
		Message:    make(chan Message),
		Quit:       make(chan bool),
	}

	return &client
}

func (server *SMQServer) ConnectionHandler(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Accepted connection for %s\n", conn.RemoteAddr())
	fixedHeader, err := packets.ReadFixedHeader(conn)

	if err != nil {
		fmt.Println(err)
		return
	}

	connect, err := packets.ReadConnect(conn, fixedHeader.RemainingLen)

	if err != nil {
		fmt.Println(err)
		return
	}
	// Create new client
	client := NewSMQClient(connect.ClientName, conn)
	// Register client to all mandatory topics
	for k := range server.MandatoryTopics {
		topic := server.MandatoryTopics[k]
		topic.Clients = append(topic.Clients, client)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go client.ClientMessageHandler(&wg)
	quit := false

	for quit != true {
		fixedHeader, err := packets.ReadFixedHeader(client.Conn)

		if err != nil {
			quit = true
			client.Quit <- quit
		} else {
			switch fixedHeader.PacketType {
			case publish:
				fmt.Println("Publish")
				break
			case subscribe:
				fmt.Println("Subscribe")
				break
			case disconnect:
				quit = true
				client.Quit <- quit
			default:
				fmt.Printf("invlaid packet type:'%d'\n", fixedHeader.PacketType)
				quit = true
				client.Quit <- quit
			}
		}
	}

	wg.Wait()
	fmt.Printf("Closed connection for %s\n", conn.RemoteAddr())
}

func (client *SMQClient) ClientMessageHandler(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case msg := <-client.Message:
			fmt.Println("Message recieved: ", msg.Topic, " / ", msg.Message)
		case <-client.Quit:
			return
		}
	}
}

func (server *SMQServer) RegisterMandatoryTopic(name string) error {
	if _, ok := server.MandatoryTopics[name]; ok {
		errStr := fmt.Sprintf("could not register topic '%s': already exists", name)
		return errors.New(errStr)
	}

	topic := Topic{
		In: make(chan<- Message),
	}

	server.MandatoryTopics[name] = topic
	fmt.Printf("New mandatory topic '%s' registered.\n", name)
	return nil
}

// END OF SOURCE
