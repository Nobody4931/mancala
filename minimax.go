package main

// TODO: Use hashing and transposition tables to optimize scan times

import "fmt" // TEST
import (
	"context"
	"math"
	"math/rand"
	"time"
)

const CalculationTime time.Duration = time.Millisecond * 5000


type Zobrist struct {
	BoardTable [][]uint64
	TurnTable []uint64
}

func NewZobrist(game *Game) *Zobrist {
	totalStones := 0
	for _, stones := range game.Board {
		totalStones += stones
	}

	boardTable := make([][]uint64, BoardSize)
	for hole := range boardTable {
		boardTable[hole] = make([]uint64, totalStones + 1)
		for stones := range boardTable[hole] {
			boardTable[hole][stones] = rand.Uint64()
		}
	}

	turnTable := make([]uint64, 2)
	turnTable[PlayerOne - 1] = rand.Uint64()
	turnTable[PlayerTwo - 1] = rand.Uint64()

	return &Zobrist{
		BoardTable: boardTable,
		TurnTable: turnTable,
	}
}

func (zobr *Zobrist) Hash(game *Game) uint64 {
	var hash uint64 = 0
	for hole, stones := range game.Board {
		hash ^= zobr.BoardTable[hole][stones]
	}
	hash ^= zobr.TurnTable[game.Turn - 1]
	return hash
}


type Transposition struct {
	Hasher *Zobrist
	Table map[uint64]int
}

func NewTransposition(game *Game) *Transposition {
	return &Transposition{
		Hasher: NewZobrist(game),
		Table: make(map[uint64]int),
	}
}

func (trans *Transposition) Get(game *Game) int {
	return trans.Table[trans.Hasher.Hash(game)]
}

func (trans *Transposition) Set(game *Game, score int) {
	trans.Table[trans.Hasher.Hash(game)] = score
}


type Node struct {
	Move int
	Next *Node
}

func (game *Game) Minimax() *Node {
	ctx, cancel := context.WithTimeout(context.Background(), CalculationTime)
	defer cancel()

	totalStones := 0
	for _, stones := range game.Board {
		totalStones += stones
	}

	depth := 1
	var lastRoot *Node
	var lastEval int

	for {
		root := &Node{ Move: -1 }
		eval := game.minimax(ctx, root, depth, -totalStones - 1, totalStones + 1)
		if chanReady(ctx.Done()) {
			break
		}
		lastRoot = root
		lastEval = eval
		depth++
	}

	fmt.Printf("scanned %d moves into the future, eval %d\n", depth - 1, lastEval)
	return lastRoot
}

// Typical minimax algorithm
func (game *Game) minimax(ctx context.Context, node *Node, depth int, alpha, beta int) int {
	if depth == 0 || game.GameOver() {
		return game.score(ctx, alpha, beta)
	}

	bestEval := math.MinInt

	for move := 0; move < HoleCount; move++ {
		if chanReady(ctx.Done()) {
			break
		}

		if !game.CanMove(move) {
			continue
		}

		nextGame := game.Clone()
		nextGame.MakeMove(move)
		nextNode := &Node{ Move: move };

		var eval int
		if nextGame.Turn == game.Turn {
			eval = nextGame.minimax(ctx, nextNode, depth, alpha, beta)
		} else {
			eval = -nextGame.minimax(ctx, nextNode, depth - 1, -beta, -alpha)
		}

		if eval >= beta {
			return beta
		}
		alpha = max(alpha, eval)
		if eval > bestEval {
			bestEval = eval
			node.Next = nextNode
		}
	}

	return bestEval
}

// Score calculation is done through a heavily limited minature version of the minimax
// algorithm that only performs moves that capture the opponent's stones
func (game *Game) score(ctx context.Context, alpha, beta int) int {
	bestEval := game.Board[storeIndex(game.Turn)] - game.Board[storeIndex(opponent(game.Turn))]
	if bestEval >= beta {
		return beta
	}
	alpha = max(alpha, bestEval)

	for move := 0; move < HoleCount; move++ {
		if chanReady(ctx.Done()) {
			break
		}

		hole := holeIndex(game.Turn, move)
		holeIfMoved := hole + game.Board[hole]
		if !game.CanMove(move) || !(holeIfMoved < storeIndex(game.Turn) && game.Board[holeIfMoved] == 0) {
			continue
		}

		nextGame := game.Clone()
		nextGame.MakeMove(move)

		var eval int
		if nextGame.Turn == game.Turn {
			eval = nextGame.score(ctx, alpha, beta)
		} else {
			eval = -nextGame.score(ctx, -beta, -alpha)
		}

		if eval >= beta {
			return beta
		}
		alpha = max(alpha, eval)
		bestEval = max(bestEval, eval)
	}

	return bestEval
}


func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func chanReady[T any](channel <-chan T) bool {
	select {
	case <-channel:
		return true
	default:
		return false
	}
}
