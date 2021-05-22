package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"runtime"
	"time"
)

func main() {
	// Connect to a server
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	// Simple Publisher
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
		err := nc.Publish("lk", []byte(fmt.Sprintf("Hello World, %d", i)))
		print(i, "\n")
		if err != nil {
			print(err)
		}
	}
	runtime.Goexit()
}
