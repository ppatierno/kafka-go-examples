package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ppatierno/kafka-go-examples/util"
	kafka "github.com/segmentio/kafka-go"
)

func main() {

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

	ctx, cancel := context.WithCancel(context.Background())

	// go routine for getting signals asynchronously
	go func() {
		sig := <-signals
		fmt.Println("Got signal: ", sig)
		cancel()
	}()

	bootstrapServers := strings.Split(util.GetEnv(util.BootstrapServers, "localhost:9092"), ",")
	topic := util.GetEnv(util.Topic, "my-topic")
	delayMs, _ := strconv.Atoi(util.GetEnv(util.DelayMs, strconv.Itoa(1000)))

	config := kafka.WriterConfig{
		Brokers:      bootstrapServers,
		Topic:        topic,
		BatchTimeout: 1 * time.Millisecond}

	w := kafka.NewWriter(config)

	fmt.Println("Producer configuration: ", config)

	i := 1

	defer func() {
		err := w.Close()
		if err != nil {
			fmt.Println("Error closing producer: ", err)
			return
		}
		fmt.Println("Producer closed")
	}()

	for {
		message := fmt.Sprintf("Message-%d", i)
		err := w.WriteMessages(ctx, kafka.Message{Value: []byte(message)})
		if err == nil {
			fmt.Println("Sent message: ", message)
		} else if err == context.Canceled {
			fmt.Println("Context canceled: ", err)
			break
		} else {
			fmt.Println("Error sending message: ", err)
		}
		i++

		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}
}
