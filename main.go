package main

// package routes

import (
	"encoding/json"
	"fmt"
	"gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Book struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Books []Book

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func BookIndex(w http.ResponseWriter, r *http.Request) {
	books := Books{
		Book{ID: 1, Name: "Cat in the hat"},
		Book{ID: 2, Name: "The Hobbit"},
	}

	if err := json.NewEncoder(w).Encode(books); err != nil {
		panic(err)
	}
}

func BookShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	json.NewEncoder(w).Encode(vars)
	bookID := vars["bookID"]
	fmt.Fprintln(w, "Book:", bookID)
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"BookIndex",
		"GET",
		"/books",
		BookIndex,
	},
	Route{
		"BookShow",
		"GET",
		"/books/{bookID}",
		BookShow,
	},
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	log.Fatal(http.ListenAndServe(":8000", router))
}
