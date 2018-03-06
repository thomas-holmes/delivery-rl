package messages

import "log"

var defaultBroker Broker

// M is the message payload struct. ID is an opaque identifier and
// Data can be anything!
type M struct {
	ID   interface{}
	Data interface{}
}

type Broker struct {
	lastSubID     int
	subscriptions []subscription
}

type Unsubscribe func()

type Listener func(M)

type subscription struct {
	id int
	Listener
}

// Subscribe to receive messages to the supplied listener callback. Call the returned
// unsubscribe method to unregister your listener
func (b *Broker) Subscribe(listener Listener) Unsubscribe {
	b.lastSubID++
	id := b.lastSubID

	b.subscriptions = append(b.subscriptions, subscription{id: id, Listener: listener})

	return func() {
		log.Printf("Removing listener %+v", listener)
		for i, s := range b.subscriptions {
			if id == s.id {
				b.subscriptions[i] = subscription{}
				b.subscriptions = append(b.subscriptions[:i], b.subscriptions[i+1:]...)
				return
			}
		}
	}
}

// Send a message to all subscribers on this broker
func (b *Broker) Broadcast(m M) {
	for _, s := range b.subscriptions {
		s.Listener(m)
	}
}

// Broadcast sends a message on the default broker
func Broadcast(m M) {
	defaultBroker.Broadcast(m)
}

// Subscribe to the default broker
func Subscribe(listener Listener) Unsubscribe {
	return defaultBroker.Subscribe(listener)
}
