# go-rpc-calculator

A calculator service exposed over Go's `net/rpc`, with operation history,
call-count stats, structured logging, and graceful shutdown.

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

## RPC methods

| Method                  | Args              | Notes                                           |
|--------------------------|-------------------|---------------------------------------------------|
| `Calculator.Add`          | A, B               | errors with `ErrInvalidOperand` on NaN/Inf          |
| `Calculator.Subtract`     | A, B               | errors with `ErrInvalidOperand` on NaN/Inf          |
| `Calculator.Multiply`     | A, B               | errors with `ErrInvalidOperand` on NaN/Inf          |
| `Calculator.Divide`       | A, B               | errors with `ErrDivideByZero` if B is 0; `ErrInvalidOperand` on NaN/Inf |
| `Calculator.Power`        | A, B               | A raised to the power of B; `ErrInvalidOperand` on NaN/Inf |
| `Calculator.Sqrt`         | A (B ignored)      | errors with `ErrNegativeSqrt` if A is negative; `ErrInvalidOperand` on NaN/Inf |
| `Calculator.Modulo`       | A, B               | sign follows the dividend A; errors with `ErrModuloByZero` if B is 0; `ErrInvalidOperand` on NaN/Inf |
| `Calculator.Negate`       | A (B ignored)      | returns -A; `ErrInvalidOperand` on NaN/Inf          |
| `Calculator.GetHistory`   | Limit (optional)  | returns the most recent `Limit` operations (0 means all retained); negative `Limit` returns `ErrInvalidLimit` |
| `Calculator.Stats`        | none               | returns per-operation call counts and the total operations served |

All arithmetic methods take a `calculator.Args{A, B float64}` struct and return a `float64` reply.

## History

Every call to an arithmetic method is recorded as a `calculator.HistoryEntry` (operation name, operands, result, error if any, and timestamp). The server keeps the most recent entries, up to a fixed size, and discards older ones once that size is exceeded. Call `Calculator.GetHistory` with a `calculator.HistoryArgs{Limit: n}` to retrieve up to `n` of the most recent entries.

## Stats

The server keeps a running, concurrency-safe count of how many times each operation has been called. Call `Calculator.Stats` with an empty args struct to get back a `calculator.Stats{Counts map[string]uint64, Total uint64}` snapshot.

## Logging

When built with `calculator.NewCalculator` (as `cmd/server` does), every operation call is logged with its name, operands, and result or error. Pass `-verbose` to the server to additionally log each TCP connection's remote address on accept and close.

## Graceful shutdown

The server stops accepting new connections on SIGINT (Ctrl+C) or SIGTERM, closes its listener, and waits for any in-flight requests to finish being served before the process exits.

## Flags

| Flag    | Binary          | Default     | Description                     |
|---------|-----------------|-------------|----------------------------------|
| `-host` | server, client  | `` (server, all interfaces) / `localhost` (client) | host/interface to bind or connect to |
| `-port` | server, client  | `1234`      | TCP port to listen on or connect to |
| `-max-history` | server  | `100`       | number of recent operations retained for `GetHistory` |
| `-timeout` | client      | `5s`        | timeout for dialing the server and for each RPC call   |
| `-verbose` | server      | `false`     | log each connection's remote address on accept and close |

Example running both on a non-default port:

```sh
go run ./cmd/server -port 5000
go run ./cmd/client -port 5000
```

## Testing

```sh
go test ./...          # unit tests plus TCP-level integration tests
go test -race ./...    # includes concurrency tests for stats/history
go test -bench=. ./calculator/...
```
