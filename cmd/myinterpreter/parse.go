package main

import (
	"fmt"
	"strings"
)

func Parse(tokens []string) {
	tokenParts := make([][]string, len(tokens))
	for i, token := range tokens {
		tokenParts[i] = strings.Fields(token)
	}

	if len(tokenParts) == 2 && len(tokenParts[0]) >= 2 {
		if tokenParts[0][0] == "TRUE" || tokenParts[0][0] == "FALSE" || tokenParts[0][0] == "NIL" {
			fmt.Println(strings.ToLower(tokenParts[0][1]))
			return
		}
	}

	if len(tokenParts) == 4 && len(tokenParts[0]) >= 3 && len(tokenParts[2]) >= 3 {
		if tokenParts[1][0] == "PLUS" {
			fmt.Printf("(+ %s %s)\n", tokenParts[0][2], tokenParts[2][2])
			return
		}
	}

	fmt.Println("Error: Unable to parse.")
}
