package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	response := make(map[string]string)
	response["hello"] = r.URL.Path[1:]
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("An error occured")
	}
	w.Write(jsonResponse)
	return
}
