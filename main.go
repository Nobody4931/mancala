package main

import (
	"fmt"
	"math/rand"
	"time"
)

func printGame(game *Game, zobr *Zobrist) {
	fmt.Printf("\nboard %X\n", zobr.Hash(game))
	fmt.Printf("[  %2d  ]\n", game.Board[HoleCount*2+1])
	for hole := 0; hole < HoleCount; hole++ {
		p1Hole := holeIndex(PlayerOne, hole)
		p2Hole := oppHoleIndex(PlayerOne, hole)
		fmt.Printf("[%2d][%2d]\n", game.Board[p1Hole], game.Board[p2Hole])
	}
	fmt.Printf("[  %2d  ]\n\n", game.Board[HoleCount])
}

func main() {
	rand.Seed(time.Now().UnixMicro())

	game := NewGameWithStones([]int{ 5, 2, 2, 5, 1, 4 }, []int{ 5, 2, 2, 5, 1, 4 })
	zobr := NewZobrist(game)

	for !game.GameOver() {
		printGame(game, zobr)
		recc := game.Minimax()
		var move int
		fmt.Printf("recc: %d\n", recc.Next.Move)
		fmt.Printf("p%d move: ", game.Turn)
		fmt.Scanf("%d\n", &move)
		game.MakeMove(move)
	}

	printGame(game, zobr)
}
