package core

import (
	"errors"
	"fmt"
	"net"
	"smq/packets"
	"sync"
)

const (
	ConnectType    = 1
	PublishType    = 2
	SubscribeType  = 3
	DisconnectType = 4
	ConnectAck     = 5
	PublishAck     = 6
	SubscribeAck   = 7
	DiscconectAck  = 8
)

type SMQClient struct {
	Conn       net.Conn
	ClientName string
	Publish    chan packets.Publish
	Quit       chan bool
}

type Topic struct {
	Clients []*SMQClient
}

type SMQServer struct {
	Addr    string
	Port    string
	Ln      net.Listener
	Handler *EventHandler
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

	server.Handler = NewEventHandler()
	server.Handler.Start()

	return &server, nil
}

func NewSMQClient(name string, conn net.Conn) *SMQClient {
	client := SMQClient{
		Conn:       conn,
		ClientName: name,
		Publish:    make(chan packets.Publish),
		Quit:       make(chan bool),
	}

	return &client
}

func (server *SMQServer) ConnectionHandler(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Accepted connection for %s\n", conn.RemoteAddr())
	var err error
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
	for k := range (*server.Handler).Topics {
		topic := server.Handler.Topics[k]
		topic.Clients = append(topic.Clients, client)
		server.Handler.Topics[k] = topic
	}

	var wg sync.WaitGroup
	quit := false
	var event Event

	for quit != true {
		fixedHeader, err := packets.ReadFixedHeader(client.Conn)

		if err != nil {
			quit = true
			client.Quit <- quit
		} else {
			switch fixedHeader.PacketType {
			case PublishType:
				var publish *packets.Publish
				publish, err = packets.GetPublish(conn, fixedHeader.RemainingLen)
				event.EventType = PublishType
				event.Payload = publish
				server.Handler.EventChannel <- event
				break
			case SubscribeType:
				fmt.Println("Subscribe")
				break
			case DisconnectType:
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

func (server *SMQServer) RegisterTopic(name string) error {
	if _, ok := (*server).Handler.Topics[name]; ok {
		errStr := fmt.Sprintf("could not register topic '%s': already exists", name)
		return errors.New(errStr)
	}

	topic := Topic{
		Clients: make([]*SMQClient, 0),
	}

	server.Handler.Topics[name] = &topic
	fmt.Printf("New mandatory topic '%s' registered.\n", name)
	return nil
}

// END OF SOURCE
