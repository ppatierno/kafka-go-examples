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
	done := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

	go func() {
		sig := <-signals
		fmt.Println("Got signal: ", sig)
		done <- true
	}()

	bootstrapServers := strings.Split(util.GetEnv(util.BootstrapServers, "localhost:9092"), ",")
	topic := util.GetEnv(util.Topic, "my-topic")
	delayMs, _ := strconv.Atoi(util.GetEnv(util.DelayMs, strconv.Itoa(1000)))

	producer, err := sarama.NewSyncProducer(bootstrapServers, nil)
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

		partition, offset, err := producer.SendMessage(&message)
		if err != nil {
			fmt.Println("Error sending message: ", err)
		} else {
			fmt.Printf("Sent message value='%s' at partition = %d, offset = %d\n", value, partition, offset)
		}

		i++

		select {
		case closing := <-done:
			fmt.Println("Got closing: ", closing)
			break producerLoop
		default:
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}
	}
}
