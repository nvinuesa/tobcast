package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/underscorenico/tobcast/internal/data"
)

type Consumer struct {
	port int
}

func NewConsumer(port int) *Consumer {
	return &Consumer{port}
}

func (c *Consumer) ListenBroadcasted(handler func(message data.Message)) {
	l, err := net.Listen("tcp4", ":"+strconv.Itoa(c.port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go c.handleBroadcastConnection(conn, handler)
	}
}

func (consumer *Consumer) handleBroadcastConnection(c net.Conn, handler func(message data.Message)) {
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
		var msg data.Message
		json.Unmarshal([]byte(temp), &msg)

		handler(msg)
	}
}
