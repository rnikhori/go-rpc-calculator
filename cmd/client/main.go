package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"strconv"
	"strings"

	"github.com/rnikhori/go-rpc-calculator/calculator"
)

var binaryOps = map[string]string{
	"add":      "Calculator.Add",
	"subtract": "Calculator.Subtract",
	"multiply": "Calculator.Multiply",
	"divide":   "Calculator.Divide",
	"power":    "Calculator.Power",
	"modulo":   "Calculator.Modulo",
}

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

	rest := flag.Args()
	if len(rest) == 0 {
		runDemo(client)
		return
	}

	if err := runOp(client, rest); err != nil {
		log.Fatal(err)
	}
}

func runOp(client *rpc.Client, args []string) error {
	op := strings.ToLower(args[0])
	method, ok := binaryOps[op]
	if !ok {
		return fmt.Errorf("unknown operation %q", args[0])
	}
	if len(args) != 3 {
		return fmt.Errorf("usage: client %s <a> <b>", op)
	}

	a, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return fmt.Errorf("invalid operand %q: %w", args[1], err)
	}
	b, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return fmt.Errorf("invalid operand %q: %w", args[2], err)
	}

	var reply float64
	if err := client.Call(method, calculator.Args{A: a, B: b}, &reply); err != nil {
		return fmt.Errorf("%s error: %w", op, err)
	}
	fmt.Println(reply)
	return nil
}

func runDemo(client *rpc.Client) {
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
