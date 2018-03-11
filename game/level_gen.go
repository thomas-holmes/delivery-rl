package main

import (
	"github.com/MichaelTJones/pcg"
	"github.com/thomas-holmes/delivery-rl/game/items"
)

type Room struct {
	ID          int
	X           int
	Y           int
	W           int
	H           int
	connectedTo []int
}

type CandidateTile struct {
	TileKind TileKind
	Item     Item
	X        int
	Y        int
}

type CandidateLevel struct {
	rng *pcg.PCG64

	W int
	H int

	stairCandidates []Position

	depth int // One Indexed

	nextRoomID int

	flags LevelGenFlag

	rooms map[int]*Room
	tiles []CandidateTile
}

const (
	MinRoomWidth            = 3
	MinRoomHeight           = 3
	MaxRoomWidth            = 20
	MaxRoomHeight           = 20
	MaxRoomIterations       = 200
	MaxWeaponPlacement      = 5
	MaxConsumeablePlacement = 20
	MaxArmourPlacement      = 5
)

func (level *CandidateLevel) genNextRoomID() int {
	id := level.nextRoomID
	level.nextRoomID++
	return id
}

func (level *CandidateLevel) tryAddRandomRoom() {
	widthBound := uint64(MaxRoomWidth - MinRoomWidth)
	heightBound := uint64(MaxRoomHeight - MinRoomHeight)

	randomWidth := int(level.rng.Bounded(widthBound)) + MinRoomWidth
	randomHeight := int(level.rng.Bounded(heightBound)) + MinRoomHeight

	topLeftX := int(level.rng.Bounded(uint64(level.W-randomWidth-1))) + 1
	topLeftY := int(level.rng.Bounded(uint64(level.H-randomHeight-1))) + 1

	for y := topLeftY; y < topLeftY+randomHeight; y++ {
		for x := topLeftX; x < topLeftX+randomWidth; x++ {
			if level.tiles[y*level.W+x].TileKind != Wall {
				// Just quit if we run into a non-wall feature
				return
			}
		}
	}

	// We can place our room, so make it all floors.
	for y := topLeftY; y < topLeftY+randomHeight; y++ {
		for x := topLeftX; x < topLeftX+randomWidth; x++ {
			level.tiles[y*level.W+x].TileKind = Floor
		}
	}

	room := &Room{
		ID: level.genNextRoomID(),
		X:  topLeftX,
		Y:  topLeftY,
		W:  randomWidth,
		H:  randomHeight,
	}

	level.rooms[room.ID] = room

	if room.ID > 0 {
		level.connectRooms(room.ID, room.ID-1)
	}
}

func (level *CandidateLevel) addRooms() {
	for i := 0; i < MaxRoomIterations; i++ {
		level.tryAddRandomRoom()
	}
}

func (room *Room) chooseTopWall(rng *pcg.PCG64) (int, int) {
	x := int(rng.Bounded(uint64(room.W))) + room.X
	y := room.Y
	return x, y
}
func (room *Room) chooseRightWall(rng *pcg.PCG64) (int, int) {
	x := room.X + room.W - 1
	y := int(rng.Bounded(uint64(room.H))) + room.Y
	return x, y
}
func (room *Room) chooseBottomWall(rng *pcg.PCG64) (int, int) {
	x := int(rng.Bounded(uint64(room.W))) + room.X
	y := room.Y + room.H - 1
	return x, y
}
func (room *Room) chooseLeftWall(rng *pcg.PCG64) (int, int) {
	x := room.X
	y := int(rng.Bounded(uint64(room.H))) + room.Y
	return x, y
}

func (level *CandidateLevel) connectRooms(roomId1 int, roomId2 int) {
	room1 := level.rooms[roomId1]
	room2 := level.rooms[roomId2]

	const (
		Top int = iota
		Right
		Bottom
		Left
	)

	yDelta := room1.Y - room2.Y
	xDelta := room1.X - room2.X

	var r1x, r1y, r2x, r2y int

	if yDelta > 0 { // Below
		if xDelta > 0 { // Right
			if xDelta > yDelta { // further right than below
				r1x, r1y = room1.chooseLeftWall(level.rng)
				r2x, r2y = room2.chooseRightWall(level.rng)
			} else { // further below than right
				r1x, r1y = room1.chooseTopWall(level.rng)
				r2x, r2y = room2.chooseBottomWall(level.rng)
			}
		} else { // Left
			if -xDelta > yDelta { // further left than below
				r1x, r1y = room1.chooseRightWall(level.rng)
				r2x, r2y = room2.chooseLeftWall(level.rng)
			} else { // further below than left
				r1x, r1y = room1.chooseTopWall(level.rng)
				r2x, r2y = room2.chooseBottomWall(level.rng)
			}
		}
	} else { // Above
		if xDelta > 0 { // Right
			if xDelta > -yDelta { // further right than above
				r1x, r1y = room1.chooseLeftWall(level.rng)
				r2x, r2y = room2.chooseRightWall(level.rng)
			} else { // further above than right
				r1x, r1y = room1.chooseBottomWall(level.rng)
				r2x, r2y = room2.chooseTopWall(level.rng)
			}
		} else { // Left
			if -xDelta > -yDelta { // further left than above
				r1x, r1y = room1.chooseRightWall(level.rng)
				r2x, r2y = room2.chooseLeftWall(level.rng)
			} else { // further above than left
				r1x, r1y = room1.chooseBottomWall(level.rng)
				r2x, r2y = room2.chooseTopWall(level.rng)
			}
		}
	}

	if r1x > r2x {
		r1x, r2x = r2x, r1x
		r1y, r2y = r2y, r1y
	}
	for ; r1x <= r2x; r1x++ {
		level.tiles[r1y*level.W+r1x].TileKind = Floor
	}
	r1x--
	if r1y > r2y {
		r1x, r2x = r2x, r1x
		r1y, r2y = r2y, r1y
	}
	for ; r1y <= r2y; r1y++ {
		level.tiles[r1y*level.W+r1x].TileKind = Floor
	}
}

func (c *CandidateLevel) chooseStairs() (Position, Position) {
	diagonalDistance := distance(Position{X: 0, Y: 0}, Position{X: c.W - 1, Y: c.H - 1})

	for _, p1 := range c.stairCandidates {
		for _, p2 := range c.stairCandidates {
			if p1 == p2 {
				continue
			}

			if distance(p1, p2) > diagonalDistance/3 {
				return p1, p2
			}
		}
	}

	// This *really* shouldn't happen. Should be impossible.
	panic("Something went horribly wrong with level gen, so sorry. Better luck next seed")
}

func (level *CandidateLevel) addContourStairs() {
	p1, p2 := level.chooseStairs()

	if level.flags&GenUpStairs != 0 {
		level.tiles[p1.Y*level.W+p1.X].TileKind = UpStair
	}

	if level.flags&GenDownStairs != 0 {
		level.tiles[p2.Y*level.W+p2.X].TileKind = DownStair
	}
}

func (level *CandidateLevel) addItems(collection items.Collection, max int) {
	for i := 0; i < max; i++ {
		itemDef := items.GetLevelBoundedItem(level.rng, collection, level.depth)

		randomItem := produceItem(itemDef)

		tileIndex := level.rng.Bounded(uint64(len(level.tiles)))
		if level.tiles[tileIndex].TileKind == Floor {
			level.tiles[tileIndex].Item = randomItem
		}
	}
}

func (level *CandidateLevel) encodeAsString() string {
	levelStr := ""
	for y := 0; y < level.H; y++ {
		if y != 0 {
			levelStr += "\n"
		}
		for x := 0; x < level.W; x++ {
			switch level.tiles[y*level.W+x].TileKind {
			case Wall:
				levelStr += string(WallGlyph)
			case Floor:
				item := level.tiles[y*level.W+x].Item
				if item != (Item{}) {
					levelStr += string(item.Symbol)
				} else {
					levelStr += string(FloorGlyph)
				}
			case DownStair:
				levelStr += string(DownStairGlyph)
			case UpStair:
				levelStr += string(UpStairGlyph)
			}
		}
	}

	return levelStr
}

type LevelGenFlag int

const (
	GenUpStairs = 1 << iota
	GenDownStairs
)

func GenLevel(rng *pcg.PCG64, maxX int, maxY int, depth int, flags LevelGenFlag) *CandidateLevel {
	subX := rng.Bounded(uint64(maxX / 4))
	subY := rng.Bounded(uint64(maxY / 4))

	W := maxX - int(subX)
	H := maxY - int(subY)

	tiles := make([]CandidateTile, W*H, W*H)
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			idx := y*W + x
			tiles[idx].X = x
			tiles[idx].Y = y
		}
	}

	level := &CandidateLevel{
		rng:   rng,
		W:     W,
		H:     H,
		flags: flags,
		depth: depth,

		rooms: make(map[int]*Room),
		tiles: tiles,
	}

	level.drawInitialContour()
	level.startBombing()
	level.addContourStairs()

	level.addItems(items.GetCollection("weapons"), MaxWeaponPlacement)
	level.addItems(items.GetCollection("consumeables"), MaxConsumeablePlacement)
	level.addItems(items.GetCollection("armour"), MaxArmourPlacement)

	return level
}

func (c *CandidateLevel) drawInitialContour() {
	{
		// Draw a line between left and right
		var leftX, leftY int
		var rightX, rightY int

		leftX = int(c.rng.Bounded(3) + 1)
		leftY = int(c.rng.Bounded(uint64(c.H)-1)) + 1

		rightX = c.W - 1 - int(c.rng.Bounded(3)+1)
		rightY = int(c.rng.Bounded(uint64(c.H)-1)) + 1

		for _, pos := range PlotLine(leftX, leftY, rightX, rightY) {
			index := pos.Y*c.W + pos.X
			c.tiles[index].TileKind = Floor
		}

		c.stairCandidates = append(c.stairCandidates, Position{X: leftX, Y: leftY})
		c.stairCandidates = append(c.stairCandidates, Position{X: rightX, Y: rightY})
	}
	{
		// Draw a line between top and bottom
		var topX, topY int
		var botX, botY int

		topX = int(c.rng.Bounded(uint64(c.W)-1)) + 1
		topY = int(c.rng.Bounded(3) + 1)

		botX = int(c.rng.Bounded(uint64(c.W)-1)) + 1
		botY = c.H - 1 - int(c.rng.Bounded(3)+1)

		for _, pos := range PlotLine(topX, topY, botX, botY) {
			index := pos.Y*c.W + pos.X
			c.tiles[index].TileKind = Floor
		}

		c.stairCandidates = append(c.stairCandidates, Position{X: topX, Y: topY})
		c.stairCandidates = append(c.stairCandidates, Position{X: botX, Y: botY})
	}
}

func shuffle(rng *pcg.PCG64, n int, swap func(i, j int)) {
	for i := 0; i < n; i++ {
		first, second := int(rng.Bounded(uint64(n))), int(rng.Bounded(uint64(n)))
		swap(first, second)
	}
}

const BombIterationMultiple int = 5
const BombEndProbability uint64 = 3
const BombEndSize int = 10
const BombBoostRadius int = 2
const BombBoostProbability uint64 = 20

func (c *CandidateLevel) startBombing() {
	var tileStack []*CandidateTile

	for i, t := range c.tiles {
		if t.TileKind == Floor {
			tileStack = append(tileStack, &c.tiles[i])
		}
	}

	shuffle(c.rng, len(tileStack), func(i, j int) {
		tileStack[i], tileStack[j] = tileStack[j], tileStack[i]
	})

	iterations := len(tileStack) * BombIterationMultiple

	for i := 0; i < iterations && len(tileStack) > 0; i++ {
		radius := 1

		var tile *CandidateTile
		if c.rng.Bounded(BombBoostProbability) == 0 {
			radius = 2
		}

		var index int
		if c.rng.Bounded(3) == 0 {
			selector := int(c.rng.Bounded(uint64(BombEndSize)))
			fifteenFromEnd := max(0, len(tileStack)-BombEndSize-1)
			index = min(len(tileStack)-1, fifteenFromEnd+selector)
		} else {
			index = int(c.rng.Bounded(uint64(len(tileStack))))
		}

		tile = tileStack[index]
		tileX, tileY := tile.X, tile.Y
		tileStack = append(tileStack[:index], tileStack[index+1:]...)

		minX := max(0, tileX-radius)
		minY := max(0, tileY-radius)
		maxX := min(c.W-1, tileX+radius)
		maxY := min(c.H-1, tileY+radius)

		for y := minY; y <= maxY; y++ {
			for x := minX; x <= maxX; x++ {
				tIndex := y*c.W + x
				if c.tiles[tIndex].TileKind == Floor {
					continue
				}

				if c.tiles[tIndex].X == 0 ||
					c.tiles[tIndex].X == c.W-1 ||
					c.tiles[tIndex].Y == 0 ||
					c.tiles[tIndex].Y == c.H-1 {
					continue
				}

				c.tiles[tIndex].TileKind = Floor

				tileStack = append(tileStack, &c.tiles[tIndex])
			}
		}
	}
}
