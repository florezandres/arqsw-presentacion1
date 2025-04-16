package kafka

import (
	"Taller2/Product/internal/command/domain/events"
	"Taller2/Product/internal/query/domain/models"
	"Taller2/Product/internal/query/domain/repositories"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type EventConsumer struct {
	reader      *kafka.Reader
	productRepo repositories.ProductReadRepository
}

func NewConsumer(brokers []string, topic string, repo repositories.ProductReadRepository) *EventConsumer {
	return &EventConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: "query_service",
		}),
		productRepo: repo,
	}
}

func (c *EventConsumer) Start(ctx context.Context) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		eventType := string(msg.Key)

		switch eventType {
		case "product_created":
			var event events.ProductCreatedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Error unmarshaling event: %v", err)
				continue
			}
			if err := c.handleProductCreated(ctx, event); err != nil {
				log.Printf("Error handling event: %v", err)
			}
		case "product_deleted":
			var event events.ProductDeletedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Error unmarshaling event: %v", err)
				continue
			}
			if err := c.handleProductDeleted(ctx, event); err != nil {
				log.Printf("Error handling event: %v", err)
			}
		case "product_updated":
			var event events.ProductUpdatedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Error unmarshaling update event: %v", err)
				continue
			}
			if err := c.handleProductUpdated(ctx, event); err != nil {
				log.Printf("Error handling update event: %v", err)
			}
		default:
			log.Printf("Unknown event type: %s", eventType)
		}
	}
}

func (c *EventConsumer) handleProductCreated(ctx context.Context, event events.ProductCreatedEvent) error {
	product := &models.ProductRead{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Price:       event.Price,
		Stock:       int32(event.Stock),
	}
	return c.productRepo.Save(ctx, product)
}

func (c *EventConsumer) handleProductDeleted(ctx context.Context, event events.ProductDeletedEvent) error {
	return c.productRepo.Delete(ctx, event.ProductID)
}

// Update desde el command
func (c *EventConsumer) handleProductUpdated(ctx context.Context, event events.ProductUpdatedEvent) error {
	product := &models.ProductRead{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Price:       event.Price,
		Stock:       int32(event.Stock),
	}
	return c.productRepo.Update(ctx, product)
}
