package thread

import (
	"fmt"
	"io"
	"json"
	"log"
)

type reader struct {
	thread Handler
}

func (e *reader) read(thread *Handler, msg string) {
	if !e.hasResponse(msg) {
		return
	}

	buf := make([]byte, 0, 16384)
	tmp := make([]byte, 256)

	for {
		n, err := thread.conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
				log.Print(err)
			}
			break
		}
		if n != 256 {
			buf = append(buf, tmp[:n]...)
			break
		}
		buf = append(buf, tmp[:n]...)
	}

	fmt.Print("Message from server: " + string(buf) + "\n")

}

func (e *reader) hasResponse(req string) bool {
	msg := json.Json{}
	msg.Create(req)
	if msg.Have("_pcode", "C_SAVE_GAME_DATA") {
		return false
	}
	return true
}
