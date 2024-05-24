package menuify

type MenuScreen interface {
	//Offloading for rendering the menu string, locks the menu until method returns
	Render(*MenuFrame)

	//Monospaced terminal screen size
	GetWidth() int
	GetHeight() int
}