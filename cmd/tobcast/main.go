package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/underscorenico/tobcast/internal/consumer"
	"github.com/underscorenico/tobcast/internal/data"
	"github.com/underscorenico/tobcast/internal/producer"
	"github.com/underscorenico/tobcast/internal/timestamps"
)

func main() {
	arguments := os.Args[1:]
	if len(arguments) == 0 {
		fmt.Println("Please provide a port number!")
		return
	}

	var ports []int
	for _, s := range arguments {
		i, _ := strconv.Atoi(s)
		ports = append(ports, i)
	}
	var deliverPorts []int
	for _, p := range ports {
		deliverPorts = append(deliverPorts, p+10)
	}
	prod := producer.New(ports, deliverPorts)

	myPort := ports[1]
	consumer := consumer.New(ports, prod)
	incrementTimestamp := make(chan int)
	updateTimestamp := make(chan data.SenderWithTimestamp)

	timestamps := timestamps.New(ports)
	go monitor(timestamps, incrementTimestamp, updateTimestamp)
	go consumer.ListenBroadcasted(ports[0], timestamps, updateTimestamp)
	go consumer.ListenDelivered(deliverPorts[0], timestamps)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("send: ")
		input, _ := reader.ReadString('\n')
		text := strings.TrimSuffix(input, "\n")
		incrementTimestamp <- myPort
		ts := timestamps.Get(myPort)
		message := data.Message{
			Timestamp: ts,
			Value:     text,
		}
		prod.Broadcast(message)
	}
}

func monitor(ts *timestamps.Timestamps, incrementTimestamp chan int, updateTimestamp chan data.SenderWithTimestamp) {
	for {
		select {
		case incre := <-incrementTimestamp:
			ts.Incr(incre)
		case update := <-updateTimestamp:
			ts.Update(update.Port, update.Timestamp)
		}
	}
}
