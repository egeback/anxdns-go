package anxdns

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	defaultTTL int = 3600
)

type Client struct {
	Domain  string
	ApiKey  string
	BaseUrl string `default:"https://dyn.anx.se/api/dns/"`
}

func New(Domain string, ApiKey string) *Client {
	return &Client{
		Domain: Domain,
		ApiKey: ApiKey,
	}
}

func (client *Client) _communicate(apiRequest Request) []byte {
	// Create client
	httpClient := &http.Client{}

	var request *http.Request
	var error error

	if apiRequest.JsonData == nil {
		request, error = http.NewRequest(apiRequest.Type, client.BaseUrl+apiRequest.QueryParams, nil)
	} else {
		request, error = http.NewRequest(apiRequest.Type, client.BaseUrl+apiRequest.QueryParams, bytes.NewBuffer(apiRequest.JsonData))
	}

	if error != nil {
		fmt.Println(error)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("apiKey", client.ApiKey)

	response, error := httpClient.Do(request)
	if error != nil {
		fmt.Println(error)
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			//klog.Fatal(err)
		}
	}()

	// Read response body
	respBody, _ := ioutil.ReadAll(response.Body)

	// Display results
	//fmt.Println("response Status : ", response.Status)
	//fmt.Println("response Body : ", string(respBody))

	if response.StatusCode != 200 {
		panic(errors.New("Could not communicate with server"))
	}

	return respBody
}

func (client Client) AddTxtRecord(name string, txt string, ttl int) {
	record := Data{
		Domain:  client.Domain,
		Type:    "TXT",
		Name:    name,
		TTL:     ttl,
		TxtData: txt,
	}

	jsonData, _ := json.Marshal(record)

	apiRequest := Request{
		Type:     POST,
		JsonData: jsonData,
	}

	client._communicate(apiRequest)
}

func (client Client) AddARecord(name string, address string, ttl int) {
	record := Data{
		Domain:  client.Domain,
		Type:    "A",
		Name:    name,
		TTL:     ttl,
		Address: address,
	}

	jsonData, _ := json.Marshal(record)

	apiRequest := Request{
		Type:     POST,
		JsonData: jsonData,
	}

	client._communicate(apiRequest)
}

func (client Client) AddCNameRecord(name string, address string, ttl int) {
	record := Data{
		Domain:  client.Domain,
		Type:    "CNAME",
		Name:    name,
		TTL:     ttl,
		Address: address,
	}

	jsonData, _ := json.Marshal(record)

	apiRequest := Request{
		Type:     POST,
		JsonData: jsonData,
	}

	client._communicate(apiRequest)
}

func (client Client) VerifyOrGetRecord(line int, name string, recordType string) Record {
	var record Record
	if line > 0 {
		record = client.getRecordsByLine(line)
	} else if len(name) > 0 {
		records := client.GetRecordsByName(name)
		if len(records) == 0 {
			panic(errors.New("0 records with that name"))
		} else if len(records) > 1 {
			panic(errors.New(">1 record with that name. Specify line instead of name."))
		}
		record = records[0]
	} else {
		panic(errors.New("Line or name needs to be provided"))
	}

	if len(recordType) > 0 && record.Type != recordType {
		panic(errors.New("Record is not a " + recordType))
	}

	return record

}

func (client Client) DeleteRecordsByTxt(name string) {

}

func (client *Client) GetAllRecords() []Record {
	request := Request{
		Type:        GET,
		QueryParams: "?domain=" + client.Domain,
	}
	respBody := client._communicate(request)
	response := Response{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		panic(err)
	}
	return response.DnsRecords
}

func (client Client) GetRecordsByName(name string) []Record {
	all_records := client.GetAllRecords()

	return ParseRecordsByName(all_records, name)
}

func (client Client) getRecordsByLine(line int) Record {
	all_records := client.GetAllRecords()

	return ParseRecordsByLine(all_records, line)
}

func (client Client) GetRecordsByTxt(txt string, name string) []Record {
	var records []Record
	if name != "" {
		fmt.Println("name: " + name)
		records = client.GetRecordsByName(name)
	} else {
		records = client.GetAllRecords()
	}

	return ParseRecordsByTxt(records, txt)
}

func ParseRecordsByTxt(all_records []Record, txt string) []Record {
	var records []Record

	fmt.Println("Len: " + fmt.Sprint(len(all_records)))

	for _, record := range all_records {
		fmt.Println(record.Type + " " + record.Txtdata)
		if record.Type == "TXT" && record.Txtdata == txt {
			records = append(records, record)
		}
	}

	return records
}

func ParseRecordsByName(all_records []Record, name string) []Record {
	n := name
	if !strings.HasSuffix(n, ".") {
		n = name + "."
	}

	var records []Record

	for _, record := range all_records {
		if record.Name == name {
			records = append(records, record)
		}
	}

	return records
}

func ParseRecordsByLine(all_records []Record, line int) Record {
	var records []Record

	for _, record := range all_records {
		if record.Line == strconv.Itoa(line) {
			records = append(records, record)
		}
	}
	if len(records) > 0 {
		return records[0]
	}

	panic(errors.New("No records found"))
}
