package main

import "github.com/NektarinR/godocker/pkg/server"

func main() {
	srv := server.Server{}

	srv.Run(8081)
}
