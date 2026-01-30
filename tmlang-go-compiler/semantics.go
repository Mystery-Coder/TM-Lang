package main

import (
	"errors"
	"fmt"
)

type FlatTransition struct {
	Src   string
	Read  string
	Write string
	Dir   string
	Next  string
}

type SemanticAnalyzer struct {
	IR           IntermediateRepresention
	FinalIR      []FlatTransition
	MacroCounter int
}

func (analyzer *SemanticAnalyzer) initSemanticAnalyzer(_IR IntermediateRepresention) {
	analyzer.IR = _IR
	analyzer.FinalIR = make([]FlatTransition, 0)
	analyzer.MacroCounter = 0
}

func (analyzer *SemanticAnalyzer) analyze() ([]FlatTransition, error) {

	start := analyzer.IR.Meta.Start

	if start == "" {
		return nil, errors.New("Syntax Error: No START found in Config")
	}

	for _, transiton := range analyzer.IR.Main {
		analyzer.processTransition(transiton)
	}

	return analyzer.FinalIR, nil
}

func (analyzer *SemanticAnalyzer) processTransition(transition Transition) error {

	target := transition.Target

	switch target.Type {
	case "GOTO":
		analyzer.FinalIR = append(analyzer.FinalIR, FlatTransition{
			Src:   transition.Src,
			Read:  transition.Read,
			Write: transition.Write,
			Dir:   transition.Dir,
			Next:  target.Name,
		})
	case "CALL":
		macroName := target.Name

		returnState := target.Return

		macroTranstions := analyzer.IR.Macros[macroName]
		if macroTranstions == nil {
			return fmt.Errorf("Call to undefined macro, %s", macroName)
		}
		analyzer.MacroCounter++
		prefix := fmt.Sprintf("%s_%d_", macroName, analyzer.MacroCounter) // For the states in the macro

		macroStart := macroTranstions[0].Src
		macroStartRenamed := prefix + macroStart // this will be the new start

		analyzer.FinalIR = append(analyzer.FinalIR, FlatTransition{
			transition.Src,
			transition.Read,
			transition.Write,
			transition.Dir,
			macroStartRenamed,
		})

		for _, macroTransition := range macroTranstions {
			newSrc := prefix + macroTransition.Src

			macroTransactionTarget := macroTransition.Target

			var newNext string

			switch macroTransactionTarget.Type {
			case "GOTO":
				newNext = prefix + macroTransactionTarget.Name
			case "RETURN":
				newNext = returnState
			case "CALL":
				return errors.New("MACROS cannot call other macros")
			}

			analyzer.FinalIR = append(analyzer.FinalIR, FlatTransition{
				newSrc,
				transition.Read,
				transition.Write,
				transition.Dir,
				newNext,
			})
		}

	}

	return nil
}
