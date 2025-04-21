package kafka

import (
	"Taller2/Sale/internal/query/persistence/persistence"
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"log"
)

type Consumer struct {
	reader *kafka.Reader
	repo   *persistence.SQLiteSaleReadRepository
}

func NewConsumer(brokers []string, topic string, repo *persistence.SQLiteSaleReadRepository) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  "sale-query-group",
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
		repo: repo,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return c.reader.Close()
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			var event SaleEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Error unmarshaling event: %v", err)
				continue
			}

			if err := c.handleEvent(event); err != nil {
				log.Printf("Error handling event: %v", err)
			}
		}
	}
}

func (c *Consumer) handleEvent(event SaleEvent) error {
	switch event.Type {
	case "sale_created":
		payload, ok := event.Payload.(map[string]interface{})
		if !ok {
			return errors.New("invalid payload for sale_created")
		}
		return c.repo.Create(payload)
	case "sale_updated":
		payload, ok := event.Payload.(map[string]interface{})
		if !ok {
			return errors.New("invalid payload for sale_updated")
		}
		return c.repo.Update(payload)
	case "sale_deleted":
		id, ok := event.Payload.(string)
		if !ok {
			return errors.New("invalid payload for sale_deleted")
		}
		return c.repo.Delete(id)
	default:
		return nil
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
