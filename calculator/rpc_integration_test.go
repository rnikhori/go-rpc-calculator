package calculator

import (
	"net"
	"net/rpc"
	"strings"
	"testing"
)

func TestServeOverTCP(t *testing.T) {
	calc := NewCalculator(nil, 10)
	server := rpc.NewServer()
	if err := server.Register(calc); err != nil {
		t.Fatalf("register error: %v", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen error: %v", err)
	}
	defer listener.Close()

	go server.Accept(listener)

	client, err := rpc.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer client.Close()

	var reply float64
	if err := client.Call("Calculator.Add", Args{A: 4, B: 5}, &reply); err != nil {
		t.Fatalf("Add call error: %v", err)
	}
	if reply != 9 {
		t.Errorf("Add(4, 5) over RPC = %v, want 9", reply)
	}
}

func TestErrorsPropagateOverTCP(t *testing.T) {
	calc := NewCalculator(nil, 10)
	server := rpc.NewServer()
	if err := server.Register(calc); err != nil {
		t.Fatalf("register error: %v", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen error: %v", err)
	}
	defer listener.Close()

	go server.Accept(listener)

	client, err := rpc.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer client.Close()

	var reply float64
	err = client.Call("Calculator.Divide", Args{A: 1, B: 0}, &reply)
	if err == nil || !strings.Contains(err.Error(), ErrDivideByZero.Error()) {
		t.Errorf("Divide(1, 0) over RPC error = %v, want it to mention %q", err, ErrDivideByZero.Error())
	}

	err = client.Call("Calculator.Sqrt", Args{A: -4}, &reply)
	if err == nil || !strings.Contains(err.Error(), ErrNegativeSqrt.Error()) {
		t.Errorf("Sqrt(-4) over RPC error = %v, want it to mention %q", err, ErrNegativeSqrt.Error())
	}
}
