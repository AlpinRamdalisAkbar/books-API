package main

import (
	"encoding/json"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"database/sql"

	"github.com/gorilla/mux"
	_"github.com/go-sql-driver/mysql"
)

type Books struct {
	ID	string	`json:"id"`
	Title 	string	`json:"title"`
}

var db *sql.DB
var err error

func getBooks(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	var books []Books
	result, err := db.Query("SELECT id, title from books")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var book Books
		err := result.Scan(&book.ID, &book.Title)
		if err != nil {
			panic(err.Error())
		}
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	result, err := db.Query("SELECT id, title FROM books WHERE id = ?", param["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var book Books
	for result.Next() {
		err := result.Scan(&book.ID, &book.Title)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request)  {
	stmt, err := db.Prepare("INSERT INTO books(title) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]
	_, err = stmt.Exec(title)

	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New Books was added")
}

func updateBook(w http.ResponseWriter, r *http.Request)  {
	param := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE books SET title = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newTitle := keyVal["title"]
	_, err = stmt.Exec(newTitle, param["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "books with ID = %s was updated", param["id"])
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM books WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(param["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Book with ID = %s was deleted", param["id"])
}

func main() {
	r := mux.NewRouter()

	db, err = sql.Open("mysql", "root:@/db_books")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()


	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	fmt.Println("server start at localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
