package main

import (
	"fmt"
	"go-gRPC-sso/internal/config"
)

func main() {
	config := config.MustLoad()
	fmt.Println(config)

	// TODO: init logger - slog

	// TODO: inti app

	// TODO: run gRPC-server app
}
