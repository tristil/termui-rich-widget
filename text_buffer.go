package rich

import (
	"bytes"
	"github.com/gizak/termui"
)

type TextBuffer struct {
	termui.Buffer

	CursorPosX int
	CursorPosY int
}

func NewTextBuffer() *TextBuffer {
	return &TextBuffer{
		Buffer:     termui.NewBuffer(),
		CursorPosX: 0,
		CursorPosY: 0,
	}
}

func (b *TextBuffer) String() string {
	buf := bytes.Buffer{}
	for p, c := range b.CellMap {
		if p.X == b.Area.Dx()-1 {
			buf.WriteRune('\n')
		}
		buf.WriteRune(c.Ch)
	}
	return buf.String()
}
