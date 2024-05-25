package ncurses

import (
	"fmt"
	"strings"
	"time"

	"github.com/JoshuaDoes/menuify"
	"seehuhn.de/go/ncurses"
)

var (
	leased *MenuScreen_Ncurses
	timing bool
)

type MenuScreen_Ncurses struct {
	Menu		*menuify.Menu
	Terminal	*ncurses.Window
	CachedFrame *menuify.MenuFrame

	//Padding for centered rendering, total for the count rather than one side
	paddingW int //i.e. use 6 if you want 3 lines of padding on both sides
	paddingH int
}

func NewMenuScreenNcurses(m *menuify.Menu) *MenuScreen_Ncurses {
	if m == nil {
		return nil
	}
	if leased != nil {
		ncurses.EndWin()
		leased = nil
		return NewMenuScreenNcurses(m)
	}
	if timing {
		menuify.Interval(time.Millisecond * 100, func() error {
			if !timing {
				return fmt.Errorf("timer gone")
			}
			return nil
		})
	}
	ms := &MenuScreen_Ncurses{
		Menu: m,
		Terminal: ncurses.Init(),
		paddingW: 6,
		paddingH: 6,
	}
	leased = ms

	go menuify.Interval(time.Nanosecond * 16666, func() error {
		if term := leased; term != nil {
			height, width := ms.Terminal.GetMaxYX()
			if height != ms.Menu.Engine.LinesV || width != ms.Menu.Engine.LinesH {
				ms.Menu.Engine.LinesV = height
				ms.Menu.Engine.LinesH = width
				ms.Menu.Engine.Redraw()
			}
			return nil
		}
		timing = false
		return fmt.Errorf("timer closed")
	})

	return ms
}

func (ms *MenuScreen_Ncurses) Render(frame *menuify.MenuFrame) {
	ms.CachedFrame = frame
	ms.Clear()
	if frame != nil && frame.Empty() {
		headLines := padStr(strings.Split(frame.Header, "\n"), ms.Engine.LinesH - ms.paddingW)
		head := strings.Join(headLines, "\n")
		menuLines := padStr(strings.Split(frame.Menu, "\n"), ms.Engine.LinesH - ms.paddingW)
		menu := strings.Join(menuLines, "\n")
		footLines := padStr(strings.Split(frame.Footer, "\n"), ms.Engine.LinesH - ms.paddingW)
		foot := strings.Join(footLines, "\n")

		text := fmt.Sprintf("%s\n\n%s", head, menu) //Inserted 2 new lines
		height := len(headLines) + len(menuLines) + ms.paddingH + 2

		if height < ms.Engine.LinesV {
			if height + len(footLines) + 2 < ms.Engine.LinesV {
				height += len(footLines) + 2
				pad := int(math.Floor(float64(menuEngine.LinesV - height))) + 1
				for i := 0; i < pad; i++ {
					text += "\n"
				}
				text += foot
			}
		}

		ms.Terminal.Printf(text)
	}
	ms.Terminal.Refresh()
}

func (ms *MenuScreen_Ncurses) GetFrame() *menuify.MenuFrame {
	return ms.CachedFrame
}

func (ms *MenuScreen_Ncurses) Clear() {
	ms.Terminal.Erase()
}

func (ms *MenuScreen_Ncurses) GetWidth() int {
	_, width := terminal.GetMaxYX()
	return width
}

func (ms *MenuScreen_Ncurses) GetHeight() int {
	height, _ := terminal.GetMaxYX()
	return height
}

//Close must be called by the creator, as screens could be repurposed after use
func (ms *MenuScreen_Ncurses) Close() {
	if leased == nil {
		return
	}
	ncurses.EndWin()
	leased = nil
}

//padStr pads the left side of each line with spaces, centering the multi-line text within the horizontal space while remaining left-justified
func padStr(lines []string, width int) []string {
	longest := 0
	for i := 0; i < len(lines); i++ {
		if len(lines[i]) > longest {
			longest = len(lines[i])
		}
	}
	if longest >= width {
		return lines
	}

	pad := int(math.Floor(float64(width - longest) / 2))
	padding := ""
	for i := 0; i < pad; i++ {
		padding += " "
	}

	for i := 0; i < len(lines); i++ {
		lines[i] = padding + lines[i]
	}
	return lines
}