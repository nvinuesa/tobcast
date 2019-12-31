package main

import (
	"github.com/underscorenico/tobcast/consumer"
	"github.com/underscorenico/tobcast/data"
	Producer "github.com/underscorenico/tobcast/producer"
	Timestamps "github.com/underscorenico/tobcast/timestamps"
	"strconv"
)
import "fmt"
import "bufio"
import "os"

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
	prod := Producer.New(ports)

	myPort := ports[1]
	go consumer.Listen(myPort)

	timestamps := Timestamps.New(ports)

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
		prod.Multicast(message)
	}
}

func deliver(m data.Message, conn int) {

}
