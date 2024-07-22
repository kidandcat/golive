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
	var count uint
	time.Sleep(5 * time.Second)
	fmt.Println("Hello World")
	for i := 0; i < 9999999; i++ {
		count += uint(i)
	}
	go hellower()
	time.Sleep(1 * time.Second)
	fmt.Println("Bye World", count)
}
