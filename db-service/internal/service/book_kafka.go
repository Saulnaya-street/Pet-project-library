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

type BookServiceImpl struct {
	bookRepo      repository.IBookRepository
	eventProducer kafka.IEventProducer
}

func BookService(bookRepo repository.IBookRepository, eventProducer kafka.IEventProducer) IBookService {
	return &BookServiceImpl{
		bookRepo:      bookRepo,
		eventProducer: eventProducer,
	}
}

func (s *BookServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	book, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting book by ID: %w", err)
	}
	return book, nil
}

func (s *BookServiceImpl) GetAll(ctx context.Context, author, genre string) ([]*domain.Book, error) {
	books, err := s.bookRepo.GetAll(ctx, author, genre)
	if err != nil {
		return nil, fmt.Errorf("error getting list of books: %w", err)
	}
	return books, nil
}

func (s *BookServiceImpl) Create(ctx context.Context, book *domain.Book) error {
	if book.ID == uuid.Nil {
		book.ID = uuid.New()
	}

	if err := s.bookRepo.Create(ctx, book); err != nil {
		return fmt.Errorf("error creating book: %w", err)
	}

	if err := s.eventProducer.PublishBookCreated(ctx, book); err != nil {
		log.Printf("Error publishing book creation event: %v", err)
	} else {
		log.Printf("Book creation event published: %s (%s)", book.Name, book.ID)
	}

	return nil
}

func (s *BookServiceImpl) Update(ctx context.Context, book *domain.Book) error {

	_, err := s.bookRepo.GetByID(ctx, book.ID)
	if err != nil {
		return fmt.Errorf("Book for update not found: %w", err)
	}

	if err := s.bookRepo.Update(ctx, book); err != nil {
		return fmt.Errorf("error updating book: %w", err)
	}

	if err := s.eventProducer.PublishBookUpdated(ctx, book); err != nil {
		log.Printf("Error publishing book update event: %v", err)
	} else {
		log.Printf("Book update event published: %s (%s)", book.Name, book.ID)
	}

	return nil
}

func (s *BookServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {

	book, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("book to delete not found: %w", err)
	}

	if err := s.bookRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting book: %w", err)
	}

	if err := s.eventProducer.PublishBookDeleted(ctx, id); err != nil {
		log.Printf("Error publishing book deletion event: %v", err)
	} else {
		log.Printf("Book deletion event published: %s (%s)", book.Name, id)
	}

	return nil
}
