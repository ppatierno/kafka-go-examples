package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ppatierno/kafka-go-examples/util"
	kafka "github.com/segmentio/kafka-go"
)

func main() {

	bootstrapServers := strings.Split(util.GetEnv(util.BootstrapServers, "localhost:9092"), ",")
	topic := util.GetEnv(util.Topic, "my-topic")

	config := kafka.WriterConfig{
		Brokers: bootstrapServers,
		Topic:   topic}

	w := kafka.NewWriter(config)

	i := 1

	for {
		message := fmt.Sprintf("Message-%d", i)
		err := w.WriteMessages(context.Background(), kafka.Message{Value: []byte(message)})
		if err == nil {
			fmt.Println("Sent message: ", message)
		}
		i++
		time.Sleep(time.Second)
	}

	w.Close()
}
