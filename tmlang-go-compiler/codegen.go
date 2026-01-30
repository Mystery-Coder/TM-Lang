package main

import (
	"fmt"
	"sort"
	"strings"
)

type CodeGenerator struct {
	Meta struct {
		Start  string
		Accept string
		Reject string
	}
	FinalIR []FlatTransition
}

func (cg *CodeGenerator) initCodegen(meta struct {
	Start  string
	Accept string
	Reject string
}, finalIR []FlatTransition) {
	cg.Meta = meta
	cg.FinalIR = finalIR
}

func (cg *CodeGenerator) GenerateC() string {

	// Collect Unique States
	stateSet := make(map[string]bool)
	stateSet[cg.Meta.Start] = true
	stateSet[cg.Meta.Accept] = true
	stateSet[cg.Meta.Reject] = true

	for _, t := range cg.FinalIR {
		stateSet[t.Src] = true
		stateSet[t.Next] = true
	}

	//  Sort States and Map to Integers
	var stateList []string
	for name := range stateSet {
		stateList = append(stateList, name)
	}
	sort.Strings(stateList)

	stateMap := make(map[string]int)
	stateComments := ""

	for i, name := range stateList {
		stateMap[name] = i
		stateComments += fmt.Sprintf("// %s: %d\n", name, i)
	}

	startID := stateMap[cg.Meta.Start]
	acceptID := stateMap[cg.Meta.Accept]
	rejectID := stateMap[cg.Meta.Reject]

	// Group transitions by Source State ID
	grouped := make(map[int][]FlatTransition)
	for _, t := range cg.FinalIR {
		srcID := stateMap[t.Src]
		grouped[srcID] = append(grouped[srcID], t)
	}

	// Build switch-case logic
	switchLogic := ""

	for i := 0; i < len(stateList); i++ {
		rules, exists := grouped[i]
		if !exists {
			continue
		}

		switchLogic += fmt.Sprintf("            case %d:\n", i)

		for j, rule := range rules {
			prefix := "else if"
			if j == 0 {
				prefix = "if"
			}

			moveCode := ""
			switch rule.Dir {
			case "R":
				moveCode = "head++;"
			case "L":
				moveCode = "head--;"
			}

			nextID := stateMap[rule.Next]

			switchLogic += fmt.Sprintf(`                %s (read_val == '%s') {
                    tape[head] = '%s';
                    %s
                    current_state = %d;
                    matched = 1;
                }
			`, prefix, rule.Read, rule.Write, moveCode, nextID)
		}
		switchLogic += "                break;\n"
	}

	// Final C Code
	cCode := fmt.Sprintf(`#include <stdio.h>
		#include <stdlib.h>
		#include <string.h>

		#define TAPE_SIZE 20000
		#define HEAD_START 10000

		/* --- STATE MAP --- 
		%s*/

		int current_state = %d;
		int ACCEPT_STATE = %d;
		int REJECT_STATE = %d;

		char tape[TAPE_SIZE];
		int head = HEAD_START;

		void print_tape() {
			printf("\r[ ");
			for(int i = head - 10; i <= head + 10; i++) {
				if(i == head) printf("[%%c]", tape[i]);
				else printf(" %%c ", tape[i]);
			}
			printf(" ] State: %%d  ", current_state);
			fflush(stdout); 
		}

		int main() {
			memset(tape, '_', TAPE_SIZE);
			
			printf("Enter Input: ");
			char input[100];
			scanf("%%s", input);
			
			for(int i=0; i<strlen(input); i++) {
				tape[head + i] = input[i];
			}

			printf("\n--- RUNNING ---\n");

			while(1) {
				print_tape();

				if (current_state == ACCEPT_STATE) { printf("\n\nACCEPTED!\n"); return 0; }
				if (current_state == REJECT_STATE) { printf("\n\nREJECTED!\n"); return 1; }

				char read_val = tape[head];
				int matched = 0;

				switch(current_state) {
						%s
				}
				
				if (!matched) {
					printf("\n\nCRASH: State %%d has no rule for char '%%c'\n", current_state, read_val);
					return 1;
				}
			}
		}
	`, stateComments, startID, acceptID, rejectID, switchLogic)

	return cCode
}

func (cg *CodeGenerator) GenerateDot() string {
	fmt.Println("--- Generating Diagram ---")
	var sb strings.Builder

	sb.WriteString("digraph TuringMachine {\n")
	sb.WriteString("    rankdir=LR;\n")
	sb.WriteString("    node [shape = circle];\n")

	// 1. Special Shapes
	sb.WriteString(fmt.Sprintf("    \"%s\" [shape = doublecircle, color=green];\n", cg.Meta.Accept))
	sb.WriteString(fmt.Sprintf("    \"%s\" [shape = doublecircle, color=red];\n", cg.Meta.Reject))

	// 2. Entry Point
	sb.WriteString("    entry [shape = point];\n")
	sb.WriteString(fmt.Sprintf("    entry -> \"%s\";\n", cg.Meta.Start))

	// 3. Group Transitions (Edge Merging)
	type EdgeKey struct {
		Src, Dst string
	}
	edges := make(map[EdgeKey][]string)

	for _, t := range cg.FinalIR {
		key := EdgeKey{Src: t.Src, Dst: t.Next}

		// Clean Labels for Diagram (Replace '_' with 'BLANK')
		rLbl := t.Read
		if rLbl == "_" {
			rLbl = "BLANK"
		}
		wLbl := t.Write
		if wLbl == "_" {
			wLbl = "BLANK"
		}

		label := fmt.Sprintf("%s / %s, %s", rLbl, wLbl, t.Dir)
		edges[key] = append(edges[key], label)
	}

	// 4. Generate Merged Edges
	for key, labels := range edges {
		combinedLabel := strings.Join(labels, "\\n")
		sb.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\" [label = \"%s\"];\n", key.Src, key.Dst, combinedLabel))
	}

	sb.WriteString("}\n")
	return sb.String()
}
