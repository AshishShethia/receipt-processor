package handlers

import (
	"encoding/json"
	"net/http"
	"receipt-processor/models"
	"github.com/google/uuid"
)

var LatestReceipt models.Receipt
var LatestReceiptID string

func PostReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt

	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "Invalid receipt data", http.StatusBadRequest)
		return
	}

	receiptID := uuid.New().String()
	LatestReceiptID = receiptID
	LatestReceipt = receipt

	response := models.ReceiptResponse{
		ID: receiptID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
