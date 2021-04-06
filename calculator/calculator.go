package calculator

import (
	"errors"
	"math"
	"sync"
	"time"
)

var (
	ErrDivideByZero = errors.New("division by zero")
	ErrNegativeSqrt = errors.New("square root of negative number")
	ErrModuloByZero = errors.New("modulo by zero")
)

const defaultMaxHistory = 100

type Args struct {
	A float64
	B float64
}

type HistoryEntry struct {
	Op     string
	A      float64
	B      float64
	Result float64
	Err    string
	Time   time.Time
}

type HistoryArgs struct {
	Limit int
}

type Stats struct {
	Counts map[string]uint64
	Total  uint64
}

type Calculator struct {
	mu         sync.Mutex
	history    []HistoryEntry
	maxHistory int
	opCounts   map[string]uint64
}

func (c *Calculator) recordHistory(op string, a, b, result float64, err error) {
	entry := HistoryEntry{Op: op, A: a, B: b, Result: result, Time: time.Now()}
	if err != nil {
		entry.Err = err.Error()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	max := c.maxHistory
	if max <= 0 {
		max = defaultMaxHistory
	}
	c.history = append(c.history, entry)
	if len(c.history) > max {
		c.history = c.history[len(c.history)-max:]
	}

	if c.opCounts == nil {
		c.opCounts = make(map[string]uint64)
	}
	c.opCounts[op]++
}

func (c *Calculator) GetHistory(args HistoryArgs, reply *[]HistoryEntry) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	n := len(c.history)
	limit := args.Limit
	if limit <= 0 || limit > n {
		limit = n
	}
	result := make([]HistoryEntry, limit)
	copy(result, c.history[n-limit:])
	*reply = result
	return nil
}

func (c *Calculator) Stats(args struct{}, reply *Stats) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	counts := make(map[string]uint64, len(c.opCounts))
	var total uint64
	for op, n := range c.opCounts {
		counts[op] = n
		total += n
	}
	*reply = Stats{Counts: counts, Total: total}
	return nil
}

func (c *Calculator) Add(args Args, reply *float64) error {
	*reply = args.A + args.B
	c.recordHistory("Add", args.A, args.B, *reply, nil)
	return nil
}

func (c *Calculator) Subtract(args Args, reply *float64) error {
	*reply = args.A - args.B
	c.recordHistory("Subtract", args.A, args.B, *reply, nil)
	return nil
}

func (c *Calculator) Multiply(args Args, reply *float64) error {
	*reply = args.A * args.B
	c.recordHistory("Multiply", args.A, args.B, *reply, nil)
	return nil
}

func (c *Calculator) Divide(args Args, reply *float64) error {
	if args.B == 0 {
		c.recordHistory("Divide", args.A, args.B, 0, ErrDivideByZero)
		return ErrDivideByZero
	}
	*reply = args.A / args.B
	c.recordHistory("Divide", args.A, args.B, *reply, nil)
	return nil
}

func (c *Calculator) Power(args Args, reply *float64) error {
	*reply = math.Pow(args.A, args.B)
	c.recordHistory("Power", args.A, args.B, *reply, nil)
	return nil
}

func (c *Calculator) Sqrt(args Args, reply *float64) error {
	if args.A < 0 {
		c.recordHistory("Sqrt", args.A, args.B, 0, ErrNegativeSqrt)
		return ErrNegativeSqrt
	}
	*reply = math.Sqrt(args.A)
	c.recordHistory("Sqrt", args.A, args.B, *reply, nil)
	return nil
}

func (c *Calculator) Modulo(args Args, reply *float64) error {
	if args.B == 0 {
		c.recordHistory("Modulo", args.A, args.B, 0, ErrModuloByZero)
		return ErrModuloByZero
	}
	*reply = math.Mod(args.A, args.B)
	c.recordHistory("Modulo", args.A, args.B, *reply, nil)
	return nil
}

func (c *Calculator) Negate(args Args, reply *float64) error {
	*reply = -args.A
	c.recordHistory("Negate", args.A, args.B, *reply, nil)
	return nil
}
