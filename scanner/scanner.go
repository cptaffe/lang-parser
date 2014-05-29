package scanner

import (
	"errors"
	"fmt"
	"github.com/cptaffe/parser/token"
	"io"
	"os"
	"unicode/utf8"
)

type Scanner struct {
	// immutable state
	file *token.File
	c    chan *token.Token

	// scanning state
	ch  rune
	pos *token.Pos
	src []byte

	// errors
	errors errorHandler
}

// New Scanner from file
func Init(f *token.File) (s *Scanner) {
	s = new(Scanner)
	s.file = f
	s.pos = new(token.Pos)
	s.errors = printError
	return
}

type errorHandler func(t *token.Token, err error)

// prints errors
// default errorHandler
func printError(t *token.Token, err error) {
	fmt.Printf("%s \x1b[31merror\x1b[0m: %s\n", t, err)
	return
}

// unrecoverable error
func (sc *Scanner) die(err error) {
	if err != io.EOF {
		sc.errors(token.NewToken(0, sc.ch, sc.pos, sc.pos), err)
	}
	//sc.c <- nil
	os.Exit(1)
}

// returns next character
// line by line management
func (sc *Scanner) next() (r rune, err error) {
	if len(sc.src) > 1 && sc.ch != '\n' && len(sc.src) > sc.pos.Ch {
		sc.pos.Ch++
		sc.ch = r
	} else {
		src, _, e := sc.file.Rd.ReadLine()
		err = e
		if err != nil {
			return r, err
		}
		sc.pos.Ln++
		sc.pos.Ch = 0
		sc.src = src
	}
	sc.ch, _ = utf8.DecodeRune(sc.src[sc.pos.Ch:])
	return sc.ch, nil
}

// returns previous character
func (sc *Scanner) backup() (r rune, err error) {
	if sc.ch > 0 {
		sc.pos.Ch--
		sc.ch, _ = utf8.DecodeRune(sc.src[sc.pos.Ch:])
		return sc.ch, nil
	} else if sc.pos.Ch == 0 {
		return '\n', nil
	} else {
		return sc.ch, errors.New("cannot backup")
	}
}

// peeks at next character
func (sc *Scanner) peek() (r rune, err error) {
	r, err = sc.next()
	if err != nil {
		return
	}
	_, err = sc.backup()
	if err != nil {
		return
	}
	return
}

// lexer type
type lexer func(sc *Scanner) (t *token.Token, err error)

func (sc *Scanner) lexType(l lexer) (err error) {
	_, err = l(sc)
	if err != nil {
		sc.die(err)
	}
	return
}

// lex lexes the file
func (sc *Scanner) lex() {
	sc.next()
	for {
		r := sc.ch
		if r == 0xFFFD || r == ' ' {
			// unknown character, ignore.
			_, err := sc.next()
			if err != nil {
				sc.die(err)
				return
			}
		} else if '0' <= r && r <= '9' {
			// numbers
			err := sc.lexType(lexNumber)
			if err != nil {
				return
			}
		} else if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
			// characters
			err := sc.lexType(lexCharacter)
			if err != nil {
				return
			}
		} else if r == '(' {
			// in a set
			err := sc.lexType(lexSet)
			if err != nil {
				return
			}
		} else if r == '+' || r == '-' || r == '*' || r == '/' {
			// characters
			err := sc.lexType(lexOperator)
			if err != nil {
				return
			}
		} else {
			e := errors.New("unexpected rune")
			sc.errors(&token.Token{token.ERR, []rune{sc.ch}, *sc.pos, *sc.pos}, e)
			_, err := sc.next()
			if err != nil {
				sc.die(err)
				return
			}
		}
	}
}

// testing (basically)
func (sc *Scanner) Scan(c chan *token.Token) {
	sc.c = c
	sc.lex()
}
