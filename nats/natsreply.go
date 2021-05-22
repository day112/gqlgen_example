package main

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"runtime"
)

func main() {
	var url = "nats://127.0.0.1:4222"
	nc, err := nats.Connect(url, nats.Name("lk-nats"))
	if err != nil {
		log.Fatal("connect error")
	}

	nc.Subscribe("lk", func(mess *nats.Msg) {
		log.Println(string(mess.Data), "from nats")
		result, _ := json.Marshal(mess)
		log.Println("the reply info is ", string(result))
		nc.Publish(mess.Reply, []byte("lk can help you"))
	})

	//nc.Publish("lk", []byte("lk can help you"))

	runtime.Goexit()
}
