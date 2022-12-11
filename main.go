package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// receipts
type receipt struct {
	Retailer     string          `json:"retailer"`
	PurchaseDate string          `json:"purchaseDate"`
	PurchaseTime string          `json:"purchaseTime"`
	Items        []receipt_items `json:"items"`
	Total        string          `json:"total"`
	points       string
}

// receipt items
type receipt_items struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// receipt items
type id_response struct {
	Id string `json:"id"`
}

// albums slice to seed record album data.
//
//	var receipts = []receipt{
//		{id: "someidhere", retailer: "test", purchaseDate: "2022-01-01", purchaseTime: "13:01", total: "10", items: []receipt_items{{shortDescription: "someitem", price: "10"}}},
//	}
var receipts map[string]receipt

func main() {
	receipts = make(map[string]receipt)
	receipts["someidhere"] = receipt{
		Retailer: "test", PurchaseDate: "2022-01-01", PurchaseTime: "13:01", Total: "10", Items: []receipt_items{{ShortDescription: "someitem", Price: "10"}},
	}
	router := gin.Default()
	router.GET("/receipts/:id/points", getReceiptByID)
	router.POST("/receipts/process", postReceipts)

	router.Run("localhost:8080")
}

// postRecept adds an receipt from JSON received in the request body.
func postReceipts(c *gin.Context) {
	var newReceipt receipt
	var idResponse id_response
	uuidWithHyphen := uuid.New()
	// Call BindJSON to bind the received JSON to
	// newReceipt.
	if err := c.BindJSON(&newReceipt); err != nil {
		return
	}
	receipts[uuidWithHyphen.String()] = newReceipt
	idResponse.Id = uuidWithHyphen.String()
	c.IndentedJSON(http.StatusCreated, idResponse)
}

// getReceiptByID locates the receipt whose ID value matches the id
// parameter sent by the client, then returns that receipt as a response.
func getReceiptByID(c *gin.Context) {
	id := c.Param("id")

	val, ok := receipts[id]
	// If the key exists
	if ok {
		c.IndentedJSON(http.StatusOK, val)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "receipt not found"})
}
