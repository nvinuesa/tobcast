package consumer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/underscorenico/tobcast/internal/data"

	"github.com/underscorenico/tobcast/internal/producer"
	"github.com/underscorenico/tobcast/internal/timestamps"
)

type Consumer struct {
	Ports       []int
	Connections []net.Conn
	producer    *producer.Producer
}

func New(ports []int, producer *producer.Producer) *Consumer {
	return &Consumer{Ports: ports, Connections: []net.Conn{}, producer: producer}
}

func (consumer *Consumer) ListenBroadcasted(port int, timestamps *timestamps.Timestamps, timestampsChan chan data.Message) {
	l, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())

	// var received []data.Message
	// var delivered []data.Message

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go consumer.handleBroadcastConnection(port, c, timestampsChan)
	}
}

func (consumer *Consumer) ListenDelivered(port int, timestamps *timestamps.Timestamps) {
	l, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go consumer.handleDeliveredConnection(port, c)
	}
}

func (consumer *Consumer) handleDeliveredConnection(port int, c net.Conn) {
	f, err := os.OpenFile(strconv.Itoa(port)+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	if _, err = f.WriteString("serving on port " + strconv.Itoa(port) + "\n"); err != nil {
		panic(err)
	}

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		switch {
		case err == io.EOF:
			log.Println("Reached EOF - close this connection.\n   ---")
			return
		case err != nil:
			log.Println("\nError reading command. Got: '"+netData+"'\n", err)
			return
		}

		temp := strings.TrimSpace(netData)
		if _, err = f.WriteString(temp + "\n"); err != nil {
			panic(err)
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (consumer *Consumer) handleBroadcastConnection(port int, c net.Conn, timestampsChan chan data.Message) {
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		switch {
		case err == io.EOF:
			log.Println("Reached EOF - close this connection.\n   ---")
			return
		case err != nil:
			log.Println("\nError reading command. Got: '"+netData+"'\n", err)
			return
		}

		_, pp, err := net.SplitHostPort(c.RemoteAddr().String())
		if err != nil {
			panic(err)
		}

		log.Println("incoming port: " + pp)
		temp := strings.TrimSpace(netData)
		var msg data.Message
		json.Unmarshal([]byte(temp), &msg)
		err = consumer.producer.Deliver(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}
