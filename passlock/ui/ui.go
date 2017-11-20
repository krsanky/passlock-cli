package ui

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	runewidth "github.com/mattn/go-runewidth"
)

var row = 0
var style = tcell.StyleDefault

func putln(s tcell.Screen, str string) {
	puts(s, style, 1, row, str)
	row++
}

func puts(s tcell.Screen, style tcell.Style, x, y int, str string) {
	i := 0
	var deferred []rune
	dwidth := 0
	for _, r := range str {
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
}

func Ui() {
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	encoding.Register()

	if err = s.Init(); err != nil {
		panic(err)
	}

	plain := tcell.StyleDefault
	bold := style.Bold(true)

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	s.Clear()

	quit := make(chan struct{})

	style = bold
	putln(s, "Press ESC to Exit")
	putln(s, "Character set: "+s.CharacterSet())
	style = plain

	putln(s, "English:   October")
	putln(s, "Icelandic: október")
	putln(s, "Arabic:    أكتوبر")
	putln(s, "Russian:   октября")
	putln(s, "Greek:     Οκτωβρίου")
	putln(s, "Chinese:   十月 (note, two double wide characters)")
	putln(s, "Combining: A\u030a (should look like Angstrom)")
	putln(s, "Emoticon:  \U0001f618 (blowing a kiss)")
	putln(s, "Airplane:  \u2708 (fly away)")
	putln(s, "Command:   \u2318 (mac clover key)")
	putln(s, "Enclose:   !\u20e3 (should be enclosed exclamation)")
	putln(s, "")
	putln(s, "Box:")
	putln(s, string([]rune{
		tcell.RuneULCorner,
		tcell.RuneHLine,
		tcell.RuneTTee,
		tcell.RuneHLine,
		tcell.RuneURCorner,
	}))
	putln(s, string([]rune{
		tcell.RuneVLine,
		tcell.RuneBullet,
		tcell.RuneVLine,
		tcell.RuneLantern,
		tcell.RuneVLine,
	})+"  (bullet, lantern/section)")
	putln(s, string([]rune{
		tcell.RuneLTee,
		tcell.RuneHLine,
		tcell.RunePlus,
		tcell.RuneHLine,
		tcell.RuneRTee,
	}))
	putln(s, string([]rune{
		tcell.RuneVLine,
		tcell.RuneDiamond,
		tcell.RuneVLine,
		tcell.RuneUArrow,
		tcell.RuneVLine,
	})+"  (diamond, up arrow)")
	putln(s, string([]rune{
		tcell.RuneLLCorner,
		tcell.RuneHLine,
		tcell.RuneBTee,
		tcell.RuneHLine,
		tcell.RuneLRCorner,
	}))

	s.Show()
	go func() {
		for {
			ev := s.PollEvent()
			switch ev_ := ev.(type) {
			case *tcell.EventKey:
				switch ev_.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	<-quit

	s.Fini()
}
