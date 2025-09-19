package main

import "fmt"

func main() {
	server := NewSaver()
	err := server.Start()

	fmt.Print("TEst", server, err)
}
