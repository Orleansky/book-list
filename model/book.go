package model

type Book struct {
	Id     string `json:"id" bson:"id"`
	Title  string `json:"title" bson:"title"`
	Author string `json:"author" bson:"author"`
}
