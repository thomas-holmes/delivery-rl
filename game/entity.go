package main

type BasicEntity struct {
	ID int
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
}
