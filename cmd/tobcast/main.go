package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/underscorenico/tobcast/internal/config"
	"github.com/underscorenico/tobcast/internal/consumer"
	"github.com/underscorenico/tobcast/internal/data"
	"github.com/underscorenico/tobcast/internal/producer"
	"github.com/underscorenico/tobcast/internal/timestamps"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	var config config.Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file, %s", err)
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("bad formatted configuration, %v", err)
	}
	broadcastPorts := config.Cluster.Broadcast.Ports
	deliverPorts := config.Cluster.Deliver.Ports
	prod := producer.New(broadcastPorts, deliverPorts)

	myDeliverPort := config.Listen.Deliver.Port
	myBroadcastPort := config.Listen.Broadcast.Port
	consumer := consumer.New(prod)
	incrementTimestamp := make(chan int)
	updateTimestamp := make(chan data.SenderWithTimestamp)

	timestamps := timestamps.New(broadcastPorts)
	go monitor(timestamps, incrementTimestamp, updateTimestamp)
	go consumer.ListenBroadcasted(myBroadcastPort, timestamps, updateTimestamp)
	go consumer.ListenDelivered(myDeliverPort, timestamps)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("send: ")
		input, _ := reader.ReadString('\n')
		text := strings.TrimSuffix(input, "\n")
		incrementTimestamp <- myDeliverPort
		ts := timestamps.Get(myDeliverPort)
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
