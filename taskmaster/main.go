package main

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry/yagnats"
	//"github.com/vito/gordon"
	"github.com/vito/taskmaster"
	"os"
	"os/signal"
	//"time"
)

type sshStartMessage struct {
	Session   string `json:"session"`
	PublicKey string `json:"public_key"`
}

type sshStopMessage struct {
	Session string `json:"session"`
}

func main() {
	fmt.Println("Waiting for connections...")

	sea := taskmaster.NewSEA("/tmp/warden.sock")

	nats := yagnats.NewClient()

	nats.Connect(&yagnats.ConnectionInfo{
		Addr: "127.0.0.1:4222",
	})

	nats.Subscribe("ssh.start", func(msg *yagnats.Message) {
		var start sshStartMessage

		err := json.Unmarshal([]byte(msg.Payload), &start)
		if err != nil {
			fmt.Println("Error unmarshalling:", err)
			return
		}

		err = sea.Start(start.Session, start.PublicKey)
		if err != nil {
			fmt.Println("Failed to start SSH session:", err)
			return
		}
	})

	nats.Subscribe("ssh.stop", func(msg *yagnats.Message) {
		var stop sshStopMessage

		err := json.Unmarshal([]byte(msg.Payload), &stop)
		if err != nil {
			fmt.Println("Error unmarshalling:", err)
			return
		}

		err = sea.Stop(stop.Session)
		if err != nil {
			fmt.Println("Tailed to stop SSH session:", err)
			return
		}
	})

	//go func() {
	//for {
	//for id, session := range sessions {
	//fmt.Println(id, session.Port)
	//}

	//time.Sleep(1 * time.Second)
	//}
	//}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	//for _, session := range sessions {
	//session.Container.Destroy()
	//}
}
