package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"receipt-processor/models"
	"strconv"
	"strings"
	"time"
)

func GetPoints(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	var receiptID string
	if len(parts) >= 4 {
		receiptID = parts[3]
	}

	if receiptID == "" {
		if LatestReceiptID == "" {
			http.Error(w, "No receipt processed yet", http.StatusNotFound)
			return
		}
		receiptID = LatestReceiptID
	}

	if receiptID != LatestReceiptID {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	var points int
	breakdown := []string{}

	alphanumericChars := 0
	for _, ch := range LatestReceipt.Retailer {
		if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
			alphanumericChars++
		}
	}
	points += alphanumericChars
	breakdown = append(breakdown, fmt.Sprintf("%d points - retailer name (%s) has %d alphanumeric characters", alphanumericChars, LatestReceipt.Retailer, alphanumericChars))

	total, err := strconv.ParseFloat(LatestReceipt.Total, 64)
	if err != nil {
		http.Error(w, "Invalid total value", http.StatusBadRequest)
		return
	}
	if total == float64(int(total)) {
		points += 50
		breakdown = append(breakdown, "50 points - total is a round dollar amount")
	}

	if int(total*4)%4 == 0 {
		points += 25
		breakdown = append(breakdown, "25 points - total is a multiple of 0.25")
	}

	itemsPoints := len(LatestReceipt.Items) / 2 * 5
	points += itemsPoints
	breakdown = append(breakdown, fmt.Sprintf("%d points - %d items (%d pairs @ 5 points each)", itemsPoints, len(LatestReceipt.Items), len(LatestReceipt.Items)/2))

	for _, item := range LatestReceipt.Items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLen%3 == 0 {
			itemPrice, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				http.Error(w, "Invalid item price", http.StatusBadRequest)
				return
			}
			itemPoints := int(itemPrice * 0.2)
			if itemPoints == 0 {
				itemPoints = 1
			}
			if float64(itemPoints) < itemPrice*0.2 {
				itemPoints++
			}
			points += itemPoints
			breakdown = append(breakdown, fmt.Sprintf("Item \"%s\" has trimmed length %d (a multiple of 3), item price of %s * 0.2 = %.2f, rounded up is %d points", item.ShortDescription, trimmedLen, item.Price, itemPrice*0.2, itemPoints))
		}
	}

	purchaseDate, err := time.Parse("2006-01-02", LatestReceipt.PurchaseDate)
	if err != nil {
		http.Error(w, "Invalid purchase date", http.StatusBadRequest)
		return
	}
	if purchaseDate.Day()%2 != 0 {
		points += 6
		breakdown = append(breakdown, "6 points - purchase date has an odd day")
	}

	purchaseTime, err := parseTime(LatestReceipt.PurchaseTime)
	if err != nil {
		http.Error(w, "Invalid purchase time", http.StatusBadRequest)
		return
	}
	if purchaseTime.After(time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)) && purchaseTime.Before(time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)) {
		points += 10
		breakdown = append(breakdown, "10 points - purchase time is between 2:00 PM and 4:00 PM")
	}

	response := models.PointsResponse{
		Points: points,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func parseTime(timeStr string) (time.Time, error) {
	if t, err := time.Parse("03:04pm", timeStr); err == nil {
		return t, nil
	}
	if t, err := time.Parse("15:04", timeStr); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("invalid time format")
}
