package main

import (
	"anxdns-go/anxdns"
	"anxdns-go/test"
	"testing"
)

var (
	client test.TestClient
)

func init() {
	client = *test.NewTestClient()
}

func TestGetAllNr(t *testing.T) {
	all_records, error := client.GetAllRecords()
	if error != nil {
		panic(error)
	}
	if len(all_records) != 7 {
		t.Errorf("len(client.GetAllRecords()) = %d; want 7", len(all_records))
	}
}

func TestGetByName1(t *testing.T) {
	all_records, error := client.GetRecordsByName("test.test.com.")
	if error != nil {
		panic(error)
	}
	if len(all_records) != 1 {
		t.Errorf("len(client.GetRecordsByTxt(\"test\",\"test.test.com.\")) = %d; want 1", len(all_records))
	}
}

func TestGetByName2(t *testing.T) {
	all_records, error := client.GetRecordsByName("test.com.")
	if error != nil {
		panic(error)
	}
	if len(all_records) != 2 {
		t.Errorf("len(client.GetRecordsByTxt(\"test\",\"test.test.com.\")) = %d; want 2", len(all_records))
	}
	if all_records[0].Type != "A" {
		t.Errorf("record.Type == %s; want A", all_records[0].Type)
	}
	if all_records[1].Type != "TXT" {
		t.Errorf("record.Type == %s; want TXT", all_records[1].Type)
	}

}

func TestGetByTxt(t *testing.T) {
	all_records, error := client.GetRecordsByTxt("test", "test.test.com.")
	if error != nil {
		panic(error)
	}
	if len(all_records) != 1 {
		t.Errorf("len(client.GetRecordsByTxt(\"test\",\"test.test.com.\")) = %d; want 1", len(all_records))
	}
	if all_records[0].Txtdata != "test" {
		t.Errorf("record.Txtdata is not test")
	}
	if all_records[0].Name != "test.test.com." {
		t.Errorf("record.Name is not test.test.com.")
	}
}

func TestAddTxtRecord(t *testing.T) {
	client.AddTxtRecord("test.test.com", "1234567890", anxdns.DefaultTTL)
}

func TestDeletTxteRecord(t *testing.T) {
	client.DeleteRecordsByTxt("test.test.com", "test")
}
