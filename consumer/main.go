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
	groupID := util.GetEnv(util.GroupID, "my-group")

	config := kafka.ReaderConfig{
		Brokers:  bootstrapServers,
		GroupID:  groupID,
		Topic:    topic,
		MaxWait:  500 * time.Millisecond,
		MinBytes: 1,
		MaxBytes: 1e6}

	r := kafka.NewReader(config)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("Received message from %v-%v [%v]: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}

	r.Close()
}
