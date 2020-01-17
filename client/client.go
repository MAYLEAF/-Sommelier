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

type ConnInfo struct {
	RAddr string `json:"serverAddress"`
	Proto string `json:"protocol"`
}

var Threadcount = 0

func New() *Client {
	newClient := Client{}
	newClient.rAddr, newClient.proto = connInfo()
	return &newClient
}

func connInfo() (string, string) {
	json := &ConnInfo{}
	json = jsonLib.ReadJsonFile("connect.json", json).(*ConnInfo)
	rAddr := json.RAddr
	proto := json.Proto
	logger.Info("Conn Info: %v", json)
	return rAddr, proto
}

func (e *Client) CreateThreads(values [][]string) {
	logger.Info("Create Threads")
	defer logger.Info("Create ThreadsEND")

	wg := &sync.WaitGroup{}
	for _, value := range values {
		Threadcount++
		wg.Add(1)
		go func(value []string) {
			thread := thread.Handler{}
			thread.Create(e.rAddr, value)
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
	go thread.Attack(actions)

	thread.Schedule.Wait()
	Threadcount--
	logger.Info("usercount %d", Threadcount)

	e.wg.Done()
}
