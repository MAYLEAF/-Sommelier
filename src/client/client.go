//Package Client implements a controller for TCP communication.
package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"thread"
	"time"
)

type Client struct {
	threads    []thread.Handler
	serverAddr string
	wg         sync.WaitGroup
	protocol   string
	Err        error
}

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
	log.Printf("Logger: SetConnection with %v", result["serverAddress"])

}

func (e *Client) CreateThreads(values [][]string) {
	log.Print("Logger: Create Threads")
	for _, value := range values {
		e.wg.Add(1)
		thread := thread.Handler{}
		thread.Create(e.serverAddr, value)
		e.threads = append(e.threads, thread)
		e.wg.Done()
	}
	e.wg.Wait()
}
func (e *Client) MakeTest(messages []string) {
	log.Print("Logger: MakeTest")
	defer log.Print("Logger: MakeTestEnd")

	for _, thread := range e.threads {
		e.wg.Add(1)
		go e.test(messages, thread)
	}
	e.wg.Wait()
}

func (e *Client) test(messages []string, thread thread.Handler) {
	log.Printf("Logger: Test A Thread;  Handler=%v", thread)
	defer log.Printf("Logger: TestEnd A Thread; Handler=%v\n\n", thread)

	go thread.RequestMaker()
	thread.Send = make(chan string, 10)

	for _, message := range messages {
		thread.Schedule.Add(1)
		thread.Send <- message
	}

	thread.Schedule.Wait()
	e.wg.Done()
}

func (e *Client) MakeRequest(Message string) {
	for _, thread := range e.threads {
		time.Sleep(1000 * time.Millisecond)
		thread.Write(Message)
	}
}
