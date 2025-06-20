package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GroverSoans/quantcrate-backend/services/market-data-service/internal/database"
	"github.com/GroverSoans/quantcrate-backend/services/market-data-service/internal/providers"
)


func helloHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hello, Http!")
}

//Health checkpoint
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func InitializeServer() {
	//Initialize postgres database
	db := database.InitializeDatabase()
	defer db.Close()
	fmt.Println("Connected to database successfully")
	
	go providers.InitializePolygon(db)
	
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/health", healthCheckHandler)


	fmt.Println("Server is running at http://localhost:8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Sever failed to start: ", err)
	}

	
}