
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Data from Service 1")
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
