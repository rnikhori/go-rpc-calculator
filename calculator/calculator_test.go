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

func TestDivideByZero(t *testing.T) {
	c := &Calculator{}
	var reply float64
	err := c.Divide(Args{A: 10, B: 0}, &reply)
	if err == nil {
		t.Fatal("expected error for division by zero, got nil")
	}
}
