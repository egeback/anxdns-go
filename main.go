package main

import (
	"anxdns-go/anxdns"
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

	//domain := "egeback.se"
	// label := "_acme-challenge.test.egeback.se."
	// label := "test.egeback.se."
	baseURL := "https://dyn.anx.se/api/dns/"
	// value := "-JU6HZrE5ajDeOwtIsS60nkRpMa2Kl2zrAQyKFo0kug"
	value := zone
	label := "test.egeback.se."

	var client = anxdns.Client{
		BaseUrl: baseURL,
		Domain:  zone,
		ApiKey:  apiKey,
	}

	all_records := client.GetAllRecords()
	fmt.Println(len(all_records))

	allTxtRecords := client.GetRecordsByTxt(value, label)
	fmt.Println(allTxtRecords[0])
}
