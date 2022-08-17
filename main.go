package main

import (
	"fmt"
	"os"
)

var (
	zone   = os.Getenv("TEST_ZONE_NAME")
	apiKey = os.Getenv("API_KEY")
)

func main() {
	fmt.Println("TEST_ZONE_NAME: '" + zone + "'")
	fmt.Println("API_KEY: '" + apiKey + "'")

	domain := "egeback.se"
	// label := "_acme-challenge.test.egeback.se."
	// label := "test.egeback.se."
	baseURL := "https://dyn.anx.se/api/dns/"
	// value := "-JU6HZrE5ajDeOwtIsS60nkRpMa2Kl2zrAQyKFo0kug"
	value := "test"
	label := "test.egeback.se."

	var client = anxdns.newClient{
		Domain: domain,
		ApiKey: apiKey,
	}

	all_records := client.getAllRecords()
	fmt.Println(len(all_records))

	allTxtRecords := client.getRecordsByTxt(value, label)
	fmt.Println(allTxtRecords[0])
}
