package calculator

import (
	"errors"
	"log"
	"math"
	"sync"
	"time"
)

var (
	ErrDivideByZero   = errors.New("division by zero")
	ErrNegativeSqrt   = errors.New("square root of negative number")
	ErrModuloByZero   = errors.New("modulo by zero")
	ErrInvalidOperand = errors.New("operand must be a finite number (not NaN or Inf)")
	ErrInvalidLimit   = errors.New("history limit must not be negative")
)

func invalidOperand(x float64) bool {
	return math.IsNaN(x) || math.IsInf(x, 0)
}

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

type opHistory struct {
	mu      sync.Mutex
	max     int
	entries []HistoryEntry
}

func (h *opHistory) add(entry HistoryEntry) {
	h.mu.Lock()
	defer h.mu.Unlock()

	max := h.max
	if max <= 0 {
		max = defaultMaxHistory
	}
	h.entries = append(h.entries, entry)
	if len(h.entries) > max {
		h.entries = h.entries[len(h.entries)-max:]
	}
}

func (h *opHistory) recent(limit int) []HistoryEntry {
	h.mu.Lock()
	defer h.mu.Unlock()

	n := len(h.entries)
	if limit <= 0 || limit > n {
		limit = n
	}
	result := make([]HistoryEntry, limit)
	copy(result, h.entries[n-limit:])
	return result
}

type Calculator struct {
	hist     opHistory
	statsMu  sync.Mutex
	opCounts map[string]uint64
	logger   *log.Logger
}

func NewCalculator(logger *log.Logger, maxHistory int) *Calculator {
	return &Calculator{logger: logger, hist: opHistory{max: maxHistory}}
}

func (c *Calculator) recordHistory(op string, a, b, result float64, err error) {
	entry := HistoryEntry{Op: op, A: a, B: b, Result: result, Time: time.Now()}
	if err != nil {
		entry.Err = err.Error()
	}
	c.hist.add(entry)

	c.statsMu.Lock()
	if c.opCounts == nil {
		c.opCounts = make(map[string]uint64)
	}
	c.opCounts[op]++
	c.statsMu.Unlock()

	if c.logger != nil {
		if err != nil {
			c.logger.Printf("%s(%v, %v) error: %v", op, a, b, err)
		} else {
			c.logger.Printf("%s(%v, %v) = %v", op, a, b, result)
		}
	}
}

func (c *Calculator) GetHistory(args HistoryArgs, reply *[]HistoryEntry) error {
	if args.Limit < 0 {
		return ErrInvalidLimit
	}
	*reply = c.hist.recent(args.Limit)
	return nil
}

func (c *Calculator) Stats(args struct{}, reply *Stats) error {
	c.statsMu.Lock()
	defer c.statsMu.Unlock()

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
	if invalidOperand(args.A) || invalidOperand(args.B) {
		c.recordHistory("Add", args.A, args.B, 0, ErrInvalidOperand)
		return ErrInvalidOperand
	}
	*reply = args.A + args.B
	c.recordHistory("Add", args.A, args.B, *reply, nil)
	return nil
}

func (c *Calculator) Subtract(args Args, reply *float64) error {
	if invalidOperand(args.A) || invalidOperand(args.B) {
		c.recordHistory("Subtract", args.A, args.B, 0, ErrInvalidOperand)
		return ErrInvalidOperand
	}
	*reply = args.A - args.B
	c.recordHistory("Subtract", args.A, args.B, *reply, nil)
	return nil
}

func (c *Calculator) Multiply(args Args, reply *float64) error {
	if invalidOperand(args.A) || invalidOperand(args.B) {
		c.recordHistory("Multiply", args.A, args.B, 0, ErrInvalidOperand)
		return ErrInvalidOperand
	}
	*reply = args.A * args.B
	c.recordHistory("Multiply", args.A, args.B, *reply, nil)
	return nil
}

func (c *Calculator) Divide(args Args, reply *float64) error {
	if invalidOperand(args.A) || invalidOperand(args.B) {
		c.recordHistory("Divide", args.A, args.B, 0, ErrInvalidOperand)
		return ErrInvalidOperand
	}
	if args.B == 0 {
		c.recordHistory("Divide", args.A, args.B, 0, ErrDivideByZero)
		return ErrDivideByZero
	}
	*reply = args.A / args.B
	c.recordHistory("Divide", args.A, args.B, *reply, nil)
	return nil
}

func (c *Calculator) Power(args Args, reply *float64) error {
	if invalidOperand(args.A) || invalidOperand(args.B) {
		c.recordHistory("Power", args.A, args.B, 0, ErrInvalidOperand)
		return ErrInvalidOperand
	}
	*reply = math.Pow(args.A, args.B)
	c.recordHistory("Power", args.A, args.B, *reply, nil)
	return nil
}

func (c *Calculator) Sqrt(args Args, reply *float64) error {
	if invalidOperand(args.A) {
		c.recordHistory("Sqrt", args.A, args.B, 0, ErrInvalidOperand)
		return ErrInvalidOperand
	}
	if args.A < 0 {
		c.recordHistory("Sqrt", args.A, args.B, 0, ErrNegativeSqrt)
		return ErrNegativeSqrt
	}
	*reply = math.Sqrt(args.A)
	c.recordHistory("Sqrt", args.A, args.B, *reply, nil)
	return nil
}

// Modulo returns the remainder of A/B using math.Mod, which takes the sign
// of the dividend A (matching Go's % operator and C's fmod), not the sign
// of the divisor B as in some languages (e.g. Python's %).
func (c *Calculator) Modulo(args Args, reply *float64) error {
	if invalidOperand(args.A) || invalidOperand(args.B) {
		c.recordHistory("Modulo", args.A, args.B, 0, ErrInvalidOperand)
		return ErrInvalidOperand
	}
	if args.B == 0 {
		c.recordHistory("Modulo", args.A, args.B, 0, ErrModuloByZero)
		return ErrModuloByZero
	}
	*reply = math.Mod(args.A, args.B)
	c.recordHistory("Modulo", args.A, args.B, *reply, nil)
	return nil
}

func (c *Calculator) Negate(args Args, reply *float64) error {
	if invalidOperand(args.A) {
		c.recordHistory("Negate", args.A, args.B, 0, ErrInvalidOperand)
		return ErrInvalidOperand
	}
	*reply = -args.A
	c.recordHistory("Negate", args.A, args.B, *reply, nil)
	return nil
}
