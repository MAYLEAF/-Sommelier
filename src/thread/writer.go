package thread

import (
	"io"
	"json"
	"log"
	"os"
	"sync"
)

type writer struct {
	lock sync.Mutex
}

func (e *writer) write(thread *Handler, message string) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	fpLog, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fpLog.Close()
	multiWriter := io.MultiWriter(fpLog, os.Stdout)
	log.SetOutput(multiWriter)

	msg := json.Json{}
	msg.Create(message)
	msg.Update("uid", thread.value[0])
	content := msg.Read()

	if _, err := thread.conn.Write(content); nil != err {
		log.Printf("failed to write err: %v", err)
		return err
	}
	log.Print("Request Message:" + string(content))

	return nil
}
