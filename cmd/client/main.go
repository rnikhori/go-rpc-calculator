package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"

	"github.com/rnikhori/go-rpc-calculator/calculator"
)

func main() {
	host := flag.String("host", "localhost", "server host to connect to")
	port := flag.Int("port", 1234, "server port to connect to")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer client.Close()

	args := calculator.Args{A: 10, B: 3}

	var sum float64
	if err := client.Call("Calculator.Add", args, &sum); err != nil {
		log.Fatal("Add error:", err)
	}
	log.Printf("%.2f + %.2f = %.2f\n", args.A, args.B, sum)

	var diff float64
	if err := client.Call("Calculator.Subtract", args, &diff); err != nil {
		log.Fatal("Subtract error:", err)
	}
	log.Printf("%.2f - %.2f = %.2f\n", args.A, args.B, diff)

	var product float64
	if err := client.Call("Calculator.Multiply", args, &product); err != nil {
		log.Fatal("Multiply error:", err)
	}
	log.Printf("%.2f * %.2f = %.2f\n", args.A, args.B, product)

	var quotient float64
	if err := client.Call("Calculator.Divide", args, &quotient); err != nil {
		log.Fatal("Divide error:", err)
	}
	log.Printf("%.2f / %.2f = %.2f\n", args.A, args.B, quotient)
}
