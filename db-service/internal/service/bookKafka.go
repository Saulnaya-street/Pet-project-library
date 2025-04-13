package service

import (
	"awesomeProject22/db-service/internal/domain"
	"awesomeProject22/db-service/internal/kafka"
	"awesomeProject22/db-service/internal/repository"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
)

type BookServiceWithKafka struct {
	bookRepo      repository.IBookRepository
	eventProducer kafka.IEventProducer
}

func NewBookServiceWithKafka(bookRepo repository.IBookRepository, eventProducer kafka.IEventProducer) IBookService {
	return &BookServiceWithKafka{
		bookRepo:      bookRepo,
		eventProducer: eventProducer,
	}
}

func (s *BookServiceWithKafka) GetByID(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	return s.bookRepo.GetByID(ctx, id)
}

func (s *BookServiceWithKafka) GetAll(ctx context.Context, author, genre string) ([]*domain.Book, error) {
	return s.bookRepo.GetAll(ctx, author, genre)
}

func (s *BookServiceWithKafka) Create(ctx context.Context, book *domain.Book) error {
	if book.ID == uuid.Nil {
		book.ID = uuid.New()
	}

	if err := s.bookRepo.Create(ctx, book); err != nil {
		return err
	}

	if err := s.eventProducer.PublishBookCreated(ctx, book); err != nil {

		log.Printf("Ошибка публикации события создания книги: %v", err)
	} else {
		log.Printf("Опубликовано событие создания книги: %s (%s)", book.Name, book.ID)
	}

	return nil
}

func (s *BookServiceWithKafka) Update(ctx context.Context, book *domain.Book) error {

	_, err := s.bookRepo.GetByID(ctx, book.ID)
	if err != nil {
		return fmt.Errorf("книга для обновления не найдена: %w", err)
	}

	if err := s.bookRepo.Update(ctx, book); err != nil {
		return err
	}

	if err := s.eventProducer.PublishBookUpdated(ctx, book); err != nil {

		log.Printf("Ошибка публикации события обновления книги: %v", err)
	} else {
		log.Printf("Опубликовано событие обновления книги: %s (%s)", book.Name, book.ID)
	}

	return nil
}

func (s *BookServiceWithKafka) Delete(ctx context.Context, id uuid.UUID) error {

	book, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("книга для удаления не найдена: %w", err)
	}

	if err := s.bookRepo.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.eventProducer.PublishBookDeleted(ctx, id); err != nil {

		log.Printf("Ошибка публикации события удаления книги: %v", err)
	} else {
		log.Printf("Опубликовано событие удаления книги: %s (%s)", book.Name, id)
	}

	return nil
}
