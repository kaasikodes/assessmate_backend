package main

import (
	"log"

	httpserver "github.com/kaasikodes/assessmate_backend/internal/adapters/http_server"
)

func main() {
	log.Fatal(httpserver.Start())

}
