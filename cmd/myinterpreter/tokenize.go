package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

func Tokenize(contents string) []string {
	n := len(contents)
	line := 1
	var tokens []string
	var errors []string

	var reserved = map[string]string{
		"and":    "AND",
		"class":  "CLASS",
		"else":   "ELSE",
		"false":  "FALSE",
		"for":    "FOR",
		"fun":    "FUN",
		"if":     "IF",
		"nil":    "NIL",
		"or":     "OR",
		"print":  "PRINT",
		"return": "RETURN",
		"super":  "SUPER",
		"this":   "THIS",
		"true":   "TRUE",
		"var":    "VAR",
		"while":  "WHILE",
	}

	for i := 0; i < n; i++ {
		char := contents[i]
		switch char {
		case '\n':
			line++
		case ' ', '\r', '\t':
			// Ignore whitespace
		case '(':
			tokens = append(tokens, "LEFT_PAREN ( null")
		case ')':
			tokens = append(tokens, "RIGHT_PAREN ) null")
		case '{':
			tokens = append(tokens, "LEFT_BRACE { null")
		case '}':
			tokens = append(tokens, "RIGHT_BRACE } null")
		case ',':
			tokens = append(tokens, "COMMA , null")
		case '.':
			tokens = append(tokens, "DOT . null")
		case '-':
			tokens = append(tokens, "MINUS - null")
		case '+':
			tokens = append(tokens, "PLUS + null")
		case ';':
			tokens = append(tokens, "SEMICOLON ; null")
		case '*':
			tokens = append(tokens, "STAR * null")
		case '=':
			if i+1 < n && contents[i+1] == '=' {
				tokens = append(tokens, "EQUAL_EQUAL == null")
				i++
			} else {
				tokens = append(tokens, "EQUAL = null")
			}
		case '!':
			if i+1 < n && contents[i+1] == '=' {
				tokens = append(tokens, "BANG_EQUAL != null")
				i++
			} else {
				tokens = append(tokens, "BANG ! null")
			}
		case '<':
			if i+1 < n && contents[i+1] == '=' {
				tokens = append(tokens, "LESS_EQUAL <= null")
				i++
			} else {
				tokens = append(tokens, "LESS < null")
			}
		case '>':
			if i+1 < n && contents[i+1] == '=' {
				tokens = append(tokens, "GREATER_EQUAL >= null")
				i++
			} else {
				tokens = append(tokens, "GREATER > null")
			}
		case '/':
			if i+1 < n && contents[i+1] == '/' {
				for i < n && contents[i] != '\n' {
					i++
				}
				i--
			} else {
				tokens = append(tokens, "SLASH / null")
			}
		case '"':
			start := i
			for i+1 < n && contents[i+1] != '"' {
				if contents[i+1] == '\n' {
					line++
				}
				i++
			}
			if i+1 >= n {
				errors = append(errors, fmt.Sprintf("[line %d] Error: Unterminated string.", line))
			} else {
				i++
				lex := string(contents[start : i+1]) 
				lit := string(contents[start+1 : i]) 
				tokens = append(tokens, fmt.Sprintf("STRING %s %s", lex, lit))
			}
		default:
			if unicode.IsLetter(rune(char)) || char == '_' {
				start := i
				for i+1 < n && (unicode.IsLetter(rune(contents[i+1])) || unicode.IsDigit(rune(contents[i+1])) || contents[i+1] == '_') {
					i++
				}
				lex := string(contents[start : i+1])

				if tokenType, isReserved := reserved[lex]; isReserved {
					tokens = append(tokens, fmt.Sprintf("%s %s null", tokenType, lex))
				} else {
					tokens = append(tokens, fmt.Sprintf("IDENTIFIER %s null", lex))
				}
			} else if unicode.IsDigit(rune(char)) {
				start := i
				isFloat := false

				for i+1 < n && unicode.IsDigit(rune(contents[i+1])) {
					i++
				}

				if i+1 < n && contents[i+1] == '.' {
					isFloat = true
					i++
					for i+1 < n && unicode.IsDigit(rune(contents[i+1])) {
						i++
					}
				}

				lex := string(contents[start : i+1])

				var lit string
				if isFloat {
					lexValue := lex
					if dotIndex := strings.Index(lexValue, "."); dotIndex != -1 {
						integerPart := lexValue[:dotIndex]
						fractionalPart := strings.TrimRight(lexValue[dotIndex+1:], "0")

						if fractionalPart == "" {
							lit = fmt.Sprintf("%s.0", integerPart)
						} else {
							lit = fmt.Sprintf("%s.%s", integerPart, fractionalPart)
						}
					}
				} else {
					lit = fmt.Sprintf("%s.0", lex)
				}

				tokens = append(tokens, fmt.Sprintf("NUMBER %s %s", lex, lit))
			} else {
				errors = append(errors, fmt.Sprintf("[line %d] Error: Unexpected character: %c", line, char))
			}
		}
	}

	for _, e := range errors {
		fmt.Fprintln(os.Stderr, e)
	}

	if len(errors) > 0 {
		os.Exit(65)
	}

	tokens = append(tokens, "EOF  null")
	return tokens
}
