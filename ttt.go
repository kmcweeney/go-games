package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type PlayField struct {
	Player int
	Board  [][]string
}

func NewPlayField() (pf *PlayField) {
	pf = &PlayField{
		Player: 1,
		Board: [][]string{
			{"1", "2", "3"},
			{"4", "5", "6"},
			{"7", "8", "9"},
		},
	}
	return pf
}

func (pf *PlayField) printBoard() {
	fmt.Println()
	fmt.Println("-------")
	for _, row := range pf.Board {
		fmt.Printf("|%s|%s|%s|\n", row[0], row[1], row[2])
		fmt.Println("-------")
	}

}

func (pf *PlayField) makePlay(pos string) error {
	for i, row := range pf.Board {
		for j := range row {
			if pf.Board[i][j] == pos {
				if pf.Player == 1 {
					pf.Board[i][j] = "X"
					return nil
				}
				pf.Board[i][j] = "O"
				return nil
			}
		}
	}
	return fmt.Errorf("%s is not a legal play", pos)
}

func (pf *PlayField) takeTurn(r *bufio.Reader) (winner bool, err error) {
	pf.printBoard()
	fmt.Printf("Ready player %d!\n", pf.Player)
	fmt.Print("> ")
	text, _ := r.ReadString('\n')
	text = strings.TrimSpace(text)
	err = pf.makePlay(text)
	if err != nil {
		return false, nil
	}
	winner, err = pf.checkForWinner()
	if err != nil {
		return false, err
	}
	if winner {
		return true, nil
	}

	if pf.Player == 1 {
		pf.Player = 2
		return false, nil
	}
	pf.Player = 1
	return false, nil
}

func (pf *PlayField) checkForWinner() (winner bool, err error) {

	var rows string
	var cols string
	var diags string
	for i, row := range pf.Board {
		for j := range row {
			rows = rows + pf.Board[i][j]
			cols = cols + pf.Board[j][i]
		}
		rows = rows + ","
		cols = cols + ","
	}
	diags = diags + pf.Board[0][0] + pf.Board[1][1] + pf.Board[2][2] + ","
	diags = diags + pf.Board[0][2] + pf.Board[1][1] + pf.Board[2][0]
	if strings.Contains(rows, "XXX") || strings.Contains(rows, "OOO") {
		return true, nil
	}
	if strings.Contains(cols, "XXX") || strings.Contains(cols, "OOO") {
		return true, nil
	}
	if strings.Contains(diags, "XXX") || strings.Contains(diags, "OOO") {
		return true, nil
	}
	all := rows + cols
	re := regexp.MustCompile("[0-9]+")
	if re.MatchString(all) {
		return false, nil
	}
	return false, errors.New("no moves left")
}

func ticTacToe() {
	pf := NewPlayField()
	r := bufio.NewReader(os.Stdin)
	for {
		winner, err := pf.takeTurn(r)
		if err != nil {
			fmt.Println("No more moves remaining, starting new game!")
			pf = NewPlayField()
		}
		if winner {
			pf.printBoard()
			fmt.Printf("Player %d won!\n", pf.Player)
			fmt.Println("Enter p to play again, anything else to quit")
			t, _ := r.ReadString('\n')
			t = strings.TrimSpace(t)
			if t == "p" {
				pf = NewPlayField()
			} else {
				return
			}
		}
	}
}
