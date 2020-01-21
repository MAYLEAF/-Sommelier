package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/MAYLEAF/Sommelier/client"
	"github.com/MAYLEAF/Sommelier/logger"
	"github.com/MAYLEAF/Sommelier/thread"
)

var request = flag.String("request", "request.json", "<file_name>.json")
var value = flag.String("value", "value.csv", "<file_name>.csv")

func main() {
	flag.Parse()
	thread.Actions = readJson(*request)
	threadList := readCsv(*value)
	defer logger.Logger().Close()

	TestClient := client.New()
	TestClient.CreateThreads(threadList)

	TestClient.Test()
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
