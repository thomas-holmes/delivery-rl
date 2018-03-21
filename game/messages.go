package main

type Message int

const (
	MoveCreature Message = iota
	AttackCreature
	PlayerDead
	KillCreature
	PlayerFloorChange
	ShowMenu
	EquipItem
	ShowFullGameLog
	GameWon
	FoodSpoiled
	PlayerQuaffPotion
	PlayerActivateItem
	PlayerThrowItem
	PlayerWarp
	PlayerDropItem
	SplashGrease
	TryMoveCreature
)

type MoveCreatureMessage struct {
	Creature *Creature
	OldX     int
	OldY     int
	NewX     int
	NewY     int
}

type AttackCreatureMessage struct {
	Attacker *Creature
	Defender *Creature
}

type KillCreatureMessage struct {
	Attacker *Creature
	Defender *Creature
}

type PlayerFloorChangeMessage struct {
	Stair
}

type ShowMenuMessage struct {
	Menu Menu
}

type EquipItemMessage struct {
	Item
}

type PlayerQuaffPotionMessage struct {
	Potion Item
}

type PlayerActivateItemMessage struct {
	Item Item
}

type PlayerThrowItemMessage struct {
	Item

	TargetX int
	TargetY int

	*World
}

type PlayerWarpMessage struct {
	TargetX int
	TargetY int

	Cost int

	*World
}

type PlayerDropItemMessage struct {
	Item

	*World
}

type SplashGreaseMessage struct {
	Item

	X int
	Y int
}

type TryMoveCreatureMessage struct {
	*Creature
	NewX int
	NewY int
}
