package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/viper"
	"github.com/underscorenico/tobcast/pkg/config"
	"github.com/underscorenico/tobcast/pkg/tobcast"
)

func main() {

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		signal.Stop(signals)
	}()

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

	tobcast := tobcast.New(&config)

	msg := make(chan string, 1)
	go func() {
		for {
			fmt.Print("send: ")
			var s string
			fmt.Scan(&s)
			msg <- s
		}
	}()

	go func() {
		<-signals
		tobcast.Stop()
		os.Exit(1)
	}()

	for {
		s := <-msg
		text := strings.TrimSuffix(s, "\n")
		tobcast.Multicast(text)
	}
}
