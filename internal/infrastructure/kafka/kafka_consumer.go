package kafka

import (
	"context"
	"encoding/json"
	"log"

	"messaging/internal/domain/event"
	"messaging/internal/usecase"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	router *usecase.MessageRouter
}

func NewConsumer(brokers []string, topic, groupID string, router *usecase.MessageRouter) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader: reader,
		router: router,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	log.Println("Kafka consumer started...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer stopping...")
			return c.reader.Close()
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return c.reader.Close()
				}
				log.Printf("Error reading kafka message: %v", err)
				continue
			}

			// Unmarshal directly into the domain Event struct
			var ev event.Event
			if err := json.Unmarshal(m.Value, &ev); err != nil {
				log.Printf("Error unmarshalling message into Event: %v", err)
				continue
			}

			// Pass the entire Event to the router
			if err := c.router.Route(ctx, ev); err != nil {
				log.Printf("Error routing event type %s: %v", ev.EventType, err)
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
