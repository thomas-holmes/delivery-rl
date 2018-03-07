package main

type Message int

const (
	PlayerUpdate Message = iota
	ClearRegion
	MoveEntity
	AttackEntity
	SpellTarget
	SpellLaunch
	PlayerDead
	KillEntity
	PlayerFloorChange
	ShowMenu
	EquipItem
	ShowFullGameLog
	SaveGame
	GameWon
	FoodSpoiled
	QuaffPotion
	ItemDetailClosed
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

type SpellTargetMessage struct {
	Spell Spell
	World *World
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

type QuaffPotionMessage struct {
	Potion Item
}

type ItemDetailClosedMessage struct {
	CloseInventory bool
}
