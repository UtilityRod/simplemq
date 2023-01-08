package main

import (
	"fmt"
	"smq/core"
)

func main() {
	server, err := core.NewSMQServer("127.0.0.1", "44567")

	if err != nil {
		fmt.Println(err)
		return
	}

	server.RegisterMandatoryTopic("test")

	for {
		conn, err := server.Ln.Accept()

		if err != nil {
			fmt.Println(err)
			break
		}

		go server.ConnectionHandler(conn)
	}
}
