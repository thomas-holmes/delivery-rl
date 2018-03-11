package main

type Message int

const (
	MoveEntity Message = iota
	AttackEntity
	PlayerDead
	KillEntity
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

type MoveEntityMessage struct {
	ID   int
	OldX int
	OldY int
	NewX int
	NewY int
}

type AttackEntityMesasge struct {
	Attacker Entity
	Defender Entity
}

type KillEntityMessage struct {
	Attacker Entity
	Defender Entity
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
