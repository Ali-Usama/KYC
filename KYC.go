package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"io"
	"os"
	"strconv"
)

type KYC struct {
	contractapi.Contract
	NextClientID int
	NextBankID   int
}

type CustomerData struct {
	Name         string         `json:"name"`
	DateOfBirth  string         `json:"dateOfBirth"`
	Address      string         `json:"address"`
	IdNumber     int            `json:"idNumber"`
	PhoneNumber  string         `json:"phoneNumber"`
	RegisteredBy OrgCredentials `json:"registeredBy"`
}

type BankData struct {
	Name           string         `json:"name"`
	IdNumber       int            `json:"idNumber"`
	OrgCredentials OrgCredentials `json:"orgCredentials"`
}

type OrgCredentials struct {
	OrgName string `json:"orgName"`
	OrgNum  int    `json:"orgNum"`
}

// InitLedger adds initial customers and financial institutions to the ledger
func (s *KYC) InitLedger(ctx contractapi.TransactionContextInterface) error {
	file, err := os.OpenFile("data/customers.json", os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	// Read the contents of the file
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var customers []CustomerData
	err = json.Unmarshal(content, &customers)
	if err != nil {
		return err
	}

	for _, customer := range customers {
		customerJSON, err := json.Marshal(customer)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(strconv.Itoa(customer.IdNumber), customerJSON)
		if err != nil {
			return fmt.Errorf("failed to insert the customer into world state: #{err}")
		}
	}
	return nil
}

func main() {

}
