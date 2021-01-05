package calculator

import "errors"

type Args struct {
	A float64
	B float64
}

type Calculator struct{}

func (c *Calculator) Add(args Args, reply *float64) error {
	*reply = args.A + args.B
	return nil
}

func (c *Calculator) Subtract(args Args, reply *float64) error {
	*reply = args.A - args.B
	return nil
}

func (c *Calculator) Multiply(args Args, reply *float64) error {
	*reply = args.A * args.B
	return nil
}

func (c *Calculator) Divide(args Args, reply *float64) error {
	if args.B == 0 {
		return errors.New("division by zero")
	}
	*reply = args.A / args.B
	return nil
}
