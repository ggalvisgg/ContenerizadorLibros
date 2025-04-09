package services

import (
    "fmt"
    "example.com/go-mongo-app/models"
    "example.com/go-mongo-app/repositories"
)

type BookServiceInterface interface {
    GetBooks() ([]models.Book, error)
    GetBookByID(id string) (*models.Book, error)
    AddBook(book models.Book) (*models.Book, error)
    UpdateBook(book *models.Book) (*models.Book, error)
    DeleteBookByID(id string) error
}

type BookService struct {
    repo *repositories.BookRepository
}

func NewBookService(repo *repositories.BookRepository) *BookService {
    return &BookService{repo}
}

func (s *BookService) GetBooks() ([]models.Book, error) {
    return s.repo.GetAllBooks()
}

func (s *BookService) GetBookByID(id string) (*models.Book, error) {
    return s.repo.GetBookByID(id)
}

func (s *BookService) AddBook(book models.Book) (*models.Book, error) {
    return s.repo.CreateBook(book)
}

func (s *BookService) UpdateBook(book *models.Book) (*models.Book, error) {
    return s.repo.UpdateBook(book)
}

func (s *BookService) DeleteBookByID(id string) error {
    deleted, err := s.repo.RemoveBookByID(id)
    if err != nil {
        return err
    }
    if !deleted {
        return fmt.Errorf("libro no encontrado")
    }
    return nil
}