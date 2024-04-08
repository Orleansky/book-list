package repository

import (
	"Anastasia/books-list/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepo struct {
	MongoCollection *mongo.Collection
}

func (r *BookRepo) CreateBook(b *model.Book) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), b)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *BookRepo) GetBookByID(bookId string) (*model.Book, error) {
	var book model.Book

	err := r.MongoCollection.FindOne(context.Background(),
		bson.D{{Key: "id", Value: bookId}}).Decode(&book)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *BookRepo) GetBooks() ([]model.Book, error) {
	result, err := r.MongoCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var books []model.Book
	err = result.All(context.Background(), &books)
	if err != nil {
		return nil, fmt.Errorf("decoding error %s", err.Error())
	}

	return books, nil
}

func (r *BookRepo) UpdateBookById(bookId, bookTitle, bookAuthor string, updateBook *model.Book) (int64, error) {
	var result *mongo.UpdateResult
	var err error

	if updateBook.Title == "" {
		updateBook.Title = bookTitle
	}

	if updateBook.Author == "" {
		updateBook.Author = bookAuthor
	}

	result, err = r.MongoCollection.UpdateOne(context.Background(),
		bson.D{{Key: "id", Value: bookId}},
		bson.D{{Key: "$set", Value: updateBook}})
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func (r *BookRepo) DeleteBookById(bookId string) (int64, error) {
	result, err := r.MongoCollection.DeleteOne(context.Background(),
		bson.D{{Key: "id", Value: bookId}})
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (r *BookRepo) DeleteBooks() (int64, error) {
	result, err := r.MongoCollection.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
