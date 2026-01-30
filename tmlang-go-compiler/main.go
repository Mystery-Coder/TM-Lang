package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Error: No input file provided.")
		fmt.Println("Usage: tmlang <file.tm>")
		os.Exit(1)
	}

	filepathArg := os.Args[1]

	code, err := os.ReadFile(filepathArg)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	cCode, dotCode, err := Compile(string(code))
	if err != nil {
		fmt.Printf("Compilation Failed: %v\n", err)
		os.Exit(1)
	}

	ext := filepath.Ext(filepathArg)
	baseName := strings.TrimSuffix(filepath.Base(filepathArg), ext)
	outputDir := "build"

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating build dir: %v\n", err)
		os.Exit(1)
	}

	cPath := filepath.Join(outputDir, baseName+".c")
	if err := os.WriteFile(cPath, []byte(cCode), 0644); err != nil {
		fmt.Printf("Error writing C file: %v\n", err)
		os.Exit(1)
	}

	dotPath := filepath.Join(outputDir, baseName+".dot")
	if err := os.WriteFile(dotPath, []byte(dotCode), 0644); err != nil {
		fmt.Printf("Error writing DOT file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("--- Converting to SVG ---")
	svgPath := filepath.Join(outputDir, baseName+".svg")

	// Check if 'dot' is installed
	if _, err := exec.LookPath("dot"); err == nil {
		cmd := exec.Command("dot", "-Tsvg", dotPath, "-o", svgPath)
		if err := cmd.Run(); err != nil {
			fmt.Println("Warning: Failed to generate SVG (is GraphViz working?)")
		}
	} else {
		fmt.Println("Note: GraphViz ('dot') not found. Skipping SVG generation.")
	}

	fmt.Printf("\n Output saved to '%s/'\n", outputDir)
}

func Compile(sourceCode string) (string, string, error) {

	var lexer Lexer
	lexer.initLexer(sourceCode)
	tokens := lexer.tokenizeSource()

	var parser Parser
	parser.initParser(tokens)
	ir, err := parser.parse()
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
