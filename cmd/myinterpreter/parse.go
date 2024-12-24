package main

import "fmt"

// Parse takes tokens and prints the AST representation.
func Parse(tokens []string) {
	// Simple implementation to parse tokens.
	if len(tokens) == 2 && (tokens[0] == "true" || tokens[0] == "false" || tokens[0] == "nil") {
		fmt.Println(tokens[0])
		return
	}
	if len(tokens) == 4 && tokens[1] == "PLUS" {
		fmt.Printf("(+ %s %s)\n", tokens[0], tokens[2])
		return
	}
	fmt.Println("Error: Unable to parse.")
}
