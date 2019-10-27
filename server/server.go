package server

import (
	"log"
	"net"
	"sync"

	"github.com/mtldglass/miris/protocol"
)

func Serve(listenAddr string) error {
	clients := map[net.Conn]struct{}{}
	clientsLock := sync.Mutex{}

	list, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := list.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		clientsLock.Lock()
		clients[conn] = struct{}{}
		clientsLock.Unlock()

		go func(conn net.Conn) {
			defer func() {
				conn.Close()
				clientsLock.Lock()
				delete(clients, conn)
				clientsLock.Unlock()
			}()

			for {
				msg, err := protocol.Read(conn)
				if err != nil {
					log.Println(err)
					protocol.Write(conn, &protocol.ErrorMessage{
						Message: err.Error(),
					})
					return
				}

				switch msg.Type() {
				case protocol.MessageTypeError:
					errorMsg := msg.(*protocol.ErrorMessage)
					log.Println(errorMsg.Error())
				case protocol.MessageTypeChat:
					chatMsg := msg.(*protocol.ChatMessage)
					log.Printf("%s > %s", chatMsg.UserName, chatMsg.Message)

					clientsLock.Lock()
					for client := range clients {
						err := protocol.Write(client, msg)
						if err != nil {
							log.Println(err)
						}
					}
					clientsLock.Unlock()
				}
			}
		}(conn)
	}
}
