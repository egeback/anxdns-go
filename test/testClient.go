package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/egeback/anxdns-go/anxdns"
)

func _communicate(apiRequest anxdns.Request) ([]byte, error) {
	if apiRequest.Type == anxdns.GET { //GET
		if apiRequest.QueryParams != "?domain=test.com" {
			panic("QueryParameters not correct")
		}
		if len(apiRequest.BaseUrl) == 0 {
			panic("BaseUrl not set")
		}

		file, err := os.Open("./test/data.json")
		// if we os.Open returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}
		// defer the closing of our jsonFile so that we can parse it later on
		defer file.Close()

		byteValue, _ := ioutil.ReadAll(file)

		return byteValue, nil
	} else if apiRequest.Type == anxdns.POST { //ADD
		data := anxdns.Data{}
		if err := json.Unmarshal(apiRequest.JsonData, &data); err != nil {
			panic(err)
		}
		//TODO: Add more checks
		if data.Type == "TXT" {
			if data.Domain != "test.com" {
				panic("Wrong DOMAIN in request to add txt record")
			}
			if data.TxtData != "1234567890" {
				panic("Wrong TXTDATA in request to add txt record")
			}
			if data.Name != "test.test.com." {
				panic("Wrong ADDRESS in request to add txt record")
			}
			if data.TTL != anxdns.DefaultTTL {
				panic("Wrong TTL in request to add txt record")
			}
		} else {
			return nil, fmt.Errorf("Wrong data provided in api call to add record")
		}
	} else if apiRequest.Type == anxdns.DELETE { //DELETE
		data := anxdns.Data{}
		if err := json.Unmarshal(apiRequest.JsonData, &data); err != nil {
			panic(err)
		}
		if data.Type == "TXT" {
			if data.Domain != "test.com" {
				panic("Wrong DOMAIN in request to add txt record")
			}
			if data.TxtData != "1234567890" {
				panic("Wrong TXTDATA in request to add txt record")
			}
			if data.Name != "test.test.com." {
				panic("Wrong ADDRESS in request to add txt record")
			}
		}
	}
	return []byte{}, nil
}

type TestClient struct {
	anxdns.Client
}

func NewTestClient() *TestClient {
	return &TestClient{
		anxdns.Client{
			Domain:      "test.com",
			Communicate: _communicate,
			BaseUrl:     "https://dyn.anx.se/api/dns/",
		},
	}
}
