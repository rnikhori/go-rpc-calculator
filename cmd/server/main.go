package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/rnikhori/go-rpc-calculator/calculator"
)

const defaultMaxHistory = 100

func main() {
	host := flag.String("host", "", "host/interface to listen on (empty means all interfaces)")
	port := flag.Int("port", 1234, "port to listen on")
	flag.Parse()

	calc := calculator.NewCalculator(log.Default(), defaultMaxHistory)
	if err := rpc.Register(calc); err != nil {
		log.Fatal("register error:", err)
	}

	addr := fmt.Sprintf("%s:%d", *host, *port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	log.Printf("calculator RPC server listening on %s\n", listener.Addr())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("shutdown signal received, closing listener")
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				log.Println("server stopped accepting connections, shutting down")
				return
			default:
				log.Println("accept error:", err)
				continue
			}
		}
		go rpc.ServeConn(conn)
	}
}
