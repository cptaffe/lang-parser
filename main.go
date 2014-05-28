package main

import (
	"github.com/cptaffe/parser/parser"
	"github.com/cptaffe/parser/scanner"
	"github.com/cptaffe/parser/token"
	"log"
	"os"
	"sync"
)

func main() {
	// channel
	c := make(chan *token.Token)
	// set up scanner
	f, err := token.NewFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	s := scanner.Init(f)

	// set up parser
	p := new(parser.Parser)
	//Scan and Parse
	var done sync.WaitGroup
	// go lexer
	done.Add(1)
	go func() {
		defer done.Done()
		s.Scan(c)
	}()
	// go parser
	done.Add(1)
	go func() {
		defer done.Done()
		p.Parse(c)
	}()
	done.Wait()
}
