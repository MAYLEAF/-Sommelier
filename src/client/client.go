//Package Client implements a controller for TCP communication.
package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logger"
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
	logger := logger.Logger()
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
	logger.Info("SetConnection with %v", result["serverAddress"])

}

func (e *Client) CreateThreads(values [][]string) {
	logger := logger.Logger()
	logger.Info("Create Threads")
	defer logger.Info("Create ThreadsEND")
	wg := &sync.WaitGroup{}
	thread.Usercount = 0
	for _, value := range values {
		thread.Usercount++
		wg.Add(1)
		go func(value []string) {
			thread := thread.Handler{}
			thread.Create(e.serverAddr, value)
			e.threads = append(e.threads, thread)
			logger.Info("Logger: Create thread:", thread)
			wg.Done()
		}(value)
	}
	wg.Wait()
}

func (e *Client) MakeTest(actions map[string]interface{}) {

	logger := logger.Logger()
	logger.Info("MakeTest", e.wg)
	defer logger.Info("MakeTestEnd")

	e.wg.Wait()
	for _, thread := range e.threads {
		e.wg.Add(1)
		go e.test(actions, thread)
	}

	e.wg.Wait()
}

func (e *Client) test(actions map[string]interface{}, thread thread.Handler) {
	logger := logger.Logger()
	logger.Info("Test A Thread;  Handler=%v", thread)
	defer logger.Info("TestEnd A Thread; Handler=%v\n\n", thread)

	go thread.RequestMaker(actions)
	thread.Send = make(chan string, 10)
	thread.Schedule.Add(1)
	thread.Schedule.Wait()
	e.wg.Done()
}

func (e *Client) MakeRequest(Message string) {
	for _, thread := range e.threads {
		time.Sleep(1000 * time.Millisecond)
		thread.Write(Message)
	}
}
