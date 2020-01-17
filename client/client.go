//Package Client implements a controller for TCP communication.
package client

import (
	"sync"

	jsonLib "github.com/MAYLEAF/Sommelier/json"
	"github.com/MAYLEAF/Sommelier/logger"
	"github.com/MAYLEAF/Sommelier/thread"
)

type Client struct {
	threads     []thread.Handler
	threadCount int
	rAddr       string
	waitGroup   sync.WaitGroup
	proto       string
	Err         error
}

type connInfo struct {
	RAddr string `json:"serverAddress"`
	Proto string `json:"protocol"`
}

var Threadcount = 0

func New() *Client {
	newClient := Client{}
	newClient.rAddr, newClient.proto = ConnInfo()
	return &newClient
}

func ConnInfo() (string, string) {
	json := &connInfo{}
	json = jsonLib.ReadJsonFile("connect.json", json).(*connInfo)
	return json.RAddr, json.Proto
}

func (Test *Client) CreateThreads(values [][]string) {
	logger.Info("Create Threads")
	for _, value := range values {
		Threadcount++
		Test.waitGroup.Add(1)
		go func(value []string) {
			newThread := thread.New(Test.rAddr, value)
			Test.threads = append(Test.threads, *newThread)
			Test.waitGroup.Done()
		}(value)
	}
	Test.waitGroup.Wait()
}

func (Test *Client) Test(actions map[string]interface{}) {
	logger.Info("MakeTest", nil)

	for _, thread := range Test.threads {
		Test.waitGroup.Add(1)
		go Test.test(actions, thread)
	}

	Test.waitGroup.Wait()
}

func (e *Client) test(actions map[string]interface{}, thread thread.Handler) {
	logger.Info("Test A Thread;  Handler=%v", thread)
	defer logger.Info("TestEnd A Thread; Handler=%v\n\n", thread)

	thread.Schedule.Add(1)
	go thread.Attack(actions)

	thread.Schedule.Wait()
	Threadcount--
	logger.Info("usercount %d", Threadcount)

	e.waitGroup.Done()
}
