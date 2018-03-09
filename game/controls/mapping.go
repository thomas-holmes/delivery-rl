package controls

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Mapping struct {
	Name    string
	Keys    []string
	SdlKeys []uint32
}

var Up = Mapping{Name: "Up", Keys: []string{"K", "KP_8", "UP"}, SdlKeys: []uint32{sdl.K_k, sdl.K_KP_8, sdl.K_KP_8}}
var UpRight = Mapping{Name: "Up-Right", Keys: []string{"U", "KP_9"}, SdlKeys: []uint32{sdl.K_u, sdl.K_KP_9}}
var Right = Mapping{Name: "Right", Keys: []string{"L", "KP_6", "RIGHT"}, SdlKeys: []uint32{sdl.K_l, sdl.K_KP_6, sdl.K_RIGHT}}
var DownRight = Mapping{Name: "Down-Right", Keys: []string{"N", "KP_3"}, SdlKeys: []uint32{sdl.K_n, sdl.K_KP_3}}
var Down = Mapping{Name: "Down", Keys: []string{"J", "KP_2", "DOWN"}, SdlKeys: []uint32{sdl.K_j, sdl.K_KP_2, sdl.K_DOWN}}
var DownLeft = Mapping{Name: "DownLeft", Keys: []string{"B", "KP_1"}, SdlKeys: []uint32{sdl.K_b, sdl.K_KP_1}}
var Left = Mapping{Name: "Left", Keys: []string{"H", "KP_4", "LEFT"}, SdlKeys: []uint32{sdl.K_h, sdl.K_KP_4, sdl.K_LEFT}}
var UpLeft = Mapping{Name: "UpLeft", Keys: []string{"Y", "KP_7"}, SdlKeys: []uint32{sdl.K_y, sdl.K_KP_7}}
var Wait = Mapping{Name: "Wait", Keys: []string{".", "KP_.", "5"}, SdlKeys: []uint32{sdl.K_PERIOD, sdl.K_5, sdl.K_KP_5, sdl.K_KP_PERIOD}}

var Messages = Mapping{Name: "Messages", Keys: []string{"M"}, SdlKeys: []uint32{sdl.K_m}}

var Cast = Mapping{Name: "Cast", Keys: []string{"Z"}, SdlKeys: []uint32{sdl.K_z}}
var Equip = Mapping{Name: "Equip", Keys: []string{"E"}, SdlKeys: []uint32{sdl.K_e}}
var Inventory = Mapping{Name: "Inventory", Keys: []string{"I"}, SdlKeys: []uint32{sdl.K_i}}
var Quaff = Mapping{Name: "Quaff", Keys: []string{"Q"}, SdlKeys: []uint32{sdl.K_q}}
var Activate = Mapping{Name: "Activate", Keys: []string{"A"}, SdlKeys: []uint32{sdl.K_a}}
var Inspect = Mapping{Name: "Inspect", Keys: []string{"X"}, SdlKeys: []uint32{sdl.K_x}}

var AllMappings = []*Mapping{
	&Up,
	&UpRight,
	&Right,
	&DownRight,
	&Down,
	&DownLeft,
	&Left,
	&UpLeft,
	&Wait,
	&Messages,
	&Cast,
	&Equip,
	&Inventory,
	&Quaff,
	&Activate,
	&Inspect,
}
