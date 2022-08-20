package anxdns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
)

const (
	DefaultTTL int = 3600
)

type Communicate func(apiRequest Request) ([]byte, error)

type Client struct {
	Domain      string
	ApiKey      string
	BaseUrl     string `default:"https://dyn.anx.se/api/dns/"`
	Communicate Communicate
}

func NewClient(Domain string, ApiKey string) *Client {
	return &Client{
		Domain:      Domain,
		ApiKey:      ApiKey,
		Communicate: _communicate,
		BaseUrl:     "https://dyn.anx.se/api/dns/",
	}
}

func _communicate(apiRequest Request) ([]byte, error) {
	// Create client
	httpClient := &http.Client{}

	var request *http.Request
	var error error

	if apiRequest.JsonData == nil {
		request, error = http.NewRequest(apiRequest.Type, apiRequest.BaseUrl+apiRequest.QueryParams, nil)
	} else {
		request, error = http.NewRequest(apiRequest.Type, apiRequest.BaseUrl+apiRequest.QueryParams, bytes.NewBuffer(apiRequest.JsonData))
	}

	if error != nil {
		return nil, error
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("apiKey", apiRequest.ApiKey)

	response, error := httpClient.Do(request)
	if error != nil {
		return nil, error
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	// Read response body
	respBody, _ := ioutil.ReadAll(response.Body)

	// Display results
	//fmt.Println("response Status : ", response.Status)
	//fmt.Println("response Body : ", string(respBody))

	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		return nil, fmt.Errorf("Could not communicate with server")
	}

	return respBody, nil
}

func checkSubdomainPartOfDomain(domain string, name string) bool {
	var n string

	if strings.HasSuffix(name, ".") {
		n = string(name[0 : len(name)-1])
	} else {
		n = name
	}

	if domainutil.Domain(n) != domain {
		return false
	}
	return true
}

func (client Client) AddTxtRecord(name string, txt string, ttl int) {
	if !checkSubdomainPartOfDomain(client.Domain, name) {
		panic("Name not part of domain")
	}

	record := Data{
		Domain:  client.Domain,
		Type:    "TXT",
		Name:    secureSuffixxDot(name),
		TTL:     ttl,
		TxtData: txt,
	}

	jsonData, _ := json.Marshal(record)

	apiRequest := Request{
		Type:     POST,
		JsonData: jsonData,
		BaseUrl:  client.BaseUrl,
		ApiKey:   client.ApiKey,
	}

	client.Communicate(apiRequest)
}

func (client Client) AddARecord(name string, address string, ttl int) {
	if !checkSubdomainPartOfDomain(client.Domain, name) {
		panic("Name not part of domain")
	}

	record := Data{
		Domain:  client.Domain,
		Type:    "A",
		Name:    secureSuffixxDot(name),
		TTL:     ttl,
		Address: address,
	}

	jsonData, _ := json.Marshal(record)

	apiRequest := Request{
		Type:     POST,
		JsonData: jsonData,
		BaseUrl:  client.BaseUrl,
		ApiKey:   client.ApiKey,
	}

	client.Communicate(apiRequest)
}

func (client Client) AddCNameRecord(name string, address string, ttl int) {
	if !checkSubdomainPartOfDomain(client.Domain, name) {
		panic("Name not part of domain")
	}

	record := Data{
		Domain:  client.Domain,
		Type:    "CNAME",
		Name:    secureSuffixxDot(name),
		TTL:     ttl,
		Address: address,
	}

	jsonData, _ := json.Marshal(record)

	apiRequest := Request{
		Type:     POST,
		JsonData: jsonData,
		BaseUrl:  client.BaseUrl,
		ApiKey:   client.ApiKey,
	}

	client.Communicate(apiRequest)
}

func (client Client) VerifyOrGetRecord(line int, name string, recordType string) (*Record, error) {
	var record *Record
	var error error
	if line > 0 {
		record, error = client.getRecordsByLine(line)
		if error != nil {
			return nil, error
		}
	} else if len(name) > 0 {
		records, error := client.GetRecordsByName(secureSuffixxDot(name))
		if error != nil {
			return nil, error
		}
		if len(records) == 0 {
			return nil, fmt.Errorf(("0 records with that name"))
		} else if len(records) > 1 {
			return nil, fmt.Errorf((">1 record with that name. Specify line instead of name."))
		}
		record = &records[0]
	} else {
		return nil, fmt.Errorf("Line or name needs to be provided")
	}

	if len(recordType) > 0 && record.Type != recordType {
		return nil, fmt.Errorf("Record is not a " + recordType)
	}

	return record, nil

}

func (client Client) DeleteRecordsByLine(line int) error {
	// Find line
	record, error := client.VerifyOrGetRecord(line, "", "")

	if error != nil {
		return error
	}

	data := Data{
		Domain: client.Domain,
		Type:   "CNAME",
		Name:   record.Name,
		Line:   record.Line,
	}

	jsonData, _ := json.Marshal(data)

	apiRequest := Request{
		Type:     DELETE,
		JsonData: jsonData,
		BaseUrl:  client.BaseUrl,
		ApiKey:   client.ApiKey,
	}

	_, error = client.Communicate(apiRequest)
	return error
}

func (client Client) DeleteRecordsByName(name string) error {
	// TODO: Add to verify also type of record
	record, error := client.VerifyOrGetRecord(-1, secureSuffixxDot(name), "")
	if error != nil {
		return error
	}

	data := Data{
		Domain: client.Domain,
		Type:   "CNAME",
		Name:   record.Name,
		Line:   record.Line,
	}

	jsonData, _ := json.Marshal(data)

	apiRequest := Request{
		Type:     DELETE,
		JsonData: jsonData,
		BaseUrl:  client.BaseUrl,
		ApiKey:   client.ApiKey,
	}

	_, error = client.Communicate(apiRequest)
	return error
}

func (client Client) DeleteRecordsByTxt(name string, txt string) error {
	// Find name
	records, error := client.GetRecordsByTxt(txt, secureSuffixxDot(name))
	if error != nil {
		return error
	}

	if len(records) == 0 {
		return fmt.Errorf("0 records with that name.")
	} else if len(records) > 1 {
		return fmt.Errorf(">1 record with that name. Specify line instead of name.")
	}

	record := records[0]

	data := Data{
		Domain: client.Domain,
		Type:   "CNAME",
		Name:   record.Name,
		Line:   record.Line,
	}

	jsonData, _ := json.Marshal(data)

	apiRequest := Request{
		Type:     DELETE,
		JsonData: jsonData,
		BaseUrl:  client.BaseUrl,
		ApiKey:   client.ApiKey,
	}

	_, error = client.Communicate(apiRequest)
	return error
}

func (client *Client) GetAllRecords() ([]Record, error) {
	request := Request{
		Type:        GET,
		QueryParams: "?domain=" + client.Domain,
		BaseUrl:     client.BaseUrl,
		ApiKey:      client.ApiKey,
	}
	respBody, error := client.Communicate(request)

	if error != nil {
		return nil, error
	}

	response := Response{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, error
	}

	return response.DnsRecords, nil
}

func (client Client) GetRecordsByName(name string) ([]Record, error) {
	all_records, error := client.GetAllRecords()
	if error != nil {
		return nil, error
	}

	return ParseRecordsByName(all_records, secureSuffixxDot(name)), nil
}

func (client Client) getRecordsByLine(line int) (*Record, error) {
	all_records, error := client.GetAllRecords()
	if error != nil {
		return nil, error
	}

	return ParseRecordsByLine(all_records, line), nil
}

func (client Client) GetRecordsByTxt(txt string, name string) ([]Record, error) {
	var records []Record
	var error error
	if name != "" {
		records, error = client.GetRecordsByName(secureSuffixxDot(name))
	} else {
		records, error = client.GetAllRecords()
	}

	if error != nil {
		return nil, error
	}

	return ParseRecordsByTxt(records, txt), nil
}

func ParseRecordsByTxt(all_records []Record, txt string) []Record {
	var records []Record

	for _, record := range all_records {
		if record.Type == "TXT" && record.Txtdata == txt {
			records = append(records, record)
		}
	}

	return records
}

func secureSuffixxDot(name string) string {
	if !strings.HasSuffix(name, ".") {
		return name + "."
	}
	return name
}

func ParseRecordsByName(all_records []Record, name string) []Record {
	var records []Record

	for _, record := range all_records {
		if record.Name == name {
			records = append(records, record)
		}
	}

	return records
}

func ParseRecordsByLine(all_records []Record, line int) *Record {
	var records []Record

	for _, record := range all_records {
		if record.Line == strconv.Itoa(line) {
			records = append(records, record)
		}
	}
	if len(records) > 0 {
		return &records[0]
	}

	return nil
}
