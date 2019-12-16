package thread

import (
	"fmt"
	"json"
	"time"
)

type context struct {
	actions map[string]interface{}
}

func (e *context) create(actions map[string]interface{}) {
	e.actions = actions
}

func (e *context) react(thread *Handler) {
	threadwriter := writer{}
	threadreader := reader{}

	msg := json.Json{}

	msg.SetJson(e.actions)

	for {
		time.Sleep(1000 * time.Millisecond)
		response := threadreader.read(&e.thread)
		res := json.Json{}
		res.Create(response)
		if res.Has("key", "something") {
			fmt.Print("Do something")
		}

	}
}
