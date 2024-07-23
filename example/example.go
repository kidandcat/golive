package main

import (
	"fmt"
	"net/http"

	_ "github.com/kidandcat/golive"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World example")
	})

	http.ListenAndServe(":8080", nil)
}
