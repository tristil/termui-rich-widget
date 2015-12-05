package rich

import (
	"fmt"
	"log"
	"os"
)

var logging = make(chan string, 1)

func Log(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	logging <- s

	go func() {
		message := <-logging
		f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)
		log.Println(message)
	}()
}
