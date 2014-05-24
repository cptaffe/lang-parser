package token

import (
	"bufio"
	"fmt"
	"os"
)

// Language id constants
const (
	ERR = iota
	begin_types
	INTEGER
	FLOAT
	end_types

	begin_operators
	PLUS     // +()
	MINUS    // -()
	MULTIPLY // *()
	DIVIDE   // /()
	end_operators

	COMMENT
	VARIABLE
)

var token = [...]string{
	ERR:      "parse error",
	INTEGER:  "int",
	FLOAT:    "float",
	VARIABLE: "var",
}

// File
// File for source code
type File struct {
	File *os.File
	Rd   *bufio.Reader
}

func NewFile(fn string) (f *File, err error) {
	f = new(File)
	file, err := os.Open(fn) // opens first argument
	f.File = file
	if err != nil {
		return
	}
	f.Rd = bufio.NewReader(f.File) // buffered reader
	return
}

// Token
// Tokens for lexing
type Token struct {
	Id    int
	Ch    []rune
	Begin *Pos
	End   *Pos
}

func (t *Token) String() (s string) {
	return fmt.Sprintf("%s \x1b[36m%s\x1b[0m(%s)", t.End, token[t.Id], string(t.Ch))
}

func (t *Token) Add(r rune, p *Pos) {
	//fmt.Printf("%c\n", r)
	t.Ch = append(t.Ch, r)
	t.End = p
}

func NewToken(id int, c rune, b *Pos, e *Pos) (t *Token) {
	//r := make([]rune, 10, 100)
	//r[0] = c
	t = &Token{id, []rune{c}, b, e}
	return
}

// Pos
// Simple Positioning w/ line and character struct
type Pos struct {
	Ln int
	Ch int
}

func (p *Pos) String() string {
	return fmt.Sprintf("(%d,%d)", p.Ln+1, p.Ch)
}
