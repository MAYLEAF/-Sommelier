package main

import (
	"bufio"
	"client"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	//TODO Separate Start Process
	first := flag.String("start", "start.json", "<file_name>.json")
	request := flag.String("request", "request.json", "<file_name>.json")
	finish := flag.String("finish", "finish.json", "<file_name>.json")
	value := flag.String("value", "value.csv", "<file_name>.csv")
	flag.Parse()

	documents := readJson(*request)
	starts := readJson(*first)
	lasts := readJson(*finish)
	rows := readCsv(*value)

	var e = client.Client{}
	e.SetConnection()
	e.CreateThreads(rows)

	for _, start := range starts {
		message, _ := json.Marshal(start)
		e.MakeRequest(string(message))
	}
	for _, document := range documents {
		message, _ := json.Marshal(document)
		e.MakeRequest(string(message))
	}
	for _, last := range lasts {
		message, _ := json.Marshal(last)
		e.MakeRequest(string(message))
	}

}

func readJson(file_name string) map[string]interface{} {
	jsonFile, err := os.Open(file_name)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	var result map[string]interface{}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &result)
	return result
}

func readCsv(file_name string) [][]string {
	file, err := os.Open(file_name)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))
	csv, _ := reader.ReadAll()

	return csv
}
