package main

import (
	"log"
	"net"
	"net/rpc"

	"github.com/rnikhori/go-rpc-calculator/calculator"
)

func main() {
	if err := rpc.Register(new(calculator.Calculator)); err != nil {
		log.Fatal("register error:", err)
	}

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer listener.Close()

	log.Println("calculator RPC server listening on :1234")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
