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
	"sync"
	"syscall"

	"github.com/rnikhori/go-rpc-calculator/calculator"
)

const defaultMaxHistory = 100

func main() {
	host := flag.String("host", "", "host/interface to listen on (empty means all interfaces)")
	port := flag.Int("port", 1234, "port to listen on")
	maxHistory := flag.Int("max-history", defaultMaxHistory, "number of recent operations to retain for GetHistory")
	flag.Parse()

	calc := calculator.NewCalculator(log.Default(), *maxHistory)
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

	var conns sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				log.Println("waiting for in-flight connections to finish")
				conns.Wait()
				log.Println("server shut down")
				return
			default:
				log.Println("accept error:", err)
				continue
			}
		}
		conns.Add(1)
		go func() {
			defer conns.Done()
			rpc.ServeConn(conn)
		}()
	}
}
