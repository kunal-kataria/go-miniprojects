package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"github.com/kunal/go-bookstore/pkg/utils"
	"github.com/kunal/go-bookstore/pkg/models"
)

var NewBook models.Book

func GetBook( w http.ResponseWriter, r *http.Request) {
	newBooks := models.GetAllBooks()
	res, _ := json.Marshal(newBooks)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookId := params["bookId"]
	Id,err := strconv.ParseInt(bookId,0,0)
	if err != nil {
		fmt.Printf("Error while parsing")
	}
	book,_ := models.GetBookById(Id);
	res,_ := json.Marshal(book)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	book := &models.Book{}
	utils.ParseBody(r, book)
	b := book.CreateBook()
	res,_ := json.Marshal(&b)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookId := params["bookId"]
	Id,err := strconv.ParseInt(bookId,0,0)
	if err != nil {
		fmt.Printf("Error while parsing")
	}
	book := models.DeleteBook(Id)
	res,_ := json.Marshal(book)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	updatebook := &models.Book{}
	utils.ParseBody(r, updatebook)
	params := mux.Vars(r)
	bookId := params["bookId"]
	Id,err := strconv.ParseInt(bookId,0,0)
	if err != nil {
		fmt.Printf("Error while parsing")
	}
	bookDet, db := models.GetBookById(Id)
	if updatebook.Name != "" {
		bookDet.Name = updatebook.Name
	}
	if updatebook.Author != "" {
		bookDet.Author = updatebook.Author
	}
	if updatebook.Publication != "" {
		bookDet.Publication = updatebook.Publication
	}
	db.Save(bookDet)
	res,_ := json.Marshal(bookDet)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}