package messages

import (
	"testing"
)

func TestMessageSend(t *testing.T) {
	var broker Broker

	var sideEffect int

	broker.Subscribe(func(m M) {
		sideEffect = 5
	})

	broker.Broadcast(M{ID: 1, Data: nil})

	if sideEffect != 5 {
		t.Error("Expected side effect to be set to 5, instead got", sideEffect)
	}
}

func TestUnsubPreventsMessages(t *testing.T) {
	var broker Broker

	var sideEffect int

	unsub := broker.Subscribe(func(m M) {
		sideEffect += 1
	})

	broker.Broadcast(M{ID: 1, Data: nil})

	if sideEffect != 1 {
		t.Error("Expected side effect to be set to 1, instead got", sideEffect)
	}

	unsub()

	broker.Broadcast(M{ID: 1, Data: nil})
	broker.Broadcast(M{ID: 1, Data: nil})
	broker.Broadcast(M{ID: 1, Data: nil})

	if sideEffect != 1 {
		t.Error("Unsub didn't work. Expected sideEffect to be 1, instead got", sideEffect)
	}
}

func TestUnsubMultipleMiddleSubscriber(t *testing.T) {
	var broker Broker

	var sideEffect1, sideEffect2, sideEffect3, sideEffect4 int

	broker.Subscribe(func(m M) {
		sideEffect1++
	})
	unsub2 := broker.Subscribe(func(m M) {
		sideEffect2++
	})
	unsub3 := broker.Subscribe(func(m M) {
		sideEffect3++
	})
	broker.Subscribe(func(m M) {
		sideEffect4++
	})

	broker.Broadcast(M{ID: 1, Data: nil})

	if sideEffect1 != 1 || sideEffect2 != 1 || sideEffect3 != 1 || sideEffect4 != 1 {
		t.Error("Expected values 1, 1, 1, 1, instead got", sideEffect1, sideEffect2, sideEffect3, sideEffect4)
	}

	unsub2()

	broker.Broadcast(M{ID: 1, Data: nil})

	if sideEffect1 != 2 || sideEffect2 != 1 || sideEffect3 != 2 || sideEffect4 != 2 {
		t.Error("Expected values 2, 1, 2, 2, instead got", sideEffect1, sideEffect2, sideEffect3, sideEffect4)
	}

	unsub3()

	broker.Broadcast(M{ID: 1, Data: nil})

	if sideEffect1 != 3 || sideEffect2 != 1 || sideEffect3 != 2 || sideEffect4 != 3 {
		t.Error("Expected values 3, 1, 2, 3, instead got", sideEffect1, sideEffect2, sideEffect3, sideEffect4)
	}
}

func TestEasyComparisonWithCustomTypes(t *testing.T) {
	var broker Broker

	type MyMID int
	const (
		MID1 MyMID = iota
		MID2
		MID3
	)

	var sideEffect int

	broker.Subscribe(func(m M) {
		if m.ID == MID1 {
			sideEffect = 5
		}
	})

	broker.Broadcast(M{ID: MID1})

	if sideEffect != 5 {
		t.Error("Expected sideEffect to be 5, instead got", sideEffect)
	}
}
