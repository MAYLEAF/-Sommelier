package client

import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"log"
)

type Client struct {
	threads []handler
	serverAddr string
	protocol string
	Err    error
}

func (e* Client) SetConnection(){
	jsonFile, err := os.Open("connect.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	e.serverAddr = fmt.Sprintf("%v", result["serverAddress"])
	e.protocol = fmt.Sprintf("%v", result["protocol"])
	log.Print("Connection set")
	log.Print(e)
}

func (e *Client) CreateThreads() {
	thread := handler{}
	thread.Create(e.serverAddr,e.protocol)
	e.threads = append(e.threads, thread)
}

func (e* Client) MakeRequest(Message string) {
	for _, thread := range e.threads {
		thread.MakeRequest(Message)
	}
}

func (e* Client) MakeCommunication(Messages []string) {
	for _, Message := range Messages {
		e.MakeRequest(Message)
	}
}

