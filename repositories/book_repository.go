package repositories

import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "example.com/go-mongo-app/models"
    "log"
    "os"  

)


type BookRepository struct {
    collection *mongo.Collection
}

func NewBookRepository() *BookRepository {
    mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        log.Fatal("MONGO_URI not set in environment")
    }

    clientOptions := options.Client().ApplyURI(mongoURI)
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    collection := client.Database("library").Collection("books")
    return &BookRepository{collection}
}



// repositories/book_repository.go
func (r *BookRepository) GetAllBooks() ([]models.Book, error) {
    var books []models.Book
    cursor, err := r.collection.Find(context.TODO(), bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var book models.Book
        if err := cursor.Decode(&book); err != nil {
            return nil, err
        }
        books = append(books, book)
    }

    return books, nil
}

func (r *BookRepository) CreateBook(book models.Book) (*models.Book, error) {
    book.ID = primitive.NewObjectID()
    _, err := r.collection.InsertOne(context.TODO(), book)
    if err != nil {
        return nil, err
    }
    return &book, nil
}

func (r *BookRepository) GetBookByID(id string) (*models.Book, error) {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    var book models.Book
    filter := bson.M{"_id": objectID}
    err = r.collection.FindOne(context.TODO(), filter).Decode(&book)
    if err != nil {
        return nil, err
    }

    return &book, nil
}

func (r *BookRepository) UpdateBook(book *models.Book) (*models.Book, error) {
    filter := bson.M{"_id": book.ID}
    update := bson.M{"$set": book}

    _, err := r.collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        return nil, err
    }

    return book, nil
}

func (r *BookRepository) RemoveBookByID(id string) (bool, error) {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return false, err
    }

    result, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
    if err != nil {
        return false, err
    }

    return result.DeletedCount > 0, nil
}

func (r *BookRepository) GetBooksByAuthor(author string) ([]models.Book, error) {
    var books []models.Book
    filter := bson.M{"author": author}
    cursor, err := r.collection.Find(context.TODO(), filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var book models.Book
        if err := cursor.Decode(&book); err != nil {
            return nil, err
        }
        books = append(books, book)
    }

    return books, nil
}