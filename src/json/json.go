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
	if err := json.Unmarshal([]byte(msg), e.json); err != nil {
		log.Printf("Fail to create json err: %v \n\n", err)
	}
	return err
}

func (e *Json) Read() (string, err) {
	if msg, err := json.Marshal(e.json); err != nil {
		log.Printf("Fail to read json err:%v \n\n", err)
		return err
	}
	return string(msg)
}

func (e *Json) Update(key string, value interface{}) {
	e.json[key] = value
}

//TODO update,create, delete, read
