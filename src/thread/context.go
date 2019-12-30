package thread

import (
	"bytes"
	"json"
	"logger"
)

type context struct {
	msg              json.Json
	ping_count       int
	turn_seconds     int
	is_finish_throw  bool
	is_request_fired bool
	is_host          bool
	is_my_turn       bool
}

var Actions map[string]interface{}
var Usercount int

func (e *context) Initialize() {
	e.msg = json.Json{}
	e.msg.SetJson(Actions)
	e.ping_count = 5
	e.turn_seconds = 15
	e.is_finish_throw = false
	e.is_request_fired = false
	e.is_host = false
	e.is_my_turn = false
}

func (e *context) react(thread *Handler) {
	threadwriter := writer{}
	threadreader := reader{}

	//TODO get refeat number
	threadwriter.write(thread, json.Read(e.msg.Load("C_LOGIN_REQ")))

	for {
		response := threadreader.read(thread)
		res := make(map[string]interface{})
		jsonObj := json.Json{}
		bytesReader := bytes.NewReader(response)
		if err := json.Decode(bytesReader, res); err != nil {
			jsonObj.SetJson(res)
		}

		if jsonObj.Load("_pcode") == "S_LOGIN_RES" {
			threadwriter.write(thread, json.Read(e.msg.Load("C_READY_TO_START")))
			continue
		}
		if jsonObj.Contains("other_uid", "AI_") {
			threadwriter.write(thread, json.Read(e.msg.Load("C_FINISH_GAME")))
			continue
		}
		if jsonObj.Load("_pcode") == "S_GAME_CREATED" && !jsonObj.Contains("other_uid", "AI_") {
			threadwriter.write(thread, json.Read(e.msg.Load("C_LOADING_COMPLETE")))
			continue
		}
		if jsonObj.Load("_pcode") == "S_GAME_START" {
			if jsonObj.Load("hostUid") == thread.value[0] {
				e.is_host = true
				threadwriter.write(thread, json.Read(e.msg.Load("C_GAME_DATA")))
			}
			continue
		}
		if jsonObj.Load("_pcode") == "C_GAME_DATA" && jsonObj.Load("uid") != thread.value[0] && e.ping_count > 0 {
			err := threadwriter.write(thread, json.Read(e.msg.Load("C_GAME_DATA")))
			if err != nil {
				logger.Info("User count %d", Usercount)
				thread.Schedule.Done()
				return
			}
			e.is_my_turn = true
			e.turn_seconds = 15
			e.ping_count--

		} else if jsonObj.Load("_pcode") == "C_GAME_DATA" && jsonObj.Load("uid") == thread.value[0] && e.ping_count > 0 {
			e.is_request_fired = false
			e.is_my_turn = false

		} else if jsonObj.Load("_pcode") == "S_PING" && e.ping_count > 0 && !e.is_request_fired && e.turn_seconds <= 0 {
			if e.is_my_turn && e.is_host {
				threadwriter.write(thread, json.Read(e.msg.Load("C_GAME_DATA")))
				e.turn_seconds = 15
			}
			if !e.is_my_turn && !e.is_host {
				threadwriter.write(thread, json.Read(e.msg.Load("C_GAME_DATA")))
				e.turn_seconds = 15
			}
			e.is_request_fired = true

		} else if jsonObj.Load("_pcode") == "C_GAME_DATA" && jsonObj.Load("uid") != thread.value[0] && !e.is_finish_throw && e.ping_count <= 0 {
			e.is_finish_throw = true
			e.turn_seconds = 15
			threadwriter.write(thread, json.Read(e.msg.Load("C_FINISH_GAME")))

		} else if jsonObj.Load("_pcode") == "S_GAME_RESULT" {
			threadwriter.write(thread, json.Read(e.msg.Load("C_BACK_TO_LOBBY")))
			Usercount--
			logger.Info("User count %d", Usercount)
			thread.Schedule.Done()
			return
		} else if jsonObj.Load("_pcode") == "S_MATCHING_FAIL" {
			threadwriter.write(thread, json.Read(e.msg.Load("C_BACK_TO_LOBBY")))
			Usercount--
			logger.Info("User count %d", Usercount)
			thread.Schedule.Done()
			return
		} else if jsonObj.Load("_pcode") == "S_PING" {
			e.turn_seconds--
			logger.Info("User "+thread.value[0]+"-turn: %d seconds, ping_count: %v left", e.turn_seconds, e.ping_count)
		} else if e.turn_seconds < 0 {
			e.turn_seconds = 15
		} else {
			continue
		}
	}
}

func (e *context) TurnChange() {
	e.is_my_turn = !e.is_my_turn
	e.turn_seconds = 15
	e.ping_count--
}
