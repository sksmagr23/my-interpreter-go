package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	contents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	processContents(string(contents))
}

func processContents(contents string) {
	reserved := map[string]string{
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

	line := 1
	tokens := []string{}
	errors := []string{}
	i := 0
	n := len(contents)

	for i < n {
		char := contents[i]
		switch char {
		case '\n':
			line++
		case ' ', '\r', '\t':
			// Ignore whitespace
		case '(', ')', '{', '}', ',', '.', '-', '+', ';', '*':
			tokens = append(tokens, singleCharToken(char))
		case '=', '!', '<', '>':
			i, tokens = appendTwoCharTokens(i, contents, char, tokens)
		case '/':
			i = handleSlash(i, contents, &tokens)
		case '"':
			i, errors = handleString(i, contents, line, &tokens, errors)
		default:
			if unicode.IsLetter(rune(char)) || char == '_' {
				i, tokens = handleIdentifier(i, contents, reserved, tokens)
			} else if unicode.IsDigit(rune(char)) {
				i, tokens = handleNumber(i, contents, tokens)
			} else {
				errors = append(errors, fmt.Sprintf("[line %d] Error: Unexpected character: %c", line, char))
			}
		}
		if char == '\n' {
			line++ // Increment after processing newline
		}
		i++
	}

	printErrorsAndTokens(errors, tokens)
	if len(errors) > 0 {
		os.Exit(65)
	}
}

func singleCharToken(char byte) string {
	return map[byte]string{
		'(': "LEFT_PAREN ( null",
		')': "RIGHT_PAREN ) null",
		'{': "LEFT_BRACE { null",
		'}': "RIGHT_BRACE } null",
		',': "COMMA , null",
		'.': "DOT . null",
		'-': "MINUS - null",
		'+': "PLUS + null",
		';': "SEMICOLON ; null",
		'*': "STAR * null",
	}[char]
}

func appendTwoCharTokens(i int, contents string, char byte, tokens []string) (int, []string) {
	nextChar := map[byte]string{
		'=': "EQUAL_EQUAL",
		'!': "BANG_EQUAL",
		'<': "LESS_EQUAL",
		'>': "GREATER_EQUAL",
	}[char]

	token := map[byte]string{
		'=': "EQUAL",
		'!': "BANG",
		'<': "LESS",
		'>': "GREATER",
	}[char]

	if i+1 < len(contents) && contents[i+1] == '=' {
		tokens = append(tokens, fmt.Sprintf("%s %c= null", nextChar, char))
		i++
	} else {
		tokens = append(tokens, fmt.Sprintf("%s %c null", token, char))
	}

	return i, tokens
}

func handleSlash(i int, contents string, tokens *[]string) int {
	if i+1 < len(contents) && contents[i+1] == '/' {
		for i < len(contents) && contents[i] != '\n' {
			i++
		}
	} else {
		*tokens = append(*tokens, "SLASH / null")
	}
	return i
}

func handleString(i int, contents string, line int, tokens *[]string, errors []string) (int, []string) {
	start := i
	for i+1 < len(contents) && contents[i+1] != '"' {
		if contents[i+1] == '\n' {
			line++
		}
		i++
	}

	if i+1 >= len(contents) {
		errors = append(errors, fmt.Sprintf("[line %d] Error: Unterminated string.", line))
	} else {
		i++
		lex := string(contents[start : i+1])
		lit := string(contents[start+1 : i])
		*tokens = append(*tokens, fmt.Sprintf("STRING %s %s", lex, lit))
	}

	return i, errors
}

func handleIdentifier(i int, contents string, reserved map[string]string, tokens []string) (int, []string) {
	start := i
	for i+1 < len(contents) && (unicode.IsLetter(rune(contents[i+1])) || unicode.IsDigit(rune(contents[i+1])) || contents[i+1] == '_') {
		i++
	}
	lex := string(contents[start : i+1])
	if tokenType, isReserved := reserved[lex]; isReserved {
		tokens = append(tokens, fmt.Sprintf("%s %s null", tokenType, lex))
	} else {
		tokens = append(tokens, fmt.Sprintf("IDENTIFIER %s null", lex))
	}
	return i, tokens
}

func handleNumber(i int, contents string, tokens []string) (int, []string) {
	start := i
	isFloat := false

	for i+1 < len(contents) && unicode.IsDigit(rune(contents[i+1])) {
		i++
	}

	if i+1 < len(contents) && contents[i+1] == '.' {
		isFloat = true
		i++
		for i+1 < len(contents) && unicode.IsDigit(rune(contents[i+1])) {
			i++
		}
	}

	lex := string(contents[start : i+1])
	var lit string
	if isFloat {
		lit = strings.TrimRight(lex, "0")
		if strings.HasSuffix(lit, ".") {
			lit += "0"
		}
	} else {
		lit = fmt.Sprintf("%s.0", lex)
	}

	tokens = append(tokens, fmt.Sprintf("NUMBER %s %s", lex, lit))
	return i, tokens
}

func printErrorsAndTokens(errors, tokens []string) {
	for _, e := range errors {
		fmt.Fprintln(os.Stderr, e)
	}
	for _, t := range tokens {
		fmt.Println(t)
	}
	fmt.Println("EOF  null")
}
