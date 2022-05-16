package main

import (
	p "github.com/achilio/mv-manager-sse/broadcaster"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	p.Serve()
}
