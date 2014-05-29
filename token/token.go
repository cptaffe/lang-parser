package token

import (
	"bufio"
	"fmt"
	"os"
)

// Language id constants
const (
	ERR = iota
	// types
	begin_types
	INTEGER
	FLOAT
	end_types

	// type collections
	begin_typeCols
	SET // (innards)
	end_typeCols

	begin_operators
	PLUS     // +
	MINUS    // -
	MULTIPLY // *
	DIVIDE   // /
	end_operators

	COMMENT
	VARIABLE
)

var token = [...]string{
	// error
	ERR: "parse error",
	// types
	INTEGER: "int",
	FLOAT:   "float",
	// type sets
	SET: "set",
	// variable
	VARIABLE: "var",
	// operators
	PLUS:     "add",
	MINUS:    "sub",
	MULTIPLY: "mult",
	DIVIDE:   "div",
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
	Begin Pos
	End   Pos
}

func (t *Token) String() (s string) {
	return fmt.Sprintf("%s \x1b[36m%s\x1b[0m(%s)", t.Begin, token[t.Id], string(t.Ch))
}

func (t *Token) Add(r rune, p *Pos) {
	t.Ch = append(t.Ch, r)
	t.End = *p
}

func NewToken(id int, c rune, b *Pos, e *Pos) (t *Token) {
	t = &Token{id, []rune{c}, *b, *e}
	return
}

// Pos
// Simple Positioning w/ line and character struct
type Pos struct {
	Ln int
	Ch int
}

func (p Pos) String() string {
	return fmt.Sprintf("(%d,%d)", p.Ln, p.Ch)
}
