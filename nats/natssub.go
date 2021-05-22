package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"runtime"
)

func main() {
	// Connect to a server
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	_, err := nc.Subscribe("lk", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	if err != nil {
		fmt.Println(err)
	}
	runtime.Goexit()
}
