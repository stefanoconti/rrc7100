package uiterm

import (
	"strings"
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

type InputFunc func(ui *Ui, textbox *Textbox, text string)

type Textbox struct {
	Text string
	Fg   Attribute
	Bg   Attribute

	Input InputFunc

	ui *Ui
	active         bool
	x0, y0, x1, y1 int
}

func (t *Textbox) uiInitialize(ui *Ui) {
	t.ui = ui
}

func (t *Textbox) setBounds(x0, y0, x1, y1 int) {
	t.x0 = x0
	t.y0 = y0
	t.x1 = x1
	t.y1 = y1
}

func (t *Textbox) setActive(active bool) {
	t.active = active
}

func (t *Textbox) draw() {
	var setCursor = false
	reader := strings.NewReader(t.Text)
	for y := t.y0; y < t.y1; y++ {
		for x := t.x0; x < t.x1; x++ {
			var chr rune
			if ch, _, err := reader.ReadRune(); err != nil {
				if t.active && !setCursor {
					termbox.SetCursor(x, y)
					setCursor = true
				}
				chr = ' '
			} else {
				chr = ch
			}
			termbox.SetCell(x, y, chr, termbox.Attribute(t.Fg), termbox.Attribute(t.Bg))
		}
	}
}

func (t *Textbox) keyEvent(mod Modifier, key Key) {
	redraw := false
	switch key {
	case KeyCtrlC:
		t.Text = ""
		redraw = true
	case KeyEnter:
		if t.Input != nil {
			t.Input(t.ui, t, t.Text)
		}
		t.Text = ""
		redraw = true
	case KeySpace:
		t.Text = t.Text + " "
		redraw = true
	case KeyBackspace:
	case KeyBackspace2:
		if len(t.Text) > 0 {
			if r, size := utf8.DecodeLastRuneInString(t.Text); r != utf8.RuneError {
				t.Text = t.Text[:len(t.Text)-size]
				redraw = true
			}
		}
	}
	if redraw {
		t.draw()
		termbox.Flush()
	}
}

func (t *Textbox) characterEvent(chr rune) {
	t.Text = t.Text + string(chr)
	t.draw()
	termbox.Flush()
}
