package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "example.com/go-mongo-app/controllers"
    "example.com/go-mongo-app/services"
    "example.com/go-mongo-app/repositories"
)

func init() {
    godotenv.Load()
}
func main() {
  
    repo := repositories.NewBookRepository()
    service := services.NewBookService(repo)
    controller := controllers.NewBookController(service)

    router := mux.NewRouter()
    router.HandleFunc("/books", controller.GetBooks).Methods("GET")
    router.HandleFunc("/books", controller.CreateBook).Methods("POST")
    router.HandleFunc("/books/{id}", controller.GetBookByID).Methods("GET")
    router.HandleFunc("/books/{id}", controller.UpdateBook).Methods("PUT")
    router.HandleFunc("/books/{id}", controller.DeleteBookByID).Methods("DELETE")

    fmt.Println("Servidor en el puerto 8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}