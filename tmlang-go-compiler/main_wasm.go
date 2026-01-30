//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"syscall/js"
)

// --- Structs for JSON Output ---

type CompileResponse struct {
	Status string `json:"status"` // "success" or "error"
	CCode  string `json:"c_code"`
	Dot    string `json:"dot"`
	Error  string `json:"error"`
}

type SimulationResult struct {
	Status  string           `json:"status"` // "ACCEPTED", "REJECTED", "TIMEOUT", "CRASH"
	History []SimulationStep `json:"history"`
}

type SimulationStep struct {
	StepCount int    `json:"step"`
	Tape      string `json:"tape"`
	Head      int    `json:"head"`
	State     string `json:"state"`
}

// JS Usage: const result = JSON.parse(window.tmCompile(sourceCode));
func compileWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return errorJson("Missing source code")
	}

	cCode, dotCode, err := Compile(args[0].String())

	if err != nil {
		return errorJson(err.Error())
	}

	resp := CompileResponse{
		Status: "success",
		CCode:  cCode,
		Dot:    dotCode,
	}
	b, _ := json.Marshal(resp)
	return string(b)
}

// JS Usage: const result = JSON.parse(window.tmRun(sourceCode, inputString));
func runWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return errorJson("Usage: tmRun(code, input)")
	}

	sourceCode := args[0].String()
	tapeInput := args[1].String()

	// 1. Re-run Pipeline to get Logic (IR)
	// We need the raw data structures (IR), not the C string.
	var lexer Lexer
	lexer.initLexer(sourceCode)

	var parser Parser
	parser.initParser(lexer.tokenizeSource())
	ir, err := parser.parse()
	if err != nil {
		return errorJson("Parse Error: " + err.Error())
	}

	var analyzer SemanticAnalyzer
	analyzer.initSemanticAnalyzer(ir)
	finalIR, err := analyzer.analyze()
	if err != nil {
		return errorJson("Semantic Error: " + err.Error())
	}

	// 2. Execute Simulation
	// finalIR is the list of transitions, ir.Meta contains Start/Accept/Reject
	result := runSimulationInternal(finalIR, ir.Meta, tapeInput, 5000)

	b, _ := json.Marshal(result)
	return string(b)
}

// This logic lives here because only the Web UI needs step-by-step history.
func runSimulationInternal(transitions []FlatTransition, meta Meta, input string, maxSteps int) SimulationResult {
	// 1. Setup Tape
	const TAPE_SIZE = 20000
	const HEAD_START = 10000

	tape := make([]rune, TAPE_SIZE)
	for i := range tape {
		tape[i] = '_'
	}

	// Load Input
	for i, char := range input {
		tape[HEAD_START+i] = char
	}

	head := HEAD_START
	currentState := meta.Start
	history := []SimulationStep{}

	// 2. Execution Loop
	for step := 0; step < maxSteps; step++ {

		// DYNAMIC VIEWPORT:
		// Calculate a window around the head so the user always sees the action.
		// We show 15 chars to the left and 15 to the right.
		viewStart := head - 15
		viewEnd := head + 15

		// Bounds safety checks
		if viewStart < 0 {
			viewStart = 0
		}
		if viewEnd > TAPE_SIZE {
			viewEnd = TAPE_SIZE
		}

		tapeWindow := string(tape[viewStart:viewEnd])

		history = append(history, SimulationStep{
			StepCount: step,
			Tape:      tapeWindow,
			Head:      head - viewStart, // The head index relative to the window string
			State:     currentState,
		})

		// Check End Conditions
		if currentState == meta.Accept {
			return SimulationResult{Status: "ACCEPTED", History: history}
		}
		if currentState == meta.Reject {
			return SimulationResult{Status: "REJECTED", History: history}
		}

		// Logic Lookup
		var match *FlatTransition
		charUnderHead := string(tape[head])

		for _, t := range transitions {
			if t.Src == currentState && t.Read == charUnderHead {
				match = &t
				break
			}
		}

		if match == nil {
			return SimulationResult{Status: "CRASH", History: history}
		}

		if len(match.Write) > 0 {
			tape[head] = rune(match.Write[0])
		}

		switch match.Dir {
		case "R":
			head++
		case "L":
			head--
		}
		currentState = match.Next
	}

	return SimulationResult{Status: "TIMEOUT", History: history}
}

func errorJson(msg string) string {
	b, _ := json.Marshal(CompileResponse{Status: "error", Error: msg})
	return string(b)
}

func main() {
	c := make(chan struct{})

	js.Global().Set("tmCompile", js.FuncOf(compileWrapper))
	js.Global().Set("tmRun", js.FuncOf(runWrapper))

	<-c
}
