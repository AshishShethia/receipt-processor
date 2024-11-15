package main

import (
	"fmt"
	"log"
	"net/http"
	"receipt-processor/handlers"
)

func main() {
	http.HandleFunc("/receipts/process", handlers.PostReceipt)
	http.HandleFunc("/receipts/points", handlers.GetPoints)
	http.HandleFunc("/receipts/points/", handlers.GetPoints)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
