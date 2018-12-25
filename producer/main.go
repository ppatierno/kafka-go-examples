package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ppatierno/kafka-go-examples/util"
	kafka "github.com/segmentio/kafka-go"
)

func main() {

	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

	// go routine for getting signals asynchronously
	go func() {
		sig := <-signals
		fmt.Println(sig)
		done <- true
	}()

	bootstrapServers := strings.Split(util.GetEnv(util.BootstrapServers, "localhost:9092"), ",")
	topic := util.GetEnv(util.Topic, "my-topic")

	config := kafka.WriterConfig{
		Brokers: bootstrapServers,
		Topic:   topic}

	w := kafka.NewWriter(config)

	i := 1

	var closing bool
	for {
		message := fmt.Sprintf("Message-%d", i)
		err := w.WriteMessages(context.Background(), kafka.Message{Value: []byte(message)})
		if err == nil {
			fmt.Println("Sent message: ", message)
		}
		i++

		// checking close or sleep
		select {
		case closing = <-done:
			fmt.Println("closing =", closing)
		default:
			time.Sleep(time.Second)
		}

		if closing {
			break
		}
	}

	fmt.Println("Closing")
	w.Close()
}
