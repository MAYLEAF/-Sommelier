package thread

import (
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
	logger := logger.Logger()

	//TODO get refeat number
	threadwriter.write(thread, string(e.msg.Select("C_LOGIN_REQ").Read()))

	for {
		response := threadreader.read(thread)
		res := json.Json{}
		res.Create(response)
		if res.Has("_pcode", "S_LOGIN_RES") {
			threadwriter.write(thread, string(e.msg.Select("C_READY_TO_START").Read()))
			continue
		}
		if res.Contains("other_uid", "AI_") {
			threadwriter.write(thread, string(e.msg.Select("C_FINISH_GAME").Read()))
			continue
		}
		if res.Has("_pcode", "S_GAME_CREATED") && !res.Contains("other_uid", "AI_") {
			threadwriter.write(thread, string(e.msg.Select("C_LOADING_COMPLETE").Read()))
			continue
		}
		if res.Has("_pcode", "S_GAME_START") {
			if res.Has("hostUid", thread.value[0]) {
				e.is_host = true
				threadwriter.write(thread, string(e.msg.Select("C_GAME_DATA").Read()))
			}
			continue
		}
		if res.Has("_pcode", "C_GAME_DATA") && !res.Has("uid", thread.value[0]) && e.ping_count > 0 {
			err := threadwriter.write(thread, string(e.msg.Select("C_GAME_DATA").Read()))
			if err != nil {
				logger.Info("User count %d", Usercount)
				thread.Schedule.Done()
				return
			}
			e.is_my_turn = true
			e.turn_seconds = 15
			e.ping_count--

		} else if res.Has("_pcode", "C_GAME_DATA") && res.Has("uid", thread.value[0]) && e.ping_count > 0 {
			e.is_request_fired = false
			e.is_my_turn = false

		} else if res.Has("_pcode", "S_PING") && e.ping_count > 0 && !e.is_request_fired && e.turn_seconds <= 0 {
			if e.is_my_turn && e.is_host {
				threadwriter.write(thread, string(e.msg.Select("C_GAME_DATA").Read()))
				e.turn_seconds = 15
			}
			if !e.is_my_turn && !e.is_host {
				threadwriter.write(thread, string(e.msg.Select("C_GAME_DATA").Read()))
				e.turn_seconds = 15
			}
			e.is_request_fired = true

		} else if res.Has("_pcode", "C_GAME_DATA") && !res.Has("uid", thread.value[0]) && !e.is_finish_throw && e.ping_count <= 0 {
			e.is_finish_throw = true
			e.turn_seconds = 15
			threadwriter.write(thread, string(e.msg.Select("C_FINISH_GAME").Read()))

		} else if res.Has("_pcode", "S_GAME_RESULT") {
			threadwriter.write(thread, string(e.msg.Select("C_BACK_TO_LOBBY").Read()))
			Usercount--
			logger.Info("User count %d", Usercount)
			thread.Schedule.Done()
			return
		} else if res.Has("_pcode", "S_MATCHING_FAIL") {
			threadwriter.write(thread, string(e.msg.Select("C_BACK_TO_LOBBY").Read()))
			Usercount--
			logger.Info("User count %d", Usercount)
			thread.Schedule.Done()
			return
		} else if res.Has("_pcode", "S_PING") {
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
