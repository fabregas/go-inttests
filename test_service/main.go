package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("starting service...")
	os.Environ()
	val, ok := os.LookupEnv("MYENV")
	if !ok {
		fmt.Println("no MYENV ;(")
		os.Exit(1)
	}
	if val != "ok" {
		fmt.Println("invalid MYENV value")
		os.Exit(2)
	}

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world")
	})

	panic(http.ListenAndServe(":5555", nil))
}
