package main

import "fmt"

func main() {
	server := NewServer("127.0.0.1",9999)
	fmt.Println("Server Starting......")
	server.Start()
}