package main

import (
	"flag"
	"log"

	"github.com/mtldglass/miris/client"
	"github.com/mtldglass/miris/server"
)

func main() {
	serverMode := flag.Bool("server", false, "Run in server mode")
	addr := flag.String("addr", ":1337", "In server mode: listen address, in client mode: the address of the server")
	userName := flag.String("user", "", "In client mode: your username")
	flag.Parse()

	if *serverMode {
		err := server.Serve(*addr)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		if *userName == "" {
			log.Fatal("please specify -user")
		}

		err := client.Client(*addr, *userName)
		if err != nil {
			log.Fatal(err)
		}
	}
}
