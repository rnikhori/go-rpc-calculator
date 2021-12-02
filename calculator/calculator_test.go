package calculator

import (
	"bytes"
	"errors"
	"log"
	"math"
	"strings"
	"sync"
	"testing"
)

func TestAdd(t *testing.T) {
	c := &Calculator{}
	var reply float64
	if err := c.Add(Args{A: 2, B: 3}, &reply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != 5 {
		t.Errorf("Add(2, 3) = %v, want 5", reply)
	}
}

func TestSubtract(t *testing.T) {
	c := &Calculator{}
	var reply float64
	if err := c.Subtract(Args{A: 5, B: 3}, &reply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != 2 {
		t.Errorf("Subtract(5, 3) = %v, want 2", reply)
	}
}

func TestMultiply(t *testing.T) {
	c := &Calculator{}
	var reply float64
	if err := c.Multiply(Args{A: 4, B: 3}, &reply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != 12 {
		t.Errorf("Multiply(4, 3) = %v, want 12", reply)
	}
}

func TestDivide(t *testing.T) {
	c := &Calculator{}
	var reply float64
	if err := c.Divide(Args{A: 10, B: 2}, &reply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != 5 {
		t.Errorf("Divide(10, 2) = %v, want 5", reply)
	}
}

func TestPower(t *testing.T) {
	c := &Calculator{}
	var reply float64
	if err := c.Power(Args{A: 2, B: 10}, &reply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != 1024 {
		t.Errorf("Power(2, 10) = %v, want 1024", reply)
	}
}

func TestSqrt(t *testing.T) {
	c := &Calculator{}
	var reply float64
	if err := c.Sqrt(Args{A: 9}, &reply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != 3 {
		t.Errorf("Sqrt(9) = %v, want 3", reply)
	}
}

func TestSqrtNegative(t *testing.T) {
	c := &Calculator{}
	var reply float64
	err := c.Sqrt(Args{A: -4}, &reply)
	if !errors.Is(err, ErrNegativeSqrt) {
		t.Fatalf("Sqrt(-4) error = %v, want ErrNegativeSqrt", err)
	}
}

func TestModulo(t *testing.T) {
	c := &Calculator{}
	var reply float64
	if err := c.Modulo(Args{A: 10, B: 3}, &reply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != 1 {
		t.Errorf("Modulo(10, 3) = %v, want 1", reply)
	}
}

func TestModuloSignFollowsDividend(t *testing.T) {
	c := &Calculator{}
	var reply float64

	if err := c.Modulo(Args{A: -7, B: 3}, &reply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != -1 {
		t.Errorf("Modulo(-7, 3) = %v, want -1 (sign follows dividend)", reply)
	}

	if err := c.Modulo(Args{A: 7, B: -3}, &reply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != 1 {
		t.Errorf("Modulo(7, -3) = %v, want 1 (sign follows dividend)", reply)
	}
}

func TestModuloByZero(t *testing.T) {
	c := &Calculator{}
	var reply float64
	err := c.Modulo(Args{A: 10, B: 0}, &reply)
	if !errors.Is(err, ErrModuloByZero) {
		t.Fatalf("Modulo(10, 0) error = %v, want ErrModuloByZero", err)
	}
}

func TestNegate(t *testing.T) {
	c := &Calculator{}
	var reply float64
	if err := c.Negate(Args{A: 7}, &reply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != -7 {
		t.Errorf("Negate(7) = %v, want -7", reply)
	}
	if err := c.Negate(Args{A: -7}, &reply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != 7 {
		t.Errorf("Negate(-7) = %v, want 7", reply)
	}
}

func TestGetHistoryRecordsOperations(t *testing.T) {
	c := &Calculator{}
	var reply float64
	if err := c.Add(Args{A: 1, B: 2}, &reply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := c.Divide(Args{A: 1, B: 0}, &reply); !errors.Is(err, ErrDivideByZero) {
		t.Fatalf("Divide(1, 0) error = %v, want ErrDivideByZero", err)
	}

	var history []HistoryEntry
	if err := c.GetHistory(HistoryArgs{}, &history); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("len(history) = %d, want 2", len(history))
	}
	if history[0].Op != "Add" || history[0].Result != 3 || history[0].Err != "" {
		t.Errorf("history[0] = %+v, want Add with result 3 and no error", history[0])
	}
	if history[1].Op != "Divide" || history[1].Err != ErrDivideByZero.Error() {
		t.Errorf("history[1] = %+v, want Divide with ErrDivideByZero", history[1])
	}
}

func TestHistoryRecordsEveryErrorKind(t *testing.T) {
	c := &Calculator{}
	var reply float64

	_ = c.Sqrt(Args{A: -9}, &reply)
	_ = c.Modulo(Args{A: 1, B: 0}, &reply)
	_ = c.Add(Args{A: math.NaN(), B: 1}, &reply)

	var history []HistoryEntry
	if err := c.GetHistory(HistoryArgs{}, &history); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(history) != 3 {
		t.Fatalf("len(history) = %d, want 3", len(history))
	}

	wantErrs := []string{ErrNegativeSqrt.Error(), ErrModuloByZero.Error(), ErrInvalidOperand.Error()}
	for i, want := range wantErrs {
		if history[i].Err != want {
			t.Errorf("history[%d].Err = %q, want %q", i, history[i].Err, want)
		}
		if history[i].Result != 0 {
			t.Errorf("history[%d].Result = %v, want 0 for a failed call", i, history[i].Result)
		}
	}
}

func TestGetHistoryRespectsLimit(t *testing.T) {
	c := &Calculator{}
	var reply float64
	for i := 0; i < 5; i++ {
		_ = c.Add(Args{A: float64(i), B: 1}, &reply)
	}

	var history []HistoryEntry
	if err := c.GetHistory(HistoryArgs{Limit: 2}, &history); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("len(history) = %d, want 2", len(history))
	}
	if history[0].A != 3 || history[1].A != 4 {
		t.Errorf("history = %+v, want the two most recent entries", history)
	}
}

func TestGetHistoryRejectsNegativeLimit(t *testing.T) {
	c := &Calculator{}
	var history []HistoryEntry
	err := c.GetHistory(HistoryArgs{Limit: -1}, &history)
	if !errors.Is(err, ErrInvalidLimit) {
		t.Fatalf("GetHistory(Limit: -1) error = %v, want ErrInvalidLimit", err)
	}
}

func TestHistoryTrimsToMaxSize(t *testing.T) {
	c := &Calculator{hist: opHistory{max: 3}}
	var reply float64
	for i := 0; i < 10; i++ {
		_ = c.Add(Args{A: float64(i), B: 0}, &reply)
	}

	var history []HistoryEntry
	if err := c.GetHistory(HistoryArgs{}, &history); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(history) != 3 {
		t.Fatalf("len(history) = %d, want 3 (trimmed to maxHistory)", len(history))
	}
	if history[0].A != 7 || history[2].A != 9 {
		t.Errorf("history = %+v, want the 3 most recent entries", history)
	}
}

func TestInvalidOperandsRejected(t *testing.T) {
	c := &Calculator{}
	var reply float64

	cases := []struct {
		name string
		call func() error
	}{
		{"Add NaN", func() error { return c.Add(Args{A: math.NaN(), B: 1}, &reply) }},
		{"Subtract +Inf", func() error { return c.Subtract(Args{A: math.Inf(1), B: 1}, &reply) }},
		{"Multiply -Inf", func() error { return c.Multiply(Args{A: 1, B: math.Inf(-1)}, &reply) }},
		{"Divide NaN", func() error { return c.Divide(Args{A: math.NaN(), B: 1}, &reply) }},
		{"Power Inf", func() error { return c.Power(Args{A: 2, B: math.Inf(1)}, &reply) }},
		{"Sqrt NaN", func() error { return c.Sqrt(Args{A: math.NaN()}, &reply) }},
		{"Modulo Inf", func() error { return c.Modulo(Args{A: 1, B: math.Inf(1)}, &reply) }},
		{"Negate NaN", func() error { return c.Negate(Args{A: math.NaN()}, &reply) }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); !errors.Is(err, ErrInvalidOperand) {
				t.Errorf("error = %v, want ErrInvalidOperand", err)
			}
		})
	}
}

func TestStatsTracksCountsPerOperation(t *testing.T) {
	c := &Calculator{}
	var reply float64
	_ = c.Add(Args{A: 1, B: 1}, &reply)
	_ = c.Add(Args{A: 2, B: 2}, &reply)
	_ = c.Subtract(Args{A: 3, B: 1}, &reply)

	var stats Stats
	if err := c.Stats(struct{}{}, &stats); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Total != 3 {
		t.Errorf("stats.Total = %d, want 3", stats.Total)
	}
	if stats.Counts["Add"] != 2 {
		t.Errorf("stats.Counts[Add] = %d, want 2", stats.Counts["Add"])
	}
	if stats.Counts["Subtract"] != 1 {
		t.Errorf("stats.Counts[Subtract] = %d, want 1", stats.Counts["Subtract"])
	}
}

func TestStatsIsConcurrencySafe(t *testing.T) {
	c := &Calculator{}
	const goroutines = 50
	const callsEach = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			var reply float64
			for j := 0; j < callsEach; j++ {
				_ = c.Add(Args{A: 1, B: 1}, &reply)
			}
		}()
	}
	wg.Wait()

	var stats Stats
	if err := c.Stats(struct{}{}, &stats); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := uint64(goroutines * callsEach)
	if stats.Total != want {
		t.Errorf("stats.Total = %d, want %d", stats.Total, want)
	}
	if stats.Counts["Add"] != want {
		t.Errorf("stats.Counts[Add] = %d, want %d", stats.Counts["Add"], want)
	}
}

func TestNewCalculatorLogsOperations(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	c := NewCalculator(logger, 0)

	var reply float64
	if err := c.Add(Args{A: 1, B: 2}, &reply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Add(1, 2) = 3") {
		t.Errorf("log output = %q, want it to contain %q", buf.String(), "Add(1, 2) = 3")
	}

	buf.Reset()
	if err := c.Divide(Args{A: 1, B: 0}, &reply); !errors.Is(err, ErrDivideByZero) {
		t.Fatalf("Divide(1, 0) error = %v, want ErrDivideByZero", err)
	}
	if !strings.Contains(buf.String(), "Divide(1, 0) error") {
		t.Errorf("log output = %q, want it to mention the Divide error", buf.String())
	}
}

func BenchmarkMultiply(b *testing.B) {
	c := &Calculator{}
	var reply float64
	args := Args{A: 7, B: 6}
	for i := 0; i < b.N; i++ {
		_ = c.Multiply(args, &reply)
	}
}

func BenchmarkSqrt(b *testing.B) {
	c := &Calculator{}
	var reply float64
	args := Args{A: 2}
	for i := 0; i < b.N; i++ {
		_ = c.Sqrt(args, &reply)
	}
}

func TestDivideByZero(t *testing.T) {
	c := &Calculator{}
	var reply float64
	err := c.Divide(Args{A: 10, B: 0}, &reply)
	if err == nil {
		t.Fatal("expected error for division by zero, got nil")
	}
}
