package main

func (w *World) LoadSavedWorld(s SaveV0) {
	window := w.Window
	messageBus := w.messageBus

	nw := &World{}
	nw.Window = window
	nw.messageBus = messageBus

	s.Restore(nw)

	w = nw
}
