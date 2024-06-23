package main

import (
	"fmt"
	_ "golive"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World example")
	})

	http.ListenAndServe(":8080", nil)
}
