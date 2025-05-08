package main

import "fmt"

type AstNode struct {
	Type      int
	Value     string
	Arguments []*AstNode
}

func NewAstNode(TokenType int, value string) *AstNode {
	return &AstNode{
		Type:      TokenType,
		Value:     value,
		Arguments: []*AstNode{},
	}
}

func Parse(tokens []Token) AstNode {
	astNode := NewAstNode(0, "root")
	j := 0
	isIdentifier := false
	for i := 0; i < len(tokens); i++ {
		aNode := NewAstNode(tokens[i].tokenType, tokens[i].value)
		switch tokens[i].tokenType {
		case keywords:
			fmt.Println(tokens[i].value)
			astNode.Arguments = append(astNode.Arguments, aNode)
			continue
		case label:
			astNode.Arguments = append(astNode.Arguments, aNode)
		case newline:
			j++
			continue
		case leftquote, rightquote, eof:
			continue
		case leftpercent:
			isIdentifier = true
		case rightpercent:
			isIdentifier = false
		case equals:
			i++
			aNode.Type = text
			aNode.Value = tokens[i].value
			astNode.Arguments[j].Arguments[len(astNode.Arguments[j].Arguments)-1].Arguments = append(astNode.Arguments[j].Arguments[len(astNode.Arguments[j].Arguments)-1].Arguments, aNode)
			continue
		default:
			if isIdentifier {
				isIdentifier = false
				aNode.Type = identifer
			}
			astNode.Arguments[j].Arguments = append(astNode.Arguments[j].Arguments, aNode)
		}
	}
	fmt.Println("===========\nParser\n===========")
	printAst(astNode, 0)
	return *astNode
}
func printAst(node *AstNode, depth int) {
	prefix := ""
	for i := 0; i < depth; i++ {
		prefix += "  "
	}
	fmt.Printf("%s- %s (type %d)\n", prefix, node.Value, node.Type)
	for _, child := range node.Arguments {
		printAst(child, depth+1)
	}
}
