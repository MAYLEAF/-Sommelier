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
	Sommelier "github.com/MAYLEAF/Sommelier/lib"
	"github.com/MAYLEAF/Sommelier/logger"
)

func main() {
	request := flag.String("request", "request.json", "<file_name>.json")
	value := flag.String("value", "value.csv", "<file_name>.csv")
	flag.Parse()

	logger.Logger()
	defer logger.Close()
	Sommelier.NewRedis()

	actions := readJson(*request)
	threadValue := readCsv(*value)

	TestClient := client.Client{}
	TestClient.SetConnection()
	TestClient.CreateThreads(threadValue)

	TestClient.Test(actions)

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
