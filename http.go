package main

import (
	"encoding/json"
	"net/http"
	"html/template" 
	"log"

	"github.com/gorilla/mux"
)

// Переменная для хранения скомпилированного шаблона
var tmpl = template.Must(template.ParseFiles("index.html")) 
func SetupRoutes(cache *Cache) *mux.Router {
	r := mux.NewRouter()

		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		
				if err := tmpl.Execute(w, nil); err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}) 

		r.HandleFunc("/order/{uid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["uid"]

		order, ok := cache.Get(uid)
		if !ok {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
	}).Methods("GET")
	
	return r
}