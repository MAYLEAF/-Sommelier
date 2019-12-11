/*
Package Json implement a simple library for json CRUD.
*/
package json

import (
	"encoding/json"
	"log"
)

type Json struct {
	json map[string]interface{}
}

func (e *Json) Create(msg string) error {
	e.json = make(map[string]interface{})
	if err := json.Unmarshal([]byte(msg), &e.json); err != nil {
		log.Printf("Fail to create json err: %v \n\n", err)
		return err
	}
	return nil
}
func (e *Json) Read() ([]byte, error) {
	if msg, err := json.Marshal(e.json); err != nil {
		log.Printf("Fail to read json err:%v \n\n", err)
		return err
	}
	return msg
}

func (e *Json) Have(key string, value interface{}) bool {
	switch {
	case e.json[key] == nil:
		log.Printf("Json key %v is empty", key)
	case e.json[key] == value:
		return true
	case e.json[key] != value:
		log.Printf("Json key %v does not have value %v", key, value)
	default:
		log.Print("Unknown error")
	}
	return false
}

func (e *Json) Update(key string, value interface{}) {
	e.json[key] = value
}
