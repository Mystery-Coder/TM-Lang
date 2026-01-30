package main

import "fmt"

func Compile(sourceCode string) (string, string, error) {

	var lexer Lexer
	lexer.initLexer(sourceCode)
	tokens := lexer.tokenizeSource()

	fmt.Println(tokens)

	var parser Parser
	parser.initParser(tokens)
	ir, err := parser.parse()
	fmt.Println(ir)

	if err != nil {
		return "", "", err
	}

	var analyzer SemanticAnalyzer
	analyzer.initSemanticAnalyzer(ir)
	finalIR, err := analyzer.analyze()
	if err != nil {
		return "", "", err
	}

	var codegen CodeGenerator
	codegen.initCodegen(ir.Meta, finalIR)

	cCode := codegen.GenerateC()
	dotCode := codegen.GenerateDot()

	return cCode, dotCode, nil
}
