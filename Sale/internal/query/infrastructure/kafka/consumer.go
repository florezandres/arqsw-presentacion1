package kafka

import (
	"Taller2/Sale/internal/query/domain/models"
	"context"
	"encoding/json"
	"log"
	"time"

	"Taller2/Sale/internal/command/domain/events"
	"Taller2/Sale/internal/query/infrastructure/persistence"
	"github.com/segmentio/kafka-go"
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

			var wrapper events.SaleEventWrapper
			if err := json.Unmarshal(msg.Value, &wrapper); err != nil {
				log.Printf("Error unmarshaling event wrapper: %v", err)
				continue
			}

			if err := c.handleEvent(ctx, wrapper); err != nil {
				log.Printf("Error handling event: %v", err)
			}
		}
	}
}

func (c *Consumer) handleEvent(ctx context.Context, wrapper events.SaleEventWrapper) error {
	switch wrapper.Type {
	case "sale_created":
		payloadBytes, err := json.Marshal(wrapper.Payload)
		if err != nil {
			return err
		}

		var event events.SaleCreatedEvent
		if err := json.Unmarshal(payloadBytes, &event); err != nil {
			return err
		}

		return c.repo.Create(ctx, &models.SaleRead{
			ID:        event.ID,
			ProductID: event.ProductID,
			Quantity:  int32(event.Quantity),
			Date:      time.Now().Format(time.RFC3339),
		})

	case "sale_updated":
		payloadBytes, err := json.Marshal(wrapper.Payload)
		if err != nil {
			return err
		}

		var event events.SaleUpdatedEvent
		if err := json.Unmarshal(payloadBytes, &event); err != nil {
			return err
		}

		return c.repo.Update(ctx, &models.SaleRead{
			ID:        event.ID,
			ProductID: event.ProductID,
			Quantity:  int32(event.Quantity),
			Date:      time.Now().Format(time.RFC3339),
		})

	case "sale_deleted":
		payloadBytes, err := json.Marshal(wrapper.Payload)
		if err != nil {
			return err
		}

		var event events.SaleDeletedEvent
		if err := json.Unmarshal(payloadBytes, &event); err != nil {
			return err
		}

		return c.repo.Delete(ctx, event.SaleID)

	default:
		return nil
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
