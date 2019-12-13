package thread

import (
	"json"
	"time"
)

type context struct {
	actions map[string]interface{}
}

func (e *context) create(actions map[string]interface{}) {
	e.actions = actions
}
