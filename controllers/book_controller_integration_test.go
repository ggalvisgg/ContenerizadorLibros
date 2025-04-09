package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/go-mongo-app/models"
	"example.com/go-mongo-app/repositories"
	"example.com/go-mongo-app/services"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *mux.Router {
	repo := repositories.NewBookRepository()
	service := services.NewBookService(repo)
	controller := NewBookController(service)

	r := mux.NewRouter()
	r.HandleFunc("/books", controller.GetBooks).Methods("GET")
	r.HandleFunc("/books", controller.CreateBook).Methods("POST")
	r.HandleFunc("/books/{id}", controller.GetBookByID).Methods("GET")
	r.HandleFunc("/books/{id}", controller.DeleteBookByID).Methods("DELETE")
	r.HandleFunc("/books/{id}", controller.UpdateBook).Methods("PUT")

	return r
}

func TestBookCRUDIntegration(t *testing.T) {
	router := setupRouter()

	// 1. Create a book
	book := models.Book{
		ISBN:   "123-456-789",
		Title:  "Test Book",
		Author: "Test Author",
	}
	body, _ := json.Marshal(book)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Get the ID of the created book
	bodyResp, _ := ioutil.ReadAll(rr.Body)
	var createdBook models.Book
	json.Unmarshal(bodyResp, &createdBook)
	assert.NotEmpty(t, createdBook.ID)

	// 2. Get all books
	reqGetAll, _ := http.NewRequest("GET", "/books", nil)
	rrGetAll := httptest.NewRecorder()
	router.ServeHTTP(rrGetAll, reqGetAll)
	assert.Equal(t, http.StatusOK, rrGetAll.Code)

	// 3. Get book by ID
	urlByID := fmt.Sprintf("/books/%s", createdBook.ID.Hex())
	reqGet, _ := http.NewRequest("GET", urlByID, nil)
	rrGet := httptest.NewRecorder()
	router.ServeHTTP(rrGet, reqGet)
	assert.Equal(t, http.StatusOK, rrGet.Code)

	// 4. Update book
	createdBook.Title = "Updated Test Book"
	updateBody, _ := json.Marshal(createdBook)
	reqUpdate, _ := http.NewRequest("PUT", "/books/"+createdBook.ID.Hex(), bytes.NewBuffer(updateBody))
	reqUpdate.Header.Set("Content-Type", "application/json")
	rrUpdate := httptest.NewRecorder()
	router.ServeHTTP(rrUpdate, reqUpdate)
	assert.Equal(t, http.StatusOK, rrUpdate.Code)

	// 5. Delete book
	reqDelete, _ := http.NewRequest("DELETE", urlByID, nil)
	rrDelete := httptest.NewRecorder()
	router.ServeHTTP(rrDelete, reqDelete)
	assert.Equal(t, http.StatusNoContent, rrDelete.Code)
}
