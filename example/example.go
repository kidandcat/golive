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
	time.Sleep(5 * time.Second)
	for i := 0; i < 99999; i++ {
		fmt.Println("Hello World")
	}
	go hellower()
	time.Sleep(1 * time.Second)
	fmt.Println("Bye World")
}
