package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

type Books struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func books(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite", "baseBooks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT id, name, age FROM books")
		if err != nil {
			http.Error(w, `{"status":"ersror"}`, http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var bookse []Books
		for rows.Next() {
			var b Books
			err := rows.Scan(&b.Id, &b.Name, &b.Age)
			if err != nil {
				http.Error(w, `{"status":"scan_error"}`, http.StatusBadRequest)
				return
			}
			bookse = append(bookse, b)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bookse)
	case http.MethodPost:
		var newBooks Books
		if err := json.NewDecoder(r.Body).Decode(&newBooks); err != nil {
			http.Error(w, `231`, http.StatusBadRequest)
		}
		_, err := db.Exec("INSERT INTO books (name, age) VALUES (?,?)", newBooks.Name, newBooks.Age)
		if err != nil {
			log.Print(err)
			http.Error(w, "dsa", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(`{"status":"kaki"}`)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}
func main() {
	db, err := sql.Open("sqlite", "baseBooks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	butes, er := os.ReadFile("create_tables.sql")
	if er != nil {
		log.Fatal("Ошибка чтения файла:", er)
	}
	base := string(butes)
	_, err = db.Exec(base)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("таблица созданна")
	http.HandleFunc("/books", books)
	http.ListenAndServe(":8080", nil)
}
