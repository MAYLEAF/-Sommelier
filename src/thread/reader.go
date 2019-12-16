package thread

import (
	"fmt"
	"io"
	"log"
)

type reader struct {
	thread Handler
}

func (e *reader) read(thread *Handler) string {

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

	fmt.Print("Message from server: " + string(buf) + "\n\n")
	return string(buf)

}
