package main

import (
	"log"
	"os"
	"strings"
)

var MapOfVars = make(map[string]string)
var MapOfFuncs = make(map[string]bool)

func handleVars(astNode *AstNode) string {
	res := ""
	for i, v := range astNode.Arguments {
		if i == 0 && v.Type == identifer {
			varName := v.Value
			varValue := v.Arguments[0].Value
			return "\n" + `	char ` + varName + `[] = "` + varValue + `";` + "\n"
		}
		if v.Type == option {
			continue
		}
		MapOfVars[astNode.Arguments[i].Value] = "int"
		res = "	int " + astNode.Arguments[i].Value + " = " + astNode.Arguments[i].Arguments[0].Value + ";\n"
	}
	return res
}

func handleDate() (string, string) {
	if MapOfFuncs["get_date"] {
		return "	get_date();\n", ""
	}
	MapOfFuncs["get_date"] = true
	return "	get_date();\n", "\nvoid get_date(){\n	time_t t; \n	time(&t); \n	" + `printf("` + `\n` + `%s"` + ", ctime(&t)); \n} \n "
}

func handleEcho(astNode *AstNode) string {
	resp := ""
	var vars []string
	resp = `	printf("\n`
	respVal := "%s"
	for i := range astNode.Arguments {
		arg := astNode.Arguments[i]
		if arg.Type == identifer {
			if MapOfVars[arg.Value] == "int" {
				respVal = "%i"
			}
			resp += respVal
			vars = append(vars, arg.Value)
			continue
		}
		resp += arg.Value
	}
	resp += `"`
	if len(vars) > 0 {
		resp += `, `
	}
	for _, v := range vars {
		resp += v
	}
	resp += `);` + "\n"
	return resp
}

func checkKeyword(arg *AstNode) (string, string) {
	switch arg.Value {
	case "echo":
		return handleEcho(arg), ""
	case "set":
		return handleVars(arg), ""
	case "date":
		return handleDate()
	case "exit":
		return "exit();", ""
	case "goto":
		return "goto " + arg.Arguments[0].Value + ";\n", ""
	}
	return "", ""
}

func generateCode(astNode *AstNode) {
	f, err := os.Create("main.c")
	if err != nil {
		log.Fatal("Failed to create main.c")
	}
	defer f.Close()
	var sb strings.Builder
	var functionBuilder strings.Builder
	_, err = functionBuilder.WriteString("#include <stdio.h>\n#include <time.h> \n")
	check(err)
	_, err = sb.WriteString("\nint main() {")
	check(err)
	os.Remove("main.c")
	for _, arg := range astNode.Arguments {
		if arg.Type == keywords {
			res, res2 := checkKeyword(arg)
			if len(res2) >= 1 {
				_, err = functionBuilder.WriteString(res2)
				check(err)
			}
			_, err := sb.WriteString(res)
			check(err)
		} else if arg.Type == label {
			_, err = sb.WriteString(arg.Value + ":\n")
			check(err)
		}
	}
	sb.WriteString("\n}\n")

	_, err = f.Write([]byte(functionBuilder.String()))
	check(err)
	_, err = f.Write([]byte(sb.String()))
	check(err)
}
