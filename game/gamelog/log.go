package gamelog

import "fmt"

var defaultLog GameLog

type GameLog struct {
	messages []string
}

func (g *GameLog) Append(format string, args ...interface{}) {
	g.messages = append(g.messages, fmt.Sprintf(format, args...))
}

func (g *GameLog) Messages() []string {
	return g.messages
}

func (g *GameLog) Clear() {
	g.messages = nil
}

func Append(format string, args ...interface{}) {
	defaultLog.Append(format, args...)
}

func Messages() []string {
	return defaultLog.Messages()
}

func Clear() {
	defaultLog.Clear()
}
