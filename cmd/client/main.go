package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strconv"
	"strings"
	"time"

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

var unaryOps = map[string]string{
	"sqrt":   "Calculator.Sqrt",
	"negate": "Calculator.Negate",
}

var callTimeout time.Duration

func main() {
	host := flag.String("host", "localhost", "server host to connect to")
	port := flag.Int("port", 1234, "server port to connect to")
	timeout := flag.Duration("timeout", 5*time.Second, "timeout for dialing and for each RPC call")
	flag.Parse()
	callTimeout = *timeout

	addr := net.JoinHostPort(*host, strconv.Itoa(*port))
	conn, err := net.DialTimeout("tcp", addr, callTimeout)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	client := rpc.NewClient(conn)
	defer client.Close()

	rest := flag.Args()
	if len(rest) == 0 {
		runDemo(client)
		return
	}

	switch strings.ToLower(rest[0]) {
	case "history":
		if err := runHistory(client, rest[1:]); err != nil {
			log.Fatal(err)
		}
		return
	case "stats":
		if err := runStats(client); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := runOp(client, rest); err != nil {
		log.Fatal(err)
	}
}

func callRPC(client *rpc.Client, method string, args, reply any) error {
	call := client.Go(method, args, reply, nil)
	select {
	case res := <-call.Done:
		return res.Error
	case <-time.After(callTimeout):
		return fmt.Errorf("rpc call %s timed out after %s", method, callTimeout)
	}
}

func runStats(client *rpc.Client) error {
	var stats calculator.Stats
	if err := callRPC(client, "Calculator.Stats", struct{}{}, &stats); err != nil {
		return fmt.Errorf("stats error: %w", err)
	}

	fmt.Printf("total operations: %d\n", stats.Total)
	for op, n := range stats.Counts {
		fmt.Printf("  %s: %d\n", op, n)
	}
	return nil
}

func runHistory(client *rpc.Client, args []string) error {
	limit := 0
	if len(args) > 0 {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid history limit %q: %w", args[0], err)
		}
		limit = n
	}

	var history []calculator.HistoryEntry
	if err := callRPC(client, "Calculator.GetHistory", calculator.HistoryArgs{Limit: limit}, &history); err != nil {
		return fmt.Errorf("history error: %w", err)
	}

	if len(history) == 0 {
		fmt.Println("no history recorded")
		return nil
	}
	for _, entry := range history {
		status := "ok"
		if entry.Err != "" {
			status = entry.Err
		}
		fmt.Printf("%s  %s(%v, %v) = %v  [%s]\n", entry.Time.Format(time.RFC3339), entry.Op, entry.A, entry.B, entry.Result, status)
	}
	return nil
}

func runOp(client *rpc.Client, args []string) error {
	op := strings.ToLower(args[0])

	if method, ok := unaryOps[op]; ok {
		if len(args) != 2 {
			return fmt.Errorf("usage: client %s <a>", op)
		}
		a, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return fmt.Errorf("invalid operand %q: %w", args[1], err)
		}
		return callAndPrint(client, method, op, calculator.Args{A: a})
	}

	method, ok := binaryOps[op]
	if !ok {
		return fmt.Errorf("unknown operation %q\n%s", args[0], usage())
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

	return callAndPrint(client, method, op, calculator.Args{A: a, B: b})
}

func callAndPrint(client *rpc.Client, method, op string, args calculator.Args) error {
	var reply float64
	if err := callRPC(client, method, args, &reply); err != nil {
		return fmt.Errorf("%s error: %w", op, err)
	}
	fmt.Println(reply)
	return nil
}

func usage() string {
	return "supported operations: add, subtract, multiply, divide, power, modulo (require <a> <b>); sqrt, negate (require <a>)"
}

func runDemo(client *rpc.Client) {
	args := calculator.Args{A: 10, B: 3}

	var sum float64
	if err := callRPC(client, "Calculator.Add", args, &sum); err != nil {
		log.Fatal("Add error:", err)
	}
	log.Printf("%.2f + %.2f = %.2f\n", args.A, args.B, sum)

	var diff float64
	if err := callRPC(client, "Calculator.Subtract", args, &diff); err != nil {
		log.Fatal("Subtract error:", err)
	}
	log.Printf("%.2f - %.2f = %.2f\n", args.A, args.B, diff)

	var product float64
	if err := callRPC(client, "Calculator.Multiply", args, &product); err != nil {
		log.Fatal("Multiply error:", err)
	}
	log.Printf("%.2f * %.2f = %.2f\n", args.A, args.B, product)

	var quotient float64
	if err := callRPC(client, "Calculator.Divide", args, &quotient); err != nil {
		log.Fatal("Divide error:", err)
	}
	log.Printf("%.2f / %.2f = %.2f\n", args.A, args.B, quotient)
}
