package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Shopify/sarama"

	"github.com/ppatierno/kafka-go-examples/util"
)

func main() {

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

	bootstrapServers := strings.Split(util.GetEnv(util.BootstrapServers, "localhost:9092"), ",")
	topic := util.GetEnv(util.Topic, "my-topic")

	consumer, err := sarama.NewConsumer(bootstrapServers, nil)
	if err != nil {
		panic("Error creating the consumer")
	}

	defer func() {
		err := consumer.Close()
		if err != nil {
			fmt.Println("Error closing consumer: ", err)
			return
		}
		fmt.Println("Consumer closed")
	}()

	var partition int32 = 0
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		panic("Error creating the partition consumer")
	}

	defer func() {
		err := partitionConsumer.Close()
		if err != nil {
			fmt.Println("Error closing partition consumer: ", err)
			return
		}
		fmt.Println("Consumer partition closed")
	}()

consumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Received message from %s-%d [%d]: %s = %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		case sig := <-signals:
			fmt.Println("Got signal: ", sig)
			break consumerLoop
		}
	}
}
