package broker

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func WaitKafkaConsumersGroupReadiness(brokerAddress string, topics ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Fatal("Kafka not ready in time")
		default:
			conn, err := kafka.Dial("tcp", brokerAddress)
			if err == nil {
				_, err = conn.ReadPartitions(topics...)
				conn.Close()
				if err == nil {
					return
				}
			}
			time.Sleep(5 * time.Second)
		}
	}
}
