package main

/*
Example Code,

CONFIG:
    START: q0
    ACCEPT: done
    REJECT: fail
MACROS:
    DEF seek_blank:
        q0, 0 -> 1, R, q0
        q0, _ -> _, S, RETURN
MAIN:
    q0, 1 -> 0, R, CALL seek_blank -> done
*/
import (
	"fmt"
	"regexp"
)

type TokenType string

const (
	SECTION   TokenType = "SECTION" // CONFIG:, MACROS:, MAIN:
	KEYWORD   TokenType = "KEYWORD" // START:, ACCEPT:, REJECT:, DEF, CALL, RETURN
	ID        TokenType = "ID"      // Identifiers (q0, my_macro)
	SYMBOL    TokenType = "SYMBOL"  // 0, 1, *, _
	DIRECTION TokenType = "DIR"     // L, R, S
	ARROW     TokenType = "ARROW"   // ->
	COMMA     TokenType = "COMMA"   // ,
	COLON     TokenType = "COLON"   // :
	NEWLINE   TokenType = "NEWLINE" // \n
	SKIP      TokenType = "SKIP"
	COMMENT   TokenType = "COMMENT"
	MISMATCH  TokenType = "MISMATCH"
	EOF       TokenType = "EOF" // End of File
)

func (token *Token) isNil() bool {
	return token.Line == 0 && token.TypeOfToken == "" && token.Value == ""
}

type Token struct {
	TypeOfToken TokenType
	Line        int
	Value       string // For Ex: ->, CONFIG:
}

type Rule struct {
	TypeOfToken TokenType
	Regex       *regexp.Regexp
}

type Lexer struct {
	SourceCode  string
	Tokens      []Token
	CurrentLine int
	Rules       []Rule
}

func (lexer *Lexer) initLexer(src string) {
	lexer.CurrentLine = 1
	lexer.SourceCode = src
	lexer.Tokens = nil

	lexer.Rules = []Rule{
		{SECTION, regexp.MustCompile(`^(CONFIG:|MACROS:|MAIN:)`)},
		{KEYWORD, regexp.MustCompile(`^(START:|ACCEPT:|REJECT:)`)},
		{KEYWORD, regexp.MustCompile(`^(DEF|CALL|RETURN)\b`)},
		{ARROW, regexp.MustCompile(`^->`)},
		{COMMA, regexp.MustCompile(`^,`)},
		{COLON, regexp.MustCompile(`^:`)},
		{DIRECTION, regexp.MustCompile(`\b(L|R|S)\b`)}, //L R S are reserved
		{ID, regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*`)},
		{SYMBOL, regexp.MustCompile(`^[0-9a-zA-Z_]`)},
		{NEWLINE, regexp.MustCompile(`^\n`)},
		{SKIP, regexp.MustCompile(`^[ \t\r]+`)},
		{COMMENT, regexp.MustCompile(`^//.*`)},
		{MISMATCH, regexp.MustCompile(`^.`)},
	}
}

func (lexer *Lexer) tokenizeSource() []Token {

	pos := 0

	for pos < len(lexer.SourceCode) {

		matched := false
		for _, rule := range lexer.Rules {
			location := rule.Regex.FindStringIndex(lexer.SourceCode[pos:]) // location has start and end position of the match

			if location != nil && location[0] == 0 { // If a token type match is found
				textValue := lexer.SourceCode[pos : pos+location[1]]

				switch rule.TypeOfToken {
				case NEWLINE:
					lexer.CurrentLine++
				case SKIP, COMMENT:
					// skipping
				case MISMATCH:
					fmt.Printf(
						"Lexer mismatch %q at line %d\n",
						textValue,
						lexer.CurrentLine,
					)
					return nil
				default:
					lexer.Tokens = append(lexer.Tokens, Token{
						TypeOfToken: rule.TypeOfToken,
						Value:       textValue,
						Line:        lexer.CurrentLine,
					})
				}
				pos += location[1]
				matched = true
				break
			}

		}
		if !matched {

			fmt.Printf(
				"Lexer error: unexpected character %q at line %d\n",
				lexer.SourceCode[pos],
				lexer.CurrentLine,
			)
			return nil
		}
	}

	lexer.Tokens = append(lexer.Tokens, Token{
		TypeOfToken: EOF,
		Value:       "",
		Line:        lexer.CurrentLine,
	})

	return lexer.Tokens

}
