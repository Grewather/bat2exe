package main

import (
	"fmt"
	"log"
)

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

type ifStruct struct {
	Condition *AstNode
	Body      *AstNode
	Head      *AstNode
}

func Parse(tokens []Token) AstNode {
	astNode := NewAstNode(0, "root")
	j := 0
	isIdentifier := false
	ifstruct := ifStruct{}
	ifstruct.Head = astNode
	var currentNode *AstNode
	for i := 0; i < len(tokens); i++ {
		aNode := NewAstNode(tokens[i].tokenType, tokens[i].value)
		switch tokens[i].tokenType {
		case keywords:
			if ifstruct.Body != nil {
				ifstruct.Body.Arguments = append(ifstruct.Body.Arguments, aNode)
				currentNode = aNode
				continue
			}
			astNode.Arguments = append(astNode.Arguments, aNode)
			if tokens[i].value == "if" {
				ifstruct.Head.Arguments[j].Arguments = append(ifstruct.Head.Arguments[j].Arguments, NewAstNode(condidion, "condition"))
				ifstruct.Condition = ifstruct.Head.Arguments[j].Arguments[0]
				continue
			}
			continue
		case label:
			ifstruct.Head.Arguments = append(ifstruct.Head.Arguments, aNode)
		case newline:
			j++
			continue
		case leftquote, rightquote, eof:
			continue
		case leftpercent:
			isIdentifier = true
		case rightpercent:
			isIdentifier = false
		case lparen:
			// ifstruct.Condition = nil
			// astNode.Arguments[j].Arguments
			if ifstruct.Condition != nil {
				ifstruct.Condition = nil
				ifstruct.Head.Arguments[j].Arguments = append(ifstruct.Head.Arguments[j].Arguments, NewAstNode(body, "body"))
				ifstruct.Body = ifstruct.Head.Arguments[j].Arguments[1]
			}
			continue
		case rparen:
			if ifstruct.Body != nil {
				ifstruct.Body = nil
			}
			continue
		case equals:
			i++
			aNode.Type = text
			aNode.Value = tokens[i].value
			ifstruct.Head.Arguments[j].Arguments[len(ifstruct.Head.Arguments[j].Arguments)-1].Arguments = append(ifstruct.Head.Arguments[j].Arguments[len(ifstruct.Head.Arguments[j].Arguments)-1].Arguments, aNode)
			continue
		default:
			if isIdentifier {
				isIdentifier = false
				aNode.Type = identifer
			}
			if ifstruct.Condition != nil {
				ifstruct.Condition.Arguments = append(ifstruct.Condition.Arguments, aNode)
				continue
			}
			if ifstruct.Body != nil && currentNode != nil {
				currentNode.Arguments = append(currentNode.Arguments, aNode)
				continue
			}
			log.Println("debug", aNode.Value, aNode.Type, aNode.Arguments)
			log.Println("debug 2", tokens[i].tokenType, tokens[i].value)
			ifstruct.Head.Arguments[j].Arguments = append(ifstruct.Head.Arguments[j].Arguments, aNode)
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
