package main

type Message int

const (
	MoveEntity Message = iota
	AttackEntity
	SpellTarget
	SpellLaunch
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

type PlayerQuaffPotionMessage struct {
	Potion Item
}

type PlayerActivateItemMessage struct {
	Item Item
}

type TryMoveCreatureMessage struct {
	*Creature
	NewX int
	NewY int
}
