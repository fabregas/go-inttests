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

	http.HandleFunc("/volume/size", func(w http.ResponseWriter, r *http.Request) {
		info, err := os.Stat("/container_volume")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		} else {
			fmt.Fprintf(w, fmt.Sprintf("%d", info.Size()))
		}
	})

	panic(http.ListenAndServe(":5555", nil))
}
