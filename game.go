package main

const HoleCount int = 6
const BoardSize int = HoleCount * 2 + 2


type Player uint8
const (
	PlayerOne Player = 0b01
	PlayerTwo Player = 0b10
	PlayerXor Player = 0b11
)

type Game struct {
	Turn Player
	Board []int
}

func NewGame() *Game {
	return &Game{
		Turn: PlayerOne,
		Board: make([]int, BoardSize),
	}
}

func NewGameWithStones(p1Stones, p2Stones []int) *Game {
	game := NewGame()
	copy(game.Board[0:HoleCount], p1Stones)
	copy(game.Board[HoleCount+1:BoardSize-1], p2Stones)
	return game
}

func NewGameWithStarter(stones int) *Game {
	pStones := make([]int, HoleCount)
	for i := range pStones {
		pStones[i] = stones
	}
	return NewGameWithStones(pStones, pStones)
}


func (game *Game) Clone() *Game {
	board := make([]int, len(game.Board))
	copy(board, game.Board)
	return &Game{
		Turn: game.Turn,
		Board: board,
	}
}

func (game *Game) CanMove(hole int) bool {
	hole = holeIndex(game.Turn, hole)
	return game.Board[hole] > 0
}

func (game *Game) MakeMove(hole int) bool {
	if !game.CanMove(hole) {
		return false
	}

	hole = holeIndex(game.Turn, hole)
	stones := game.Board[hole]
	oppStore := storeIndex(opponent(game.Turn))

	// Pick up the stones
	game.Board[hole] = 0

	// Place one stone in every consecutive hole (skipping the opponent's store)
	// until there are no stones remaining
	for stones > 0 {
		hole = (hole + 1) % BoardSize
		if hole == oppStore {
			continue
		}
		game.Board[hole]++
		stones--
	}

	// If the last stone was placed in an empty hole on the player's side and
	// there are stones on the opponent's side of the board, capture all the stones
	if holeNum, holeSide := holeNumber(hole); holeNum != HoleCount && holeSide == game.Turn {
		store := storeIndex(game.Turn)
		oppHole := oppHoleIndex(game.Turn, holeNum)
		if game.Board[hole] == 1 && game.Board[oppHole] > 0 {
			game.Board[store] += game.Board[hole] + game.Board[oppHole]
			game.Board[hole] = 0
			game.Board[oppHole] = 0
		}
	}

	// If there are no more turns to be made, then each player captures all the
	// stones on their side of the board and the game ends
	if game.GameOver() {
		p1Store := storeIndex(PlayerOne)
		p2Store := storeIndex(PlayerTwo)
		for hole := 0; hole < HoleCount; hole++ {
			p1Hole := holeIndex(PlayerOne, hole)
			p2Hole := holeIndex(PlayerTwo, hole)
			game.Board[p1Store] += game.Board[p1Hole]
			game.Board[p2Store] += game.Board[p2Hole]
			game.Board[p1Hole] = 0
			game.Board[p2Hole] = 0
		}
	}

	// If the last stone was placed in the active player's store, gain an extra turn
	if hole == storeIndex(game.Turn) {
		return true
	}

	game.Turn = opponent(game.Turn)
	return true
}

func (game *Game) GameOver() bool {
	movesRemaining := func(player Player) bool {
		for hole := 0; hole < HoleCount; hole++ {
			if game.Board[holeIndex(player, hole)] > 0 {
				return true
			}
		}
		return false
	}
	return !movesRemaining(PlayerOne) || !movesRemaining(PlayerTwo)
}


func opponent(player Player) Player {
	return player ^ PlayerXor
}

func holeIndex(player Player, hole int) int {
	return int(player - 1) * (HoleCount + 1) + hole
}

func storeIndex(player Player) int {
	return holeIndex(player, HoleCount)
}

func oppHoleIndex(player Player, hole int) int {
	return holeIndex(opponent(player), HoleCount - 1 - hole)
}

func holeNumber(holeIdx int) (int, Player) {
	if holeIdx <= HoleCount {
		return holeIdx, PlayerOne
	} else {
		return holeIdx - (HoleCount + 1), PlayerTwo
	}
}
