package service

import (
	"Anastasia/books-list/model"
	"Anastasia/books-list/repository"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookService struct {
	MongoCollection *mongo.Collection
}

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func (svc *BookService) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	var newBook model.Book

	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid body ", err)
		res.Error = err.Error()
		return
	}

	newBook.Id = uuid.NewString()

	repo := repository.BookRepo{MongoCollection: svc.MongoCollection}

	createId, err := repo.CreateBook(&newBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("create error ", err)
		res.Error = err.Error()
		return
	}

	res.Data = newBook.Id
	w.WriteHeader(http.StatusOK)

	log.Println("book created with id ", createId, newBook)
}

func (svc *BookService) GetBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	params := mux.Vars(r)
	log.Println("book id: ", params["id"])

	repo := repository.BookRepo{MongoCollection: svc.MongoCollection}

	book, err := repo.GetBookByID(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error: ", err)
		res.Error = err.Error()
		return
	}

	res.Data = book
	w.WriteHeader(http.StatusOK)
}

func (svc *BookService) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	repo := repository.BookRepo{MongoCollection: svc.MongoCollection}

	books, err := repo.GetBooks()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error: ", err)
		res.Error = err.Error()
		return
	}

	res.Data = books
	w.WriteHeader(http.StatusOK)
}
func (svc *BookService) UpdateBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	params := mux.Vars(r)

	if params["id"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid employee id")
		res.Error = "invalid employee id"
		return
	}

	var book model.Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid body ", err)
		res.Error = err.Error()
		return
	}

	book.Id = params["id"]

	log.Println("employee id: ", book.Id)

	repo := repository.BookRepo{MongoCollection: svc.MongoCollection}
	updateBookPointer, err := repo.GetBookByID(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error: ", err)
		res.Error = err.Error()
		return
	}
	updateBook := *updateBookPointer

	updatesCount, err := repo.UpdateBookById(book.Id, updateBook.Title, updateBook.Author, &book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error: ", err)
		res.Error = err.Error()
		return
	}

	res.Data = updatesCount
	w.WriteHeader(http.StatusOK)
}

func (svc *BookService) DeleteBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	bookId := mux.Vars(r)["id"]
	log.Println("book id: ", bookId)

	repo := repository.BookRepo{MongoCollection: svc.MongoCollection}

	deletedCount, err := repo.DeleteBookById(bookId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error: ", err)
		res.Error = err.Error()
		return
	}

	res.Data = deletedCount
	w.WriteHeader(http.StatusOK)
}
func (svc *BookService) DeleteBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	repo := repository.BookRepo{MongoCollection: svc.MongoCollection}

	deletedCount, err := repo.DeleteBooks()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error: ", err)
		res.Error = err.Error()
		return
	}

	res.Data = deletedCount
	w.WriteHeader(http.StatusOK)
}
