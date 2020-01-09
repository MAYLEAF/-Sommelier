package thread

import (
	"bytes"
	"github.com/MAYLEAF/Sommelier/json"
	"github.com/MAYLEAF/Sommelier/logger"
	"reflect"
)

type context struct {
	msg              json.Json
	ping_count       int
	turn_seconds     int
	is_finish_throw  bool
	is_request_fired bool
	is_host          bool
	is_my_turn       bool
	thread           *Handler
}

var Actions map[string]interface{}
var Usercount int

func (e *context) Initialize(thread *Handler) {
	e.msg = json.Json{}
	e.msg.SetJson(Actions)
	e.ping_count = 5
	e.turn_seconds = 15
	e.is_finish_throw = false
	e.is_request_fired = false
	e.is_host = false
	e.is_my_turn = false
	e.thread = thread
}

func (e *context) react(thread *Handler) {
	threadwriter := writer{}
	threadreader := reader{}

	chjson := make(chan []byte)
	threadwriter.write(thread, json.Read(e.msg.Load("C_LOGIN_REQ")))

	go func() {
		for {
			response := threadreader.read(thread)
			res := make(map[string]interface{})
			jsonObj := json.Json{}
			bytesReader := bytes.NewReader(response)
			if err := json.Decode(bytesReader, res); err != nil {
				logger.Info("%v", err)
			}
			jsonObj.SetJson(res)

			if jsonObj.Load("_pcode") == "S_LOGIN_RES" {
				chjson <- json.Read(e.msg.Load("C_READY_TO_START"))
				continue
			}
			if jsonObj.Contains("other_uid", "AI_") {
				chjson <- json.Read(e.msg.Load("C_FINISH_GAME"))
				continue
			}
			if jsonObj.Load("_pcode") == "S_GAME_CREATED" && !jsonObj.Contains("other_uid", "AI_") {
				chjson <- json.Read(e.msg.Load("C_LOADING_COMPLETE"))
				continue
			}

			if jsonObj.Load("_pcode") == "S_GAME_START" {
				if jsonObj.Load("hostUid") == thread.value[0] {
					e.is_host = true
					chjson <- json.Read(e.msg.Load("C_GAME_DATA"))
				}
				continue
			}
			if jsonObj.Load("_pcode") == "C_GAME_DATA" && !e.is_finish_throw {
				if jsonObj.Load("uid") == thread.value[0] {
					e.ping_count--
				}
				if jsonObj.Load("uid") == thread.value[0] && e.ping_count <= 0 {
					chjson <- json.Read(e.msg.Load("C_FINISH_GAME"))
					e.is_finish_throw = true
				}
				if jsonObj.Load("uid") == thread.value[0] {
					continue
				}
				chjson <- json.Read(e.msg.Load("C_GAME_DATA"))
				continue
			}

			if jsonObj.Load("_pcode") == "S_GAME_RESULT" {
				chjson <- json.Read(e.msg.Load("C_BACK_TO_LOBBY"))
				return
			}

			if jsonObj.Load("_pcode") == "S_DISCONNECT_RES" && !e.is_finish_throw {
				chjson <- json.Read(e.msg.Load("C_FINISH_GAME"))
				e.is_finish_throw = true
				continue
			}

			if jsonObj.Load("_pcode") == "S_MATCHING_FAIL" {
				chjson <- json.Read(e.msg.Load("C_BACK_TO_LOBBY"))
				thread.Schedule.Done()
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case msg := <-chjson:
				if err := threadwriter.write(thread, msg); err != nil {
					Usercount--
					logger.Info("User count %d", Usercount)
					thread.Schedule.Done()
					return
				}
				if reflect.DeepEqual(msg, json.Read(e.msg.Load("C_BACK_TO_LOBBY"))) {
					Usercount--
					logger.Info("User count %d", Usercount)
					thread.Schedule.Done()
				}
			}
		}
	}()

	thread.Schedule.Wait()
}

func (e *context) TurnChange() {
	e.is_my_turn = !e.is_my_turn
	e.turn_seconds = 15
	e.ping_count--
}

func (e *context) Send(pcode string) error {
	writer := writer{}
	if err := writer.write(e.thread, json.Read(e.msg.Load(pcode))); err != nil {
		logger.Info("Thread context send fail User - %v, pcode: %v", e.thread.value[0], pcode)
		return err
	}
	return nil
}
