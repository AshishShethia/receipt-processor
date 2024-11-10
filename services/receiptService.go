package services

import (
    "receipt-processor/models"
    "math"
    "strings"
    "time"
)

func ProcessReceipt(receipt models.Receipt) string {
    // Generate an ID for the receipt and store it (in-memory)
    // In a real application, this would save to a database or similar.
    return "some-generated-id" // Replace with real logic as needed
}

func CalculatePoints(receipt models.Receipt) int {
    points := 0

    // Rule 1: One point for every alphanumeric character in the retailer name
    points += len(removeNonAlphanumeric(receipt.Retailer))

    // Rule 2: 50 points if the total is a round dollar amount
    if receipt.Total == math.Floor(receipt.Total) {
        points += 50
    }

    // Rule 3: 25 points if the total is a multiple of 0.25
    if math.Mod(receipt.Total, 0.25) == 0 {
        points += 25
    }

    // Rule 4: 5 points for every two items
    points += 5 * (len(receipt.Items) / 2)

    // Rule 5: Item description and price points
    for _, item := range receipt.Items {
        trimmedDescription := strings.TrimSpace(item.ShortDescription)
        if len(trimmedDescription)%3 == 0 {
            points += int(math.Ceil(item.Price * 0.2))
        }
    }

    // Rule 6: 6 points if the day in the purchase date is odd
    date, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
    if date.Day()%2 != 0 {
        points += 6
    }

    // Rule 7: 10 points if purchase time is between 2:00 PM and 4:00 PM
    purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
    if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
        points += 10
    }

    return points
}

func removeNonAlphanumeric(s string) string {
    return strings.Map(func(r rune) rune {
        if ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z') || ('0' <= r && r <= '9') {
            return r
        }
        return -1
    }, s)
}
