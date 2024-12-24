package main

import (
	"fmt"
	"os"
	"unicode"
)

func main() {
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh <command> <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	filename := os.Args[2]

	contents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "tokenize":
		tokenize(string(contents))
	case "parse":
		tokens, errors := tokenize(string(contents))
		if len(errors) > 0 {
			for _, e := range errors {
				fmt.Fprintln(os.Stderr, e)
			}
			os.Exit(65)
		}
		parse(tokens)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

type Token struct {
	Type    string
	Lexeme  string
	Literal string
	Line    int
}

func tokenize(input string) ([]Token, []string) {
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

	line := 1
	var errors []string
	var tokens []Token
	n := len(input)

	for i := 0; i < n; i++ {
		char := input[i]
		switch char {
		case '\n':
			line++
		case ' ', '\r', '\t':
			// Ignore whitespace
		case '(':
			tokens = append(tokens, Token{"LEFT_PAREN", "(", "", line})
		case ')':
			tokens = append(tokens, Token{"RIGHT_PAREN", ")", "", line})
		case '{':
			tokens = append(tokens, Token{"LEFT_BRACE", "{", "", line})
		case '}':
			tokens = append(tokens, Token{"RIGHT_BRACE", "}", "", line})
		case ',':
			tokens = append(tokens, Token{"COMMA", ",", "", line})
		case '.':
			tokens = append(tokens, Token{"DOT", ".", "", line})
		case '-':
			tokens = append(tokens, Token{"MINUS", "-", "", line})
		case '+':
			tokens = append(tokens, Token{"PLUS", "+", "", line})
		case ';':
			tokens = append(tokens, Token{"SEMICOLON", ";", "", line})
		case '*':
			tokens = append(tokens, Token{"STAR", "*", "", line})
		case '=':
			if i+1 < n && input[i+1] == '=' {
				tokens = append(tokens, Token{"EQUAL_EQUAL", "==", "", line})
				i++
			} else {
				tokens = append(tokens, Token{"EQUAL", "=", "", line})
			}
		case '!':
			if i+1 < n && input[i+1] == '=' {
				tokens = append(tokens, Token{"BANG_EQUAL", "!=", "", line})
				i++
			} else {
				tokens = append(tokens, Token{"BANG", "!", "", line})
			}
		case '<':
			if i+1 < n && input[i+1] == '=' {
				tokens = append(tokens, Token{"LESS_EQUAL", "<=", "", line})
				i++
			} else {
				tokens = append(tokens, Token{"LESS", "<", "", line})
			}
		case '>':
			if i+1 < n && input[i+1] == '=' {
				tokens = append(tokens, Token{"GREATER_EQUAL", ">=", "", line})
				i++
			} else {
				tokens = append(tokens, Token{"GREATER", ">", "", line})
			}
		case '/':
			if i+1 < n && input[i+1] == '/' {
				for i < n && input[i] != '\n' {
					i++
				}
				i--
			} else {
				tokens = append(tokens, Token{"SLASH", "/", "", line})
			}
		case '"':
			start := i
			for i+1 < n && input[i+1] != '"' {
				if input[i+1] == '\n' {
					line++
				}
				i++
			}
			if i+1 >= n {
				errors = append(errors, fmt.Sprintf("[line %d] Error: Unterminated string.", line))
			} else {
				i++
				lex := string(input[start : i+1])
				lit := string(input[start+1 : i])
				tokens = append(tokens, Token{"STRING", lex, lit, line})
			}
		default:
			if unicode.IsLetter(rune(char)) || char == '_' {
				start := i
				for i+1 < n && (unicode.IsLetter(rune(input[i+1])) || unicode.IsDigit(rune(input[i+1])) || input[i+1] == '_') {
					i++
				}
				lex := string(input[start : i+1])
				if tokenType, isReserved := reserved[lex]; isReserved {
					tokens = append(tokens, Token{tokenType, lex, "", line})
				} else {
					tokens = append(tokens, Token{"IDENTIFIER", lex, "", line})
				}
			} else if unicode.IsDigit(rune(char)) {
				start := i
				for i+1 < n && unicode.IsDigit(rune(input[i+1])) {
					i++
				}
				lex := string(input[start : i+1])
				tokens = append(tokens, Token{"NUMBER", lex, lex + ".0", line})
			} else {
				errors = append(errors, fmt.Sprintf("[line %d] Error: Unexpected character: %c", line, char))
			}
		}
	}

	return tokens, errors
}

func parse(tokens []Token) {
	fmt.Println("Parsing tokens...")
	for _, token := range tokens {
		fmt.Printf("%s %s %s\n", token.Type, token.Lexeme, token.Literal)
	}
	fmt.Println("Parse complete.")
}
