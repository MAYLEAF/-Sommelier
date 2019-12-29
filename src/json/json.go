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

var logger = logger.Logger()

type Json struct {
	json    map[string]interface{}
	encoder json.Encoder
}

func Decode(r io.Reader, v interface{}) {
	dec := json.NewDecoder(r)
	msg := make(map[string]interface{})
	if err := dec.Decode(&v); err != nil {
		logger.Error("%v", err)
	}
}

func Encode(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		logger.Error("%v", err)
	}
	return err
}

func Read(v interface{}) []byte {
	msg, err := json.Marshal(v)
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

func (e *Json) Load(key string) interface{} {
	if value, ok := e.json[key]; ok {
		return value
	}
	return nil
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
