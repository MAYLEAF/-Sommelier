package thread

import (
	"json"
	"logger"
	"sync"
)

type context struct {
	actions map[string]interface{}
}

func (e *context) create(actions map[string]interface{}) {
	e.actions = actions
}

var Usercount int

func (e *context) react(thread *Handler) {
	threadwriter := writer{}
	threadreader := reader{}
	logger := logger.Logger()

	msg := json.Json{}
	game_count := 1

	msg.SetJson(e.actions)
	//TODO get refeat number
	threadwriter.write(thread, string(msg.Select("C_LOGIN_REQ").Read()))
	var once sync.Once

	for {
		response := threadreader.read(thread)
		res := json.Json{}
		res.Create(response)
		if res.Has("_pcode", "S_LOGIN_RES") {
			threadwriter.write(thread, string(msg.Select("C_READY_TO_START").Read()))
		} else if res.Contains("other_uid", "AI_") {
			threadwriter.write(thread, string(msg.Select("C_BACK_TO_LOBBY").Read()))
		} else if res.Has("_pcode", "S_GAME_CREATED") {
			threadwriter.write(thread, string(msg.Select("C_LOADING_COMPLETE").Read()))
		} else if res.Has("_pcode", "S_GAME_START") {
			threadwriter.write(thread, string(msg.Select("C_GAME_DATA").Read()))
			game_count--
		} else if res.Has("_pcode", "C_GAME_DATA") && game_count == 0 {
			once.Do(func() {
				threadwriter.write(thread, string(msg.Select("C_GAME_SAVE_DATA").Read()))
			})
		} else if res.Has("_pcode", "S_GAME_SAVE_DATA") {
			threadwriter.write(thread, string(msg.Select("C_FINISH_GAME").Read()))
		} else if res.Has("_pcode", "S_MATCHING_FAIL") {
			threadwriter.write(thread, string(msg.Select("C_BACK_TO_LOBBY").Read()))
			Usercount--
			logger.Info("User count %d", Usercount)
			thread.Schedule.Done()
			return
		} else if res.Has("_pcode", "S_GAME_RESULT") {
			threadwriter.write(thread, string(msg.Select("C_BACK_TO_LOBBY").Read()))
			Usercount--
			logger.Info("User count %d", Usercount)
			thread.Schedule.Done()
			return
		} else if res.Has("_pcode", "S_PING") {
			continue
		} else {
			continue
		}
	}
}
