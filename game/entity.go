package main

type BasicEntity struct {
	ID       int
	disabled bool
}

func (a BasicEntity) Enabled() bool {
	return !a.disabled
}

func (a *BasicEntity) Disable() {
	a.disabled = true
}

func (a *BasicEntity) Enable() {
	a.disabled = false
}

func (e *BasicEntity) SetIdentity(id int) {
	e.ID = id
}

func (e BasicEntity) Identity() int {
	return e.ID
}

type Entity interface {
	Identity() int
	SetIdentity(int)
	Enabled() bool
	Disable()
	Enable()
}
