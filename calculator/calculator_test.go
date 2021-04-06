package calculator

import (
	"errors"
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

func TestHistoryTrimsToMaxSize(t *testing.T) {
	c := &Calculator{maxHistory: 3}
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

func BenchmarkMultiply(b *testing.B) {
	c := &Calculator{}
	var reply float64
	args := Args{A: 7, B: 6}
	for i := 0; i < b.N; i++ {
		_ = c.Multiply(args, &reply)
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
