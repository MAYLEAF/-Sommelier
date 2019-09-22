package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"client"
)

func main() {
	//TODO READ JSON FILE AND TAKE REQUEST JSON INTO STRING
	documents := readRequest()

	var e = client.Client{}
	e.SetConnection()
	e.CreateThreads()

	for _, document := range documents{
		message, _ := json.Marshal(document)
		e.MakeRequest(string(message))
	}
}

func readRequest() map[string]interface{}{
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
