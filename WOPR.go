package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const boardPrint = "┌─┬─┬─┐\n│%c│%c│%c│\n├─┼─┼─┤\n│%c│%c│%c│\n├─┼─┼─┤\n│%c│%c│%c│\n└─┴─┴─┘\n"

var combos = [][]int{
	{0, 1, 2},
	{3, 4, 5},
	{6, 7, 8},
	{0, 3, 6},
	{1, 4, 7},
	{2, 5, 8},
	{0, 4, 8},
	{2, 4, 6},
}

type Game struct {
	board []rune
	w     io.Writer
	r     io.Reader
	turnX bool
}

func NewGame(w io.Writer, r io.Reader) *Game {
	g := &Game{
		board: []rune("         "),
		w:     w,
		r:     r,
		turnX: false,
	}
	return g
}

func (g *Game) Boot() {
	fmt.Fprint(g.w, "BOOTING WOPR")
	for i := 0; i < 6; i++ {
		fmt.Fprint(g.w, ".")
		time.Sleep(time.Millisecond * 10)
	}
	time.Sleep(time.Second)
	fmt.Fprintf(g.w, " COMPLETE\n")
	g.Log("[WOPR] Would you like to play a game?\n")
	g.Print()
}

func (g *Game) Reset() {
	g.board = []rune("         ")
}

func (g *Game) Print() {
	g.Log(boardPrint, g.board[0], g.board[1], g.board[2], g.board[3], g.board[4], g.board[5], g.board[6], g.board[7], g.board[8])
}

func (g *Game) Log(thing string, args ...interface{}) {
	str := fmt.Sprintf(thing, args...)

	for _, ch := range str {
		fmt.Fprintf(g.w, "%c", ch)
		time.Sleep(time.Millisecond * 10)
	}
	time.Sleep(time.Millisecond * 10)
}
func (g *Game) GoRandom() error {
	sym := 'X'
	if !g.turnX {
		sym = 'O'
	}
	options := make([]int, 0, 9)
	for i, chr := range g.board {
		if chr == ' ' {
			options = append(options, i)
		}
	}
	if len(options) == 0 {
		return fmt.Errorf(">> game over: cat-scratch")
	}
	i, _ := rand.Int(rand.Reader, big.NewInt(int64(len(options))))
	g.Log("[WOPR] It's my turn")
	time.Sleep(time.Millisecond * 20)
	g.Log(". Done\n")
	g.board[options[i.Int64()]] = sym
	if g.CheckWin() {
		return fmt.Errorf(">> game over: WOPR wins")
	}
	return nil
}

func (g *Game) GoPlayer() error {
	sym := 'X'
	if !g.turnX {
		sym = 'O'
	}
	bufReader := bufio.NewReader(g.r)
	for {
		fmt.Fprintf(g.w, "<< enter move [0-8]: ")
		input, err := bufReader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(g.w, "Error reading the content:", err)
		}

		move, err := strconv.Atoi(strings.Trim(input, "\n \t"))
		if err != nil || move > 8 || g.board[move] != ' ' {
			continue
		}
		g.board[move] = sym
		break
	}
	if g.CheckWin() {
		return fmt.Errorf(">> game over: player wins")
	}
	return nil
}

func (g *Game) CheckWin() bool {
	// Check win conditions
	for _, c := range combos {
		if g.board[c[0]] == g.board[c[1]] && g.board[c[1]] == g.board[c[2]] && g.board[c[1]] != ' ' {
			return true
		}
	}
	g.turnX = !g.turnX
	return false
}

func (g *Game) Reward() {
	cmd := exec.Command("/bin/sh", "-i")
	cmd.Stdout = g.w
	cmd.Stderr = g.w
	cmd.Stdin = g.r
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(g.w, "error giving reward: %s", err)
	}
}

func (g *Game) Play() {
	g.Boot()
	for {
		// WOPR goes first
		if isOver := g.GoRandom(); isOver != nil {
			g.Log(isOver.Error() + "\n")
			g.Log("[WOPR] let's play again\n")
			g.Reset()
		}
		g.Print()
		// Player goes
		if isOver := g.GoPlayer(); isOver != nil {
			g.Log(isOver.Error() + "\n")
			g.Log("[WOPR] Good job, you have bested me...\n")
			g.Reward()
			break
		}
	}
}

func main() {
	serv, err := net.Listen("tcp", "0.0.0.0:4000")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer serv.Close()
	fmt.Fprintf(os.Stderr, "Listening on 0.0.0.0:4000\n")
	for {
		conn, err := serv.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error recieving connection: %s\n", err)
			continue
		}
		fmt.Fprintf(os.Stderr, "starting game for: %s\n", conn.RemoteAddr())
		game := NewGame(conn, conn)
		go func(c net.Conn) {
			defer fmt.Println("Connection closed", c.RemoteAddr())
			game.Play()
		}(conn)
	}
}
