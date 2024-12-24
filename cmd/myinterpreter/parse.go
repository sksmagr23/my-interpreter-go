package main

import (
	"fmt"
	"strconv"
)

// Token represents a single token with type and value.
type Token struct {
	Type  string
	Value string
}

// Expr is the base interface for all expression types.
type Expr interface {
	Accept(visitor Visitor) string
}

// LiteralExpr represents a literal value (e.g., true, false, nil, or numbers).
type LiteralExpr struct {
	Value interface{}
}

// BinaryExpr represents a binary expression (e.g., 2 + 3).
type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

// Visitor interface for AST printing.
type Visitor interface {
	VisitLiteralExpr(expr *LiteralExpr) string
	VisitBinaryExpr(expr *BinaryExpr) string
}

// AstPrinter implements the Visitor interface for printing the AST.
type AstPrinter struct{}

func (p *AstPrinter) VisitLiteralExpr(expr *LiteralExpr) string {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (p *AstPrinter) VisitBinaryExpr(expr *BinaryExpr) string {
	return fmt.Sprintf("(%s %s %s)",
		expr.Operator.Value,
		expr.Left.Accept(p),
		expr.Right.Accept(p))
}

// Accept method for LiteralExpr.
func (l *LiteralExpr) Accept(visitor Visitor) string {
	return visitor.VisitLiteralExpr(l)
}

// Accept method for BinaryExpr.
func (b *BinaryExpr) Accept(visitor Visitor) string {
	return visitor.VisitBinaryExpr(b)
}

// Parse takes tokens and constructs the AST, then prints it.
func Parse(tokens []Token) {
	if len(tokens) == 1 {
		// Handle single literal values like `true`, `false`, `nil`, or numbers.
		token := tokens[0]
		if token.Type == "TRUE" || token.Type == "FALSE" || token.Type == "NIL" || token.Type == "NUMBER" {
			value := parseLiteralValue(token)
			expr := &LiteralExpr{Value: value}
			printer := &AstPrinter{}
			fmt.Println(expr.Accept(printer))
			return
		}
	}

	if len(tokens) == 3 {
		// Handle binary expressions like `2 + 3`.
		left := parseLiteralValue(tokens[0])
		operator := tokens[1]
		right := parseLiteralValue(tokens[2])

		if operator.Type == "PLUS" {
			expr := &BinaryExpr{
				Left:     &LiteralExpr{Value: left},
				Operator: operator,
				Right:    &LiteralExpr{Value: right},
			}
			printer := &AstPrinter{}
			fmt.Println(expr.Accept(printer))
			return
		}
	}

	fmt.Println("Error: Unable to parse input.")
}

// Helper function to parse literal values from tokens.
func parseLiteralValue(token Token) interface{} {
	switch token.Type {
	case "TRUE":
		return true
	case "FALSE":
		return false
	case "NIL":
		return nil
	case "NUMBER":
		if val, err := strconv.ParseFloat(token.Value, 64); err == nil {
			return val
		}
	}
	return token.Value
}
