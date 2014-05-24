package scanner

import (
	"errors"
	"fmt"
	"github.com/cptaffe/parser/token"
	"io"
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
}

// New Scanner from file
func NewScanner(f *token.File) (s *Scanner) {
	s = new(Scanner)
	s.file = f
	s.pos = new(token.Pos)
	return
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

// prints errors
func (sc *Scanner) errors(t *token.Token, err error) {
	fmt.Printf("%s \x1b[31merror\x1b[0m: %s\n", t, err)
	return
}

// die dies gracefully
func (sc *Scanner) die(m string) {
	if m != "" {
		sc.errors(token.NewToken(0, sc.ch, sc.pos, sc.pos), errors.New(m))
	}
	sc.c <- nil
	return
}

func (sc *Scanner) lexNumber() (t *token.Token, err error) {
	t = token.NewToken(token.INTEGER, sc.ch, sc.pos, sc.pos)
	for {
		var r rune
		deci := false
		r, err = sc.next()
		if err != nil {
			return
		}
		// lex int
		if '0' <= r && r <= '9' {
			t.Add(r, sc.pos)
			// lex float
		} else if r == '.' {
			if !deci {
				t.Id = token.FLOAT // float
				t.Add(r, sc.pos)
				deci = true
			} else {
				return t, errors.New("more than one decimal")
			}
		} else {
			// do not have to backup() b/c lex next()s at end.
			return
		}
	}
}

func (sc *Scanner) lexCharacter() (t *token.Token, err error) {
	t = token.NewToken(token.VARIABLE, sc.ch, sc.pos, sc.pos)
	for {
		var r rune
		r, err = sc.next()
		if err != nil {
			return
		}
		// lex variable
		if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
			t.Add(r, sc.pos)
		} else {
			return
		}
	}
}

// lex lexes the file
func (sc *Scanner) lex() {
	var r rune // current rune
	for {
		r = sc.ch
		if '0' <= r && r <= '9' {
			t, err := sc.lexNumber()
			if err != nil {
				if err == io.EOF {
					sc.die("eof")
				} else {
					sc.die("reading error")
				}

				return
			} else {
				sc.c <- t
			}
		} else if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
			t, err := sc.lexCharacter()
			if err != nil {
				if err == io.EOF {
					sc.die("eof")
				} else {
					sc.die("reading error")
				}

				return
			} else {
				sc.c <- t
			}
		} else {
			e := errors.New("unexpected rune")
			sc.errors(&token.Token{0, []rune{sc.ch}, sc.pos, sc.pos}, e)
		}

		// get next character as last thing in loop
		rn, err := sc.next()
		r = rn
		if err != nil {
			// reading error
			if err == io.EOF {
				sc.die("eof")
			} else {
				sc.die("reading error")
			}
			return
		}
	}
}

// testing (basically)
func (sc *Scanner) Scan(c chan *token.Token) {
	sc.c = c
	sc.lex()
}
