package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ppatierno/kafka-go-examples/util"
)

func main() {

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

	bootstrapServers := strings.Split(util.GetEnv(util.BootstrapServers, "localhost:9092"), ",")
	topic := util.GetEnv(util.Topic, "my-topic")
	delayMs, _ := strconv.Atoi(util.GetEnv(util.DelayMs, strconv.Itoa(1000)))

	config := sarama.NewConfig()
	// enabling the read from the Success() channel
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(bootstrapServers, config)
	if err != nil {
		panic("Error creating the sync producer")
	}

	i := 1

	defer func() {
		err := producer.Close()
		if err != nil {
			fmt.Println("Error closing producer: ", err)
			return
		}
		fmt.Println("Producer closed")
	}()

producerLoop:
	for {

		value := fmt.Sprintf("Message-%d", i)
		message := sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(value)}

		select {
		case producer.Input() <- &message:
			i++
		case err := <-producer.Errors():
			fmt.Println("Failed to produce message", err)
		case success := <-producer.Successes():
			fmt.Printf("Sent message value='%s' at partition = %d, offset = %d\n", success.Value, success.Partition, success.Offset)
		case sig := <-signals:
			fmt.Println("Got signal: ", sig)
			break producerLoop
		default:
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}
	}
}
