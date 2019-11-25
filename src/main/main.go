package main

import (
	"fmt"
	"bufio"
	"encoding/json"
	"encoding/csv"
	"io/ioutil"
	"os"
	"client"
)

func main() {
	documents := readJson()
	rows := readCsv()
	//TODO Separate Start Process

	var e = client.Client{}
	e.SetConnection()
	for _, row := range rows{
		e.CreateThreads(row)
	}

	for _, document := range documents{
		message, _ := json.Marshal(document)
		e.MakeRequest(string(message))
	}
}

func readJson() map[string]interface{}{
	jsonFile, err := os.Open("request.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	var result map[string]interface{}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &result)
	return result
}

func readCsv() [][]string{
	file, err := os.Open("value.csv")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))
	csv, _ := reader.ReadAll()

	return csv
}
