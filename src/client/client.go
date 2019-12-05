package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type Client struct {
	threads    []handler
	serverAddr string
	wg         sync.WaitGroup
	protocol   string
	Err        error
}

//TODO Separate Each Client By Protocol
func (e *Client) SetConnection() {
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
	log.Print(e)
}

func (e *Client) CreateThreads(values [][]string) {
	for _, value := range values {
		e.wg.Add(1)
		thread := handler{}
		thread.Create(e.serverAddr, value)
		e.threads = append(e.threads, thread)
		log.Print(thread)
		e.wg.Done()
	}
	e.wg.Wait()
}
func (e *Client) MakeTest(messages []string) {
	for _, thread := range e.threads {
		e.wg.Add(1)
		go thread.test(messages, e.wg)
	}
	e.wg.Wait()
}

func (e *Client) MakeRequest(Message string) {
	for _, thread := range e.threads {
		e.wg.Add(1)
		ch := make(chan string, 10)
		time.Sleep(1 * time.Second)
		thread.MakeRequest(Message, ch)
		select {
		case send := <-ch:
			log.Print(send)
			e.wg.Done()
		default:
			log.Print("default")
		}
	}
	e.wg.Wait()
}

func (e *Client) MakeCommunication(Messages []string) {
	for _, Message := range Messages {
		e.MakeRequest(Message)
	}
}
