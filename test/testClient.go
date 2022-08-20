package test

import (
	"anxdns-go/anxdns"
	"fmt"
	"io/ioutil"
	"os"
)

func _communicate(apiRequest anxdns.Request) []byte {
	if apiRequest.Type == anxdns.GET { //GET
		if apiRequest.QueryParams != "?domain=test.com" {
			panic("QueryParameters not correct")
		}

		file, err := os.Open("./test/data.json")
		// if we os.Open returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}
		// defer the closing of our jsonFile so that we can parse it later on
		defer file.Close()

		byteValue, _ := ioutil.ReadAll(file)

		return byteValue
	} else if apiRequest.Type == anxdns.POST { //ADD
		//"{\"domain\":\"test.com\",\"type\":\"TXT\",\"name\":\"test.test.com\",\"ttl\":3600,\"txtData\":\"1234567890\"}"
		fmt.Println("TEST")

	} else if apiRequest.Type == anxdns.PUT { //DELETE
	}
	return []byte{}
}

type TestClient struct {
	anxdns.Client
}

func NewTestClient() *TestClient {
	return &TestClient{
		anxdns.Client{
			Domain:      "test.com",
			Communicate: _communicate,
		},
	}
}
