package main

import (
	"github.com/NektarinR/godocker/pkg/server"
	"os"
	"strconv"
)

func main() {
	srv := server.Server{}
	evnvPort := os.Getenv("HTTP_PORT")
	port, err := strconv.Atoi(evnvPort)
	if err != nil {
		panic(err)
	}
	srv.Run(port)
}
