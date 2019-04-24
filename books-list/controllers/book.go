package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/inviincible/rest_api/books-list/models"
)

type Controller struct{}
type Book = models.Book

var books []Book

func logFatal(err error) {
	if err != nil {
		log.Println(err)
	}
}

func (c Controller) GetBooks(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var book Book
		books = []Book{}

		rows, err := db.Query("select * from books")
		logFatal(err)

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
			logFatal(err)
			books = append(books, book)
		}
		json.NewEncoder(w).Encode(books)
	}
}
func (c *Controller) GetBook(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var book Book
		params := mux.Vars(r)

		bookID, _ := strconv.Atoi(params["id"])

		row := db.QueryRow("select * from books where id = $1", bookID)

		err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)

		json.NewEncoder(w).Encode(&book)
	}
}

func (c *Controller) AddBook(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var book Book
		//var bookID int
		err := json.NewDecoder(r.Body).Decode(&book)
		logFatal(err)

		err = db.QueryRow("insert into books (title, author, year) values($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&book.ID)
		logFatal(err)

		json.NewEncoder(w).Encode(&book)

	}
}

func (c *Controller) UpdateBook(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var book Book

		err := json.NewDecoder(r.Body).Decode(&book)
		logFatal(err)

		result, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id;", book.Title, book.Author, book.Year, book.ID)
		logFatal(err)

		rowsUpdated, err := result.RowsAffected()
		logFatal(err)

		json.NewEncoder(w).Encode(rowsUpdated)
	}

}

func (c *Controller) RemoveBook(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		result, err := db.Exec("delete from books where id=$1", params["id"])
		logFatal(err)

		rowsDeleted, err := result.RowsAffected()
		logFatal(err)

		json.NewEncoder(w).Encode(rowsDeleted)
	}
}
