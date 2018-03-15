package messages

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

// Change this to enforce a dispose function or something
// so I can clean these up without leaking memory
type Listener func(M)

type subscription struct {
	id int
	Listener
}

// Subscribe to receive messages to the supplied listener callback. Call the returned
// unsubscribe method to unregister your listener

// Probably a gross memory leak
func (b *Broker) Subscribe(listener Listener) Unsubscribe {
	b.lastSubID++
	id := b.lastSubID

	b.subscriptions = append(b.subscriptions, subscription{id: id, Listener: listener})

	return func() {
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

func (b *Broker) UnSubAll() {
	for i, _ := range b.subscriptions {
		b.subscriptions[i] = subscription{}
	}
	b.subscriptions = nil
}

// Broadcast sends a message on the default broker
func Broadcast(m M) {
	defaultBroker.Broadcast(m)
}

// Subscribe to the default broker
func Subscribe(listener Listener) Unsubscribe {
	return defaultBroker.Subscribe(listener)
}

func UnSubAll() {
	defaultBroker.UnSubAll()
}
