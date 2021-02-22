# go-rpc-calculator

A simple calculator service exposed over Go's `net/rpc`.

## RPC methods

| Method              | Args        | Notes                                   |
|----------------------|-------------|------------------------------------------|
| `Calculator.Add`      | A, B        |                                            |
| `Calculator.Subtract` | A, B        |                                            |
| `Calculator.Multiply` | A, B        |                                            |
| `Calculator.Divide`   | A, B        | errors with `ErrDivideByZero` if B is 0   |
| `Calculator.Power`    | A, B        | A raised to the power of B                |
| `Calculator.Sqrt`     | A           | errors with `ErrNegativeSqrt` if A is negative |
| `Calculator.Modulo`   | A, B        | errors with `ErrModuloByZero` if B is 0   |
| `Calculator.Negate`   | A           | returns -A                                |

All methods take a `calculator.Args{A, B float64}` struct and return a `float64` reply. Methods that only need one operand (`Sqrt`, `Negate`) ignore `B`.

## Build

```sh
go build ./...
```

## Run

In one terminal, start the server:

```sh
go run ./cmd/server
```

In a second terminal, run the client:

```sh
go run ./cmd/client
```

The client connects to `localhost:1234` and prints the result of each operation.
