package menuify

type MenuConfig struct {
	Environment map[string]string        `json:"environment"`
	Keybinds    []*MenuKeycodeBinding    `json:"keybinds"`
	HomeMenu    string                   `json:"home"`
	Menus       map[string]*MenuItemList `json:"menus"`
}