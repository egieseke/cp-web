/*
Copyright 2016 IBM

Licensed under the Apache License, Version 2.0 (the "License")
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Licensed Materials - Property of IBM
© Copyright IBM Corp. 2016
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var vehiclePrefix = "vehicle:"
var accountPrefix = "acct:"
var licensePrefix = "license:"
var registrationPrefix = "reg:"
var violationPrefix = "vio:"
var tollPrefix = "toll:"

var titleKeys = "TitleKeys"
var licenseKeys = "LicenseKeys"
var registrationKeys = "RegistrationKeys"
var violationKeys = "TrafficViolationKeys"
var tollKeys = "TollKeys"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func generateCUSIPSuffix(issueDate string, days int) (string, error) {

	t, err := msToTime(issueDate)
	if err != nil {
		return "", err
	}

	maturityDate := t.AddDate(0, 0, days)
	month := int(maturityDate.Month())
	day := maturityDate.Day()

	suffix := seventhDigit[month] + eigthDigit[day]
	return suffix, nil

}

const (
	millisPerSecond     = int64(time.Second / time.Millisecond)
	nanosPerMillisecond = int64(time.Millisecond / time.Nanosecond)
)

func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(msInt/millisPerSecond,
		(msInt%millisPerSecond)*nanosPerMillisecond), nil
}

type VehicleRegistration struct {
	RegistrationId string `json:"registrationId"`
	PlateNum       string `json:"plateNum"`
	VIN            string `json:"vin"`
	TestId         string `json:"testId"`
	PolicyId       string `json:"policyId"`
	Owner          string `json:"owner"`
	AutoRenewal    string `json:"auto"`
	IssueDate      string `json:"issueDate"`
	ExpiryDate     string `json:"expiryDate"`
}

type DriverLicense struct {
	LicenseId  string `json:"licenseId"`
	TestId     string `json:"testId"`
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	Zip        string `json:"zip"`
	Driver     string `json:"driver"`
	AutoRenewal    string `json:"auto"`
	IssueDate  string `json:"issueDate"`
	ExpiryDate string `json:"expiryDate"`
}

type CP struct {
	VIN       string  `json:"vin"`
	Make      string  `json:"make"`
	Model     string  `json:"model"`
	Year      int     `json:"year"`
	Color     string  `json:"color"`
	Miles     int     `json:"miles"`
	Value     float64 `json:"value"`
	Owner     string `json:"owner"`
	Issuer    string  `json:"issuer"`
	State     string  `json:"state"`
	IssueDate string  `json:"issueDate"`
}

type Account struct {
	ID          string   `json:"id"`
	Prefix      string   `json:"prefix"`
	CashBalance float64  `json:"cashBalance"`
	AssetsIds   []string `json:"assetIds"`
}

type RenewLicenseTx struct {
	TxId       string `json:"txId"`
	LicenseId  string `json:"licenseId"`
	Driver     string `json:"driver"`
	IssueDate  string `json:"issueDate"`
	ExpiryDate string `json:"expiryDate"`
}

type TrafficViolationTx struct {
	TxId       string `json:"txId"`
	ViolationType string `json:"type"`
	LicenseId  string `json:"licenseId"`
	Driver     string `json:"driver"`
	IssueDate  string `json:"issueDate"`
	Fine       float64  `json:"fine"`
	Location   string `json:"location"`
}

type TollTx struct {
	TxId       string `json:"txId"`
	TollType string `json:"type"`
	RegistrationId  string `json:"registrationId"`
	Driver     string `json:"owner"`
	IssueDate  string `json:"issueDate"`
	Toll       float64  `json:"tollAmt"`
	Location   string `json:"location"`
}

type RenewRegistrationTx struct {
	TxId       string `json:"txId"`
	RegistrationId  string `json:"registrationId"`
	Owner string `json:"owner"`
	IssueDate  string `json:"issueDate"`
	ExpiryDate string `json:"expiryDate"`
}

type TransferTitleTx struct {
	VIN         string  `json:"vin"`
	FromOwner string  `json:"fromOwner"`
	ToOwner string  `json:"toOwner"`
	IssueDate  string `json:"issueDate"`
        AmountPaid  float64 `json:"amountPaid"`
}

type TerminateAssetTx struct {
	VIN         string  `json:"vin"`
	Owner string  `json:"owner"`
	IssueDate  string `json:"issueDate"`
}

func (t *SimpleChaincode) createAccounts(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Creating accounts")

	//  				0
	// "number of accounts to create"
	var err error
	numAccounts, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("error creating accounts with input")
		return nil, errors.New("createAccounts accepts a single integer argument")
	}
	//create a bunch of accounts
	var account Account
	counter := 1
	for counter <= numAccounts {
		var prefix string
		suffix := "000A"
		if counter < 10 {
			prefix = strconv.Itoa(counter) + "0" + suffix
		} else {
			prefix = strconv.Itoa(counter) + suffix
		}
		var assetIds []string
		account = Account{ID: "company" + strconv.Itoa(counter), Prefix: prefix, CashBalance: 10000.0, AssetsIds: assetIds}
		accountBytes, err := json.Marshal(&account)
		if err != nil {
			fmt.Println("error creating account" + account.ID)
			return nil, errors.New("Error creating account " + account.ID)
		}
		err = stub.PutState(accountPrefix+account.ID, accountBytes)
		counter++
		fmt.Println("created account" + accountPrefix + account.ID)
	}

	fmt.Println("Accounts created")
	return nil, nil

}

func (t *SimpleChaincode) createAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Creating account")

	// Obtain the username to associate with the account
	if len(args) != 1 {
		fmt.Println("Error obtaining username")
		return nil, errors.New("createAccount accepts a single username argument")
	}
	username := args[0]

	// Build an account object for the user
	var assetIds []string
	suffix := "000A"
	prefix := username + suffix
	var account = Account{ID: username, Prefix: prefix, CashBalance: 10000.0, AssetsIds: assetIds}
	accountBytes, err := json.Marshal(&account)
	if err != nil {
		fmt.Println("error creating account" + account.ID)
		return nil, errors.New("Error creating account " + account.ID)
	}

	fmt.Println("Attempting to get state of any existing account for " + account.ID)
	existingBytes, err := stub.GetState(accountPrefix + account.ID)
	if err == nil {

		var company Account
		err = json.Unmarshal(existingBytes, &company)
		if err != nil {
			fmt.Println("Error unmarshalling account " + account.ID + "\n--->: " + err.Error())

			if strings.Contains(err.Error(), "unexpected end") {
				fmt.Println("No data means existing account found for " + account.ID + ", initializing account.")
				err = stub.PutState(accountPrefix+account.ID, accountBytes)

				if err == nil {
					fmt.Println("created account" + accountPrefix + account.ID)
					return nil, nil
				} else {
					fmt.Println("failed to create initialize account for " + account.ID)
					return nil, errors.New("failed to initialize an account for " + account.ID + " => " + err.Error())
				}
			} else {
				return nil, errors.New("Error unmarshalling existing account " + account.ID)
			}
		} else {
			fmt.Println("Account already exists for " + account.ID + " " + company.ID)
			return nil, errors.New("Can't reinitialize existing user " + account.ID)
		}
	} else {

		fmt.Println("No existing account found for " + account.ID + ", initializing account.")
		err = stub.PutState(accountPrefix+account.ID, accountBytes)

		if err == nil {
			fmt.Println("created account" + accountPrefix + account.ID)
			return nil, nil
		} else {
			fmt.Println("failed to create initialize account for " + account.ID)
			return nil, errors.New("failed to initialize an account for " + account.ID + " => " + err.Error())
		}

	}

}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Init firing. Function will be ignored: " + function)

	// Initialize the collection of commercial paper keys
	fmt.Println("Initializing keys")
	var blank []string
	blankBytes, _ := json.Marshal(&blank)
	err := stub.PutState(titleKeys, blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize title key collection")
	}
	err = stub.PutState(licenseKeys, blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize license key collection")
	}
	err = stub.PutState(registrationKeys, blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize registration key collection")
	}
	err = stub.PutState(violationKeys, blankBytes)
        if err != nil {
                fmt.Println("Failed to initialize traffic violation key collection")
        }
        err = stub.PutState(tollKeys, blankBytes)
        if err != nil {
                fmt.Println("Failed to initialize toll key collection")
        }
	fmt.Println("Initialization complete")
	return nil, nil
}

func (t *SimpleChaincode) issueVehicleRegistration(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Creating Vehicle Registration")

	//need one arg
	if len(args) != 1 {
		fmt.Println("error invalid arguments")
		return nil, errors.New("Incorrect number of arguments. Expecting vehicle registration record")
	}

	var registration VehicleRegistration
	var err error
	var account Account

	fmt.Println("Unmarshalling VehicleRegistration")
	err = json.Unmarshal([]byte(args[0]), &registration)
	if err != nil {
		fmt.Println("error invalid registration issue")
		return nil, errors.New("Invalid VehicleRegistration issue")
	}

	//get account prefix
	fmt.Println("Getting state of - " + accountPrefix + registration.Owner)
	accountBytes, err := stub.GetState(accountPrefix + registration.Owner)
	if err != nil {
		fmt.Println("Error Getting state of - " + accountPrefix + registration.Owner)
		return nil, errors.New("Error retrieving account " + registration.Owner)
	}
	err = json.Unmarshal(accountBytes, &account)
	if err != nil {
		fmt.Println("Error Unmarshalling accountBytes")
		return nil, errors.New("Error retrieving account " + registration.Owner)
	}

	account.AssetsIds = append(account.AssetsIds, registration.RegistrationId)

	var govaccount Account
	govaccountBytes, err := stub.GetState(accountPrefix + "government")
	if err != nil {
		fmt.Println("Error Getting state of - " + accountPrefix + "government")
		return nil, errors.New("Error retrieving account acct:government ")
	}
	err = json.Unmarshal(govaccountBytes, &govaccount)
	if err != nil {
		fmt.Println("Error Unmarshalling govaccountBytes")
		return nil, errors.New("Error retrieving account acct:government")
	}
	//deduct fee from account
	account.CashBalance -= 40
	govaccount.CashBalance += 40


	fmt.Println("Getting State on Vehicle Registration" + registration.RegistrationId)
	cpRxBytes, err := stub.GetState(registrationPrefix + registration.RegistrationId)
	if cpRxBytes == nil {
		fmt.Println("Registration does not exist, creating it")
		licenseBytes, err := json.Marshal(&registration)
		if err != nil {
			fmt.Println("Error marshalling registration")
			return nil, errors.New("Error issuing vehicle registration")
		}
		err = stub.PutState(registrationPrefix+registration.RegistrationId, licenseBytes)
		if err != nil {
			fmt.Println("Error issuing registration")
			return nil, errors.New("Error issuing vehicle registration")
		}

		fmt.Println("Marshalling account bytes to write")
		accountBytesToWrite, err := json.Marshal(&account)
		if err != nil {
			fmt.Println("Error marshalling account")
			return nil, errors.New("Error issuing vehicle registration")
		}
		err = stub.PutState(accountPrefix+registration.Owner, accountBytesToWrite)
		if err != nil {
			fmt.Println("Error putting state on accountBytesToWrite")
			return nil, errors.New("Error issuing vehicle registration")
		}

		govaccountBytesToWrite, err := json.Marshal(&govaccount)
		if err != nil {
			fmt.Println("Error marshalling govt account")
			return nil, errors.New("Error issuing vehicle registration")
		}
		err = stub.PutState(accountPrefix+"government", govaccountBytesToWrite)
		if err != nil {
			fmt.Println("Error putting state on govaccountBytesToWrite")
			return nil, errors.New("Error issuing vehicle registration")
		}

		// Update the registration keys by adding the new key
		fmt.Println("Getting Registration Keys")
		keysBytes, err := stub.GetState(registrationKeys)
		if err != nil {
			fmt.Println("Error retrieving registration keys")
			return nil, errors.New("Error retrieving registration keys")
		}
		var keys []string
		err = json.Unmarshal(keysBytes, &keys)
		if err != nil {
			fmt.Println("Error unmarshel keys")
			return nil, errors.New("Error unmarshalling registration keys ")
		}

		fmt.Println("Appending the new key to Registration Keys")
		foundKey := false
		for _, key := range keys {
			if key == registrationPrefix+registration.RegistrationId {
				foundKey = true
			}
		}
		if foundKey == false {
			keys = append(keys, registrationPrefix+registration.RegistrationId)
			keysBytesToWrite, err := json.Marshal(&keys)
			if err != nil {
				fmt.Println("Error marshalling keys")
				return nil, errors.New("Error marshalling the keys")
			}
			fmt.Println("Put state on Registration Keys")
			err = stub.PutState(registrationKeys, keysBytesToWrite)
			if err != nil {
				fmt.Println("Error writting keys back")
				return nil, errors.New("Error writing the keys back")
			}
		}

		fmt.Println("Issue vehicle registration", registration)
	}
	return nil, nil
}

func (t *SimpleChaincode) issueDriverLicense(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Creating Driver License")

	//need one arg
	if len(args) != 1 {
		fmt.Println("error invalid arguments")
		return nil, errors.New("Incorrect number of arguments. Expecting driver license record")
	}

	var license DriverLicense
	var err error
	var account Account

	fmt.Println("Unmarshalling DriverLicense")
	err = json.Unmarshal([]byte(args[0]), &license)
	if err != nil {
		fmt.Println("error invalid license issue")
		return nil, errors.New("Invalid driver license issue")
	}

	//get account prefix
	fmt.Println("Getting state of - " + accountPrefix + license.Driver)
	accountBytes, err := stub.GetState(accountPrefix + license.Driver)
	if err != nil {
		fmt.Println("Error Getting state of - " + accountPrefix + license.Driver)
		return nil, errors.New("Error retrieving account " + license.Driver)
	}
	err = json.Unmarshal(accountBytes, &account)
	if err != nil {
		fmt.Println("Error Unmarshalling accountBytes")
		return nil, errors.New("Error retrieving account " + license.Driver)
	}

	account.AssetsIds = append(account.AssetsIds, license.LicenseId)

	var govaccount Account
	govaccountBytes, err := stub.GetState(accountPrefix + "government")
	if err != nil {
		fmt.Println("Error Getting state of - " + accountPrefix + "government")
		return nil, errors.New("Error retrieving account acct:government ")
	}
	err = json.Unmarshal(govaccountBytes, &govaccount)
	if err != nil {
		fmt.Println("Error Unmarshalling govaccountBytes")
		return nil, errors.New("Error retrieving account acct:government")
	}
	//deduct fee from account
	account.CashBalance -= 50
	govaccount.CashBalance += 50


	fmt.Println("Getting State on Driver License" + license.LicenseId)
	cpRxBytes, err := stub.GetState(licensePrefix + license.LicenseId)
	if cpRxBytes == nil {
		fmt.Println("License does not exist, creating it")
		licenseBytes, err := json.Marshal(&license)
		if err != nil {
			fmt.Println("Error marshalling license")
			return nil, errors.New("Error issuing driver license")
		}
		err = stub.PutState(licensePrefix+license.LicenseId, licenseBytes)
		if err != nil {
			fmt.Println("Error issuing license")
			return nil, errors.New("Error issuing driver license")
		}

		fmt.Println("Marshalling account bytes to write")
		accountBytesToWrite, err := json.Marshal(&account)
		if err != nil {
			fmt.Println("Error marshalling account")
			return nil, errors.New("Error issuing driver license")
		}
		err = stub.PutState(accountPrefix+license.Driver, accountBytesToWrite)
		if err != nil {
			fmt.Println("Error putting state on accountBytesToWrite")
			return nil, errors.New("Error issuingdriver licensepaper")
		}

		govaccountBytesToWrite, err := json.Marshal(&govaccount)
		if err != nil {
			fmt.Println("Error marshalling govt account")
			return nil, errors.New("Error issuing driver license")
		}
		err = stub.PutState(accountPrefix+"government", govaccountBytesToWrite)
		if err != nil {
			fmt.Println("Error putting state on govaccountBytesToWrite")
			return nil, errors.New("Error issuing driver license")
		}

		// Update the License keys by adding the new key
		fmt.Println("Getting License Keys")
		keysBytes, err := stub.GetState(licenseKeys)
		if err != nil {
			fmt.Println("Error retrieving license keys")
			return nil, errors.New("Error retrieving license keys")
		}
		var keys []string
		err = json.Unmarshal(keysBytes, &keys)
		if err != nil {
			fmt.Println("Error unmarshel keys")
			return nil, errors.New("Error unmarshalling paper keys ")
		}

		fmt.Println("Appending the new key to License Keys")
		foundKey := false
		for _, key := range keys {
			if key == licensePrefix+license.LicenseId {
				foundKey = true
			}
		}
		if foundKey == false {
			keys = append(keys, licensePrefix+license.LicenseId)
			keysBytesToWrite, err := json.Marshal(&keys)
			if err != nil {
				fmt.Println("Error marshalling keys")
				return nil, errors.New("Error marshalling the keys")
			}
			fmt.Println("Put state on License Keys")
			err = stub.PutState(licenseKeys, keysBytesToWrite)
			if err != nil {
				fmt.Println("Error writting keys back")
				return nil, errors.New("Error writing the keys back")
			}
		}

		fmt.Println("Issue driver license", license)
	}
	return nil, nil
}

/* 
asset: new title object: owner=issuer
tx fee deducted from issuer and added to gov
Account: assetId list updated
Error: if vin already exists - throw error 
Todo: Account has enough money
*/
func (t *SimpleChaincode) issueVehicleTitle(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Creating vehicle title")
        var titleFee = 30.0
	//need one arg
	if len(args) != 1 {
		fmt.Println("error invalid arguments")
		return nil, errors.New("Incorrect number of arguments. Expecting title record")
	}

	var cp CP
	var err error
	var account Account

	fmt.Println("Unmarshalling CP")
	err = json.Unmarshal([]byte(args[0]), &cp)
	if err != nil {
		fmt.Println("error invalid paper issue")
		return nil, errors.New("Invalid commercial paper issue")
	}

	//get account prefix
	fmt.Println("Getting state of - " + accountPrefix + cp.Issuer)
	accountBytes, err := stub.GetState(accountPrefix + cp.Issuer)
	if err != nil {
		fmt.Println("Error Getting state of - " + accountPrefix + cp.Issuer)
		return nil, errors.New("Error retrieving account " + cp.Issuer)
	}
	err = json.Unmarshal(accountBytes, &account)
	if err != nil {
		fmt.Println("Error Unmarshalling accountBytes")
		return nil, errors.New("Error retrieving account " + cp.Issuer)
	}

	account.AssetsIds = append(account.AssetsIds, cp.VIN)

	var govaccount Account
	govaccountBytes, err := stub.GetState(accountPrefix + "government")
	if err != nil {
		fmt.Println("Error Getting state of - " + accountPrefix + "government")
		return nil, errors.New("Error retrieving account acct:government ")
	}
	err = json.Unmarshal(govaccountBytes, &govaccount)
	if err != nil {
		fmt.Println("Error Unmarshalling govaccountBytes")
		return nil, errors.New("Error retrieving account acct:government")
	}

	//deduct fee from account
	account.CashBalance -= titleFee
	govaccount.CashBalance += titleFee

	cp.Owner = cp.Issuer

	fmt.Println("Getting State on CP " + cp.VIN)
	cpRxBytes, err := stub.GetState(vehiclePrefix + cp.VIN)
	if cpRxBytes == nil {
		fmt.Println("VIN does not exist, creating it")
		cpBytes, err := json.Marshal(&cp)
		if err != nil {
			fmt.Println("Error marshalling cp")
			return nil, errors.New("Error issuing title")
		}
		err = stub.PutState(vehiclePrefix+cp.VIN, cpBytes)
		if err != nil {
			fmt.Println("Error issuing title")
			return nil, errors.New("Error issuing title")
		}

		fmt.Println("Marshalling account bytes to write")
		accountBytesToWrite, err := json.Marshal(&account)
		if err != nil {
			fmt.Println("Error marshalling account")
			return nil, errors.New("Error issuing title")
		}
		err = stub.PutState(accountPrefix+cp.Issuer, accountBytesToWrite)
		if err != nil {
			fmt.Println("Error putting state on accountBytesToWrite")
			return nil, errors.New("Error issuing title")
		}

		govaccountBytesToWrite, err := json.Marshal(&govaccount)
		if err != nil {
			fmt.Println("Error marshalling govt account")
			return nil, errors.New("Error issuing title")
		}
		err = stub.PutState(accountPrefix+"government", govaccountBytesToWrite)
		if err != nil {
			fmt.Println("Error putting state on govaccountBytesToWrite")
			return nil, errors.New("Error issuing title")
		}

		// Update the title keys by adding the new key
		fmt.Println("Getting Title Keys")
		keysBytes, err := stub.GetState(titleKeys)
		if err != nil {
			fmt.Println("Error retrieving title keys")
			return nil, errors.New("Error retrieving title keys")
		}
		var keys []string
		err = json.Unmarshal(keysBytes, &keys)
		if err != nil {
			fmt.Println("Error unmarshel keys")
			return nil, errors.New("Error unmarshalling title keys ")
		}

		fmt.Println("Appending the new key to Title Keys")
		foundKey := false
		for _, key := range keys {
			if key == vehiclePrefix+cp.VIN{
				foundKey = true
			}
		}
		if foundKey == false {
			keys = append(keys, vehiclePrefix+cp.VIN)
			keysBytesToWrite, err := json.Marshal(&keys)
			if err != nil {
				fmt.Println("Error marshalling keys")
				return nil, errors.New("Error marshalling the keys")
			}
			fmt.Println("Put state on Title Keys")
			err = stub.PutState(titleKeys, keysBytesToWrite)
			if err != nil {
				fmt.Println("Error writting keys back")
				return nil, errors.New("Error writing the keys back")
			}
		}

		fmt.Println("Issue commercial paper", cp)
		return nil, nil
	} else {
		fmt.Println("Error VIN exists")
	        return nil, errors.New("Error issuing title")
	}
}
func GetAllDriverLicenses(stub shim.ChaincodeStubInterface) ([]DriverLicense, error) {

	var allCPs []DriverLicense

	// Get list of all the license keys
	keysBytes, err := stub.GetState(licenseKeys)
	if err != nil {
		fmt.Println("Error retrieving license keys")
		return nil, errors.New("Error retrieving license keys")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling license keys")
		return nil, errors.New("Error unmarshalling license keys")
	}

	// Get all the cps
	for _, value := range keys {
		cpBytes, err := stub.GetState(value)

		var cp DriverLicense
		err = json.Unmarshal(cpBytes, &cp)
		if err != nil {
			fmt.Println("Error retrieving license " + value)
			return nil, errors.New("Error retrieving license " + value)
		}

		fmt.Println("Appending license" + value)
		allCPs = append(allCPs, cp)
	}

	return allCPs, nil
}
func GetAllVehicleRegistrations(stub shim.ChaincodeStubInterface) ([]VehicleRegistration, error) {

	var allCPs []VehicleRegistration

	// Get list of all the registration keys
	keysBytes, err := stub.GetState(registrationKeys)
	if err != nil {
		fmt.Println("Error retrieving registration keys")
		return nil, errors.New("Error retrieving registration keys")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling registration keys")
		return nil, errors.New("Error unmarshalling registration keys")
	}

	// Get all the cps
	for _, value := range keys {
		cpBytes, err := stub.GetState(value)

		var cp VehicleRegistration
		err = json.Unmarshal(cpBytes, &cp)
		if err != nil {
			fmt.Println("Error retrieving cp " + value)
			return nil, errors.New("Error retrieving cp " + value)
		}

		fmt.Println("Appending CP" + value)
		allCPs = append(allCPs, cp)
	}

	return allCPs, nil
}

func GetAllTolls(stub shim.ChaincodeStubInterface) ([]TollTx, error) {

        var allCPs []TollTx

        // Get list of all the toll keys
        keysBytes, err := stub.GetState(tollKeys)
        if err != nil {
                fmt.Println("Error retrieving toll keys")
                return nil, errors.New("Error retrieving toll keys")
        }
        var keys []string
        err = json.Unmarshal(keysBytes, &keys)
        if err != nil {
                fmt.Println("Error unmarshalling toll keys")
                return nil, errors.New("Error unmarshalling toll keys")
        }

        // Get all the cps
        for _, value := range keys {
                cpBytes, err := stub.GetState(value)

                var cp TollTx
                err = json.Unmarshal(cpBytes, &cp)
                if err != nil {
                        fmt.Println("Error retrieving toll " + value)
                        return nil, errors.New("Error retrieving toll " + value)
                }

                fmt.Println("Appending Toll" + value)
                allCPs = append(allCPs, cp)
        }

        return allCPs, nil
}

func GetAllViolations(stub shim.ChaincodeStubInterface) ([]TrafficViolationTx, error) {

        var allCPs []TrafficViolationTx

        // Get list of all the violation keys
        keysBytes, err := stub.GetState(violationKeys)
        if err != nil {
                fmt.Println("Error retrieving violation keys")
                return nil, errors.New("Error retrieving violation keys")
        }
        var keys []string
        err = json.Unmarshal(keysBytes, &keys)
        if err != nil {
                fmt.Println("Error unmarshalling violation keys")
                return nil, errors.New("Error unmarshalling violation keys")
        }

        // Get all the cps
        for _, value := range keys {
                cpBytes, err := stub.GetState(value)

                var cp TrafficViolationTx
                err = json.Unmarshal(cpBytes, &cp)
                if err != nil {
                        fmt.Println("Error retrieving violation " + value)
                        return nil, errors.New("Error retrieving traffic violation" + value)
                }

                fmt.Println("Appending traffic violation" + value)
                allCPs = append(allCPs, cp)
        }

        return allCPs, nil
}
func GetAllTitles(stub shim.ChaincodeStubInterface) ([]CP, error) {

	var allCPs []CP

	// Get list of all the title keys
	keysBytes, err := stub.GetState(titleKeys)
	if err != nil {
		fmt.Println("Error retrieving title keys")
		return nil, errors.New("Error retrieving title keys")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling title keys")
		return nil, errors.New("Error unmarshalling title keys")
	}

	// Get all the cps
	for _, value := range keys {
		cpBytes, err := stub.GetState(value)

		var cp CP
		err = json.Unmarshal(cpBytes, &cp)
		if err != nil {
			fmt.Println("Error retrieving title " + value)
			return nil, errors.New("Error retrieving title " + value)
		}

		fmt.Println("Appending title" + value)
		allCPs = append(allCPs, cp)
	}

	return allCPs, nil
}

func GetCP(cpid string, stub shim.ChaincodeStubInterface) (CP, error) {
	var cp CP

	cpBytes, err := stub.GetState(cpid)
	if err != nil {
		fmt.Println("Error retrieving cp " + cpid)
		return cp, errors.New("Error retrieving cp " + cpid)
	}

	err = json.Unmarshal(cpBytes, &cp)
	if err != nil {
		fmt.Println("Error unmarshalling cp " + cpid)
		return cp, errors.New("Error unmarshalling cp " + cpid)
	}

	return cp, nil
}

func GetCompany(companyID string, stub shim.ChaincodeStubInterface) (Account, error) {
	var company Account
	companyBytes, err := stub.GetState(accountPrefix + companyID)
	if err != nil {
		fmt.Println("Account not found " + companyID)
		return company, errors.New("Account not found " + companyID)
	}

	err = json.Unmarshal(companyBytes, &company)
	if err != nil {
		fmt.Println("Error unmarshalling account " + companyID + "\n err:" + err.Error())
		return company, errors.New("Error unmarshalling account " + companyID)
	}

	return company, nil
}

func (t *SimpleChaincode) renewRegistration(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var registrationRenewalFee float64 = 40
	fmt.Println("Renewing Registration")
	//need one arg
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting renew registration record")
	}

	var tr RenewRegistrationTx

	fmt.Println("Unmarshalling Transaction")
	err := json.Unmarshal([]byte(args[0]), &tr)
	if err != nil {
		fmt.Println("Error Unmarshalling Transaction")
		return nil, errors.New("Invalid renew registration record")
	}

	cpBytes, err := stub.GetState(registrationPrefix + tr.RegistrationId)
	if err != nil {
		fmt.Println("RegistrationId not found")
		return nil, errors.New("RegistrationId not found " + tr.RegistrationId)
	}

	var cp VehicleRegistration
	fmt.Println("Unmarshalling Registration " + tr.RegistrationId)
	err = json.Unmarshal(cpBytes, &cp)
	if err != nil {
		fmt.Println("Error unmarshalling registration  " + tr.RegistrationId)
		return nil, errors.New("Error unmarshalling registration " + tr.RegistrationId)
	}

	var driver Account
	fmt.Println("Getting State on Owner " + tr.Owner)
	driverBytes, err := stub.GetState(accountPrefix + tr.Owner)
	if err != nil {
		fmt.Println("Account not found " + tr.Owner)
		return nil, errors.New("Account not found " + tr.Owner)
	}

	fmt.Println("Unmarshalling Driver")
	err = json.Unmarshal(driverBytes, &driver)
	if err != nil {
		fmt.Println("Error unmarshalling account " + tr.Owner)
		return nil, errors.New("Error unmarshalling account " + tr.Owner)
	}

	var toCompany Account
	fmt.Println("Getting State on ToCompany " + "government")
	toCompanyBytes, err := stub.GetState(accountPrefix + "government")
	if err != nil {
		fmt.Println("Account not found " + "government")
		return nil, errors.New("Account not found " + "government")
	}

	fmt.Println("Unmarshalling tocompany")
	err = json.Unmarshal(toCompanyBytes, &toCompany)
	if err != nil {
		fmt.Println("Error unmarshalling account " + "government")
		return nil, errors.New("Error unmarshalling account " + "government")
	}


	// If toCompany doesn't have enough cash to buy the papers
	if driver.CashBalance < registrationRenewalFee {
		fmt.Println("The owner " + tr.Owner + "doesn't have enough cash to renew the registration")
		return nil, errors.New("The owner " + tr.Owner+ "doesn't have enough cash to renew the registration")
	} else {
		fmt.Println("The owner has enough money to renew the registration")
	}

	toCompany.CashBalance += registrationRenewalFee
	driver.CashBalance -= registrationRenewalFee

	// update the license renewal date

	cp.IssueDate = tr.IssueDate
	cp.ExpiryDate = tr.ExpiryDate


	// Write everything back
	// To Company
	toCompanyBytesToWrite, err := json.Marshal(&toCompany)
	if err != nil {
		fmt.Println("Error marshalling the government")
		return nil, errors.New("Error marshalling the government")
	}
	fmt.Println("Put state on toCompany")
	err = stub.PutState(accountPrefix+"government", toCompanyBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the government back")
		return nil, errors.New("Error writing the government back")
	}

	// Save the Driver state
	driverBytesToWrite, err := json.Marshal(&driver)
	if err != nil {
		fmt.Println("Error marshalling the driver")
		return nil, errors.New("Error marshalling the driver")
	}
	fmt.Println("Put state on driver")
	err = stub.PutState(accountPrefix+tr.Owner, driverBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the driver back")
		return nil, errors.New("Error writing the driver back")
	}

	// save the updated registration
	cpBytesToWrite, err := json.Marshal(&cp)
	if err != nil {
		fmt.Println("Error marshalling the cp")
		return nil, errors.New("Error marshalling the cp")
	}
	fmt.Println("Put state on vehicle registration")
	err = stub.PutState(registrationPrefix+tr.RegistrationId, cpBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the drivers license back")
		return nil, errors.New("Error writing the owner registration back")
	}

	fmt.Println("Successfully completed Invoke of renew vehicle registration")
	return nil, nil
}

func (t *SimpleChaincode) issueTollTicket(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Issue Toll ticket")
	//need one arg
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting  toll record")
	}

	var tr TollTx

	fmt.Println("Unmarshalling Transaction")
	err := json.Unmarshal([]byte(args[0]), &tr)
	if err != nil {
		fmt.Println("Error Unmarshalling Transaction")
		return nil, errors.New("Invalid toll record")
	}

	cpBytes, err := stub.GetState(registrationPrefix + tr.RegistrationId)
	if err != nil {
		fmt.Println("RegistrationId not found")
		return nil, errors.New("RegistrationId not found " + tr.RegistrationId)
	}

	var cp VehicleRegistration
	fmt.Println("Unmarshalling Registration" + tr.RegistrationId)
	err = json.Unmarshal(cpBytes, &cp)
	if err != nil {
		fmt.Println("Error unmarshalling registration " + tr.RegistrationId)
		return nil, errors.New("Error unmarshalling registration " + tr.RegistrationId)
	}

	var driver Account
	fmt.Println("Getting State on Driver " + tr.Driver)
	driverBytes, err := stub.GetState(accountPrefix + tr.Driver)
	if err != nil {
		fmt.Println("Account not found " + tr.Driver)
		return nil, errors.New("Account not found " + tr.Driver)
	}

	fmt.Println("Unmarshalling Driver")
	err = json.Unmarshal(driverBytes, &driver)
	if err != nil {
		fmt.Println("Error unmarshalling account " + tr.Driver)
		return nil, errors.New("Error unmarshalling account " + tr.Driver)
	}

	var toCompany Account
	fmt.Println("Getting State on ToCompany " + "government")
	toCompanyBytes, err := stub.GetState(accountPrefix + "government")
	if err != nil {
		fmt.Println("Account not found " + "government")
		return nil, errors.New("Account not found " + "government")
	}

	fmt.Println("Unmarshalling tocompany")
	err = json.Unmarshal(toCompanyBytes, &toCompany)
	if err != nil {
		fmt.Println("Error unmarshalling account " + "government")
		return nil, errors.New("Error unmarshalling account " + "government")
	}

	// If toCompany doesn't have enough cash to buy the papers
/*
	if driver.CashBalance < licenseRenewalFee {
		fmt.Println("The driver " + tr.Driver + "doesn't have enough cash to pay fine")
		return nil, errors.New("The driver " + tr.Driver + "doesn't have enough cash to pay fine")
	} else {
		fmt.Println("The driver has enough money to pay fine")
	}
*/
	toCompany.CashBalance += tr.Toll
	driver.CashBalance -= tr.Toll

        // email driver

	// Write everything back
	// To Company
	toCompanyBytesToWrite, err := json.Marshal(&toCompany)
	if err != nil {
		fmt.Println("Error marshalling the government")
		return nil, errors.New("Error marshalling the government")
	}
	fmt.Println("Put state on toCompany")
	err = stub.PutState(accountPrefix+"government", toCompanyBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the government back")
		return nil, errors.New("Error writing the government back")
	}

	// Save the Driver state
	driverBytesToWrite, err := json.Marshal(&driver)
	if err != nil {
		fmt.Println("Error marshalling the driver")
		return nil, errors.New("Error marshalling the driver")
	}
	fmt.Println("Put state on driver")
	err = stub.PutState(accountPrefix+tr.Driver, driverBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the driver back")
		return nil, errors.New("Error writing the driver back")
	}

        fmt.Println("creating toll")
        tollBytes, err := json.Marshal(&tr)
        if err != nil {
                        fmt.Println("Error marshalling toll")
                        return nil, errors.New("Error issuing toll")
        }
        err = stub.PutState(tollPrefix+tr.TxId, tollBytes)
        if err != nil {
                        fmt.Println("Error issuing toll")
                        return nil, errors.New("Error issuing toll")
        }

        // Update the toll keys by adding the new key
        fmt.Println("Getting Toll Keys")
        keysBytes, err := stub.GetState(tollKeys)
        if err != nil {
        	fmt.Println("Error retrieving toll keys")
                return nil, errors.New("Error retrieving toll keys")
        }
        var keys []string
        err = json.Unmarshal(keysBytes, &keys)
        if err != nil {
               fmt.Println("Error unmarshel keys")
               return nil, errors.New("Error unmarshalling toll keys ")
        }

        fmt.Println("Appending the new key to Toll Keys")
        keys = append(keys, tollPrefix+tr.TxId)
        keysBytesToWrite, err := json.Marshal(&keys)
        if err != nil {
               fmt.Println("Error marshalling keys")
               return nil, errors.New("Error marshalling the keys")
        }
        fmt.Println("Put state on Toll Keys")
        err = stub.PutState(tollKeys, keysBytesToWrite)
        if err != nil {
               fmt.Println("Error writting keys back")
               return nil, errors.New("Error writing the keys back")
        }
	fmt.Println("Successfully completed Invoke of issue toll ticket")
	return nil, nil
}

func (t *SimpleChaincode) issueTrafficViolation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Issue Traffic Violation ticket")
	//need one arg
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting  traffic violation record")
	}

	var tr TrafficViolationTx

	fmt.Println("Unmarshalling Transaction")
	err := json.Unmarshal([]byte(args[0]), &tr)
	if err != nil {
		fmt.Println("Error Unmarshalling Transaction")
		return nil, errors.New("Invalid traffic violation record")
	}

	cpBytes, err := stub.GetState(licensePrefix + tr.LicenseId)
	if err != nil {
		fmt.Println("LicenseId not found")
		return nil, errors.New("LicenseId not found " + tr.LicenseId)
	}

	var cp DriverLicense
	fmt.Println("Unmarshalling License " + tr.LicenseId)
	err = json.Unmarshal(cpBytes, &cp)
	if err != nil {
		fmt.Println("Error unmarshalling cp " + tr.LicenseId)
		return nil, errors.New("Error unmarshalling license " + tr.LicenseId)
	}

	var driver Account
	fmt.Println("Getting State on Driver " + tr.Driver)
	driverBytes, err := stub.GetState(accountPrefix + tr.Driver)
	if err != nil {
		fmt.Println("Account not found " + tr.Driver)
		return nil, errors.New("Account not found " + tr.Driver)
	}

	fmt.Println("Unmarshalling Driver")
	err = json.Unmarshal(driverBytes, &driver)
	if err != nil {
		fmt.Println("Error unmarshalling account " + tr.Driver)
		return nil, errors.New("Error unmarshalling account " + tr.Driver)
	}

	var toCompany Account
	fmt.Println("Getting State on ToCompany " + "government")
	toCompanyBytes, err := stub.GetState(accountPrefix + "government")
	if err != nil {
		fmt.Println("Account not found " + "government")
		return nil, errors.New("Account not found " + "government")
	}

	fmt.Println("Unmarshalling tocompany")
	err = json.Unmarshal(toCompanyBytes, &toCompany)
	if err != nil {
		fmt.Println("Error unmarshalling account " + "government")
		return nil, errors.New("Error unmarshalling account " + "government")
	}

	// If toCompany doesn't have enough cash to buy the papers
/*
	if driver.CashBalance < licenseRenewalFee {
		fmt.Println("The driver " + tr.Driver + "doesn't have enough cash to pay fine")
		return nil, errors.New("The driver " + tr.Driver + "doesn't have enough cash to pay fine")
	} else {
		fmt.Println("The driver has enough money to pay fine")
	}
*/
	toCompany.CashBalance += tr.Fine
	driver.CashBalance -= tr.Fine

	//TODO introduce violation array into license
        // email driver

	// Write everything back
	// To Company
	toCompanyBytesToWrite, err := json.Marshal(&toCompany)
	if err != nil {
		fmt.Println("Error marshalling the government")
		return nil, errors.New("Error marshalling the government")
	}
	fmt.Println("Put state on toCompany")
	err = stub.PutState(accountPrefix+"government", toCompanyBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the government back")
		return nil, errors.New("Error writing the government back")
	}

	// Save the Driver state
	driverBytesToWrite, err := json.Marshal(&driver)
	if err != nil {
		fmt.Println("Error marshalling the driver")
		return nil, errors.New("Error marshalling the driver")
	}
	fmt.Println("Put state on driver")
	err = stub.PutState(accountPrefix+tr.Driver, driverBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the driver back")
		return nil, errors.New("Error writing the driver back")
	}

	// save the updated drivers license
	cpBytesToWrite, err := json.Marshal(&cp)
	if err != nil {
		fmt.Println("Error marshalling the cp")
		return nil, errors.New("Error marshalling the cp")
	}
	fmt.Println("Put state on drivers license")
	err = stub.PutState(licensePrefix+tr.LicenseId, cpBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the drivers license back")
		return nil, errors.New("Error writing the drivers license back")
	}

        fmt.Println("creating violation")
        violationBytes, err := json.Marshal(&tr)
        if err != nil {
                        fmt.Println("Error marshalling violation")
                        return nil, errors.New("Error issuing violation")
        }
        err = stub.PutState(violationPrefix+tr.TxId, violationBytes)
        if err != nil {
                        fmt.Println("Error issuing violation")
                        return nil, errors.New("Error issuing violation")
        }

        // Update the violation keys by adding the new key
        fmt.Println("Getting violation Keys")
        keysBytes, err := stub.GetState(violationKeys)
        if err != nil {
                fmt.Println("Error retrieving violation keys")
                return nil, errors.New("Error retrieving violation keys")
        }
        var keys []string
        err = json.Unmarshal(keysBytes, &keys)
        if err != nil {
               fmt.Println("Error unmarshel keys")
               return nil, errors.New("Error unmarshalling violation keys ")
        }

        fmt.Println("Appending the new key to violation Keys")
        keys = append(keys, violationPrefix+tr.TxId)
        keysBytesToWrite, err := json.Marshal(&keys)
        if err != nil {
               fmt.Println("Error marshalling keys")
               return nil, errors.New("Error marshalling the keys")
        }
        fmt.Println("Put state on violation Keys")
        err = stub.PutState(violationKeys, keysBytesToWrite)
        if err != nil {
               fmt.Println("Error writting keys back")
               return nil, errors.New("Error writing the keys back")
        }

	fmt.Println("Successfully completed Invoke of issue traffic violation")
	return nil, nil
}

func (t *SimpleChaincode) renewLicense(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var licenseRenewalFee float64 = 50
	fmt.Println("Renewing License")
	//need one arg
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting renew licenserecord")
	}

	var tr RenewLicenseTx

	fmt.Println("Unmarshalling Transaction")
	err := json.Unmarshal([]byte(args[0]), &tr)
	if err != nil {
		fmt.Println("Error Unmarshalling Transaction")
		return nil, errors.New("Invalid renew license record")
	}

	cpBytes, err := stub.GetState(licensePrefix + tr.LicenseId)
	if err != nil {
		fmt.Println("LicenseId not found")
		return nil, errors.New("LicenseId not found " + tr.LicenseId)
	}

	var cp DriverLicense
	fmt.Println("Unmarshalling License " + tr.LicenseId)
	err = json.Unmarshal(cpBytes, &cp)
	if err != nil {
		fmt.Println("Error unmarshalling cp " + tr.LicenseId)
		return nil, errors.New("Error unmarshalling license " + tr.LicenseId)
	}

	var driver Account
	fmt.Println("Getting State on Driver " + tr.Driver)
	driverBytes, err := stub.GetState(accountPrefix + tr.Driver)
	if err != nil {
		fmt.Println("Account not found " + tr.Driver)
		return nil, errors.New("Account not found " + tr.Driver)
	}

	fmt.Println("Unmarshalling Driver")
	err = json.Unmarshal(driverBytes, &driver)
	if err != nil {
		fmt.Println("Error unmarshalling account " + tr.Driver)
		return nil, errors.New("Error unmarshalling account " + tr.Driver)
	}

	var toCompany Account
	fmt.Println("Getting State on ToCompany " + "government")
	toCompanyBytes, err := stub.GetState(accountPrefix + "government")
	if err != nil {
		fmt.Println("Account not found " + "government")
		return nil, errors.New("Account not found " + "government")
	}

	fmt.Println("Unmarshalling tocompany")
	err = json.Unmarshal(toCompanyBytes, &toCompany)
	if err != nil {
		fmt.Println("Error unmarshalling account " + "government")
		return nil, errors.New("Error unmarshalling account " + "government")
	}

	// If toCompany doesn't have enough cash to buy the papers
	if driver.CashBalance < licenseRenewalFee {
		fmt.Println("The driver " + tr.Driver + "doesn't have enough cash to renew the license")
		return nil, errors.New("The driver " + tr.Driver + "doesn't have enough cash to renew the license")
	} else {
		fmt.Println("The driver has enough money to renew the license")
	}

	toCompany.CashBalance += licenseRenewalFee
	driver.CashBalance -= licenseRenewalFee

	// update the license renewal date

	cp.IssueDate = tr.IssueDate
	cp.ExpiryDate = tr.ExpiryDate


	// Write everything back
	// To Company
	toCompanyBytesToWrite, err := json.Marshal(&toCompany)
	if err != nil {
		fmt.Println("Error marshalling the government")
		return nil, errors.New("Error marshalling the government")
	}
	fmt.Println("Put state on toCompany")
	err = stub.PutState(accountPrefix+"government", toCompanyBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the government back")
		return nil, errors.New("Error writing the government back")
	}

	// Save the Driver state
	driverBytesToWrite, err := json.Marshal(&driver)
	if err != nil {
		fmt.Println("Error marshalling the driver")
		return nil, errors.New("Error marshalling the driver")
	}
	fmt.Println("Put state on driver")
	err = stub.PutState(accountPrefix+tr.Driver, driverBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the driver back")
		return nil, errors.New("Error writing the driver back")
	}

	// save the updated drivers license
	cpBytesToWrite, err := json.Marshal(&cp)
	if err != nil {
		fmt.Println("Error marshalling the cp")
		return nil, errors.New("Error marshalling the cp")
	}
	fmt.Println("Put state on drivers license")
	err = stub.PutState(licensePrefix+tr.LicenseId, cpBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the drivers license back")
		return nil, errors.New("Error writing the drivers license back")
	}

	fmt.Println("Successfully completed Invoke of renew drivers license")
	return nil, nil
}

func (t *SimpleChaincode) terminateAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Terminating Asset")

	//need one arg
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting terminate asset record")
	}

	var tr TerminateAssetTx

	fmt.Println("Unmarshalling Transaction")
	err := json.Unmarshal([]byte(args[0]), &tr)
	if err != nil {
		fmt.Println("Error Unmarshalling Transaction")
		return nil, errors.New("Invalid transfer title issue")
	}

	fmt.Println("Getting State on title " + tr.VIN)
	cpBytes, err := stub.GetState(vehiclePrefix + tr.VIN)
	if err != nil {
		fmt.Println("VIN not found")
		return nil, errors.New("VIN not found " + tr.VIN)
	}

	var cp CP
	fmt.Println("Unmarshalling Title " + tr.VIN)
	err = json.Unmarshal(cpBytes, &cp)
	if err != nil {
		fmt.Println("Error unmarshalling title" + tr.VIN)
		return nil, errors.New("Error unmarshalling title " + tr.VIN)
	}

	var fromCompany Account
	fmt.Println("Getting State on fromOwner " + tr.Owner)
	fromCompanyBytes, err := stub.GetState(accountPrefix + tr.Owner)
	if err != nil {
		fmt.Println("Account not found " + tr.Owner)
		return nil, errors.New("Account not found " + tr.Owner)
	}

	fmt.Println("Unmarshalling FromOwner")
	err = json.Unmarshal(fromCompanyBytes, &fromCompany)
	if err != nil {
		fmt.Println("Error unmarshalling account " + tr.Owner)
		return nil, errors.New("Error unmarshalling account " + tr.Owner)
	}


	// If fromCompany doesn't own this paper
	if cp.Owner != tr.Owner {
		fmt.Println("The owner " + tr.Owner+ " doesn't own this title/vehicle")
		return nil, errors.New("The owner " + tr.Owner+ "doesn't own this title/vehicle")
	} else {
		fmt.Println("The FromOwner does own this title/vehicle")
	}


        cp.Owner = ""
        cp.State="INACTIVE"
	// cp
	cpBytesToWrite, err := json.Marshal(&cp)
	if err != nil {
		fmt.Println("Error marshalling the asset")
		return nil, errors.New("Error marshalling the asset")
	}
	fmt.Println("Put state on Vehicle Asset")
	err = stub.PutState(vehiclePrefix+tr.VIN, cpBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the asset back")
		return nil, errors.New("Error writing the asset back")
	}

	fmt.Println("Successfully completed Invoke")
	return nil, nil
}
/*
Asset: title  owner=toOwner   
fee deducted from FromOwner and sent to government
Account AssetId - Remove from fromOwner & add to toOwner
Checks: fromOwner is current owner and has enough balance & toOwner exists
ToDo: remove item from original owner AssetIds array
*/
func (t *SimpleChaincode) transferTitle(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Transferring Title")
        var amountToBeTransferred float64 = 30

	//need one arg
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting transfer title record")
	}

	var tr TransferTitleTx

	fmt.Println("Unmarshalling Transaction")
	err := json.Unmarshal([]byte(args[0]), &tr)
	if err != nil {
		fmt.Println("Error Unmarshalling Transaction")
		return nil, errors.New("Invalid transfer title issue")
	}

	fmt.Println("Getting State on title " + tr.VIN)
	cpBytes, err := stub.GetState(vehiclePrefix + tr.VIN)
	if err != nil {
		fmt.Println("VIN not found")
		return nil, errors.New("VIN not found " + tr.VIN)
	}

	var cp CP
	fmt.Println("Unmarshalling Title " + tr.VIN)
	err = json.Unmarshal(cpBytes, &cp)
	if err != nil {
		fmt.Println("Error unmarshalling title" + tr.VIN)
		return nil, errors.New("Error unmarshalling title " + tr.VIN)
	}

	var fromCompany Account
	fmt.Println("Getting State on fromOwner " + tr.FromOwner)
	fromCompanyBytes, err := stub.GetState(accountPrefix + tr.FromOwner)
	if err != nil {
		fmt.Println("Account not found " + tr.FromOwner)
		return nil, errors.New("Account not found " + tr.FromOwner)
	}

	fmt.Println("Unmarshalling FromOwner")
	err = json.Unmarshal(fromCompanyBytes, &fromCompany)
	if err != nil {
		fmt.Println("Error unmarshalling account " + tr.FromOwner)
		return nil, errors.New("Error unmarshalling account " + tr.FromOwner)
	}

	var toCompany Account
	fmt.Println("Getting State on ToOwner " + tr.ToOwner)
	toCompanyBytes, err := stub.GetState(accountPrefix + tr.ToOwner)
	if err != nil {
		fmt.Println("Account not found " + tr.ToOwner)
		return nil, errors.New("Account not found " + tr.ToOwner)
	}

	fmt.Println("Unmarshalling toOwner")
	err = json.Unmarshal(toCompanyBytes, &toCompany)
	if err != nil {
		fmt.Println("Error unmarshalling account " + tr.ToOwner)
		return nil, errors.New("Error unmarshalling account " + tr.ToOwner)
	}

	// If fromCompany doesn't own this paper
	if cp.Owner != tr.FromOwner {
		fmt.Println("The owner " + tr.FromOwner+ " doesn't own this title/vehicle")
		return nil, errors.New("The owner " + tr.FromOwner+ "doesn't own this title/vehicle")
	} else {
		fmt.Println("The FromOwner does own this title/vehicle")
	}


	if fromCompany.CashBalance < amountToBeTransferred {
		fmt.Println("The owner " + tr.FromOwner+ "doesn't have enough cash to transfer the title")
		return nil, errors.New("The owner " + tr.FromOwner + "doesn't have enough cash to transfer the title")
	} else {
		fmt.Println("The From Owner has enough money to transfer title")
	}

        var govaccount Account
        govaccountBytes, err := stub.GetState(accountPrefix + "government")
        if err != nil {
                fmt.Println("Error Getting state of - " + accountPrefix + "government")
                return nil, errors.New("Error retrieving account acct:government ")
        }
        err = json.Unmarshal(govaccountBytes, &govaccount)
        if err != nil {
                fmt.Println("Error Unmarshalling govaccountBytes")
                return nil, errors.New("Error retrieving account acct:government")
        }

	fromCompany.CashBalance -= amountToBeTransferred
	govaccount.CashBalance += amountToBeTransferred

        cp.Owner = tr.ToOwner
        //adjust assetIds
        toCompany.AssetsIds = append(toCompany.AssetsIds, tr.VIN)
        //TODO fromCompany.AssetsIds = remove(fromCompany.AssetsIds, tr.VIN)
       
	// Write everything back
	// To Company
	toCompanyBytesToWrite, err := json.Marshal(&toCompany)
	if err != nil {
		fmt.Println("Error marshalling the toCompany")
		return nil, errors.New("Error marshalling the toCompany")
	}
	fmt.Println("Put state on toOwner")
	err = stub.PutState(accountPrefix+tr.ToOwner, toCompanyBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the toCompany back")
		return nil, errors.New("Error writing the toCompany back")
	}

	// From company
	fromCompanyBytesToWrite, err := json.Marshal(&fromCompany)
	if err != nil {
		fmt.Println("Error marshalling the fromCompany")
		return nil, errors.New("Error marshalling the fromCompany")
	}
	fmt.Println("Put state on fromOwner")
	err = stub.PutState(accountPrefix+tr.FromOwner, fromCompanyBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the fromCompany back")
		return nil, errors.New("Error writing the fromCompany back")
	}
        //government
	govaccountBytesToWrite, err := json.Marshal(&govaccount)
                if err != nil {
                        fmt.Println("Error marshalling govt account")
                        return nil, errors.New("Error issuing driver license")
                }
                err = stub.PutState(accountPrefix+"government", govaccountBytesToWrite)
                if err != nil {
                        fmt.Println("Error putting state on govaccountBytesToWrite")
                        return nil, errors.New("Error issuing driver license")
                }

	// cp
	cpBytesToWrite, err := json.Marshal(&cp)
	if err != nil {
		fmt.Println("Error marshalling the cp")
		return nil, errors.New("Error marshalling the cp")
	}
	fmt.Println("Put state on CP")
	err = stub.PutState(vehiclePrefix+tr.VIN, cpBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the cp back")
		return nil, errors.New("Error writing the cp back")
	}

	fmt.Println("Successfully completed Invoke")
	return nil, nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Query running. Function: " + function)

	if function == "GetAllTitles" {
		fmt.Println("Getting all titles")
		allCPs, err := GetAllTitles(stub)
		if err != nil {
			fmt.Println("Error from getall titles")
			return nil, err
		} else {
			allCPsBytes, err1 := json.Marshal(&allCPs)
			if err1 != nil {
				fmt.Println("Error marshalling allcps")
				return nil, err1
			}
			fmt.Println("All success, returning all titles")
			return allCPsBytes, nil
		}
	} else if function == "GetAllDriverLicenses" {
		fmt.Println("Getting all Driver Licenses")
		allCPs, err := GetAllDriverLicenses(stub)
		if err != nil {
			fmt.Println("Error from getall licenses")
			return nil, err
		} else {
			allCPsBytes, err1 := json.Marshal(&allCPs)
			if err1 != nil {
				fmt.Println("Error marshalling all licenses")
				return nil, err1
			}
			fmt.Println("All success, returning all licenses")
			return allCPsBytes, nil
		}
	} else if function == "GetAllVehicleRegistrations" {
		fmt.Println("Getting all Vehicle Registrations")
		allCPs, err := GetAllVehicleRegistrations(stub)
		if err != nil {
			fmt.Println("Error from getall registrations")
			return nil, err
		} else {
			allCPsBytes, err1 := json.Marshal(&allCPs)
			if err1 != nil {
				fmt.Println("Error marshalling all registrations")
				return nil, err1
			}
			fmt.Println("All success, returning all registrations")
			return allCPsBytes, nil
		}
        } else if function == "GetAllTolls" {
                fmt.Println("Getting all Tolls")
                allCPs, err := GetAllTolls(stub)
                if err != nil {
                        fmt.Println("Error from getall tolls")
                        return nil, err
                } else {
                        allCPsBytes, err1 := json.Marshal(&allCPs)
                        if err1 != nil {
                                fmt.Println("Error marshalling all tolls")
                                return nil, err1
                        }
                        fmt.Println("All success, returning all tolls")
                        return allCPsBytes, nil
                }
         } else if function == "GetAllViolations" {
                fmt.Println("Getting all Violations")
                allCPs, err := GetAllViolations(stub)
                if err != nil {
                        fmt.Println("Error from getall violations")
                        return nil, err
                } else {
                        allCPsBytes, err1 := json.Marshal(&allCPs)
                        if err1 != nil {
                                fmt.Println("Error marshalling all violations")
                                return nil, err1
                        }
                        fmt.Println("All success, returning all violations")
                        return allCPsBytes, nil
                }
	} else if function == "GetCP" {
		fmt.Println("Getting particular cp")
		cp, err := GetCP(args[0], stub)
		if err != nil {
			fmt.Println("Error Getting particular cp")
			return nil, err
		} else {
			cpBytes, err1 := json.Marshal(&cp)
			if err1 != nil {
				fmt.Println("Error marshalling the cp")
				return nil, err1
			}
			fmt.Println("All success, returning the cp")
			return cpBytes, nil
		}
	} else if function == "GetCompany" {
		fmt.Println("Getting the company")
		company, err := GetCompany(args[0], stub)
		if err != nil {
			fmt.Println("Error from getCompany")
			return nil, err
		} else {
			companyBytes, err1 := json.Marshal(&company)
			if err1 != nil {
				fmt.Println("Error marshalling the company")
				return nil, err1
			}
			fmt.Println("All success, returning the company")
			return companyBytes, nil
		}
	} else {
		fmt.Println("Generic Query call")
		bytes, err := stub.GetState(args[0])

		if err != nil {
			fmt.Println("Some error happenend: " + err.Error())
			return nil, err
		}

		fmt.Println("All success, returning from generic")
		return bytes, nil
	}
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Invoke running. Function: " + function)

	if function == "issueVehicleTitle" {
		return t.issueVehicleTitle(stub, args)
	} else if function == "issueDriverLicense" {
		return t.issueDriverLicense(stub, args)
	} else if function == "issueVehicleRegistration" {
		return t.issueVehicleRegistration(stub, args)
	} else if function == "transferTitle" {
		return t.transferTitle(stub, args)
	} else if function == "terminateAsset" {
		return t.terminateAsset(stub, args)
	} else if function == "renewLicense" {
		return t.renewLicense(stub, args)
        } else if function == "issueTrafficViolation" {
                return t.issueTrafficViolation(stub, args)
        } else if function == "issueTollTicket" {
                return t.issueTollTicket(stub, args)
	} else if function == "renewRegistration" {
		return t.renewRegistration(stub, args)
	} else if function == "createAccounts" {
		return t.createAccounts(stub, args)
	} else if function == "createAccount" {
		return t.createAccount(stub, args)
	}

	return nil, errors.New("Received unknown function invocation: " + function)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Println("Error starting Simple chaincode: %s", err)
	}
}

//lookup tables for last two digits of CUSIP
var seventhDigit = map[int]string{
	1:  "A",
	2:  "B",
	3:  "C",
	4:  "D",
	5:  "E",
	6:  "F",
	7:  "G",
	8:  "H",
	9:  "J",
	10: "K",
	11: "L",
	12: "M",
	13: "N",
	14: "P",
	15: "Q",
	16: "R",
	17: "S",
	18: "T",
	19: "U",
	20: "V",
	21: "W",
	22: "X",
	23: "Y",
	24: "Z",
}

var eigthDigit = map[int]string{
	1:  "1",
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "A",
	11: "B",
	12: "C",
	13: "D",
	14: "E",
	15: "F",
	16: "G",
	17: "H",
	18: "J",
	19: "K",
	20: "L",
	21: "M",
	22: "N",
	23: "P",
	24: "Q",
	25: "R",
	26: "S",
	27: "T",
	28: "U",
	29: "V",
	30: "W",
	31: "X",
}
