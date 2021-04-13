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
| `Calculator.GetHistory` | Limit (optional) | returns the most recent `Limit` operations (or all retained if `Limit` <= 0) |
| `Calculator.Stats`    | none        | returns per-operation call counts and the total operations served |

All arithmetic methods take a `calculator.Args{A, B float64}` struct and return a `float64` reply. Methods that only need one operand (`Sqrt`, `Negate`) ignore `B`.

## History

## Stats

The server keeps a running, concurrency-safe count of how many times each operation has been called. Call `Calculator.Stats` with an empty args struct to get back a `calculator.Stats{Counts map[string]uint64, Total uint64}` snapshot.

## History

Every call to an arithmetic method is recorded as a `calculator.HistoryEntry` (operation name, operands, result, error if any, and timestamp). The server keeps the most recent entries, up to a fixed size, and discards older ones once that size is exceeded. Call `Calculator.GetHistory` with a `calculator.HistoryArgs{Limit: n}` to retrieve up to `n` of the most recent entries.

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
