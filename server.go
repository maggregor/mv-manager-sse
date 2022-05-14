package main

import (
	p "github.com/achilio/mv-manager-sse/broadcaster"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	p.Serve()
}
