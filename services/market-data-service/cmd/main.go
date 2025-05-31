package main

import (
	"fmt"
	"net/http"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func main() {
	http.HandleFunc("/healthz", healthz)
	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}