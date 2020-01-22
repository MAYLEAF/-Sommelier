package thread

import (
	"bytes"
	"github.com/MAYLEAF/Sommelier/json"
	"github.com/MAYLEAF/Sommelier/logger"
	"reflect"
)

type context struct {
	is_finish_throw bool
	ping_count      int
	actions         map[string]interface{}

	msg json.Json
}

var Actions map[string]interface{}
var GameCount = 0

func (ctx *context) Initialize(thread *Handler) {
	ctx.actions = Actions
	ctx.actions["uid"] = thread.id
	ctx.msg = json.Json{}
	ctx.msg.SetJson(Actions)
	ctx.ping_count = 100
	ctx.is_finish_throw = false
}

func (ctx *context) react(user *Handler) {
	threadwriter := writer{}
	threadreader := reader{}

	chjson := make(chan []byte)
	threadwriter.write(user, json.Read(ctx.msg.Load("C_LOGIN_REQ")))

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
				chjson <- json.Read(ctx.msg.Load("C_READY_TO_START"))
				continue
			}
			if jsonObj.Contains("other_uid", "AI_") {
				chjson <- json.Read(ctx.msg.Load("C_FINISH_GAME"))
				continue
			}
			if jsonObj.Load("_pcode") == "S_GAME_CREATED" && !jsonObj.Contains("other_uid", "AI_") {
				chjson <- json.Read(ctx.msg.Load("C_LOADING_COMPLETE"))
				continue
			}

			if jsonObj.Load("_pcode") == "S_GAME_START" {
				if jsonObj.Load("hostUid") == user.value[0] {
					GameCount++
					chjson <- json.Read(ctx.msg.Load("C_GAME_DATA"))
				}
				continue
			}
			if jsonObj.Load("_pcode") == "C_GAME_DATA" && !ctx.is_finish_throw {
				if jsonObj.Load("uid") == user.id {
					ctx.ping_count--
				}
				if jsonObj.Load("uid") == user.id && ctx.ping_count <= 0 {
					chjson <- json.Read(ctx.msg.Load("C_FINISH_GAME"))
					ctx.is_finish_throw = true
				}
				if jsonObj.Load("uid") == user.id {
					continue
				}
				chjson <- json.Read(ctx.msg.Load("C_GAME_DATA"))
				continue
			}

			if jsonObj.Load("_pcode") == "S_GAME_RESULT" {
				chjson <- json.Read(ctx.msg.Load("C_BACK_TO_LOBBY"))
				return
			}

			if jsonObj.Load("_pcode") == "S_DISCONNECT_RES" && !ctx.is_finish_throw {
				chjson <- json.Read(ctx.msg.Load("C_FINISH_GAME"))
				ctx.is_finish_throw = true
				continue
			}

			if jsonObj.Load("_pcode") == "S_MATCHING_FAIL" {
				chjson <- json.Read(ctx.msg.Load("C_BACK_TO_LOBBY"))
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
				if reflect.DeepEqual(msg, json.Read(ctx.msg.Load("C_BACK_TO_LOBBY"))) {
					user.Schedule.Done()
				}
			}
		}
	}()
	user.Schedule.Wait()
}
