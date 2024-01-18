package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting serving content on port 8080")
	http.Handle("/", http.FileServer(http.Dir("wasm/static")))
	http.ListenAndServe(":8080", nil)
}
