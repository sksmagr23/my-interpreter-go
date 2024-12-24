package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

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

func main() {
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	if command != "tokenize" {
		printUnknownCommand(command)
		os.Exit(1)
	}

	filename := os.Args[2]
	contents, err := readFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	tokens, errors := tokenize(contents)

	for _, e := range errors {
		fmt.Fprintln(os.Stderr, e)
	}
	for _, t := range tokens {
		fmt.Println(t)
	}
	fmt.Println("EOF  null")

	if len(errors) > 0 {
		os.Exit(65)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
}

func printUnknownCommand(command string) {
	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
}

func readFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func tokenize(contents []byte) ([]string, []string) {
	line := 1
	n := len(contents)
	var errors []string
	var tokens []string

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
			i = handleDoubleCharToken(contents, i, '=', "EQUAL_EQUAL == null", "EQUAL = null", &tokens)
		case '!':
			i = handleDoubleCharToken(contents, i, '=', "BANG_EQUAL != null", "BANG ! null", &tokens)
		case '<':
			i = handleDoubleCharToken(contents, i, '=', "LESS_EQUAL <= null", "LESS < null", &tokens)
		case '>':
			i = handleDoubleCharToken(contents, i, '=', "GREATER_EQUAL >= null", "GREATER > null", &tokens)
		case '/':
			i = handleSlash(contents, i, &tokens, &line)
		case '"':
			_, err := handleString(contents, i, &tokens, &line)
			if err != "" {
				errors = append(errors, err)
			}
		default:
			if unicode.IsLetter(rune(char)) || char == '_' {
				i = handleIdentifierOrKeyword(contents, i, &tokens)
			} else if unicode.IsDigit(rune(char)) {
				i = handleNumber(contents, i, &tokens)
			} else {
				errors = append(errors, fmt.Sprintf("[line %d] Error: Unexpected character: %c", line, char))
			}
		}
	}

	return tokens, errors
}

func handleDoubleCharToken(contents []byte, i int, nextChar byte, doubleToken, singleToken string, tokens *[]string) int {
	if i+1 < len(contents) && contents[i+1] == nextChar {
		*tokens = append(*tokens, doubleToken)
		return i + 1
	}
	*tokens = append(*tokens, singleToken)
	return i
}

func handleSlash(contents []byte, i int, tokens *[]string, line *int) int {
	if i+1 < len(contents) && contents[i+1] == '/' {
		for i < len(contents) && contents[i] != '\n' {
			i++
		}
		return i - 1
	}
	*tokens = append(*tokens, "SLASH / null")
	return i
}

func handleString(contents []byte, i int, tokens *[]string, line *int) (int, string) {
	start := i
	for i+1 < len(contents) && contents[i+1] != '"' {
		if contents[i+1] == '\n' {
			(*line)++
		}
		i++
	}
	if i+1 >= len(contents) {
		return i, fmt.Sprintf("[line %d] Error: Unterminated string.", *line)
	}
	i++
	lex := string(contents[start : i+1])
	lit := string(contents[start+1 : i])
	*tokens = append(*tokens, fmt.Sprintf("STRING %s %s", lex, lit))
	return i, ""
}

func handleIdentifierOrKeyword(contents []byte, i int, tokens *[]string) int {
	start := i
	for i+1 < len(contents) && (unicode.IsLetter(rune(contents[i+1])) || unicode.IsDigit(rune(contents[i+1])) || contents[i+1] == '_') {
		i++
	}
	lex := string(contents[start : i+1])
	if tokenType, isReserved := reserved[lex]; isReserved {
		*tokens = append(*tokens, fmt.Sprintf("%s %s null", tokenType, lex))
	} else {
		*tokens = append(*tokens, fmt.Sprintf("IDENTIFIER %s null", lex))
	}
	return i
}

func handleNumber(contents []byte, i int, tokens *[]string) int {
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
	lit := lex
	if isFloat {
		if dotIndex := strings.Index(lex, "."); dotIndex != -1 {
			integerPart := lex[:dotIndex]
			fractionalPart := strings.TrimRight(lex[dotIndex+1:], "0")
			if fractionalPart == "" {
				lit = fmt.Sprintf("%s.0", integerPart)
			} else {
				lit = fmt.Sprintf("%s.%s", integerPart, fractionalPart)
			}
		}
	} else {
		lit = fmt.Sprintf("%s.0", lex)
	}
	*tokens = append(*tokens, fmt.Sprintf("NUMBER %s %s", lex, lit))
	return i
}
