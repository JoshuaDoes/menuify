package menuify

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/JoshuaDoes/json"
)

var (
	keyCalibration map[string][]*MenuKeycodeBinding = make(map[string][]*MenuKeycodeBinding)
)

type MenuKeycodeBinding struct {
	Keycode   uint16 `json:"keycode"`
	Action    string `json:"action"`
	OnRelease bool   `json:"onRelease"`
}

func (me *MenuEngine) BindKeys() {
	for keyboard, bindings := range keyCalibration {
		kl, err := NewKeycodeListener(keyboard)
		if err != nil {
			panic(fmt.Smenuify.ScreenPrintf(me.Screen, "error listening to keyboard %s: %v", keyboard, err))
		}
		for _, binding := range bindings {
			var action func()
			switch binding.Action {
			case "prevItem":
				action = me.PrevItem
			case "nextItem":
				action = me.NextItem
			case "selectItem":
				action = me.Action
			default:
				panic("unknown action: " + binding.Action)
			}
			kl.Bind(binding.Keycode, binding.OnRelease, action)
		}
		go kl.Run()
	}
}

type KeyCalibration struct {
	Ready  bool
	Cancel bool
	Action string
	KLs    []*KeycodeListener
}

func (kc *KeyCalibration) Input(keyboard string, keycode uint16, onRelease bool) {
	if kc.Cancel {
		return
	}
	if !kc.Ready {
		kc.Cancel = true
		return
	}
	if kc.Action == "" || kc.Action == "cancel" {
		kc.Action = ""
		return
	}
	if onRelease {
		return
	}
	if keyCalibration[keyboard] == nil {
		keyCalibration[keyboard] = make([]*MenuKeycodeBinding, 0)
	}
	keyCalibration[keyboard] = append(keyCalibration[keyboard], &MenuKeycodeBinding{
		Keycode:   keycode,
		Action:    kc.Action,
		OnRelease: true,
	})
	kc.Action = ""
}

func (me *MenuEngine) Calibrate(keyCalibrationFile string) error {
	if keyCalibrationFile == "" {
		keyCalibrationFile = "./keyCalibration.json"
	}

	//Generate a key calibration file if one doesn't exist yet
	calibrator := &KeyCalibration{KLs: make([]*KeycodeListener, 0)}

	//Get a list of keyboards
	keyboards := make([]string, 0)
	err := filepath.Walk("/dev/input", func(path string, info os.FileInfo, err error) error {
		if len(path) < 16 || string(path[:16]) != "/dev/input/event" {
			return nil
		}
		keyboards = append(keyboards, path)
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking inputs: %v", err)
	}

	//Bind all keyboards to calibrator input
	for _, keyboard := range keyboards {
		kl, err := NewKeycodeListener(keyboard)
		if err != nil {
			return fmt.Errorf("error listening to walked keyboard %s: %v", keyboard, err)
		}
		kl.RootBind = calibrator.Input
		calibrator.KLs = append(calibrator.KLs, kl)
		go kl.Run()
	}

	//Start calibrating!
	stages := 6
	for stage := 0; stage < stages; stage++ {
		switch stage {
		case 0:
			keyCalibrationJSON, err := ioutil.ReadFile(keyCalibrationFile)
			if err == nil {
				keyCalibration = make(map[string][]*MenuKeycodeBinding)
				err = json.Unmarshal(keyCalibrationJSON, &keyCalibration)
				if err != nil {
					stage = 1
					continue
				}

				me.Screen.Clear()
				menuify.ScreenPrintln(me.Screen, "Press any key within\n5 seconds to recalibrate.\n")
				calibrator.Ready = true
				calibrator.Action = "cancel"
				timeout := time.Now()
				for calibrator.Action != "" {
					if time.Now().Sub(timeout).Seconds() > 5 {
						break
					}
					time.Sleep(time.Millisecond * 100)
				}
				if time.Now().Sub(timeout).Seconds() < 5 {
					calibrator.Action = ""
					menuify.ScreenPrintln(me.Screen, "Recalibration time!")
					time.Sleep(time.Second * 2)
					continue
				}
				stage = stages-1 //Skip to the end of the stages
			}
		case 1:
			calibrator.Ready = false
			keyCalibration = make(map[string][]*MenuKeycodeBinding)
			me.Screen.Clear()
			menuify.ScreenPrintln(me.Screen, "Welcome to the calibrator!\n")
			menuify.ScreenPrintln(me.Screen, "Press any key to cancel.\n")
			time.Sleep(time.Second * 2)
			if calibrator.Cancel { return ERR_CANCELLED }
			menuify.ScreenPrintln(me.Screen, "Controllers and remotes\nare also supported.\n")
			time.Sleep(time.Second * 2)
			if calibrator.Cancel { return ERR_CANCELLED }
			menuify.ScreenPrintln(me.Screen, "This is a guided process.\n")
			time.Sleep(time.Second * 2)
			if calibrator.Cancel { return ERR_CANCELLED }
			menuify.ScreenPrintln(me.Screen, "Get ready!\n")
			if calibrator.Cancel { return ERR_CANCELLED }
			time.Sleep(time.Second * 3)
			if calibrator.Cancel { return ERR_CANCELLED }
		case 2:
			me.Screen.Clear()
			calibrator.Ready = true
			calibrator.Action = "nextItem"
			menuify.ScreenPrintf(me.Screen, "\n")
			menuify.ScreenPrintln(me.Screen, "Press any key to use to\nnavigate down in a menu.\n")
			menuify.ScreenPrintln(me.Screen, "Recommended: volume down")
			for calibrator.Action != "" {
			}
		case 3:
			calibrator.Action = "prevItem"
			menuify.ScreenPrintf(me.Screen, "\n")
			menuify.ScreenPrintln(me.Screen, "Press any key to use to\nnavigate up in a menu.\n")
			menuify.ScreenPrintln(me.Screen, "Recommended: volume up")
			for calibrator.Action != "" {
			}
		case 4:
			calibrator.Action = "selectItem"
			menuify.ScreenPrintf(me.Screen, "\n")
			menuify.ScreenPrintln(me.Screen, "Press any key to use to\nselect a menu item.\n")
			menuify.ScreenPrintln(me.Screen, "Recommended: touch screen")
			for calibrator.Action != "" {
			}
		case 5:
			me.Screen.Clear()
			menuify.ScreenPrintln(me.Screen, "Saving results...\n")
			keyboards, err := json.Marshal(keyCalibration, true)
			if err != nil {
				return fmt.Errorf("error encoding calibration results: %v", err)
			}
			keyboardsFile, err := os.Create(keyCalibrationFile)
			if err != nil {
				return fmt.Errorf("error creating calibration file: %v", err)
			}
			defer keyboardsFile.Close()
			_, err = keyboardsFile.Write(keyboards)
			if err != nil {
				return fmt.Errorf("error writing calibration file: %v", err)
			}
			//menuify.ScreenPrintln(me.Screen, string(keyboards))
			//menuify.ScreenPrintln(me.Screen, "Calibration complete!")
			//time.Sleep(time.Second * 2)
			//calibrator.Ready = false
		}
	}

	for i := 0; i < len(calibrator.KLs); i++ {
		calibrator.KLs[i].Close()
	}
	return nil
}
