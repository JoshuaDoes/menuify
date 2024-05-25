package menuify

import (
	"fmt"
	"io/ioutil"

	"github.com/JoshuaDoes/json"
)

type Menu struct {
	Config *MenuConfig
	Engine *MenuEngine
	Screen *MenuScreen
	Keysrv []*KeycodeListener
}

func NewMenu() *Menu {
	return &Menu{Engine: NewMenuEngine(), Keysrv: make([]*KeycodeListener, 0)}
}

func (m *Menu) SetScreen(screen MenuScreen) {
	m.Engine.SetScreen(screen)
}

func (m *Menu) Load(configPath string) error {
	if m.Engine == nil {
		return fmt.Errorf("menu: need engine to load config")
	}

	configJSON, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	cfg := &MenuConfig{}
	if err := json.Unmarshal(configJSON, cfg); err != nil {
		return err
	}

	for key, val := range cfg.Environment {
		m.Engine.Environment[key] = val
	}

	m.Engine.ClearMenus()
	for id, itemList := range cfg.Menus {
		m.Engine.AddMenu(id, itemList)
	}
	m.Engine.HomeMenu = cfg.HomeMenu

	//TODO: Load keybinds from cfg.Keybinds, apply to m.Keysrv
	//If necessary, unload all keybinds first in case of hot reloading config

	return nil
}
