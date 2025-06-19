package main

import (
	"log"

	httpserver "github.com/kaasikodes/assessmate_backend/internal/adapters/http_server"
)

func main() {
	// start rest api
	log.Fatal(httpserver.Start())

	// TODO: Start grpc api if need be ...

}
