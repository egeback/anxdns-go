package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/egeback/anxdns-go/anxdns"

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
		Type string `short:"t" default:"A" enum:"A,CNAME,TXT" help:"Ticket type. Values [A, CNAME, TXT}. Default A]"`
		Ttl  string `short:"l" format:"int" default:"3600" help:"Ticket type. Values [A, CNAME, TXT}. Default A]"`
		Name string `arg:"" help:"Name of the records to get"`
		Data string `arg:"" optional:"" help:"TXT, CNAME, A data to add to record"`
	} `cmd help:"Add Record"`
	Update struct {
	} `cmd help:"Update Record"`
	Delete struct {
		Type string `short:"t" default:"A" enum:"A,CNAME,TXT" help:"Ticket type. Values [A, CNAME, TXT}. Default A]"`
		Ttl  string `short:"l" format:"int" default:"3600" help:"Ticket type. Values [A, CNAME, TXT}. Default A]"`
		Name string `arg:"" help:"Name of the records to get"`
		Data string `arg:"" optional:"" help:"TXT, CNAME, A data to add to record"`
	} `cmd help:"Delete Record"`
	Zone    string `short:"d" help:"Domain/Zone to update. Can also be provided with env arg ANXDNS_DOMAIN"`
	Apikey  string `short:"k" help:"API key used in request header. Can also be provided with env arg ANXDNS_APIKEY"`
	Verbose bool   `short:"v" help:"Verbose"`
	Baseurl string `short:"b" help:"Url to API" type:"url"`
}

func main() {
	ctx := kong.Parse(&cli)

	if cli.Verbose {
		fmt.Println("ANXDNS_DOMAIN: '" + zone + "'")
		fmt.Println("ANXDNS_APIKEY: '" + apiKey + "'")
	}

	var baseURL string

	if cli.Baseurl != "" {
		baseURL = cli.Baseurl
	} else {
		baseURL = "https://dyn.anx.se/api/dns/"
	}

	if len(cli.Zone) != 0 {
		zone = cli.Zone
	}
	if len(cli.Apikey) != 0 {
		apiKey = cli.Apikey
	}

	if len(zone) == 0 || len(apiKey) == 0 {
		fmt.Println("Zone or ApiKey not set")
		ctx.PrintUsage(false)
		os.Exit(-1)
	}

	var client = anxdns.NewClient(zone, apiKey)
	client.BaseUrl = baseURL

	switch ctx.Command() {
	case "get":
		getAllRecords(*client)
	case "get <name>":
		if cli.Get.Txt != "" {
			getTxtRecord(*client, cli.Get.Name, cli.Get.Txt)
		} else {
			getRecord(*client, cli.Get.Name)
		}
	case "add <name> <data>":
		addRecord(*client, cli.Add.Name, cli.Add.Type, cli.Add.Data, cli.Add.Ttl, ctx)
	case "update":
		fmt.Println("Not implemented yet")
	case "delete <name> <data>":
		deleteRecord(*client, cli.Delete.Name, cli.Delete.Type, cli.Delete.Data, cli.Delete.Ttl, ctx)
	default:
		ctx.PrintUsage(false)
		os.Exit(-1)
	}

}

func getAllRecords(client anxdns.Client) {
	all_records, error := client.GetAllRecords()
	if error != nil {
		fmt.Println(error)
		os.Exit(-1)
	}
	fmt.Println("Number of records: " + fmt.Sprint(len(all_records)))
	spew.Dump(all_records)
}

func getRecord(client anxdns.Client, name string) {
	if len(name) == 0 {
		fmt.Println("Name of record not specified")
		os.Exit(1)
	}

	all_records, error := client.GetRecordsByName(name)
	if error != nil {
		fmt.Println(error)
		os.Exit(-1)
	}
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

	all_records, error := client.GetRecordsByTxt(txt, name)
	if error != nil {
		fmt.Println(error)
		os.Exit(-1)
	}
	fmt.Println("Number of records: " + fmt.Sprint(len(all_records)))
	spew.Dump(all_records)
}

func addRecord(client anxdns.Client, name string, recordType string, data string, ttlString string, ctx *kong.Context) {
	// fmt.Println(name, recordType, data, ttlString)
	ttl, err := strconv.Atoi(ttlString)
	if err != nil {
		fmt.Println("Could not parse ttl value")
		ctx.PrintUsage(false)
		os.Exit(-1)
	}

	if recordType == "A" {
		client.AddARecord(name, data, ttl)
	} else if recordType == "CNAME" {
		client.AddCNameRecord(name, data, ttl)
	} else if recordType == "TXT" {
		client.AddTxtRecord(name, data, ttl)
	} else {
		fmt.Println("Wrong record type")
		ctx.PrintUsage(false)
		os.Exit(-1)
	}
	fmt.Println("Record added successfully")
}

func deleteRecord(client anxdns.Client, name string, recordType string, data string, ttlString string, ctx *kong.Context) {
	ttl, err := strconv.Atoi(ttlString)
	if err != nil {
		fmt.Println("Could not parse ttl value")
		ctx.PrintUsage(false)
		os.Exit(-1)
	}

	if recordType == "A" {
		fmt.Println(ttl)
		client.DeleteRecordsByName(data)
	} else if recordType == "CNAME" {
		client.DeleteRecordsByName(data)
	} else if recordType == "TXT" {
		client.DeleteRecordsByTxt(name, data)
	} else {
		fmt.Println("Wrong record type")
		ctx.PrintUsage(false)
		os.Exit(-1)
	}
	fmt.Println("Record removed successfully")
}
