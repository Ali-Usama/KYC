package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
)

type KYC struct {
	contractapi.Contract
	NextClientID int `default:"1"`
	NextBankID   int `default:"1"`
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
		customerID := s.NextClientID
		customerJSON, err := json.Marshal(customer)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(strconv.Itoa(customerID), customerJSON)
		if err != nil {
			return fmt.Errorf("failed to insert the customer into world state: #{err}")
		}
		s.NextClientID++
	}

	bankFile, err := os.OpenFile("data/bankData.json", os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer func(bankFile *os.File) {
		err := bankFile.Close()
		if err != nil {
			panic(err)
		}
	}(bankFile)

	bankContents, err := io.ReadAll(bankFile)
	if err != nil {
		return err
	}
	var banks []BankData
	err = json.Unmarshal(bankContents, &banks)
	if err != nil {
		return err
	}
	for _, bank := range banks {
		bankID := s.NextBankID
		bankJSON, err := json.Marshal(bank)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(strconv.Itoa(bankID), bankJSON)
		if err != nil {
			return fmt.Errorf("failed to insert the bank into world state: #{err}")
		}
		s.NextBankID++
	}
	return nil
}

func main() {
	KYCchaincode, err := contractapi.NewChaincode(&KYC{
		NextBankID:   1,
		NextClientID: 1,
	})
	if err != nil {
		log.Panicf("Error creating KYC chaincode: #{err}")
	}

	if err := KYCchaincode.Start(); err != nil {
		log.Panicf("Error starting KYC chaincode: #{err}")
	}

}
