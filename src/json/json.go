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

func (e *Json) Update(key string, value interface{}) {
	e.json[key] = value
}

//TODO update,create, delete, read
