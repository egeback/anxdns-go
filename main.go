package main

import (
	"anxdns-go/anxdns"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

var (
	zone   = os.Getenv("TEST_ZONE_NAME")
	apiKey = os.Getenv("API_KEY")
)

var CLI struct {
	Get struct {
		all bool `help:"Get all records."`
	} `cmd help:"Get Records"`
}

func main() {
	fmt.Println("TEST_ZONE_NAME: '" + zone + "'")
	fmt.Println("API_KEY: '" + apiKey + "'")

	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "get":
	default:
		panic(ctx.Command())
	}

	//domain := "egeback.se"
	// label := "_acme-challenge.test.egeback.se."
	// label := "test.egeback.se."
	baseURL := "https://dyn.anx.se/api/dns/"
	// value := "-JU6HZrE5ajDeOwtIsS60nkRpMa2Kl2zrAQyKFo0kug"
	value := "test"
	label := "test.egeback.se."

	var client = anxdns.Client{
		BaseUrl: baseURL,
		Domain:  zone,
		ApiKey:  apiKey,
	}

	//all_records := client.GetAllRecords()
	//fmt.Println(len(all_records))

	allTxtRecords := client.GetRecordsByTxt(value, label)
	if len(allTxtRecords) > 0 {
		fmt.Println(allTxtRecords[0])
	} else {
		fmt.Println("No recods found")
	}

}
