package main

import (
	"fmt"
	"github.com/cptaffe/parser/parser"
	"github.com/cptaffe/parser/scanner"
	"github.com/cptaffe/parser/token"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	// channel
	c := make(chan *token.Token)
	// set up scanner
	f, err := token.NewFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	s := scanner.NewScanner(f)

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
	t0 := time.Now()
	done.Wait()
	fmt.Printf("Waited %v\n", time.Since(t0))
}
