package main

import "log"

type ErrorMsg struct {
	message string
	line    int
}

func (e *ErrorMsg) throw() {
	log.Fatalf("Error %s (line %d)\n", e.message, e.line)
}

func check(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}
