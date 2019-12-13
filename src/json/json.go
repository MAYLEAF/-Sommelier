/*
Package Json implement a simple library for json CRUD.
*/
package json

import (
	"encoding/json"
	"log"
	"strings"
)

type Json struct {
	json map[string]interface{}
}

func (e *Json) Create(msg string) error {
	e.json = make(map[string]interface{})
	dec := json.NewDecoder(strings.NewReader(msg))
	err := dec.Decode(&e.json)

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (e *Json) Read() ([]byte, error) {
	msg, err := json.Marshal(e.json)
	if err != nil {
		log.Printf("Fail to read json err:%v \n\n", err)
	}
	return msg, err
}

func (e *Json) Have(key string, value interface{}) bool {
	switch {
	case e.json[key] == nil:
		log.Printf("Json key %v is empty", key)
	case e.json[key] == value:
		return true
	case e.json[key] != value:
		return false
	default:
		log.Print("Unknown error")
	}
	return false
}

func (e *Json) Has(key string, value interface{}) bool {
	switch {
	case e.json[key] == nil:
		log.Printf("Json key %v is empty", key)
	case e.json[key] == value:
		return true
	case e.json[key] != value:
		log.Printf("Json key $v have not value %v", key, value)
	default:
		log.Print("Unknown error")
	}
	return false
}

func (e *Json) Update(key string, value interface{}) {
	e.json[key] = value
}
