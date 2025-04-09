package controllers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "fmt"

    "github.com/gorilla/mux"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "example.com/go-mongo-app/models"
)

// ---------------------- MOCK DEL SERVICIO ----------------------


type MockBookService struct {
    mock.Mock
}

func (m *MockBookService) GetBooks() ([]models.Book, error) {
    args := m.Called()
    return args.Get(0).([]models.Book), args.Error(1)
}
func (m *MockBookService) GetBookByID(id string) (*models.Book, error) {

    args := m.Called(id)
    return args.Get(0).(*models.Book), args.Error(1)
}
func (m *MockBookService) AddBook(book models.Book) (*models.Book, error) {
    args := m.Called(book)
    return args.Get(0).(*models.Book), args.Error(1)
}
func (m *MockBookService) UpdateBook(book *models.Book) (*models.Book, error) {
    args := m.Called(book)
    return args.Get(0).(*models.Book), args.Error(1)
}
func (m *MockBookService) DeleteBookByID(id string) error {
    args := m.Called(id)
    return args.Error(0)
}



// ---------------------- TESTS ----------------------


func TestCreateBook_Success(t *testing.T) {
    mockService := new(MockBookService)
    controller := NewBookController(mockService)

    book := models.Book{
        Title:  "Go Programming",
        Author: "John Doe",
        ISBN:   "1234567890",
    }

    body, _ := json.Marshal(book)
    req := httptest.NewRequest("POST", "/books", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    resp := httptest.NewRecorder()

    mockService.On("AddBook", mock.MatchedBy(func(b models.Book) bool {
        return b.Title == book.Title && b.Author == book.Author && b.ISBN == book.ISBN
    })).Return(&models.Book{
        ID:     primitive.NewObjectID(),
        Title:  book.Title,
        Author: book.Author,
        ISBN:   book.ISBN,
    }, nil)

    controller.CreateBook(resp, req)

    assert.Equal(t, http.StatusCreated, resp.Code)
    mockService.AssertExpectations(t)
}

func TestGetBooks_Success(t *testing.T) {
    mockService := new(MockBookService)
    controller := NewBookController(mockService)

    books := []models.Book{
        {Title: "Book 1", Author: "Author 1", ISBN: "1111111111"},
        {Title: "Book 2", Author: "Author 2", ISBN: "2222222222"},
    }

    mockService.On("GetBooks").Return(books, nil)

    req := httptest.NewRequest("GET", "/books", nil)
    resp := httptest.NewRecorder()

    controller.GetBooks(resp, req)

    assert.Equal(t, http.StatusOK, resp.Code)
    mockService.AssertExpectations(t)
}

func TestGetBookByID_Success(t *testing.T) {
    mockService := new(MockBookService)
    controller := NewBookController(mockService)

    id := primitive.NewObjectID().Hex()
    book := &models.Book{ID: primitive.NewObjectID(), Title: "Book Title", Author: "Author Name", ISBN: "1234567890"}

    mockService.On("GetBookByID", id).Return(book, nil)

    req := httptest.NewRequest("GET", "/books/"+id, nil)
    req = mux.SetURLVars(req, map[string]string{"id": id})
    resp := httptest.NewRecorder()

    controller.GetBookByID(resp, req)

    assert.Equal(t, http.StatusOK, resp.Code)
    mockService.AssertExpectations(t)
}

func TestDeleteBook_Success(t *testing.T) {
    mockService := new(MockBookService)
    controller := NewBookController(mockService)

    id := primitive.NewObjectID().Hex()

    mockService.On("DeleteBookByID", mock.Anything).Return(nil)

    req := httptest.NewRequest("DELETE", "/books/"+id, nil)
    req = mux.SetURLVars(req, map[string]string{"id": id})
    resp := httptest.NewRecorder()

    controller.DeleteBook(resp, req)

    assert.Equal(t, http.StatusNoContent, resp.Code)
    mockService.AssertExpectations(t)
}

func TestUpdateBook_Success(t *testing.T) {
    mockService := new(MockBookService)
    controller := NewBookController(mockService)

    id := primitive.NewObjectID().Hex()
    updatedBook := &models.Book{
        ID:     primitive.NewObjectID(),
        Title:  "Updated Title",
        Author: "Updated Author",
        ISBN:   "0987654321",
    }

    body, _ := json.Marshal(updatedBook)
    req := httptest.NewRequest("PUT", "/books/"+id, bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req = mux.SetURLVars(req, map[string]string{"id": id})
    resp := httptest.NewRecorder()

    mockService.On("UpdateBook", mock.MatchedBy(func(b *models.Book) bool {
        return b.Title == updatedBook.Title && b.Author == updatedBook.Author && b.ISBN == updatedBook.ISBN
    })).Return(updatedBook, nil)

    controller.UpdateBook(resp, req)

    assert.Equal(t, http.StatusOK, resp.Code)
    mockService.AssertExpectations(t)
}

func TestCreateBook_InvalidData(t *testing.T) {
    mockService := new(MockBookService)
    controller := NewBookController(mockService)

    req := httptest.NewRequest("POST", "/books", bytes.NewReader([]byte("invalid json")))
    req.Header.Set("Content-Type", "application/json")
    resp := httptest.NewRecorder()

    controller.CreateBook(resp, req)

    assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestUpdateBook_InvalidID(t *testing.T) {
    mockService := new(MockBookService)
    controller := NewBookController(mockService)

    req := httptest.NewRequest("PUT", "/books/invalid-id", bytes.NewReader([]byte("{}")))
    req.Header.Set("Content-Type", "application/json")
    req = mux.SetURLVars(req, map[string]string{"id": "invalid-id"})
    resp := httptest.NewRecorder()

    controller.UpdateBook(resp, req)

    assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestUpdateBook_InvalidData(t *testing.T) {
    mockService := new(MockBookService)
    controller := NewBookController(mockService)

    id := primitive.NewObjectID().Hex()

    req := httptest.NewRequest("PUT", "/books/"+id, bytes.NewReader([]byte("invalid json")))
    req.Header.Set("Content-Type", "application/json")
    req = mux.SetURLVars(req, map[string]string{"id": id})
    resp := httptest.NewRecorder()

    controller.UpdateBook(resp, req)

    assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestDeleteBook_Error(t *testing.T) {
    mockService := new(MockBookService)
    controller := NewBookController(mockService)

    id := primitive.NewObjectID().Hex()

    mockService.On("DeleteBookByID", id).Return(assert.AnError)

    req := httptest.NewRequest("DELETE", "/books/"+id, nil)
    req = mux.SetURLVars(req, map[string]string{"id": id})
    resp := httptest.NewRecorder()

    controller.DeleteBook(resp, req)

    assert.Equal(t, http.StatusInternalServerError, resp.Code)
    mockService.AssertExpectations(t)
}

func TestGetBookByID_MissingID(t *testing.T) {
    mockService := new(MockBookService)
    controller := NewBookController(mockService)

    req := httptest.NewRequest("GET", "/books/", nil)
    req = mux.SetURLVars(req, map[string]string{"id": ""})
    resp := httptest.NewRecorder()

    controller.GetBookByID(resp, req)

    assert.Equal(t, http.StatusBadRequest, resp.Code)
    assert.Contains(t, resp.Body.String(), "ID no proporcionado")
}

func TestGetBookByID_NotFound(t *testing.T) {
    mockService := new(MockBookService)
    controller := NewBookController(mockService)

    id := "507f191e810c19729de860ea"

    mockService.On("GetBookByID", id).Return((*models.Book)(nil), fmt.Errorf("no encontrado"))

    req := httptest.NewRequest("GET", "/books/"+id, nil)
    req = mux.SetURLVars(req, map[string]string{"id": id})
    resp := httptest.NewRecorder()

    controller.GetBookByID(resp, req)

    assert.Equal(t, http.StatusNotFound, resp.Code)
    assert.Contains(t, resp.Body.String(), "Libro no encontrado")
}