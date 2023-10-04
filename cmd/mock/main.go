package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
			}

			fmt.Println(string(body))

		}
	})

	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Fatal(err)
	}
}
