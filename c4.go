package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

//7x6

type Cell struct {
	Color  int
	Name   string
	Below  *Cell
	Left   *Cell
	Right  *Cell
	DLeft  *Cell
	DRight *Cell
	ULeft  *Cell
	URight *Cell
}
type GameBoard struct {
	Player         int
	Board          [][]*Cell
	ComputerPlayer bool
}

const defaultColor = "\x1b[30;1m"

const redColor = "\x1b[31;1m"

const blueColor = "\x1b[34;1m"

func (c Cell) String() string {
	out := fmt.Sprintf("o:%s", c.Name)
	if c.Left != nil {
		out += fmt.Sprintf(",l:%s", c.Left.Name)
	}
	if c.Right != nil {
		out += fmt.Sprintf(",r:%s", c.Right.Name)
	}
	if c.Below != nil {
		out += fmt.Sprintf(",b:%s", c.Below.Name)
	}
	if c.ULeft != nil {
		out += fmt.Sprintf(",ul:%s", c.ULeft.Name)
	}
	if c.URight != nil {
		out += fmt.Sprintf(",ur:%s", c.URight.Name)
	}
	if c.DLeft != nil {
		out += fmt.Sprintf(",dl%s", c.DLeft.Name)
	}
	if c.DRight != nil {
		out += fmt.Sprintf(",dr%s", c.DRight.Name)
	}
	return out
}

func genRow() []*Cell {
	var out []*Cell
	for i := 0; i < 7; i++ {
		cCell := &Cell{Color: 0, Name: strconv.Itoa(i)}
		if i != 0 {
			out[i-1].Right = cCell
			cCell.Left = out[i-1]
		}
		out = append(out, cCell)
	}
	return out
}

func (gb *GameBoard) printBoard() {
	fmt.Println()
	fmt.Println(" 0  1  2  3  4  5  6")
	fmt.Println("---------------------")
	for i := 0; i <= len(gb.Board)-1; i++ {
		for j := 0; j <= len(gb.Board[i])-1; j++ {
			if gb.Board[i][j].Color == 0 {
				fmt.Printf("| |")
			} else if gb.Board[i][j].Color == 1 {
				fmt.Printf("|%sO%s|", redColor, defaultColor)
			} else {
				fmt.Printf("|%sO%s|", blueColor, defaultColor)
			}
			//fmt.Printf("|%d|", gb.Board[i][j].Color)
		}
		fmt.Println("\n---------------------")
	}

}

func (gb *GameBoard) newBoard() {
	gb.Player = 1
	var lastRow []*Cell
	for j := 0; j < 6; j++ {
		row := genRow()
		gb.Board = append(gb.Board, row)
	}
	for j := 0; j <= len(gb.Board)-1; j++ {
		if lastRow != nil {
			for i := 6; i >= 0; i-- {
				lr := lastRow[i]
				cr := gb.Board[j][i]
				cr.Name = cr.Name + "," + strconv.Itoa(j)
				lr.Below = cr
				if i > 0 {
					// connect left diags
					cr.ULeft = lastRow[i-1]
					lastRow[i-1].DRight = cr
				}
				if i < 6 {
					// connect right diags
					cr.URight = lastRow[i+1]
					lastRow[i+1].DLeft = cr
				}
			}
			lastRow = gb.Board[j]
		} else {
			for i := 6; i >= 0; i-- {
				cr := gb.Board[j][i]
				cr.Name = cr.Name + "," + strconv.Itoa(j)
				lastRow = gb.Board[j]
			}
		}
	}
}

func (gb *GameBoard) makePlay(column int64) (*Cell, error) {
	if column < 0 || column > 6 {
		return nil, errors.New("has to be 0-6")
	}
	if gb.Board[0][column].Color > 0 {
		return nil, errors.New("bad play, there is no room there")
	}
	var c *Cell
	//c := &gb.Board[0][column]
	for c = gb.Board[0][column]; c.Below != nil && c.Below.Color == 0; {
		c = c.Below
	}
	c.Color = gb.Player
	return c, nil
}

func (gb *GameBoard) columnFull(column int) bool {
	if gb.Board[5][column].Color == 0 {
		return false
	}
	return true
}

func (gb *GameBoard) takeTurn(r *bufio.Reader) (winner bool, err error) {
	gb.printBoard()
	var column int64
	if gb.Player == 2 && gb.ComputerPlayer {
		// computer makes a pick
		pick := rand.Intn(6)
		for gb.columnFull(pick) {
			pick = rand.Intn(6)
		}
		column = int64(pick)
	} else {
		fmt.Printf("Ready player %d!\n", gb.Player)
		fmt.Print("> ")
		play, _ := r.ReadString('\n')
		play = strings.TrimSpace(play)
		column, err = strconv.ParseInt(play, 10, 64)
	}
	if err != nil {
		return false, nil
	}
	c, err := gb.makePlay(column)
	if err != nil {
		return false, nil
	}
	if checkForWinner(c, gb.Player) {
		return true, nil
	}
	return false, nil
}

func countRight(c *Cell, lookingFor int) int {

	if c.Right != nil && c.Right.Color == lookingFor {
		return countRight(c.Right, lookingFor) + 1
	}
	return 0
}
func countLeft(c *Cell, lookingFor int) int {
	if c.Left != nil && c.Left.Color == lookingFor {

		return countLeft(c.Left, lookingFor) + 1
	}
	return 0
}

func countBelow(c *Cell, lookingFor int) int {
	count := 0
	b := c
	for b.Below != nil && b.Below.Color == lookingFor {
		fmt.Println(b.Color)
		count++
		b = b.Below
	}
	return count
}

func countUL(c *Cell, lookingFor int) int {
	if c.ULeft != nil && c.ULeft.Color == lookingFor {
		return countUL(c.ULeft, lookingFor) + 1
	}
	return 0

}
func countUR(c *Cell, lookingFor int) int {
	if c.URight != nil && c.URight.Color == lookingFor {
		return countUR(c.URight, lookingFor) + 1
	}
	return 0

}

func countDL(c *Cell, lookingFor int) int {
	if c.DLeft != nil && c.DLeft.Color == lookingFor {
		return countDL(c.DLeft, lookingFor) + 1
	}
	return 0

}
func countDR(c *Cell, lookingFor int) int {
	if c.DRight != nil && c.DRight.Color == lookingFor {
		return countDR(c.DRight, lookingFor) + 1
	}
	return 0

}

func checkForWinner(c *Cell, player int) bool {
	right := countRight(c, player)
	left := countLeft(c, player)
	if left+right+1 >= 4 {
		return true
	}
	if countBelow(c, player) >= 3 {
		return true
	}
	if countUL(c, player)+countDR(c, player)+1 >= 4 {
		return true
	}
	if countUR(c, player)+countDL(c, player)+1 >= 4 {
		return true
	}
	return false
}

func c4() {
	gb := GameBoard{Player: 1}
	gb.newBoard()
	rand.Seed(time.Now().UnixNano())
	r := bufio.NewReader(os.Stdin)
	fmt.Println("How many players 1 or 2")
	fmt.Printf(">")
	numPlayers, _ := r.ReadString('\n')
	numPlayers = strings.TrimSpace(numPlayers)
	i, err := strconv.ParseInt(numPlayers, 10, 32)
	if err != nil {
		log.Fatalf("Not a number, needs to be 1 or 2")
	}
	if i < 2 {
		gb.ComputerPlayer = true
	}
	for {
		winner, err := gb.takeTurn(r)
		if err != nil {
			fmt.Println("stalemate", err)
			gb.newBoard()
		}
		if winner {
			gb.printBoard()
			fmt.Printf("Player %d won!\n", gb.Player)
			log.Fatalln("Game Over")
		}
		if gb.Player == 1 {
			gb.Player = 2
		} else {
			gb.Player = 1
		}
	}
}
