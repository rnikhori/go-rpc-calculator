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

## Stats

The server keeps a running, concurrency-safe count of how many times each operation has been called. Call `Calculator.Stats` with an empty args struct to get back a `calculator.Stats{Counts map[string]uint64, Total uint64}` snapshot.

## History

Every call to an arithmetic method is recorded as a `calculator.HistoryEntry` (operation name, operands, result, error if any, and timestamp). The server keeps the most recent entries, up to a fixed size, and discards older ones once that size is exceeded. Call `Calculator.GetHistory` with a `calculator.HistoryArgs{Limit: n}` to retrieve up to `n` of the most recent entries; `Limit` of 0 returns everything retained, and a negative `Limit` returns `ErrInvalidLimit`.

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

Running the client with no arguments connects to `localhost:1234` and runs a fixed demo sequence (Add, Subtract, Multiply, Divide on 10 and 3).

## Client CLI usage

Pass an operation name and its operands to call that RPC method directly and print only the result:

```sh
go run ./cmd/client add 10 3        # binary ops: add, subtract, multiply, divide, power, modulo
go run ./cmd/client sqrt 9          # unary ops: sqrt, negate
go run ./cmd/client negate 5
go run ./cmd/client history 5       # print the 5 most recent operations recorded by the server
go run ./cmd/client stats           # print total and per-operation call counts
```

An unknown operation name prints a short usage summary and exits non-zero.

## Graceful shutdown

The server stops accepting new connections on SIGINT (Ctrl+C) or SIGTERM, closes its listener, and waits for any in-flight requests to finish being served before the process exits.

## Flags

| Flag    | Binary          | Default     | Description                     |
|---------|-----------------|-------------|----------------------------------|
| `-host` | server, client  | `` (server, all interfaces) / `localhost` (client) | host/interface to bind or connect to |
| `-port` | server, client  | `1234`      | TCP port to listen on or connect to |

Example running both on a non-default port:

```sh
go run ./cmd/server -port 5000
go run ./cmd/client -port 5000
```
