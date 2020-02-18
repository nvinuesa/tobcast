package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/underscorenico/tobcast/internal/data"
)

type Consumer struct {
	port     int
	listener net.Listener
}

func NewConsumer(port int) *Consumer {
	listener := createListener(port)
	return &Consumer{port, listener}
}

func (c *Consumer) Register(handler func(message data.Message)) {
	for {
		conn, err := c.listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go c.handleBroadcastConnection(conn, handler)
	}
}

func createListener(port int) net.Listener {
	l, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	// defer l.Close()
	return l
}

func (c *Consumer) handleBroadcastConnection(conn net.Conn, handler func(message data.Message)) {
	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		go c.handleMessage(netData, err, handler)
	}
}

func (c *Consumer) handleMessage(netData string, err error, handler func(message data.Message)) {
	switch {
	case err == io.EOF:
		log.Println("reached EOF - close this connection.\n   ---")
		return
	case err != nil:
		log.Println("error reading message, got '"+netData+"'\n", err)
		return
	}
	temp := strings.TrimSpace(netData)
	var msg data.Message
	json.Unmarshal([]byte(temp), &msg)

	handler(msg)
}
