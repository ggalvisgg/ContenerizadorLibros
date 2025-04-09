package controllers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "example.com/go-mongo-app/models"
    "example.com/go-mongo-app/services"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type BookController struct {
    Service services.BookServiceInterface
}

func NewBookController(service services.BookServiceInterface) *BookController {
    return &BookController{Service: service}
}

func (c *BookController) GetBooks(w http.ResponseWriter, r *http.Request) {
    books, err := c.Service.GetBooks()
    if err != nil {
        http.Error(w, "Error al obtener libros", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Libros obtenidos correctamente",
        "books":   books,
    })
}

func (c *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
    var book models.Book
    err := json.NewDecoder(r.Body).Decode(&book)
    if err != nil {
        http.Error(w, "Datos inválidos", http.StatusBadRequest)
        return
    }

    newBook, err := c.Service.AddBook(book)
    if err != nil {
        http.Error(w, "Error al insertar libro", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newBook)
}

func (c *BookController) GetBookByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    if id == "" {
        http.Error(w, "ID no proporcionado", http.StatusBadRequest)
        return
    }

    book, err := c.Service.GetBookByID(id)
    if err != nil {
        http.Error(w, "Libro no encontrado", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Libro encontrado",
        "book":    book,
    })
}

func (c *BookController) UpdateBook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var book models.Book
    err := json.NewDecoder(r.Body).Decode(&book)
    if err != nil {
        http.Error(w, "Datos inválidos", http.StatusBadRequest)
        return
    }

    book.ID, err = primitive.ObjectIDFromHex(id)
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    updatedBook, err := c.Service.UpdateBook(&book)
    if err != nil {
        http.Error(w, "Error al actualizar libro", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Libro actualizado correctamente",
        "book":    updatedBook,
    })
}

func (c *BookController) DeleteBook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    if id == "" {
        http.Error(w, "ID no proporcionado", http.StatusBadRequest)
        return
    }

    err := c.Service.DeleteBookByID(id)
    if err != nil {
        http.Error(w, "Error al eliminar libro", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}


func (c *BookController) DeleteBookByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    if id == "" {
        http.Error(w, "ID no proporcionado", http.StatusBadRequest)
        return
    }

    err := c.Service.DeleteBookByID(id)
    if err != nil {
        http.Error(w, "Error al eliminar libro", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
