package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Film struct {
	Id    int      `json:"id"`
	Title string   `json:"title"`
	Uear  int      `json:"uear"`
	Ganre []string `json:"ganre"`
}

var ID = make(map[int]Film)

func main() {
	http.HandleFunc("/movies/", InfAndDel)
	http.HandleFunc("/movies", movies)
	http.ListenAndServe(":8080", nil)
}

func InfAndDel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/movies/")
	ok2, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "ID должен быть числом"}`, http.StatusBadRequest)
		return
	}
	if _, err := ID[ok2]; !err {
		http.Error(w, `{"status" : "Not found film"}`, http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		var film = make(map[int]string)
		film[ok2] = ID[ok2].Title
		json.NewEncoder(w).Encode(film)
	case http.MethodDelete:
		delete(ID, ok2)
	}
}

func movies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		if url := r.URL.Query().Get("q"); url != "" {
			for kye, _ := range ID {
				if url == ID[kye].Title {
					type Idfilm struct {
						Id    int
						Uear  int
						Ganre []string
					}
					fim := Idfilm{kye, ID[kye].Uear, ID[kye].Ganre}
					json.NewEncoder(w).Encode(fim)
					return
				}
			}
			http.Error(w, `{"status" : "Not found film"}`, http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(ID)

	case http.MethodPost:
		var film Film
		if err := json.NewDecoder(r.Body).Decode(&film); err != nil {
			http.Error(w, `{"status" : "error"}`, http.StatusBadRequest)
			fmt.Print(err)
			return
		}
		ID[film.Id] = film
		w.Write([]byte(`{"status" : "completion}"`))
	}
}
