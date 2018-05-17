package main

import (
	"fmt"
	"regexp"
)

type token struct {
	value     string
	pos       int
	tokenType string
}

func (tok token) String() string {
	return fmt.Sprintf("{%s '%s' %d}", tok.tokenType, tok.value, tok.pos)
}

var baseDictionary = map[string]*regexp.Regexp{	
	"WORD":         regexp.MustCompile("[^\\s]+"),
	"SPACE":       regexp.MustCompile("[\\s]+"),
	"PUNCTUATION": regexp.MustCompile("[\\.,:!]+"),
}

type stateFn func(*lexer) stateFn

type lexer struct {
	dictionary map[string]*regexp.Regexp
	start      int // The position of the last emission
	pos        int // The current position of the lexer
	input      string
	tokens     chan token
	state      stateFn
}

func (l *lexer) next() (val string) {
	if l.pos >= len(l.input) {
		l.pos++
		return ""
	}

	val = l.input[l.pos : l.pos+1]

	l.pos++

	return
}

func (l *lexer) backup() {
	l.pos--
}

func (l *lexer) peek() (val string) {
	val = l.next()

	l.backup()

	return
}

func (l *lexer) emit(t string) {
	val := l.input[l.start:l.pos]
	tok := token{val, l.start, t}
	l.tokens <- tok
	l.start = l.pos
}

func (l *lexer) tokenize() {
	for l.state = lexData; l.state != nil; {
		l.state = l.state(l)
	}
}

func lexData(l *lexer) stateFn {
	s := l.peek()

	if s == "" {
		l.emit("EOF")
		return nil
	}

	if baseDictionary["PUNCTUATION"].MatchString(s) {
		return step("PUNCTUATION", baseDictionary["PUNCTUATION"])
	}

	if baseDictionary["SPACE"].MatchString(s) {
		return step("SPACE", baseDictionary["SPACE"])
	}

	for k, v := range l.dictionary {
		if v.MatchString(s) {
			return step(k, v)
		}
	}

	return step("WORD", baseDictionary["WORD"])
}

func step(t string, rgx *regexp.Regexp) func(l *lexer) stateFn {
	return func(l *lexer) stateFn {
		matched := rgx.FindString(l.input[l.pos:])
		l.pos += len(matched)
		l.emit(t)

		return lexData
	}
}

func newLexer(dictionary map[string]*regexp.Regexp, input string) *lexer {
	return &lexer{dictionary, 0, 0, input, make(chan token), nil}
}
