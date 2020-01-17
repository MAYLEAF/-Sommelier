package thread

import (
	"bytes"
	"github.com/MAYLEAF/Sommelier/json"
	"github.com/MAYLEAF/Sommelier/logger"
	"reflect"
)

type context struct {
	msg             json.Json
	ping_count      int
	is_finish_throw bool
	is_host         bool
}

var Actions map[string]interface{}
var GameCount = 0

func (e *context) Initialize(thread *Handler) {
	e.msg = json.Json{}
	e.msg.SetJson(Actions)
	e.ping_count = 100
	e.is_finish_throw = false
	e.is_host = false
}

func (e *context) react(user *Handler) {
	threadwriter := writer{}
	threadreader := reader{}

	chjson := make(chan []byte)
	threadwriter.write(user, json.Read(e.msg.Load("C_LOGIN_REQ")))

	go func() {
		for {
			response := threadreader.read(user)
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
				if jsonObj.Load("hostUid") == user.value[0] {
					e.is_host = true
					GameCount++
					chjson <- json.Read(e.msg.Load("C_GAME_DATA"))
				}
				continue
			}
			if jsonObj.Load("_pcode") == "C_GAME_DATA" && !e.is_finish_throw {
				if jsonObj.Load("uid") == user.id {
					e.ping_count--
				}
				if jsonObj.Load("uid") == user.id && e.ping_count <= 0 {
					chjson <- json.Read(e.msg.Load("C_FINISH_GAME"))
					e.is_finish_throw = true
				}
				if jsonObj.Load("uid") == user.id {
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
				user.Schedule.Done()
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case msg := <-chjson:
				if err := threadwriter.write(user, msg); err != nil {
					user.Schedule.Done()
					return
				}
				if reflect.DeepEqual(msg, json.Read(e.msg.Load("C_BACK_TO_LOBBY"))) {
					user.Schedule.Done()
				}
			}
		}
	}()
	user.Schedule.Wait()
}
