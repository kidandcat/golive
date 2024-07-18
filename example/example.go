package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/kidandcat/golive"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World example")
	})

	go hellower()

	http.ListenAndServe(":8080", nil)
}

func hellower() {
	for {
		time.Sleep(3 * time.Second)
		fmt.Println("Hello World")
	}
}
