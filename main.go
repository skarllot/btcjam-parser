package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	DATA_0            = `^([A-Za-z ]+)<br />([A-Za-z0-9 :,]+)`
	DATA_1_PAYMENT    = `payment (\d+) of (\d+) for <a [^ ]+ href="([^"]+)">listing (\d+)`
	DATA_1_INVESTMENT = `href="([^"]+)">listing (\d+)`
	DATA_1_DEPOSIT    = `txid: ([0-9a-f]+)`
	DATA_2            = `(-)? BTC ([0-9.,]+)`
	DATE_FORMAT       = `Jan 2, 2006 15:04:05`
)

func main() {
	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	t := TransactionOriginal{}
	r0 := regexp.MustCompile(DATA_0)
	r1Payment := regexp.MustCompile(DATA_1_PAYMENT)
	r1Investment := regexp.MustCompile(DATA_1_INVESTMENT)
	r1Deposit := regexp.MustCompile(DATA_1_DEPOSIT)
	r2 := regexp.MustCompile(DATA_2)
	loc := time.Now().Location()

	dec.Decode(&t)

	newTransactions := make([]TransactionData, 0)

	for _, item := range t.Data {
		tdata := TransactionData{}
		result := r0.FindStringSubmatch(item[0])
		tdata.Type = result[1]
		tdata.Date, _ = time.ParseInLocation(DATE_FORMAT, result[2], loc)

		switch tdata.Type {
		case "Loan Payment Received":
			result = r1Payment.FindStringSubmatch(item[1])
			tdata.PaymentNumber, _ = strconv.Atoi(result[1])
			tdata.PaymentTotal, _ = strconv.Atoi(result[2])
			tdata.ListingHref = result[3]
			tdata.ListingId, _ = strconv.ParseInt(result[4], 10, 64)
		case "Investment", "Return of Investment":
			result = r1Investment.FindStringSubmatch(item[1])
			tdata.ListingHref = result[1]
			tdata.ListingId, _ = strconv.ParseInt(result[2], 10, 64)
		case "Deposit", "Withdrawal":
			result = r1Deposit.FindStringSubmatch(item[1])
			tdata.BTCTransaction = result[1]
		default:
			panic(fmt.Sprintf("Unexpected transaction type: %s", tdata.Type))
		}

		result = r2.FindStringSubmatch(item[2])
		result[2] = strings.Replace(result[2], ",", "", -1)
		result[2] = strings.Replace(result[2], ".", "", -1)
		if result[1] == "-" {
			result[2] = "-" + result[2]
		}

		tdata.Value, _ = strconv.ParseInt(result[2], 10, 64)

		newTransactions = append(newTransactions, tdata)
	}

	enc.Encode(newTransactions)
}

type TransactionOriginal struct {
	Echo                string     `json:"sEcho"`
	TotalRecords        string     `json:"iTotalRecords"`
	TotalDisplayRecords string     `json:"iTotalDisplayRecords"`
	Data                [][]string `json:"aaData"`
}

type TransactionNew struct {
	TotalRecords    int               `json:"totalRecords"`
	ReturnedRecords int               `json:"returnedRecords`
	Data            []TransactionData `json:"data"`
}

type TransactionData struct {
	Type           string    `json:"type"`
	Date           time.Time `json:"date"`
	ListingId      int64     `json:"listingId,omitempty"`
	ListingHref    string    `json:"listingHref,omitempty"`
	PaymentNumber  int       `json:"paymentNumber,omitempty"`
	PaymentTotal   int       `json:"paymentTotal,omitempty"`
	BTCTransaction string    `json:"btcTransaction,omitempty"`
	Value          int64     `json:"value"`
}
