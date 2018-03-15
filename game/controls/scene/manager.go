package scene

import (
	"log"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/gterm"
)

type Scene interface {
	Name() string
	Update(input controls.InputEvent, deltaT uint32)
	Render(wwindow *gterm.Window, deltaT uint32)
}

type Manager struct {
	scenes      map[string]Scene
	activeScene string
}

var defaultManager = NewManager()

func NewManager() *Manager {
	return &Manager{
		scenes:      make(map[string]Scene),
		activeScene: "",
	}
}

func (m *Manager) AddScene(s Scene) {
	if _, ok := m.scenes[s.Name()]; ok {
		log.Panicf("Can't add scene %s more than once", s.Name())
	}

	m.scenes[s.Name()] = s
}

func (m *Manager) RemoveScene(name string) {
	if _, ok := m.scenes[name]; !ok {
		log.Printf("Tried to remove non-registered scene %s", name)
	}

	delete(m.scenes, name)
}

func (m *Manager) SetActiveScene(name string) {
	if _, ok := m.scenes[name]; !ok {
		log.Panicf("Tried to change scenes to %s but it is not registered.", name)
	}
	m.activeScene = name
}

func (m *Manager) UpdateActiveScene(input controls.InputEvent, deltaT uint32) {
	if s, ok := m.scenes[m.activeScene]; !ok {
		log.Panicf("Tried to update active scene %s but it does not exist", s)
	} else {
		s.Update(input, deltaT)
	}
}

func (m *Manager) RenderActiveScene(window *gterm.Window, deltaT uint32) {
	if s, ok := m.scenes[m.activeScene]; !ok {
		log.Panicf("Tried to render active scene %s but it does not exist", s)
	} else {
		s.Render(window, deltaT)
	}
}

func (m *Manager) ClearScenes() {
	m.scenes = make(map[string]Scene)
}

func AddScene(s Scene) {
	defaultManager.AddScene(s)
}

func RemoveScene(name string) {
	defaultManager.RemoveScene(name)
}

func UpdateActiveScene(input controls.InputEvent, deltaT uint32) {
	defaultManager.UpdateActiveScene(input, deltaT)
}

func SetActiveScene(name string) {
	defaultManager.SetActiveScene(name)
}

func RenderActiveScene(window *gterm.Window, deltaT uint32) {
	defaultManager.RenderActiveScene(window, deltaT)
}

func ClearScenes() {
	defaultManager.ClearScenes()
}
