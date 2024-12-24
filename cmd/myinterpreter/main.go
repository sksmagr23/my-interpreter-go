package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\r' || char == '\t'
}

func isLetterOrUnderscore(char byte) bool {
	return unicode.IsLetter(rune(char)) || char == '_'
}

func isDigit(char byte) bool {
	return unicode.IsDigit(rune(char))
}

func handleWhitespace(tokens *[]string, i int) {
	// Ignore whitespace
}

func handleSingleCharToken(tokens *[]string, i int, char byte) {
	tokenMap := map[byte]string{
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
	}
	if token, ok := tokenMap[char]; ok {
		*tokens = append(*tokens, token)
	}
}

func handleComparisonOperator(tokens *[]string, i int, char byte, contents []byte) int {
	if i+1 < len(contents) && contents[i+1] == char {
		*tokens = append(*tokens, "EQUAL_EQUAL == null")
		return i + 1
	}
	*tokens = append(*tokens, "EQUAL = null")
	return i
}

func handleStringLiteral(tokens *[]string, i int, char byte, contents []byte) int {
	start := i
	for i+1 < len(contents) && contents[i+1] != char {
		if contents[i+1] == '\n' {
			line++
		}
		i++
	}
	if i+1 >= len(contents) {
		errors = append(errors, fmt.Sprintf("[line %d] Error: Unterminated string.", line))
		return i
	}
	i++
	lexeme := string(contents[start : i+1])
	literal := string(contents[start+1 : i])
	*tokens = append(*tokens, fmt.Sprintf("STRING %s %s", lexeme, literal))
	return i
}

func handleIdentifierOrNumber(tokens *[]string, i int, char byte, contents []byte) int {
	start := i
	if isLetterOrUnderscore(char) {
		for i+1 < len(contents) && (isLetterOrUnderscore(contents[i+1]) || isDigit(contents[i+1])) {
			i++
		}
		lexeme := string(contents[start : i+1])
		if tokenType, ok := reservedKeywords[lexeme]; ok {
			*tokens = append(*tokens, fmt.Sprintf("%s %s null", tokenType, lexeme))
		} else {
			*tokens = append(*tokens, fmt.Sprintf("IDENTIFIER %s null", lexeme))
		}
	} else if isDigit(char) {
		isFloat := false
		for i+1 < len(contents) && isDigit(contents[i+1]) {
			i++
		}
		if i+1 < len(contents) && contents[i+1] == '.' {
			isFloat = true
			i++
			for i+1 < len(contents) && isDigit(contents[i+1]) {
				i++
			}
		}
		lexeme := string(contents[start : i+1])
		literal := lexeme
		if isFloat {
			dotIndex := strings.Index(lexeme, ".")
			integerPart := lexeme[:dotIndex]
			fractionalPart := strings.TrimRight(lexeme[dotIndex+1:], "0")
			literal = integerPart + "." + fractionalPart
		}
		*tokens = append(*tokens, fmt.Sprintf("NUMBER %s %s", lexeme, literal))
	} else {
		errors = append(errors, fmt.Sprintf("[line %d] Error: Unexpected character: %c", line, char))
	}
	return i
}

var (
	reservedKeywords = map[string]string{
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
	errors []string
	line   int
)

func tokenize(filename string) ([]string, []string) {
	contents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	tokens := make([]string, 0)
	errors = []string{}
	line = 1
	for i := 0; i < len(contents); i++ {
		char := contents[i]
		switch char {
		case '\n':
			line++
		case ' ', '\r', '\t':
			handleWhitespace(&tokens, i)
		case '(', ')', '{', '}', ',', '.', '-', '+', ';', '*':
			handleSingleCharToken(&tokens, i, char)
		case '=':
			i = handleComparisonOperator(&tokens, i, char, contents)
		case '!':
			i = handleComparisonOperator(&tokens, i, char, contents)
		case '<':
			i = handleComparisonOperator(&tokens, i, char, contents)
		case '>':
			i = handleComparisonOperator(&tokens, i, char, contents)
		case '/':
			if i+1 < len(contents) && contents[i+1] == '/' {
				for i < len(contents) && contents[i] != '\n' {
					i++
				}
				i--
			} else {
				handleSingleCharToken(&tokens, i, char)
			}
		case '"':
			i = handleStringLiteral(&tokens, i, char, contents)
		default:
			i = handleIdentifierOrNumber(&tokens, i, char, contents)
		}
	}
	tokens = append(tokens, "EOF null")
	return tokens, errors
}

func main() {
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
	tokens, errors := tokenize(filename)
	for _, e := range errors {
		fmt.Fprintln(os.Stderr, e)
	}
	for _, t := range tokens {
		fmt.Println(t)
	}
	if len(errors) > 0 {
		os.Exit(65)
	}
}
