package main

import (
	"encoding/json"
	"net/http"
)

var data = make(map[string]string)

func main() {
	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			json.NewEncoder(w).Encode(data)
		case http.MethodPost:
			var Data = make(map[string]string)
			if err := json.NewDecoder(r.Body).Decode(&Data); err != nil {
				http.Error(w, "ошибка", http.StatusBadRequest)
			}
			for kye, result := range Data {
				data[kye] = result
			}
			w.Write([]byte(`{"status" : "данные обновленны"`))
		case http.MethodDelete:
			city := r.URL.Query().Get("city")
			if _, err := data[city]; !err {
				http.Error(w, `{"status": "нет такого"}`, http.StatusBadRequest)
				return
			}
			delete(data, city)
			w.Write([]byte(`{"status" : "данные удалены"}`))
		}
	})
	http.ListenAndServe(":8080", nil)
}
