package main

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

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

// receipt data
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
	isreceiptvalid, validationerror := validateReceipt(newReceipt)
	if isreceiptvalid {
		receipts[uuidWithHyphen.String()] = newReceipt
	} else {

		c.IndentedJSON(http.StatusCreated, gin.H{"message": "Validaiton Error: " + validationerror})
		return
	}

	idResponse.Id = uuidWithHyphen.String()
	c.IndentedJSON(http.StatusCreated, idResponse)
}

// validate the request
func validateReceipt(pReceipt receipt) (bool, string) {
	var itemnumber int64 = 0
	var itemnumber_str string = ""
	if strings.TrimSpace(pReceipt.PurchaseDate) == "" {
		return false, "No Purchase Date Provided"
	} else if strings.TrimSpace(pReceipt.PurchaseTime) == "" {
		return false, "No Purchase Time Provided"
	} else if strings.TrimSpace(pReceipt.Retailer) == "" {
		return false, "No Retailer Provided"
	} else if strings.TrimSpace(pReceipt.Total) == "" {
		return false, "No Total Provided"
	}
	for _, item := range pReceipt.Items {
		itemnumber++
		itemnumber_str = "Item Number " + strconv.FormatInt(int64(itemnumber), 10)
		if strings.TrimSpace(item.Price) == "" {
			return false, itemnumber_str + " - No Price Provided"

		} else if strings.TrimSpace(item.ShortDescription) == "" {
			return false, itemnumber_str + " - No Short Description Provided"
		}
		_, err := strconv.ParseFloat(item.Price, 32)
		if err != nil {
			return false, itemnumber_str + " - Invalid Price Amount: " + item.Price
		}
	}
	datevalidate := regexp.MustCompile("^[0-9]{4}-(1[0-2]|0[1-9])-(3[01]|[12][0-9]|0[1-9])$")
	if !datevalidate.MatchString(pReceipt.PurchaseDate) {
		return false, "Invalid Purchase Date - Not In Expected Format YYYY-MM-DD: " + pReceipt.PurchaseDate
	}
	_, err := time.Parse("2006-01-02", pReceipt.PurchaseDate)
	if err != nil {
		return false, "Invalid Purchase Date: " + pReceipt.PurchaseDate
	}
	timevalidate := regexp.MustCompile("^(((0|1)[0-9])|(2[0-3])):((0|1|2|3|4|5)[0-9])$")
	if !timevalidate.MatchString(pReceipt.PurchaseTime) {
		return false, "Invalid Purchase Time - Not In Expected 24-hour Format HH:MI : " + pReceipt.PurchaseTime
	}
	_, err = strconv.ParseFloat(pReceipt.Total, 32)

	if err != nil {
		return false, "Invalid Total"

	}
	return true, ""
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
