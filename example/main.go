package main

import (
	"log"

	".."
	"github.com/gizak/termui"
)

func main() {
	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialise UI: %s", err)
	}
	defer termui.Close()

	w := rich.New()

	w.X = 1
	w.Y = 1
	w.Resize(50, 3)

	termui.Handle("/sys/kbd", func(e termui.Event) {
		key := e.Data.(termui.EvtKbd)

		switch key.KeyStr {
		case "<left>":
			w.MoveCursor(-1)
		case "<right>":
			w.MoveCursor(1)
		case "<home>":
			// Start of line
		case "<end>":
			// End of line
		case "<up>":
			// Up one line
		case "<down>":
			// Down one line
		case "<delete>":
			w.Delete(1)
		case "<backspace>":
			w.Delete(-1)
		case "<enter>":
			w.WriteChar('\n')
		case "<space>":
			w.WriteChar(' ')
		case "C-q":
			termui.StopLoop()
		default:
			w.WriteChar(rune(key.KeyStr[0]))
		}
	})

	termui.Loop()
}
