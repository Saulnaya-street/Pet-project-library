package kafka

import (
	"awesomeProject22/db-service/internal/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

type IEventProducer interface {
	PublishBookCreated(ctx context.Context, book *domain.Book) error
	PublishBookUpdated(ctx context.Context, book *domain.Book) error
	PublishBookDeleted(ctx context.Context, id uuid.UUID) error
	PublishUserCreated(ctx context.Context, user *domain.User) error
	PublishUserUpdated(ctx context.Context, user *domain.User) error
	PublishUserDeleted(ctx context.Context, id uuid.UUID) error
	PublishUserLoggedIn(ctx context.Context, userId uuid.UUID, username string) error
}

type EventProducer struct {
	client IKafkaClient
}

func NewEventProducer(client IKafkaClient) IEventProducer {
	return &EventProducer{
		client: client,
	}
}

func (p *EventProducer) PublishBookCreated(ctx context.Context, book *domain.Book) error {
	payload := BookEvent{
		Book: *book,
	}

	event := NewEvent(BookCreated, payload)
	return p.publishEvent(ctx, event)
}

func (p *EventProducer) PublishBookUpdated(ctx context.Context, book *domain.Book) error {
	payload := BookEvent{
		Book: *book,
	}

	event := NewEvent(BookUpdated, payload)
	return p.publishEvent(ctx, event)
}

func (p *EventProducer) PublishBookDeleted(ctx context.Context, id uuid.UUID) error {
	payload := map[string]string{
		"id": id.String(),
	}

	event := NewEvent(BookDeleted, payload)
	return p.publishEvent(ctx, event)
}

func (p *EventProducer) PublishUserCreated(ctx context.Context, user *domain.User) error {

	safeUser := *user
	safeUser.PasswordHash = ""

	payload := UserEvent{
		User: safeUser,
	}

	event := NewEvent(UserCreated, payload)
	return p.publishEvent(ctx, event)
}

func (p *EventProducer) PublishUserUpdated(ctx context.Context, user *domain.User) error {

	safeUser := *user
	safeUser.PasswordHash = ""

	payload := UserEvent{
		User: safeUser,
	}

	event := NewEvent(UserUpdated, payload)
	return p.publishEvent(ctx, event)
}

func (p *EventProducer) PublishUserDeleted(ctx context.Context, id uuid.UUID) error {
	payload := map[string]string{
		"id": id.String(),
	}

	event := NewEvent(UserDeleted, payload)
	return p.publishEvent(ctx, event)
}

func (p *EventProducer) PublishUserLoggedIn(ctx context.Context, userId uuid.UUID, username string) error {
	payload := LoginEvent{
		UserID:    userId,
		Username:  username,
		Timestamp: time.Now(),
	}

	event := NewEvent(UserLoggedIn, payload)
	return p.publishEvent(ctx, event)
}

func (p *EventProducer) publishEvent(ctx context.Context, event Event) error {
	eventData, err := event.Serialize()
	if err != nil {
		return fmt.Errorf("ошибка сериализации события: %w", err)
	}

	err = p.client.Publish(ctx, string(event.Type), eventData)
	if err != nil {
		return fmt.Errorf("ошибка публикации события: %w", err)
	}

	log.Printf("Опубликовано событие %s с ID %s", event.Type, event.ID)
	return nil
}
