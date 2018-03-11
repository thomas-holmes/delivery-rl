package controls

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Mapping struct {
	Name     string
	Action   Action
	Keys     []string
	SdlKeys  []sdl.Keycode
	Shift    bool
	Control  bool
	HideHelp bool
}

var KeyUp = Mapping{Name: "Up", Action: Up, Keys: []string{"K", "KP_8", "UP"}, SdlKeys: []sdl.Keycode{sdl.K_k, sdl.K_KP_8, sdl.K_UP}}
var KeyUpRight = Mapping{Name: "Up-Right", Action: UpRight, Keys: []string{"U", "KP_9"}, SdlKeys: []sdl.Keycode{sdl.K_u, sdl.K_KP_9}}
var KeyRight = Mapping{Name: "Right", Action: Right, Keys: []string{"L", "KP_6", "RIGHT"}, SdlKeys: []sdl.Keycode{sdl.K_l, sdl.K_KP_6, sdl.K_RIGHT}}
var KeyDownRight = Mapping{Name: "Down-Right", Action: DownRight, Keys: []string{"N", "KP_3"}, SdlKeys: []sdl.Keycode{sdl.K_n, sdl.K_KP_3}}
var KeyDown = Mapping{Name: "Down", Action: Down, Keys: []string{"J", "KP_2", "DOWN"}, SdlKeys: []sdl.Keycode{sdl.K_j, sdl.K_KP_2, sdl.K_DOWN}}
var KeyDownLeft = Mapping{Name: "DownLeft", Action: DownLeft, Keys: []string{"B", "KP_1"}, SdlKeys: []sdl.Keycode{sdl.K_b, sdl.K_KP_1}}
var KeyLeft = Mapping{Name: "Left", Action: Left, Keys: []string{"H", "KP_4", "LEFT"}, SdlKeys: []sdl.Keycode{sdl.K_h, sdl.K_KP_4, sdl.K_LEFT}}
var KeyUpLeft = Mapping{Name: "UpLeft", Action: UpLeft, Keys: []string{"Y", "KP_7"}, SdlKeys: []sdl.Keycode{sdl.K_y, sdl.K_KP_7}}
var KeyFive = Mapping{Name: "Wait", Action: Wait, Keys: []string{".", "KP_.", "5"}, SdlKeys: []sdl.Keycode{sdl.K_PERIOD, sdl.K_5, sdl.K_KP_5, sdl.K_KP_PERIOD}}
var KeyGreater = Mapping{Name: "Descend", Action: Descend, Keys: []string{">"}, SdlKeys: []sdl.Keycode{sdl.K_PERIOD}, Shift: true}
var KeyLesser = Mapping{Name: "Ascend", Action: Ascend, Keys: []string{"<"}, SdlKeys: []sdl.Keycode{sdl.K_COMMA}, Shift: true}

var KeyQuestion = Mapping{Name: "Help", Action: Help, Keys: []string{"?"}, SdlKeys: []sdl.Keycode{sdl.K_SLASH}, Shift: true}

var KeyM = Mapping{Name: "Messages", Action: Messages, Keys: []string{"M"}, SdlKeys: []sdl.Keycode{sdl.K_m}}
var KeyZ = Mapping{Name: "Cast", Action: Cast, Keys: []string{"Z"}, SdlKeys: []sdl.Keycode{sdl.K_z}}
var KeyE = Mapping{Name: "Equip", Action: Equip, Keys: []string{"E"}, SdlKeys: []sdl.Keycode{sdl.K_e}, HideHelp: true}
var KeyI = Mapping{Name: "Inventory", Action: Inventory, Keys: []string{"I"}, SdlKeys: []sdl.Keycode{sdl.K_i}}
var KeyQ = Mapping{Name: "Quaff", Action: Quaff, Keys: []string{"Q"}, SdlKeys: []sdl.Keycode{sdl.K_q}, HideHelp: true}
var KeyA = Mapping{Name: "Activate", Action: Activate, Keys: []string{"A"}, SdlKeys: []sdl.Keycode{sdl.K_a}, HideHelp: true}
var KeyX = Mapping{Name: "Examine", Action: Examine, Keys: []string{"X"}, SdlKeys: []sdl.Keycode{sdl.K_x}}
var KeyG = Mapping{Name: "Get", Action: Get, Keys: []string{"G"}, SdlKeys: []sdl.Keycode{sdl.K_g}}
var KeyT = Mapping{Name: "Throw", Action: Throw, Keys: []string{"T"}, SdlKeys: []sdl.Keycode{sdl.K_t}, HideHelp: true}
var KeyD = Mapping{Name: "Drop", Action: Drop, Keys: []string{"D"}, SdlKeys: []sdl.Keycode{sdl.K_d}, HideHelp: true}

var KeyCQ = Mapping{Name: "Quit", Action: Quit, Keys: []string{"CTRL-Q"}, SdlKeys: []sdl.Keycode{sdl.K_q}, Control: true}

var KeyEnter = Mapping{Name: "Confirm", Action: Confirm, Keys: []string{"enter", "return"}, SdlKeys: []sdl.Keycode{sdl.K_RETURN, sdl.K_KP_ENTER}, HideHelp: true}
var KeyEsc = Mapping{Name: "Cancel", Action: Cancel, Keys: []string{"ESC"}, SdlKeys: []sdl.Keycode{sdl.K_ESCAPE}, HideHelp: true}
var KeyPgUp = Mapping{Name: "Page Up", Action: SkipUp, Keys: []string{"PGUP"}, SdlKeys: []sdl.Keycode{sdl.K_PAGEUP}, HideHelp: true}
var KeyPgDown = Mapping{Name: "Page Down", Action: SkipDown, Keys: []string{"PGDN"}, SdlKeys: []sdl.Keycode{sdl.K_PAGEDOWN}, HideHelp: true}
var KeyHome = Mapping{Name: "Home", Action: Top, Keys: []string{"HOME"}, SdlKeys: []sdl.Keycode{sdl.K_HOME}, HideHelp: true}
var KeyEnd = Mapping{Name: "End", Action: Bottom, Keys: []string{"END"}, SdlKeys: []sdl.Keycode{sdl.K_END}, HideHelp: true}

var AllMappings = []*Mapping{
	&KeyUp,
	&KeyUpRight,
	&KeyRight,
	&KeyDownRight,
	&KeyDown,
	&KeyDownLeft,
	&KeyLeft,
	&KeyUpLeft,
	&KeyGreater,
	&KeyLesser,
	&KeyFive,
	&KeyEsc,
	&KeyM,
	&KeyZ,
	&KeyE,
	&KeyI,
	&KeyQ,
	&KeyA,
	&KeyX,
	&KeyG,
	&KeyT,
	&KeyD,
	&KeyEnter,
	&KeyQuestion,
	&KeyCQ,
	&KeyPgUp,
	&KeyPgDown,
	&KeyHome,
	&KeyEnd,
}

type Action int

const (
	None Action = iota
	Up
	UpRight
	Right
	DownRight
	Down
	DownLeft
	Ascend
	Descend
	Left
	UpLeft
	Wait
	Cancel
	Messages
	Cast
	Equip
	Inventory
	Quaff
	Activate
	Throw
	Drop
	Examine
	Get
	Confirm
	Help
	Quit
	SkipUp
	SkipDown
	Top
	Bottom
	FinalUnused
)

type InputEvent struct {
	sdl.Event
	sdl.Keymod
}

func (i InputEvent) Action() Action {
	shiftPressed := i.Keymod&sdl.KMOD_SHIFT > 0
	controlPressed := i.Keymod&sdl.KMOD_CTRL > 0
	switch e := i.Event.(type) {
	case sdl.KeyDownEvent:
		keyPressed := e.Keysym.Sym

		for _, m := range AllMappings {
			for _, mappedKey := range m.SdlKeys {
				if mappedKey == keyPressed && shiftPressed == m.Shift && controlPressed == m.Control {
					return m.Action
				}
			}
		}
	}

	return None
}
