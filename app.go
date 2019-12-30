package main

import (
	"github.com/underscorenico/tobcast/consumer"
	"github.com/underscorenico/tobcast/data"
	producer "github.com/underscorenico/tobcast/producer"
	"strconv"
	"time"
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

	prod := producer.New(ports[2:])

	go consumer.Listen(ports[1])
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("send: ")
		text, _ := reader.ReadString('\n')
		prod.Multicast(time.Now(), text)
	}
}

func deliver(m data.Message, conn int) {

}
