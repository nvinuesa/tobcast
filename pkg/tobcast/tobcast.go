package tobcast

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/underscorenico/tobcast/internal/utils"

	"github.com/underscorenico/tobcast/internal/data"
	"github.com/underscorenico/tobcast/internal/tcp"
	"github.com/underscorenico/tobcast/internal/timestamps"
	"github.com/underscorenico/tobcast/pkg/config"
)

type Tobcast struct {
	producer      *tcp.Producer
	timestamps    *timestamps.Timestamps
	consumer      *tcp.Consumer
	deliverFile   *os.File
	config        *config.Config
	received      []data.Message
	delivered     []data.Message
	keepAliveFreq time.Duration
}

func New(config *config.Config) *Tobcast {
	clusterPorts := config.Cluster.Broadcast.Ports

	consumer := tcp.NewConsumer(config.Listen.Port)
	keepAliveFreq, err := time.ParseDuration(config.KeepAliveFreq)
	if err != nil {
		log.Println(config.KeepAliveFreq)
		log.Fatalln(err)
		panic("bad formatted keep alive frequency in config")
	}
	producer := tcp.NewProducer(clusterPorts)
	timestamps := timestamps.New(clusterPorts)

	f, err := os.OpenFile(strconv.Itoa(config.Listen.Port)+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	instance := &Tobcast{producer, timestamps, consumer, f, config, []data.Message{}, []data.Message{}, keepAliveFreq}
	consumer.Start(instance.handler)
	go instance.keepAlive()

	return instance
}

func (p *Tobcast) Stop() {
	p.consumer.Stop()
	p.producer.Stop()
	p.deliverFile.Close()
}

func (p *Tobcast) Multicast(text interface{}) error {
	myPort := p.config.Listen.Port
	ts := p.timestamps.Incr(myPort)
	message := &data.Message{
		Timestamp: ts,
		Value:     text,
		Sender:    myPort,
	}
	return p.producer.Broadcast(message)
}

/**
Method needed to keep the algorithm live. Without it, it would only deliver the messages
if and when received a broadcasted message.
This method sends empty messages that are not delivered to the user but ensure the delivery
of deliverable messages.
*/
func (p *Tobcast) keepAlive() {
	for {
		<-time.After(p.keepAliveFreq)
		p.Multicast("")
	}
}

func (p *Tobcast) handler(message data.Message) {
	p.timestamps.IncrOrSet(p.config.Listen.Port, message.Timestamp)
	p.timestamps.Set(message.Sender, message.Timestamp)

	p.received = append(p.received, message)

	p.order()
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

func (p *Tobcast) deliver(deliverable []data.Message) {
	for _, message := range deliverable {
		if message.Value != "" {
			bytes, err := json.Marshal(message)
			if err != nil {
				panic(err)
			}
			if _, err = p.deliverFile.WriteString(string(bytes) + "\n"); err != nil {
				panic(err)
			}
		}
	}
}

func (p *Tobcast) appendToDelivered(deliverable []data.Message) {
	p.delivered = append(p.delivered, deliverable...)
}
