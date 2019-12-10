/*
Package Json implement a simple library for json CRUD.
*/
package json

import (
	"encoding/json"
)

type Json struct {
	json map[string]interface{}
}

func (e *Json) Create(msg string) error {
	err := json.Unmarshal([]byte(msg), e.json)
	return err
}

//TODO update,create, delete, read
