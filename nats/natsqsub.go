package main

import (
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"runtime"
	"time"
)

func main() {
	var url = "nats://127.0.0.1:4222"
	nc, err := nats.Connect(url, nats.Name("lk-nats"), nats.Timeout(time.Second*10), nats.MaxReconnects(3))
	if err != nil {
		log.Fatal("connect error")
	}

	i := 0
	nc.QueueSubscribe("lk", "queue", func(msg *nats.Msg) {
		i++
		printMsg(msg, i)
	})
	nc.Flush()
	runtime.Goexit()
}

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s] Queue[%s] Pid[%d]: '%s'", i, m.Subject, m.Sub.Queue, os.Getpid(), string(m.Data))
}
