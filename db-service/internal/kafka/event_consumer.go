package kafka

import (
	"context"
	"encoding/json"
	"log"
)

type EventHandler func(ctx context.Context, event Event) error

type EventConsumer struct {
	client   IKafkaClient
	handlers map[EventType][]EventHandler
}

func NewEventConsumer(client IKafkaClient) *EventConsumer {
	return &EventConsumer{
		client:   client,
		handlers: make(map[EventType][]EventHandler),
	}
}

func (c *EventConsumer) RegisterHandler(eventType EventType, handler EventHandler) {
	if _, exists := c.handlers[eventType]; !exists {
		c.handlers[eventType] = []EventHandler{}
	}
	c.handlers[eventType] = append(c.handlers[eventType], handler)
	log.Printf("Зарегистрирован обработчик для события типа: %s", eventType)
}

func (c *EventConsumer) Start(ctx context.Context) error {
	log.Println("Запуск потребителя событий Kafka...")

	return c.client.Subscribe(ctx, func(key string, value []byte) error {
		log.Printf("Получено сообщение из Kafka. Ключ: %s", key)

		var event Event
		if err := json.Unmarshal(value, &event); err != nil {
			log.Printf("Ошибка десериализации события: %v. Содержимое: %s", err, string(value))
			return nil
		}

		log.Printf("Обработка события. Тип: %s, ID: %s", event.Type, event.ID)

		handlers, exists := c.handlers[event.Type]
		if !exists || len(handlers) == 0 {
			log.Printf("Нет зарегистрированных обработчиков для типа события: %s", event.Type)
			return nil
		}

		log.Printf("Найдено %d обработчиков для типа события %s", len(handlers), event.Type)
		for i, handler := range handlers {
			log.Printf("Запуск обработчика #%d для события %s", i+1, event.Type)
			if err := handler(ctx, event); err != nil {
				log.Printf("Ошибка обработки события %s: %v", event.Type, err)
				// Продолжаем обработку другими обработчиками
			} else {
				log.Printf("Обработчик #%d успешно обработал событие %s", i+1, event.Type)
			}
		}

		log.Printf("Завершена обработка события %s с ID %s", event.Type, event.ID)
		return nil
	})
}
