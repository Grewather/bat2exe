package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	eof = iota + 1
	keywords
	mathOperation
	operator
	text
	newline
	leftquote
	rightquote
	leftpercent
	rightpercent
	equals
	and
	identifer
	option
	label
	lparen
	rparen
	condidion
	body
)

func main() {
	fileDirectoryPtr := flag.String("dir", "", "Directory to .bat file")
	flag.Parse()
	dat, err := os.ReadFile(*fileDirectoryPtr)
	check(err)
	tkn := startLexer(string(dat))
	for _, tok := range tkn {
		fmt.Printf("TokenType: %d, Value: '%s'\n", tok.tokenType, tok.value)
	}

	astNode := Parse(tkn)
	generateCode(&astNode)
}
