package thread

import (
	"json"
)

type game_room struct {
	ping_count   int
	turn_seconds int
	is_finish    bool
	is_host_turn bool
	host_user    *Handler
	guest_user   *Handler
	actions      json.Json
}

var game_rooms map[string]game_room

func (e *game_room) Create() {
	e.ping_count = 5
	e.turn_seconds = 15
	e.is_finish = false
	e.is_host_turn = true
	actions := json.Json{}
	actions.SetJson(Actions)
	e.actions = actions
}

func (e *game_room) AddUser(host_user *Handler, guest_user *Handler) {
	e.host_user = host_user
	e.guest_user = guest_user
}

func (e *game_room) GameStart() {
	writer := writer{}
	writer.write(e.host_user, e.Action("C_GAME_DATA"))
}

func (e *game_room) SendGameData() {
	writer := writer{}
	defer e.Update()

	if e.is_host_turn {
		writer.write(e.host_user, e.Action("C_GAME_DATA"))
		return
	}
	writer.write(e.guest_user, e.Action("C_GAME_DATA"))
	e.ping_count--
	return
}

func (e *game_room) Update() {
	e.turn_seconds = 15
	e.is_host_turn = !e.is_host_turn
}

func (e *game_room) Action(action string) []byte {
	return json.Read(e.actions.Select(action))
}
