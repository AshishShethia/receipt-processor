# Receipt Processor - A Go-based API for Reward Points Calculation

## Overview

The **Receipt Processor** is a RESTful API built in **Go** that processes receipts, calculates reward points based on several rules, and allows users to retrieve the calculated points for a specific receipt or the latest processed one. The system is designed to handle receipt data with multiple items and apply a series of rules to compute points that users can redeem.

## Features

- **POST /receipts/process**: Submit a receipt for processing, generating a unique receipt ID.
- **GET /receipts/points**: Fetch points for the most recent processed receipt.
- **GET /receipts/points/{receiptID}**: Fetch points for a specific receipt using its unique ID.

## Rules for Points Calculation

The points are awarded based on the following rules:

1. **Retailer Name**: Points are awarded based on the number of alphanumeric characters in the retailer's name.
2. **Total Amount**: Points are awarded if the total is a round dollar amount (without cents).
3. **Multiple of 0.25**: Points are awarded if the total is a multiple of 0.25.
4. **Number of Items**: Points are awarded for every 2 items on the receipt.
5. **Item Description Length**: Points are awarded based on the trimmed length of item descriptions being a multiple of 3.
6. **Purchase Date**: Points are awarded if the purchase date falls on an odd day of the month.
7. **Purchase Time**: Points are awarded if the purchase time is between 2:00 PM and 4:00 PM.

## Installation

### Prerequisites

- **Go 1.18+** installed on your machine

### Getting Started

1. Clone the repository to your local machine:
```bash
    git clone https://github.com/ashishshethia/receipt-processor.git
   cd receipt-processor
```

   

2. Install Go dependencies:

```bash
Copy code
go mod tidy
```

3. Run the API server:

```bash
Copy code
go run main.go
```

### Scope

Since this project was intended as an exercise, some advanced features such as user authentication, database integration, and receipt image processing were not implemented. However, these enhancements can easily be added in future iterations of the project.