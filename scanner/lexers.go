package scanner

import (
	"errors"
	"github.com/cptaffe/parser/token"
)

// Lexer functions must ALWAYS return with the Scanner's current character being the character after the lexed item.

// Lexers that comply to the type lexer: func(Scanner) (t *token.Token, err error)

func lexNumber(sc *Scanner) (err error) {
	t := token.NewToken(token.INTEGER, sc.ch, sc.pos, sc.pos)
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
				return errors.New("more than one decimal")
			}
		} else {
			sc.c <- t
			return
		}
	}
}

func lexCharacter(sc *Scanner) (err error) {
	t := token.NewToken(token.VARIABLE, sc.ch, sc.pos, sc.pos)
	for {
		var r rune
		r, err = sc.next()
		if err != nil {
			return
		}
		// lex variable
		if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
			t.Add(r, sc.pos)
		} else if r == '(' {
			// is a function call

		} else {
			sc.c <- t
			return
		}
	}
}

func lexList(sc *Scanner) (err error) {
	// send list beginning/end
	sc.c <- token.NewToken(token.LIST, sc.ch, sc.pos, sc.pos)
	_, err = sc.next()
	return
}

// Lexers that comply to the type lexerStr: func(Scanner, string) (t *token.Token, err error)

type stringLexer func(sc *Scanner, s string) (t *token.Token, err error)
