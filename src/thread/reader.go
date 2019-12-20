package thread

import (
	"encoding/json"
	"logger"
	"sync"
)

type reader struct {
	thread Handler
	lock   sync.Mutex
}

func (e *reader) read(thread *Handler) string {
	logger := logger.Logger()
	e.lock.Lock()
	defer e.lock.Unlock()

	/*
		buf := make([]byte, 0, 32768)
		tmp := make([]byte, 256)

		for {
			n, err := thread.conn.Read(tmp[0:])
			if err != nil {
				if err != io.EOF {
					logger.Error("Reader read error:", err)
				}
				break
			}
			if n < 256 {
				buf = append(buf, tmp[:n]...)
				break
			}
			buf = append(buf, tmp[:n]...)
		}
	*/
	d := json.NewDecoder(thread.conn)

	msg := make(map[string]interface{})
	_ = d.Decode(&msg)
	buf, _ := json.Marshal(msg)

	logger.Info("Message from server: " + string(buf) + "\n")
	return string(buf)

}
