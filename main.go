package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Account string

type Transaction struct {
	From 	Account `json: "from"`
	To 		Account `json: "to"`
	Value 	uint 	`json: "value`
	Data 	string 	`json: "data"`
}

type State struct {
	Balances map[Account]uint
	dbFile *os.File
}

// func loadToStruct(path string) (error) {
// 	jsonFile, err := os.Open("data/genesis.json")

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	bVal, _ := ioutil.ReadAll(jsonFile)

// } 

func loadToMap(path string) (map[string]interface{}, error) {
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Sprintf("Found %s", path)

	bVal, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
    json.Unmarshal([]byte(bVal), &result)

	return result, nil
}

func stateFromDisk() (*State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dbFilePath := filepath.Join(cwd, "data", "genesis.json")
	// Load genesis database
	
	if err != nil {
		return nil, err
	}

	var genFile map[string]interface{}
	genFile, err = loadToMap(dbFilePath)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(genFile["gen_date"])

	return nil, nil
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	dbFilePath := filepath.Join(cwd, "data", "genesis.json")
	// Load genesis database
	
	if err != nil {
		fmt.Println(err)
	}

	var genFile map[string]interface{}
	genFile, err = loadToMap(dbFilePath)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(genFile["balances"].(map[string]interface{})["person"])

}