package scanner

import (
	"errors"
	"github.com/cptaffe/parser/token"
)

// Lexer functions must ALWAYS return with the Scanner's current character being the character after the lexed item.

// Lexers that comply to the type lexer: func(Scanner) (t *token.Token, err error)

func lexNumber(sc *Scanner) (t *token.Token, err error) {
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
			sc.c <- t
			return
		}
	}
}

func lexCharacter(sc *Scanner) (t *token.Token, err error) {
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
		} else if r == '(' {
			// is a function call

		} else {
			sc.c <- t
			return
		}
	}
}

func lexSet(sc *Scanner) (t *token.Token, err error) {
	// send set beginning/end
	sc.c <- token.NewToken(token.SET, sc.ch, sc.pos, sc.pos)
	_, err = sc.next()
	if err != nil {
		return
	}

	for r := sc.ch; r != ')'; r = sc.ch {
		// list must be of one type
		if '0' <= r && r <= '9' {
			err = sc.lexType(lexNumber)
			if err != nil {
				return
			}
		} else if r == ' ' {
			r, err = sc.next()
			if err != nil {
				return
			}
		} else {
			return
		}
	}
	// ')' has been found
	sc.c <- token.NewToken(token.SET, sc.ch, sc.pos, sc.pos)
	_, err = sc.next()
	return
}

func lexOperator(sc *Scanner) (t *token.Token, err error) {
	var tok int
	switch sc.ch {
	case '+':
		tok = token.PLUS
	case '-':
		tok = token.MINUS
	case '*':
		tok = token.MULTIPLY
	case '/':
		tok = token.DIVIDE
	}
	sc.c <- token.NewToken(tok, sc.ch, sc.pos, sc.pos)
	_, err = sc.next()
	err = sc.lexType(lexSet)
	return
}

// Lexers that comply to the type lexerStr: func(Scanner, string) (t *token.Token, err error)

type stringLexer func(sc *Scanner, s string) (t *token.Token, err error)
