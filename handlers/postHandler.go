package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"receipt-processor/models"
	"github.com/google/uuid"
)

// Global variables to store the latest receipt and its ID
var LatestReceipt models.Receipt
var LatestReceiptID string

// PostReceipt handles POST requests to process a receipt
func PostReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt

	// Parse the receipt data from the request body
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "Invalid receipt data", http.StatusBadRequest)
		return
	}

	// Generate a unique ID for this receipt
	receiptID := uuid.New().String()
	LatestReceiptID = receiptID
	LatestReceipt = receipt

	// Create a response with the generated ID
	response := models.ReceiptResponse{
		ID: receiptID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	fmt.Printf("Receipt processed with ID: %s\n", receiptID)
}
