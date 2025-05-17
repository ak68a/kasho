package main

import "github/kasho/backend/api"

func main() {
	// api.NewServer(3000)

	server := api.NewServer(".")
	server.Start(3000)
}