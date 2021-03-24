package main

import (
	"time"

	srv "github.com/balazshorvath/go-srv"

	"CloudflareDDNS/server"
)

func main() {
	srv.CreateAndRunServer(server.New, 5*time.Second)
}
