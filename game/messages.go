package main

type Message int

const (
	PlayerUpdate Message = iota
	ClearRegion
	MoveEntity
	AttackEntity
	SpellLaunch
	PlayerDead
	KillEntity
	PlayerFloorChange
	ShowMenu
	EquipItem
	GameLogAppend
	ShowFullGameLog
	SaveGame
	GameWon
	FoodSpoiled
	QuaffPotion
)

type ClearRegionMessage struct {
	X int
	Y int
	W int
	H int
}

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

type SpellLaunchMessage struct {
	Caster Entity
	X      int
	Y      int
	Spell  Spell
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

type GameLogAppendMessage struct {
	Messages []string
}

type QuaffPotionMessage struct {
	Potion Item
}
