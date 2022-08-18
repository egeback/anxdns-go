package main

import (
	"anxdns-go/anxdns"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/davecgh/go-spew/spew"
)

var (
	zone   = os.Getenv("ANXDNS_DOMAIN")
	apiKey = os.Getenv("ANXDNS_APIKEY")
)

var cli struct {
	Get struct {
		All  bool   `short:"a" help:"Get all records"`
		Txt  string `short:"t" help:"Text value of record"`
		Name string `arg:"" optional:"" help:"Name of the records to get"`
	} `cmd help:"Get Records"`
	Add struct {
		Type string `short:"t" help:"Ticket type"`
	} `cmd help:"Add Record"`
	Update struct {
	} `cmd help:"Update Record"`
	Delete struct {
	} `cmd help:"Delete Record"`
	Apikey  bool   `short:"k" help:"API key used in request header"`
	Verbose bool   `short:"v" help:"Verbose"`
	Baseurl string `short:"b" help:"Url to API" type:"url"`
}

func main() {
	fmt.Println("ANXDNS_DOMAIN: '" + zone + "'")
	fmt.Println("ANXDNS_APIKEY: '" + apiKey + "'")

	ctx := kong.Parse(&cli)

	var baseURL string

	if cli.Baseurl != "" {
		baseURL = cli.Baseurl
	} else {
		baseURL = "https://dyn.anx.se/api/dns/"
	}

	var client = anxdns.Client{
		BaseUrl: baseURL,
		Domain:  zone,
		ApiKey:  apiKey,
	}

	switch ctx.Command() {
	case "get":
		getAllRecords(client)
	case "get <name>":
		if cli.Get.Txt != "" {
			getTxtRecord(client, cli.Get.Name, cli.Get.Txt)
		} else {
			getRecord(client, cli.Get.Name)
		}
	case "add":
		fmt.Println("Not implemented yet")
	case "update":
		fmt.Println("Not implemented yet")
	case "delete":
		fmt.Println("Not implemented yet")
	default:
		fmt.Printf("ctx: %v\n", ctx)
	}

}

func getAllRecords(client anxdns.Client) {
	all_records := client.GetAllRecords()
	fmt.Println("Number of records: " + fmt.Sprint(len(all_records)))
	spew.Dump(all_records)
}

func getRecord(client anxdns.Client, name string) {
	if len(name) == 0 {
		fmt.Println("Name of record not specified")
		os.Exit(1)
	}

	all_records := client.GetRecordsByName(name)
	fmt.Println("Number of records: " + fmt.Sprint(len(all_records)))
	spew.Dump(all_records)
}

func getTxtRecord(client anxdns.Client, name string, txt string) {
	if len(name) == 0 {
		fmt.Println("Name of record not specified")
		os.Exit(1)
	}

	if len(txt) == 0 {
		fmt.Println("Txt of record not specified")
		os.Exit(1)
	}

	all_records := client.GetRecordsByTxt(txt, name)
	fmt.Println("Number of records: " + fmt.Sprint(len(all_records)))
	spew.Dump(all_records)
}
