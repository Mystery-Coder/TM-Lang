package main

import (
	"errors"
	"fmt"
)

type IntermediateRepresention struct {
	Meta struct {
		Start  string
		Accept string
		Reject string
	}

	Macros map[string][]Transition

	Main []Transition
}

type Target struct {
	Type   string // CALL or GOTO or RETURN
	Name   string // State name
	Return string // CALL <macro> -> q0         q0 is Return state
}

type Transition struct {
	Src    string
	Read   string
	Write  string
	Dir    string
	Target Target
}

type Parser struct {
	Tokens       []Token
	Position     int
	CurrentToken Token
	IR           IntermediateRepresention
}

func (parser *Parser) initParser(_tokens []Token) {
	parser.Tokens = _tokens
	parser.Position = 0
	if len(_tokens) > 0 {
		parser.CurrentToken = parser.Tokens[0]
	} else {
		parser.CurrentToken = Token{TypeOfToken: EOF}
	}
	parser.IR = IntermediateRepresention{}
	parser.IR.Macros = make(map[string][]Transition)
}

func (parser *Parser) advance() {
	parser.Position += 1
	if parser.Position < len(parser.Tokens) {
		parser.CurrentToken = parser.Tokens[parser.Position]
	} else {
		parser.CurrentToken = Token{}
	}
}

func (parser *Parser) consume(tokenType TokenType) (string, error) {
	if !parser.CurrentToken.isNil() && parser.CurrentToken.TypeOfToken == tokenType {
		currentTextValue := parser.CurrentToken.Value
		parser.advance()
		return currentTextValue, nil
	} else {
		currentTextValue := parser.CurrentToken.Value
		if parser.CurrentToken.isNil() {
			currentTextValue = "EOF"
		}

		return "", fmt.Errorf("Syntax Error: Expected %s but got %s at Line %d", tokenType, currentTextValue, parser.CurrentToken.Line)

	}
}

func (parser *Parser) parseConfig() error {
	_, err := parser.consume(SECTION) // Handle error
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {

		configKeyword, err := parser.consume(KEYWORD)

		if err != nil {
			return err
		}

		configIdentifier, err := parser.consume(ID)

		if err != nil {
			return err
		}

		switch configKeyword {
		case "START:":
			parser.IR.Meta.Start = configIdentifier
		case "ACCEPT:":
			parser.IR.Meta.Accept = configIdentifier
		case "REJECT:":
			parser.IR.Meta.Reject = configIdentifier
		}
	}

	return nil
}

func (parser *Parser) parseMacros() error {
	if _, err := parser.consume(SECTION); err != nil { // Move parser to next token after "MACROS:", DEF <name>:
		return err
	}

	for parser.CurrentToken.TypeOfToken != SECTION {
		if _, err := parser.consume(KEYWORD); err != nil {
			return err
		}
		macroIdentifier, err := parser.consume(ID)
		if err != nil {
			return err
		}
		if _, err := parser.consume(COLON); err != nil {
			return err
		}
		var transitions []Transition
		for parser.CurrentToken.TypeOfToken == ID {
			transition, err := parser.parseTransition()

			if err != nil {
				return err
			}

			transitions = append(transitions, transition)
		}

		parser.IR.Macros[macroIdentifier] = transitions

	}

	return nil
}

func (parser *Parser) parseMain() error {
	_, err := parser.consume(SECTION) // Handle error
	if err != nil {
		return err
	}

	for parser.CurrentToken.TypeOfToken != EOF {
		transtion, err := parser.parseTransition()
		if err != nil {
			return err
		}
		parser.IR.Main = append(parser.IR.Main, transtion)
	}
	return nil
}

func (parser *Parser) parseTransition() (Transition, error) { // Parses main and macros transitions | q0, 1 -> 1, R, q0

	srcIdentifier, err := parser.consume(ID)
	if err != nil {
		return Transition{}, err

	}

	if _, err := parser.consume(COMMA); err != nil {
		return Transition{}, err

	}

	readSymbol, err := parser.consume(SYMBOL)
	if err != nil {
		return Transition{}, err

	}

	if _, err := parser.consume(ARROW); err != nil {
		return Transition{}, err

	}

	writeSymbol, err := parser.consume(SYMBOL)
	if err != nil {
		return Transition{}, err

	}

	if _, err := parser.consume(COMMA); err != nil {
		return Transition{}, err

	}

	direction, err := parser.consume(DIRECTION)
	if err != nil {
		return Transition{}, err

	}
	if _, err := parser.consume(COMMA); err != nil {
		return Transition{}, err

	}

	var target Target

	if parser.CurrentToken.TypeOfToken == KEYWORD { // CALL Keyword

		kw, err := parser.consume(KEYWORD)
		if err != nil {
			return Transition{}, err
		}
		switch kw {
		case "CALL":
			macroIdentifier, err := parser.consume(ID) // CALL move_to_end -> q1

			if err != nil {
				return Transition{}, err

			}

			if _, err := parser.consume(ARROW); err != nil {
				return Transition{}, err

			}
			returnStateIdentifier, err := parser.consume(ID)
			if err != nil {
				return Transition{}, err
			}

			target.Type = "CALL"
			target.Name = macroIdentifier
			target.Return = returnStateIdentifier

		case "RETURN":
			target.Type = "RETURN"
		}
	} else { // regular, q0, 0 -> 0, q1
		returnStateIdentifier, err := parser.consume(ID)

		if err != nil {
			return Transition{}, err

		}

		target.Type = "GOTO"
		target.Name = returnStateIdentifier
	}

	return Transition{
		srcIdentifier,
		readSymbol,
		writeSymbol,
		direction,
		target,
	}, nil

}

func (parser *Parser) parse() (IntermediateRepresention, error) {

	if parser.CurrentToken.TypeOfToken == SECTION && parser.CurrentToken.Value == "CONFIG:" {
		if err := parser.parseConfig(); err != nil {
			return IntermediateRepresention{}, err
		}
	} else {
		return parser.IR, errors.New("Program must contain a CONFIG section")
	}

	if parser.CurrentToken.TypeOfToken == SECTION && parser.CurrentToken.Value == "MACROS:" {
		parser.parseMacros()
	} // Macros are optinal

	if parser.CurrentToken.TypeOfToken == SECTION && parser.CurrentToken.Value == "MAIN:" {
		if err := parser.parseMain(); err != nil {
			return IntermediateRepresention{}, err
		}
	} else {
		return parser.IR, errors.New("Program must contain a MAIN section")
	}

	return parser.IR, nil
}
