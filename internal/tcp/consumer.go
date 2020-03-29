package tcp

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/underscorenico/tobcast/internal/data"
)

type Consumer struct {
	port     int
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
}

func NewConsumer(port int) *Consumer {
	return &Consumer{port: port, quit: make(chan interface{})}
}

func createListener(port int) net.Listener {
	l, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalln("error opening tcp connection", err)
	}
	return l
}

func (c *Consumer) Stop() {
	log.Println("closing tcp listener")
	close(c.quit)
	c.listener.Close()
	c.wg.Wait()
}

func (c *Consumer) Start(handler func(message data.Message)) {

	if c.listener != nil {
		panic("tcp consumer already started")
	}
	c.listener = createListener(c.port)
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		for {
			select {
			case <-c.quit:
				return
			default:
				conn, err := c.listener.Accept()
				if err != nil {
					select {
					case <-c.quit:
						return
					default:
						log.Println("accept error", err)
					}
				}
				c.wg.Add(1)
				go c.handleBroadcastConnection(conn, handler)
				c.wg.Done()
			}
		}
	}()
}

func (c *Consumer) handleBroadcastConnection(conn net.Conn, handler func(message data.Message)) {
	defer conn.Close()

ReadLoop:
	for {
		select {
		case <-c.quit:
			return
		default:
			conn.SetDeadline(time.Now().Add(200 * time.Millisecond))
			netData, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue ReadLoop
				} else if err != io.EOF {
					log.Println("read error", err)
					return
				}
			}
			temp := strings.TrimSpace(netData)
			var msg data.Message
			json.Unmarshal([]byte(temp), &msg)

			handler(msg)
		}
	}
}
