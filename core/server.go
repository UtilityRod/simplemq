package core

import (
	"fmt"
	"net"
	"smq/packets"
	"smq/utilities"
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
	DefaultUser    = "admin"
	DefaultPass    = "password"
)

type SMQClient struct {
	Conn       net.Conn
	ClientName string
	SubTopics  []string
}

type SMQServer struct {
	Addr    string
	Port    string
	Ln      net.Listener
	Handler *EventHandler
	Topics  map[string]*Topic
	Clients map[string]*SMQClient
	Auth    *utilities.Authenticator
}

type SubscribePayload struct {
	Payload *packets.Subscribe
	Client  *SMQClient
}

func NewSMQServer(addr, port string) (*SMQServer, error) {
	fmt.Printf("Starting SMQ Server on '%s:%s'\n", addr, port)
	ln, err := net.Listen("tcp", addr+":"+port)

	if err != nil {
		return nil, err
	}

	server := SMQServer{
		Addr:    addr,
		Port:    port,
		Ln:      ln,
		Topics:  make(map[string]*Topic),
		Clients: make(map[string]*SMQClient),
		Auth:    utilities.NewAuthenticator(DefaultUser, DefaultPass),
	}

	server.Handler = NewEventHandler()
	server.Handler.Start()

	return &server, nil
}

func NewSMQClient(name string, conn net.Conn) *SMQClient {
	client := SMQClient{
		Conn:       conn,
		ClientName: name,
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

	if server.Auth.Authenticate(connect.Username, connect.Password) != true {
		fmt.Printf("incorrect username or password from: %s\n", conn.RemoteAddr())
		return
	}

	// Create new client
	server.Clients[connect.ClientName] = NewSMQClient(connect.ClientName, conn)
	client := server.Clients[connect.ClientName]
	quit := false
	var event Event

	for quit != true {
		fixedHeader, err := packets.ReadFixedHeader(client.Conn)

		if err != nil {
			quit = true
		} else {
			switch fixedHeader.PacketType {
			case PublishType:
				event.Payload, err = packets.GetPublish(conn, fixedHeader.RemainingLen)
				event.Func = server.publishHandler
				break
			case SubscribeType:
				var subscribe *packets.Subscribe
				subscribe, err = packets.GetSubscribe(conn, fixedHeader.RemainingLen)
				event.Payload = SubscribePayload{
					Payload: subscribe,
					Client:  client,
				}
				event.Func = server.subscribeHandler
				break
			case DisconnectType:
				quit = true
				event.Func = server.disconnectHandler
				event.Payload = client
			default:
				quit = true
				fmt.Printf("invlaid packet type:'%d'\n", fixedHeader.PacketType)
			}

			if err != nil {
				fmt.Println(err)
				break
			}

			server.Handler.EventChannel <- event
			event.Func = nil
			event.Payload = nil
		}
	}

	fmt.Printf("Closed connection for %s\n", conn.RemoteAddr())
}

func (server *SMQServer) publishHandler(pubInt interface{}) {
	publish := pubInt.(*packets.Publish)
	topic := server.Topics[publish.Topic]
	for _, clientName := range topic.Clients {
		client := server.Clients[clientName]
		err := publish.Send(client.Conn)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func (server *SMQServer) subscribeHandler(paylodInt interface{}) {
	payload := paylodInt.(SubscribePayload)

	// Get the topic being subscribed too
	topicStr := payload.Payload.Topic
	topic := server.Topics[topicStr]
	// Check to see if client already subscribed to topic
	clientName := payload.Client.ClientName
	if _, ok := topic.Clients[clientName]; ok {
		// Client already subscribed
		return
	}

	// Update topic to contain client information
	topic.Clients[clientName] = payload.Client.ClientName
	// Update client's subscribed topics
	payload.Client.SubTopics = append(payload.Client.SubTopics, topicStr)
}

func (server *SMQServer) disconnectHandler(payloadInt interface{}) {
	client := payloadInt.(*SMQClient)

	// For every topic the client is subscribed too
	for _, topicStr := range client.SubTopics {
		topic := server.Topics[topicStr]
		// Delete the client from topic
		delete(topic.Clients, client.ClientName)
	}
}

// END OF SOURCE
