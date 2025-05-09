package main

import (
	"strings"
)

type Token struct {
	tokenType int
	value     string
	word      string
	line      int
}

type LexerState struct {
	isOption      bool
	isPercent     bool
	isSet         bool
	isA           bool
	countNumbers  bool
	isLabel       bool
	isGoto        bool
	isIfStatement bool
}

func (tkn *Token) checkForTokenType(lexerState *LexerState) {
	var keywordsMap = map[string]int{
		"ver":    keywords,
		"assoc":  keywords,
		"cd":     keywords,
		"cls":    keywords,
		"copy":   keywords,
		"del":    keywords,
		"dir":    keywords,
		"date":   keywords,
		"echo":   keywords,
		"exit":   keywords,
		"md":     keywords,
		"move":   keywords,
		"path":   keywords,
		"pause":  keywords,
		"prompt": keywords,
		"set":    keywords,
		"choice": keywords,
		"goto":   keywords,
		"if":     keywords,
	}
	if val, ok := keywordsMap[tkn.word]; ok {
		switch tkn.word {
		case "set":
			lexerState.isSet = true
		case "goto":
			lexerState.isGoto = true
		case "if":
			lexerState.isIfStatement = true
		}
		tkn.tokenType = val
		tkn.value = strings.TrimSpace(tkn.word)
	}
}

func (tkn *Token) checkForCharType(char string, lexerState *LexerState, j *int, isQuoteOpened *bool) bool {
	switch string(char[*j]) {
	case "%":
		if !lexerState.isPercent {
			tkn.tokenType = leftpercent
			lexerState.isPercent = true
			return true
		} else {
			tkn.tokenType = rightpercent
			lexerState.isPercent = false
			return true
		}
	case "=":
		if !lexerState.isIfStatement {
			tkn.tokenType = equals
			if lexerState.isA {
				lexerState.countNumbers = true
			}
			return true
		}
		if string(char[*j+1]) == "=" {
			tkn.tokenType = operator
			*j += 1
			return true
			// log.Println(tkn.word)
		}
	case "&":
		tkn.tokenType = and
		return true
	case "/":
		tkn.tokenType = option
		lexerState.isOption = true
		return true
	case ":":
		tkn.word = ""
		lexerState.isLabel = true
		return true
	case `"`:
		if !*isQuoteOpened {
			*isQuoteOpened = true
			tkn.tokenType = leftquote
		} else {
			*isQuoteOpened = false
			tkn.tokenType = rightquote
		}
		return true
	case "(":
		tkn.tokenType = lparen
		return true
	case ")":
		tkn.tokenType = rparen
		return true
	}
	return false
}

func appendToken(word string, tokens *[]Token, tokenType, line int) {
	if strings.TrimSpace(word) == "" && tokenType != newline {
		return
	}
	*tokens = append(*tokens, Token{
		value:     strings.TrimSpace(word),
		tokenType: tokenType,
		line:      line,
	})
}

func (tkn *Token) handleChar(splittedLine string, tokens *[]Token, line int, isIfStatement bool) {
	if splittedLine == ":" {
		return
	}
	trimmedWord := strings.TrimSpace(tkn.word)
	if trimmedWord != "" {
		appendToken(trimmedWord, tokens, text, line)
	}
	tkn.value = string(splittedLine)
	if isIfStatement && splittedLine == "=" {
		tkn.value = "=="
	}
	*tokens = append(*tokens, *tkn)
	*tkn = Token{}
}

func (tkn *Token) checkLexerState(lexerState *LexerState, splittedLine string, j int, tokens *[]Token) {
	if lexerState.isSet && string(splittedLine[j+1]) == "=" {
		appendToken(tkn.word, tokens, identifer, j)
		*tkn = Token{}
		lexerState.isSet = false
	}
	if lexerState.isOption && strings.TrimSpace(string(splittedLine[j])) == "A" {
		appendToken(tkn.word, tokens, option, j)
		lexerState.isOption = false
		*tkn = Token{}
		lexerState.isA = true
	}
}

func (tkn *Token) handleOperations(splittedLine string, j *int, line int) bool {
	mapOperations := []string{"/", "*", "-", "+"}
	numb := string(splittedLine[*j])
	for i := range splittedLine {
		if i <= *j {
			continue
		}
		if string(splittedLine[i]) == " " {
			continue
		}
		if IsNumber(string(splittedLine[i])) || numb == "-" {
			if IsNumber(string(splittedLine[i])) {
				numb += string(splittedLine[i])
			} else {
			}
		} else {
			found := false
			for _, v := range mapOperations {
				if string(splittedLine[i]) == v {
					numb += string(splittedLine[i])
					found = true
				}
			}
			if !found && i != len(splittedLine)-1 {
				err := ErrorMsg{
					message: "Unexpected character",
					line:    line + 1,
				}
				err.throw()
			}
		}
	}
	return false
}

func startLexer(file string) []Token {
	var tokens []Token
	splittedFile := strings.Split(file, "\n")
	for i := range splittedFile {
		var lexerState LexerState = LexerState{false, false, false, false, false, false, false, false}
		splittedLine := splittedFile[i]
		var tkn Token
		isQuoteOpened := false
		for j := 0; j < len(splittedLine); j++ {
			if lexerState.countNumbers && tkn.handleOperations(splittedLine, &j, i) {
				continue
			}
			// log.Println(tkn.word)
			isFound := tkn.checkForCharType(splittedLine, &lexerState, &j, &isQuoteOpened)
			if isFound && !lexerState.isOption {
				tkn.handleChar(string(splittedLine[j]), &tokens, j, lexerState.isIfStatement)
				continue
			}
			if splittedLine[j] != ' ' {
				tkn.word += string(splittedLine[j])
			}
			tkn.checkForTokenType(&lexerState)
			if tkn.tokenType == keywords && tkn.value != "" {
				appendToken(tkn.word, &tokens, tkn.tokenType, j)
				tkn = Token{}
				continue
			}
			if !lexerState.isIfStatement {
				tkn.checkLexerState(&lexerState, splittedLine, j, &tokens)
			}
			if j == len(splittedLine)-1 {
				if lexerState.isLabel {
					appendToken(tkn.word, &tokens, label, j)
				} else if lexerState.isGoto {
					appendToken(tkn.word, &tokens, identifer, j)
				} else {
					appendToken(tkn.word, &tokens, text, j)
				}
				if len(tokens) == 0 || tokens[len(tokens)-1].tokenType != newline {
					appendToken("", &tokens, newline, j)
				}
			}
		}
	}
	if tokens[len(tokens)-1].tokenType == newline {
		tokens = tokens[:len(tokens)-1]
	}
	tokens = append(tokens, Token{tokenType: eof})

	return tokens
}
