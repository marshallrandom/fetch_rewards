package main

import (
	"math"
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

// point respons
type point_response struct {
	Points int `json:"points"`
}

// receipt data
var receipts map[string]receipt

func main() {
	receipts = make(map[string]receipt)
	//test data
	receipts["someidhere"] = receipt{
		Retailer: "test", PurchaseDate: "2022-01-01", PurchaseTime: "13:01", Total: "10", Items: []receipt_items{{ShortDescription: "someitem", Price: "10"}},
	}
	router := gin.Default()
	router.GET("/receipts/:id/points", getReceiptByID)
	router.POST("/receipts/process", postReceipts)

	router.Run("localhost:8080")
}

// getReceiptByID locates the receipt whose ID value matches the id
// parameter sent by the client, then returns that receipt as a response.
func getReceiptByID(c *gin.Context) {
	id := c.Param("id")
	var returnpoints point_response

	val, ok := receipts[id]

	// If the key exists
	if ok {
		returnpoints.Points, _ = strconv.Atoi(val.points)
		c.IndentedJSON(http.StatusOK, returnpoints)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "receipt not found"})
}

// postReceipt adds an receipt from JSON received in the request body.
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
		//computes the points and stores it with the receipt
		newReceipt.points = strconv.FormatInt(int64(computePoints(newReceipt)), 10)
		receipts[uuidWithHyphen.String()] = newReceipt
	} else {
		//if a validation error occurs let them know what is invalid
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "Validaiton Error: " + validationerror})
		return
	}

	idResponse.Id = uuidWithHyphen.String()
	c.IndentedJSON(http.StatusCreated, idResponse)
}

// validate submission, checks for invalid data
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

// computes the points of a receipt
func computePoints(pReceipt receipt) int {
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	var returntotal int = 0
	returntotal = 0
	//1 point per alpha numeric character
	returntotal += len(nonAlphanumericRegex.ReplaceAllString(pReceipt.Retailer, ""))
	totalVal, err := strconv.ParseFloat(pReceipt.Total, 64)
	if err == nil {
		//50 points if total is an even dollar amount
		if totalVal == math.Round(totalVal) {
			returntotal += 50
		}
		//25 points if total is a multiple of 0.25
		if totalVal/0.25 == math.Round(totalVal/0.25) {
			returntotal += 25
		}
		//5 points per pair of items
		NumberOfPairItems := math.Floor(float64(len(pReceipt.Items)) / 2)
		returntotal += 5 * int(NumberOfPairItems)

	}
	for _, item := range pReceipt.Items {
		totalVal = 0
		totalVal, _ = strconv.ParseFloat(item.Price, 64)

		//price*0.2 rounded up if price description length is a multiple of 3
		if math.Mod(float64(len(strings.TrimSpace(item.ShortDescription))), 3) == 0 {

			returntotal += int(math.Ceil(totalVal * 0.2))

		}
	}
	//6 points if the purchase date is an odd number
	purchase_day_unit_digit := pReceipt.PurchaseDate[9:10]
	if purchase_day_unit_digit == "1" ||
		purchase_day_unit_digit == "3" ||
		purchase_day_unit_digit == "5" ||
		purchase_day_unit_digit == "7" ||
		purchase_day_unit_digit == "9" {
		returntotal += 6
	}
	//10 points if purchase time is after 2pm and before 4pm
	if (pReceipt.PurchaseTime[0:2] == "14" || pReceipt.PurchaseTime[0:2] == "15") &&
		(pReceipt.PurchaseTime != "14:00") {
		returntotal += 10
	}

	return returntotal
}
