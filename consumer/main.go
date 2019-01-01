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
	groupID := util.GetEnv(util.GroupID, "my-group")

	config := kafka.ReaderConfig{
		Brokers:  bootstrapServers,
		GroupID:  groupID,
		Topic:    topic,
		MaxWait:  500 * time.Millisecond,
		MinBytes: 1,
		MaxBytes: 1e6}

	r := kafka.NewReader(config)

	fmt.Println("Consumer configuration: ", config)

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			fmt.Println("Error reading message: ", err)
			break
		}
		fmt.Printf("Received message from %s-%d [%d]: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}

	err := r.Close()
	if err != nil {
		fmt.Println("Error closing consumer: ", err)
	}
}
