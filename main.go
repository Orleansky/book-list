package main

import (
	"Anastasia/books-list/service"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("env load error", err)
	}

	log.Println("env file loaded")

	mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err != nil {
		log.Fatal("connection error", err)
	}

	err = mongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("ping failed", err)
	}

	log.Println("mongo connected")
}

func main() {
	defer mongoClient.Disconnect(context.Background())

	coll := mongoClient.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLLECTION_NAME"))

	newService := service.BookService{MongoCollection: coll}

	r := mux.NewRouter()

	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)

	r.HandleFunc("/books", newService.CreateBook).Methods(http.MethodPost)
	r.HandleFunc("/books/{id}", newService.GetBookByID).Methods(http.MethodGet)
	r.HandleFunc("/books", newService.GetBooks).Methods(http.MethodGet)
	r.HandleFunc("/books/{id}", newService.UpdateBookById).Methods(http.MethodPut)
	r.HandleFunc("/books/{id}", newService.DeleteBookById).Methods(http.MethodDelete)
	r.HandleFunc("/books", newService.DeleteBooks).Methods(http.MethodDelete)

	log.Println("server is running on 28017")
	http.ListenAndServe(":28017", r)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("running..."))
}
