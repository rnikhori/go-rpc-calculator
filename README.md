# go-rpc-calculator

A simple calculator service exposed over Go's `net/rpc` (Add, Subtract, Multiply, Divide).

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
