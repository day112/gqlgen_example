package main

import (
	"github.com/nats-io/nats.go"
	"log"
	"runtime"
	"time"
)

func main() {
	var url = "nats://127.0.0.1:4222"
	nc, err := nats.Connect(url, nats.Name("lk-nats"))
	if err != nil {
		log.Fatal("connect error")
	}
	//nc.Subscribe("lk", func(mess *nats.Msg) {
	//	log.Println(string(mess.Data), "from nats")
	//	result, _ := json.Marshal(mess)
	//	log.Println("the reply info is ", string(result))
	//})

	message, err := nc.Request("lk", []byte("消息"), 1*time.Second)
	if err != nil {
		log.Println("get error, timeout", err)
	}

	log.Println("get data", string(message.Data))
	runtime.Goexit()
}
