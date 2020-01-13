//Package Client implements a controller for TCP communication.
package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/MAYLEAF/Sommelier/logger"
	"github.com/MAYLEAF/Sommelier/thread"
)

type Client struct {
	threads    []thread.Handler
	serverAddr string
	wg         sync.WaitGroup
	protocol   string
	Err        error
}

var Threadcount int

func (e *Client) SetConnection() {
	jsonFile, err := os.Open("connect.json")
	if err != nil {
		logger.Info("SetConnection err:%v", err)
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
	logger.Info("Create Threads")
	defer logger.Info("Create ThreadsEND")
	wg := &sync.WaitGroup{}
	Threadcount = 0
	for _, value := range values {
		Threadcount++
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

func (e *Client) Test(actions map[string]interface{}) {
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
	logger.Info("Test A Thread;  Handler=%v", thread)
	defer logger.Info("TestEnd A Thread; Handler=%v\n\n", thread)

	thread.Schedule.Add(1)
	go thread.Play(actions)

	thread.Schedule.Wait()
	Threadcount--
	logger.Info("usercount %d", Threadcount)

	e.wg.Done()
}
