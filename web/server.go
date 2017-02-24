package main

import (
	"fmt"
	"net/http"
)

func main() {

	wsserver := NewWSServer("/")
	wsserver.HandlePaths()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
