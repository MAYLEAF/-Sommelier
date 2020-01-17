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
	ThreadCount int
	rAddr       string
	wg          sync.WaitGroup
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

func (e *Client) CreateThreads(values [][]string) {
	logger.Info("Create Threads")
	defer logger.Info("Create ThreadsEND")

	wg := &sync.WaitGroup{}
	for _, value := range values {
		Threadcount++
		wg.Add(1)
		go func(value []string) {
			newThread := thread.New(e.rAddr, value)
			e.threads = append(e.threads, *newThread)
			logger.Info("Logger: Create thread: %v", newThread)
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
	go thread.Attack(actions)

	thread.Schedule.Wait()
	Threadcount--
	logger.Info("usercount %d", Threadcount)

	e.wg.Done()
}
