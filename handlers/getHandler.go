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

// GetPoints handles GET requests to return points for a specific receipt
func GetPoints(w http.ResponseWriter, r *http.Request) {
	// Extract receiptID from the URL path (if provided)
	parts := strings.Split(r.URL.Path, "/")
	var receiptID string
	if len(parts) >= 4 {
		receiptID = parts[3] // Get the receiptID from URL
	}

	// If no receiptID is passed in the URL, use the latest processed one
	if receiptID == "" {
		if LatestReceiptID == "" {
			http.Error(w, "No receipt processed yet", http.StatusNotFound)
			return
		}
		receiptID = LatestReceiptID
	}

	// Check if the requested receiptID matches the latest one
	if receiptID != LatestReceiptID {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Now, calculate points based on the receipt data
	var points int
	breakdown := []string{}

	// Rule 1: Points for retailer name (alphanumeric characters)
	alphanumericChars := 0
	for _, ch := range LatestReceipt.Retailer {
		if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
			alphanumericChars++
		}
	}
	points += alphanumericChars
	breakdown = append(breakdown, fmt.Sprintf("%d points - retailer name (%s) has %d alphanumeric characters", alphanumericChars, LatestReceipt.Retailer, alphanumericChars))
	// Log the points for this rule
	fmt.Printf("Total points after Rule 1: %d\n", points)

	// Rule 2: Points if total is a round dollar amount (no cents)
	total, err := strconv.ParseFloat(LatestReceipt.Total, 64)
	if err != nil {
		http.Error(w, "Invalid total value", http.StatusBadRequest)
		return
	}
	if total == float64(int(total)) { // Total is a round dollar amount if its fractional part is 0
		points += 50
		breakdown = append(breakdown, "50 points - total is a round dollar amount")
	}
	// Log points for this rule
	fmt.Printf("Total points after Rule 2: %d\n", points)

	// Rule 3: Points if total is a multiple of 0.25
	if int(total*4)%4 == 0 { // Multiply by 4 and check if it's divisible by 4
		points += 25
		breakdown = append(breakdown, "25 points - total is a multiple of 0.25")
	}
	// Log points for this rule
	fmt.Printf("Total points after Rule 3: %d\n", points)

	// Rule 4: Points for number of items (5 points for every 2 items)
	itemsPoints := len(LatestReceipt.Items) / 2 * 5 // 5 points for every 2 items
	points += itemsPoints
	breakdown = append(breakdown, fmt.Sprintf("%d points - %d items (%d pairs @ 5 points each)", itemsPoints, len(LatestReceipt.Items), len(LatestReceipt.Items)/2))
	// Log points for this rule
	fmt.Printf("Total points after Rule 4: %d\n", points)

	// Rule 5: Points for trimmed length of item description being a multiple of 3
	for _, item := range LatestReceipt.Items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLen%3 == 0 {
			// Convert Price to float64 for calculations
			itemPrice, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				http.Error(w, "Invalid item price", http.StatusBadRequest)
				return
			}
			// Calculate points based on item price
			itemPoints := int(itemPrice * 0.2)
			if itemPoints == 0 {
				itemPoints = 1
			}
			// Round up item points correctly
			if float64(itemPoints) < itemPrice*0.2 {
				itemPoints++
			}

			points += itemPoints
			breakdown = append(breakdown, fmt.Sprintf("Item \"%s\" has trimmed length %d (a multiple of 3), item price of %s * 0.2 = %.2f, rounded up is %d points", item.ShortDescription, trimmedLen, item.Price, itemPrice*0.2, itemPoints))
		}
	}
	// Log points for this rule
	fmt.Printf("Total points after Rule 5: %d\n", points)

	// Rule 6: Points for the day of the purchase being odd
	purchaseDate, err := time.Parse("2006-01-02", LatestReceipt.PurchaseDate)
	if err != nil {
		http.Error(w, "Invalid purchase date", http.StatusBadRequest)
		return
	}
	// Log the purchase date to check it
	fmt.Println("Purchase Date:", purchaseDate)
	if purchaseDate.Day()%2 != 0 { // Day is odd
		points += 6
		breakdown = append(breakdown, "6 points - purchase date has an odd day")
	}
	// Log points for this rule
	fmt.Printf("Total points after Rule 6: %d\n", points)

	// Rule 7: Points if the time of purchase is after 2:00pm and before 4:00pm (flexible time format)
	purchaseTime, err := parseTime(LatestReceipt.PurchaseTime)
	if err != nil {
		http.Error(w, "Invalid purchase time", http.StatusBadRequest)
		return
	}
	if purchaseTime.After(time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)) && purchaseTime.Before(time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)) {
		points += 10
		breakdown = append(breakdown, "10 points - purchase time is between 2:00 PM and 4:00 PM")
	}
	// Log points for this rule
	fmt.Printf("Total points after Rule 7: %d\n", points)

	// Return points and the breakdown
	response := models.PointsResponse{
		Points: points,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	// Log the breakdown
	fmt.Println("Breakdown:")
	for _, line := range breakdown {
		fmt.Println(line)
	}
}

// parseTime tries to parse time in multiple formats (12-hour and 24-hour)
func parseTime(timeStr string) (time.Time, error) {
	// Try 12-hour format (3:04pm)
	if t, err := time.Parse("03:04pm", timeStr); err == nil {
		return t, nil
	}
	// Try 24-hour format (14:30)
	if t, err := time.Parse("15:04", timeStr); err == nil {
		return t, nil
	}
	// If neither format works, return an error
	return time.Time{}, fmt.Errorf("invalid time format")
}
