package tobcast

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/underscorenico/tobcast/internal/utils"

	"github.com/underscorenico/tobcast/internal/data"
	"github.com/underscorenico/tobcast/internal/tcp"
	"github.com/underscorenico/tobcast/internal/timestamps"
	"github.com/underscorenico/tobcast/package/config"
)

type Tobcast struct {
	producer    *tcp.Producer
	timestamps  *timestamps.Timestamps
	consumer    *tcp.Consumer
	deliverFile *os.File
	config      *config.Config
	received    []data.Message
	delivered   []data.Message
}

func New(config *config.Config) *Tobcast {
	clusterPorts := config.Cluster.Broadcast.Ports

	producer := tcp.NewProducer(clusterPorts)
	timestamps := timestamps.New(clusterPorts)
	consumer := tcp.NewConsumer(config.Listen.Port)

	f, err := os.OpenFile(strconv.Itoa(config.Listen.Port)+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	instance := &Tobcast{producer, timestamps, consumer, f, config, []data.Message{}, []data.Message{}}

	go consumer.ListenBroadcasted(instance.handler)
	return instance
}

func (p *Tobcast) Multicast(text interface{}) error {
	myPort := p.config.Listen.Port
	ts := p.timestamps.Incr(myPort)
	message := data.Message{
		Timestamp: ts,
		Value:     text,
		Sender:    myPort,
	}
	p.producer.Broadcast(message)

	// FIXME
	return nil
}

func (p *Tobcast) handler(message data.Message) {
	log.Println("received msg from sender " + strconv.Itoa(message.Sender))
	p.timestamps.IncrOrSet(p.config.Listen.Port, message.Timestamp)
	p.timestamps.Set(message.Sender, message.Timestamp)

	p.received = append(p.received, message)

	p.order()
}

func (p *Tobcast) deliver(deliverable []data.Message) {
	for _, message := range deliverable {
		bytes, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}
		if _, err = p.deliverFile.WriteString(string(bytes) + "\n"); err != nil {
			panic(err)
		}
	}
}

func (p *Tobcast) order() {
	deliverable := []data.Message{}
	remaining := utils.ArrayDifference(p.received, p.delivered)

	for _, m := range remaining {
		if m.Timestamp <= p.timestamps.Min() {
			deliverable = append(deliverable, m)
		}
	}

	defer p.appendToDelivered(deliverable)
	p.deliver(deliverable)
}

func (p *Tobcast) appendToDelivered(deliverable []data.Message) {
	p.delivered = append(p.delivered, deliverable...)
}
