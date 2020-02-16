package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/underscorenico/tobcast/package/config"
	"github.com/underscorenico/tobcast/package/tobcast"
)

func main() {
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

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("send: ")
		input, _ := reader.ReadString('\n')
		text := strings.TrimSuffix(input, "\n")
		tobcast.Multicast(text)
	}
}
