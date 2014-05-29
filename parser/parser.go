package parser

import (
	"fmt"
	"github.com/cptaffe/parser/token"
	//"log"
	"strconv"
)

type Parser struct {
	c chan *token.Token
}

type Variable struct {
	Name  string
	Value []token.Token
}

func (p *Parser) Parse(c chan *token.Token) {
	p.c = c
	p.parse()
}

func (p *Parser) parse() {
	for {
		t := <-p.c
		if t == nil {
			//fmt.Printf("got nil.\n")
			return
		}
		if t.Id == token.FLOAT {
			f, _ := strconv.ParseFloat(string(t.Ch), 64)
			fmt.Printf("%f\n", f)
		} else if t.Id == token.INTEGER {
			i, _ := strconv.ParseInt(string(t.Ch), 10, 0)
			fmt.Printf("%d\n", i)
		}
	}
}

/*func (p *Parser) parse() {
	for {
		t := <-p.c
		if t == nil {
			//fmt.Printf("got nil.\n")
			return
		}
		if t.Id == token.FLOAT {
			f, err := strconv.ParseFloat(string(t.Ch), 64) // float 64
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s %s\n", t, stringify(fmt.Sprintf("%f", f)))
		} else if t.Id == token.INTEGER {
			i, err := strconv.ParseInt(string(t.Ch), 10, 0) // int base 10
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s %s\n", t, stringify(fmt.Sprintf("%d", i)))
		} else if t.Id == token.VARIABLE {
			s := string(t.Ch) // string
			fmt.Printf("%s %s\n", t, stringify(s))
		} else if t.Id == token.SET {
			s := string(t.Ch) // string
			fmt.Printf("%s %s\n", t, stringify(s))
		} else if t.Id == token.PLUS || t.Id == token.MINUS || t.Id == token.MULTIPLY || t.Id == token.DIVIDE {
			s := string(t.Ch) // string
			fmt.Printf("%s %s\n", t, stringify(s))
		}
	}
}

func stringify(s string) string {
	return fmt.Sprintf("\x1b[31m[\x1b[0m%s\x1b[31m]\x1b[0m", s)
}*/
