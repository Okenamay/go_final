package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	fmt.Println("Starting server at port 7540")

	r.Handle("/*", http.FileServer(http.Dir("./web")))

	err := http.ListenAndServe(":7540", r)
	if err != nil {

		panic(err)
	}
}
