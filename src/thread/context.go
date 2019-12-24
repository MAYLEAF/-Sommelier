package thread

import (
	"json"
	"logger"
)

type context struct {
	actions map[string]interface{}
}

func (e *context) create(actions map[string]interface{}) {
	e.actions = actions
}

var Actions map[string]interface{}
var Usercount int

func (e *context) react(thread *Handler) {
	threadwriter := writer{}
	threadreader := reader{}
	logger := logger.Logger()

	msg := json.Json{}
	game_room = game_room{}
	ping_count := 5
	turn_seconds := 15
	finish_throw := false

	msg.SetJson(Actions)
	//TODO get refeat number
	threadwriter.write(thread, string(msg.Select("C_LOGIN_REQ").Read()))

	for {
		response := threadreader.read(thread)
		res := json.Json{}
		res.Create(response)
		if res.Has("_pcode", "S_LOGIN_RES") {
			threadwriter.write(thread, string(msg.Select("C_READY_TO_START").Read()))
		} else if res.Contains("other_uid", "AI_") {
			threadwriter.write(thread, string(msg.Select("C_FINISH_GAME").Read()))

		} else if res.Has("_pcode", "S_GAME_CREATED") && !res.Contains("other_uid", "AI_") {
			threadwriter.write(thread, string(msg.Select("C_LOADING_COMPLETE").Read()))

		} else if res.Has("_pcode", "S_GAME_START") && res.Has("hostUid", thread.value[0]) {
			threadwriter.write(thread, string(msg.Select("C_GAME_DATA").Read()))
			turn_seconds = 15
			ping_count--

		} else if res.Has("_pcode", "C_GAME_DATA") && !res.Has("uid", thread.value[0]) && ping_count > 0 {
			threadwriter.write(thread, string(msg.Select("C_GAME_DATA").Read()))
			turn_seconds = 15
			ping_count--

		} else if res.Has("_pcode", "S_PING") && ping_count > 0 && turn_seconds <= 0 {
			threadwriter.write(thread, string(msg.Select("C_GAME_DATA").Read()))
			turn_seconds = 15
			ping_count--

		} else if res.Has("_pcode", "C_GAME_DATA") && !res.Has("uid", thread.value[0]) && ping_count <= 0 {
			threadwriter.write(thread, string(msg.Select("C_FINISH_GAME").Read()))
			turn_seconds = 15

		} else if res.Has("_pcode", "S_PING") && ping_count <= 0 && turn_seconds <= 0 && finish_throw == false {
			threadwriter.write(thread, string(msg.Select("C_FINISH_GAME").Read()))
			finish_throw = true
		} else if res.Has("_pcode", "S_GAME_RESULT") {
			threadwriter.write(thread, string(msg.Select("C_BACK_TO_LOBBY").Read()))
			Usercount--
			logger.Info("User count %d", Usercount)
			thread.Schedule.Done()
			return
		} else if res.Has("_pcode", "S_MATCHING_FAIL") {
			threadwriter.write(thread, string(msg.Select("C_BACK_TO_LOBBY").Read()))
			Usercount--
			logger.Info("User count %d", Usercount)
			thread.Schedule.Done()
			return
		} else if res.Has("_pcode", "S_PING") && ping_count > 0 && turn_seconds > 0 {
			turn_seconds--
			logger.Info("User "+thread.value[0]+"-turn: %d seconds, ping_count: %d left", turn_seconds, ping_count)
		} else if res.Has("_pcode", "S_PING") {
			turn_seconds--
			logger.Info("User "+thread.value[0]+"-turn: %d seconds, ping_count: %d left", turn_seconds, ping_count)
		} else {
			continue
		}
	}
}
