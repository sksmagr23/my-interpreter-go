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
		printUsageAndExit()
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

	tokens, errors := tokenize(string(contents))

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

func printUsageAndExit() {
	fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
	os.Exit(1)
}

func tokenize(input string) ([]string, []string) {
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

	var tokens []string
	var errors []string
	line := 1
	n := len(input)

	for i := 0; i < n; i++ {
		char := input[i]
		switch char {
		case '\n':
			line++
		case ' ', '\r', '\t':
			// Ignore whitespace
		case '(', ')', '{', '}', ',', '.', '-', '+', ';', '*':
			tokens = append(tokens, simpleToken(char))
		case '=', '!', '<', '>':
			i = handleEqualityOperators(input, i, char, &tokens)
		case '/':
			i = handleSlash(input, i, &tokens, line)
		case '"':
			_, err := handleString(input, i, &tokens, line)
			if err != "" {
				errors = append(errors, err)
			}
		default:
			if unicode.IsLetter(rune(char)) || char == '_' {
				i = handleIdentifier(input, i, &tokens, reserved)
			} else if unicode.IsDigit(rune(char)) {
				i = handleNumber(input, i, &tokens)
			} else {
				errors = append(errors, fmt.Sprintf("[line %d] Error: Unexpected character: %c", line, char))
			}
		}
	}

	return tokens, errors
}

func simpleToken(char byte) string {
	symbols := map[byte]string{
		'(': "LEFT_PAREN",
		')': "RIGHT_PAREN",
		'{': "LEFT_BRACE",
		'}': "RIGHT_BRACE",
		',': "COMMA",
		'.': "DOT",
		'-': "MINUS",
		'+': "PLUS",
		';': "SEMICOLON",
		'*': "STAR",
	}
	return fmt.Sprintf("%s %c null", symbols[char], char)
}

func handleEqualityOperators(input string, i int, char byte, tokens *[]string) int {
	n := len(input)
	if i+1 < n && input[i+1] == '=' {
		operator := map[byte]string{
			'=': "EQUAL_EQUAL",
			'!': "BANG_EQUAL",
			'<': "LESS_EQUAL",
			'>': "GREATER_EQUAL",
		}
		*tokens = append(*tokens, fmt.Sprintf("%s %c%c null", operator[char], char, '='))
		return i + 1
	}

	single := map[byte]string{
		'=': "EQUAL",
		'!': "BANG",
		'<': "LESS",
		'>': "GREATER",
	}
	*tokens = append(*tokens, fmt.Sprintf("%s %c null", single[char], char))
	return i
}

func handleSlash(input string, i int, tokens *[]string, line int) int {
	n := len(input)
	if i+1 < n && input[i+1] == '/' {
		// Skip comments
		for i < n && input[i] != '\n' {
			i++
		}
		return i - 1
	}
	*tokens = append(*tokens, "SLASH / null")
	return i
}

func handleString(input string, i int, tokens *[]string, line int) (int, string) {
    start := i
    n := len(input)
    for i+1 < n && input[i+1] != '"' {
        if input[i+1] == '\\' && i+2 < n && (input[i+2] == '"' || input[i+2] == '\\') {
            i++
        }
        if input[i+1] == '\n' {
            line++
        }
        i++
    }

    if i+1 >= n {
        return i, fmt.Sprintf("[line %d] Error: Unterminated string.", line)
    }

    i++ 
    lex := input[start : i+1]
    lit := input[start+1 : i] 
    *tokens = append(*tokens, fmt.Sprintf("STRING %s %s", lex, lit))
    return i, ""
}


func handleIdentifier(input string, i int, tokens *[]string, reserved map[string]string) int {
	start := i
	n := len(input)
	for i+1 < n && (unicode.IsLetter(rune(input[i+1])) || unicode.IsDigit(rune(input[i+1])) || input[i+1] == '_') {
		i++
	}
	lex := input[start : i+1]
	if tokenType, isReserved := reserved[lex]; isReserved {
		*tokens = append(*tokens, fmt.Sprintf("%s %s null", tokenType, lex))
	} else {
		*tokens = append(*tokens, fmt.Sprintf("IDENTIFIER %s null", lex))
	}
	return i
}

func handleNumber(input string, i int, tokens *[]string) int {
	start := i
	n := len(input)
	isFloat := false

	for i+1 < n && unicode.IsDigit(rune(input[i+1])) {
		i++
	}

	if i+1 < n && input[i+1] == '.' {
		isFloat = true
		i++
		for i+1 < n && unicode.IsDigit(rune(input[i+1])) {
			i++
		}
	}

	lex := input[start : i+1]
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
