package menuify

import (
	"fmt"
)

//MenuScreen is a wrapper for a screen manager, which could theoretically wrap multiple screens...
type MenuScreen interface {
	//Offloading for rendering the menu string, locks the menu until method returns
	Render(*MenuFrame)
	GetFrame() *MenuFrame //Returns the cached menu frame
	Clear()

	//Monospaced terminal screen size
	GetWidth() int
	GetHeight() int
}

func ScreenPrintf(ms MenuScreen, format string, args ...interface{}) {
	for {
		if len(format) == 0 {
			return
		}
		if format[len(format)-1] != '\n' {
			break
		}
		format = string(format[:len(format)-1])
	}
	line := fmt.Sprintf(format, args...)
	frame := ms.GetFrame()
	frame.Menu += line
	ms.Render(frame)
}

func ScreenPrintln(ms MenuScreen, line string) {
	ScreenPrintf(ms, "%s\n", line)
}
