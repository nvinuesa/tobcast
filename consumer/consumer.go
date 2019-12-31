package consumer

import (
	"bufio"
	"fmt"
	"github.com/underscorenico/tobcast/data"
	"github.com/underscorenico/tobcast/timestamps"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func Listen(port int, timestamps *timestamps.Timestamps) {
	l, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())

	var received []data.Message
	var delivered []data.Message

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(port, c)
	}
}

func handleConnection(port int, c net.Conn) {
	f, err := os.OpenFile(strconv.Itoa(port)+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	if _, err = f.WriteString("serving on port "+strconv.Itoa(port)+"\n"); err != nil {
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
		if _, err = f.WriteString(temp+"\n"); err != nil {
			panic(err)
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}
