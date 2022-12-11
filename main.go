package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// receipts
type receipt struct {
	retailer     string          `json:"retailer"`
	purchaseDate string          `json:"purchaseDate"`
	purchaseTime string          `json:"purchaseTime"`
	items        []receipt_items `json:"items"`
	total        string          `json:"total"`
}

// receipt items
type receipt_items struct {
	shortDescription string `json:"shortDescription"`
	price            string `json:"price"`
}

// albums slice to seed record album data.
var receipts = []receipt{
	{retailer: "test", purchaseDate: "2022-01-01", purchaseTime: "13:01", total: "10", []receipt_items{shortDescipriton: "someitem", price: "10"}}
}

func main() {
	router := gin.Default()
	router.GET("/receipts/:id/points", getReceiptByID)
	router.POST("/receipts/process", postReceipts)

	router.Run("localhost:8080")
}
// postRecept adds an receipt from JSON received in the request body.
func postReceipts(c *gin.Context) {
    var newReceipt receipt

    // Call BindJSON to bind the received JSON to
    // newReceipt.
    if err := c.BindJSON(&newReceipt); err != nil {
        return
    }

    // Add the new receipt to the slice.
    receipts = append(receipts, newReceipts)
    c.IndentedJSON(http.StatusCreated, newReceipt)
}
// getReceiptByID locates the receipt whose ID value matches the id
// parameter sent by the client, then returns that receipt as a response.
func getReceiptByID(c *gin.Context) {
	id := c.Param("id")

	// Loop through the list of receipts, looking for
	// an receipts whose ID value matches the parameter.
	for _, a := range receipts {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "receipt not found"})
}
