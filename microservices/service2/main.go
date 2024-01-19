
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Data from Service 2")
	})

	log.Fatal(http.ListenAndServe(":8082", nil))
}
