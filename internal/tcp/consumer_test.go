package tcp

import (
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/underscorenico/tobcast/internal/data"
	"gopkg.in/go-playground/assert.v1"
)

const (
	testPort = 8080
)

func TestAcceptIncoming(t *testing.T) {
	// Given
	consumer := NewConsumer(testPort)

	// When	/ Then
	consumer.Start(func(message data.Message) {})
	_, err := net.Dial("tcp", ":"+strconv.Itoa(testPort))
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
	consumer.Stop()
}

func TestCallHandlerOnMsg(t *testing.T) {
	// Given
	counter := make(chan int)
	handler := func(msg data.Message) {
		fmt.Println(msg)
		counter <- 1
	}
	cnt := 0
	consumer := NewConsumer(testPort)

	// When
	consumer.Start(handler)
	conn, err := net.Dial("tcp", ":"+strconv.Itoa(testPort))
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
	conn.Write([]byte("{\"timestamp\":0,\"value\":\"message 0\", \"sender\": 8080}\n"))
	cnt += <-counter
	conn.Write([]byte("{\"timestamp\":1,\"value\":\"message 1\", \"sender\": 8080}\n"))
	cnt += <-counter

	// Then
	assert.Equal(t, cnt, 2)
	consumer.Stop()
	conn.Close()
}
