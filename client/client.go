package client

import (
	"bufio"
	"log"
	"net"
	"os"

	"github.com/mtldglass/miris/protocol"
)

func Client(serverAddr, userName string) error {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return err
	}

	errs := make(chan error)
	go func() {
		for {
			msg, err := protocol.Read(conn)
			if err != nil {
				errs <- err
				return
			}

			switch msg.Type() {
			case protocol.MessageTypeChat:
				chatMsg := msg.(*protocol.ChatMessage)
				log.Printf("%s > %s", chatMsg.UserName, chatMsg.Message)
			default:
				log.Println(msg)
			}
		}
	}()

	for {
		select {
		case err := <-errs:
			return err
		default:
		}

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		err := protocol.Write(conn, &protocol.ChatMessage{
			UserName: userName,
			Message:  text,
		})
		if err != nil {
			return err
		}
	}
}
