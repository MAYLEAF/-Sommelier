/*
Package Json implement a simple library for json CRUD.
*/
package json

import (
	"encoding/json"
	"io"
	"logger"
	"regexp"
	"strings"
)

type Json struct {
	json    map[string]interface{}
	encoder json.Encoder
}

func (e *Json) Create(msg string) error {
	logger := logger.Logger()
	e.json = make(map[string]interface{})
	dec := json.NewDecoder(strings.NewReader(msg))
	err := dec.Decode(&e.json)

	if err != nil {
		logger.Error("Json Create Error: %v", err)
	}

	return err
}

func (e *Json) SetEncoder(w io.Writer) {
	e.encoder = *json.NewEncoder(w)
}

func (e *Json) Encode(v interface{}) error {
	err := e.encoder.Encode(v)
	return err
}

func (e *Json) Read() []byte {
	logger := logger.Logger()
	msg, err := json.Marshal(e.json)
	if err != nil {
		logger.Error("Fail to read json err:%v \n\n", err)
		return nil
	}
	return msg
}

func (e *Json) Json() map[string]interface{} {
	return e.json
}

func (e *Json) SetJson(json map[string]interface{}) {
	e.json = json
}

func (e *Json) Has(key string, value interface{}) bool {
	logger := logger.Logger()
	switch {
	case e.json[key] == nil:
		return false
	case e.json[key] == value:
		return true
	case e.json[key] != value:
		return false
	default:
		logger.Error("Json key %v have not value %v", key, value)
	}
	return false
}

func (e *Json) Contains(key string, value string) bool {
	logger := logger.Logger()
	if e.json[key] == nil {
		return false
	}
	re := regexp.MustCompile(`(.*)` + value + `(.*)`)
	msg, err := json.Marshal(e.json[key])
	if err != nil {
		logger.Info("Fail to read json err: %v \n\n", err)
	}
	if re.Find(msg) == nil {
		return false
	}
	return true
}

func (e *Json) Update(key string, value interface{}) {
	e.json[key] = value
}

func (e *Json) Select(key string) *Json {
	logger := logger.Logger()
	msg, err := json.Marshal(e.json[key])
	if err != nil {
		logger.Info("Fail to read json err: %v \n", err)
	}
	dec := json.NewDecoder(strings.NewReader(string(msg)))
	newJson := Json{}
	err = dec.Decode(&newJson.json)

	if err != nil {
		logger.Info("Json Select Error: %v", err)
	}
	return &newJson
}
