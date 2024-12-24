package main

import (
	"fmt"
	"os"
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

	line := 1
	var errors []string
	var tokens []string

	for i := 0; i < len(contents); i++ {
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
			if i+1 < len(contents) && contents[i+1] == '=' {
				tokens = append(tokens, "EQUAL_EQUAL == null")
				i++
			} else {
				tokens = append(tokens, "EQUAL = null")
			}
		case '!':
			if i+1 < len(contents) && contents[i+1] == '=' {
				tokens = append(tokens, "BANG_EQUAL != null")
				i++
			} else {
				tokens = append(tokens, "BANG ! null")
			}
		case '<':
			if i+1 < len(contents) && contents[i+1] == '=' {
				tokens = append(tokens, "LESS_EQUAL <= null")
				i++
			} else {
				tokens = append(tokens, "LESS < null")
			}
		case '>':
			if i+1 < len(contents) && contents[i+1] == '=' {
				tokens = append(tokens, "GREATER_EQUAL >= null")
				i++
			} else {
				tokens = append(tokens, "GREATER > null")
			}
		case '/':
			if i+1 < len(contents) && contents[i+1] == '/' {
				for i < len(contents) && contents[i] != '\n' {
					i++
				}
				i--
			} else {
				tokens = append(tokens, "SLASH / null")
			}
		case '"':
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
				tokens = append(tokens, fmt.Sprintf("STRING %s %s", lex, lit))
			}
		default:
			if unicode.IsDigit(rune(char)) {
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

				lexeme := string(contents[start : i+1])

				var literal string
				if isFloat {
					var floatValue float64
					fmt.Sscanf(lexeme, "%f", &floatValue)
					literal = fmt.Sprintf("%g", floatValue)
				} else {
					literal = fmt.Sprintf("%s.0", lexeme)
				}

				tokens = append(tokens, fmt.Sprintf("NUMBER %s %s", lexeme, literal))
			} else {
				errors = append(errors, fmt.Sprintf("[line %d] Error: Unexpected character: %c", line, char))
			}
		}
	}

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
