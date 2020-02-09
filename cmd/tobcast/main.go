package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/underscorenico/tobcast/internal/consumer"
	"github.com/underscorenico/tobcast/internal/data"
	"github.com/underscorenico/tobcast/internal/producer"
	"github.com/underscorenico/tobcast/internal/timestamps"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	var ports []int
	for _, s := range arguments {
		i, _ := strconv.Atoi(s)
		ports = append(ports, i)
	}
	prod := producer.New(ports)

	myPort := ports[1]
	consumer := consumer.New(ports, prod)
	timestampsChan := make(chan data.Message)
	go consumer.ListenBroadcasted(ports[0], timestamps.New(ports), timestampsChan)

	timestamps := timestamps.New(ports)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("send: ")
		text, _ := reader.ReadString('\n')

		timestamps.Incr(myPort)
		ts := timestamps.Get(myPort)
		message := data.Message{
			Timestamp: ts,
			Value:     text,
		}
		prod.Broadcast(message)
	}
}
