package lexer

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/OrcaLLC/jave/src/util/timed"
)

type lexer struct {
	input            *JaveFile
	startPos, endPos int
	curPos           Position
	tokStart         Position
}

// this is almost entirely stolen because Jave believes uniqueness
// is a flaw

func (v *lexer) errPos(pos Position, err string, stuff ...interface{}) {
	fmt.Printf("[sexer error]: [%s:%d:%d] %s",
		pos.Filename, pos.Line, pos.Char, fmt.Sprintf(err, stuff...))

	fmt.Printf("[lexer] %v", v.input.MarkPos(pos))

	os.Exit(1)
}

func (v *lexer) err(err string, stuff ...interface{}) {
	v.errPos(v.curPos, err, stuff...)
}

func (v *lexer) peek(ahead int) rune {
	if ahead < 0 {
		panic(fmt.Sprintf("[pervert] Tried to peek a negative number: %d", ahead))
	}

	if v.endPos+ahead >= len(v.input.Contents) {
		return 0
	}
	return v.input.Contents[v.endPos+ahead]
}

func (v *lexer) consume() {
	v.curPos.Char++
	if v.peek(0) == '\n' {
		v.curPos.Char = 1
		v.curPos.Line++
		v.input.NewLines = append(v.input.NewLines, v.endPos)
	}

	v.endPos++
}

func (v *lexer) expect(r rune) {
	if v.peek(0) == r {
		v.consume()
	} else {
		v.err("Expected `%c`, found `%c`", r, v.peek(0))
	}
}

func (v *lexer) discardBuffer() {
	v.startPos = v.endPos

	v.tokStart = v.curPos
}

func (v *lexer) pushToken(t TokenType) {
	tok := &Token{
		Type:     t,
		Contents: string(v.input.Contents[v.startPos:v.endPos]),
		Where:    NewSpan(v.tokStart, v.curPos),
	}

	v.input.Tokens = append(v.input.Tokens, tok)

	fmt.Printf("lexer sexer [%4d:%4d:% 11s] `%s`\n", v.startPos, v.endPos, tok.Type, tok.Contents)

	v.discardBuffer()
}

// Lex the sex up
func Lex(input *JaveFile) []*Token {
	v := &lexer{
		input:    input,
		startPos: 0,
		endPos:   0,
		curPos:   Position{Filename: input.Name, Line: 1, Char: 1},
		tokStart: Position{Filename: input.Name, Line: 1, Char: 1},
	}

	timed.Timed("lexing", input.Name, func() {
		v.lex()
	})
	//fmt.Println("Sexing up the code... ahem I mean lexing")
	//v.lex()

	return v.input.Tokens
}

func (v *lexer) lex() {
	for {
		v.skipLayoutAndComments()

		if isEOF(v.peek(0)) {
			v.input.NewLines = append(v.input.NewLines, v.endPos)
			return
		}

		fmt.Printf("Peeking: %v\n", v.peek(0))

		if isLetter(v.peek(0)) || v.peek(0) == '_' {
			v.recognizeIdentifierToken()
		} else if v.peek(0) == '<' && v.peek(1) == '<' {
			v.recognizeStringToken()
		} else if isSeparator(v.peek(0)) {
			v.recognizeSeparatorToken()
		} else {
			v.err("I don't know this token.\n")
		}

		//
		//		if isEOF(v.peek(0)) {
		//			v.input.NewLines = append(v.input.NewLines, v.endPos)
		//			return
		//		}
		//
		//		if isDecimalDigit(v.peek(0)) {
		//			v.recognizeNumberToken()
		//		} else if isLetter(v.peek(0)) || v.peek(0) == '_' {
		//			v.recognizeIdentifierToken()
		//		} else if v.peek(0) == '"' {
		//			v.recognizeStringToken()
		//		} else if v.peek(0) == '\'' {
		//			v.recognizeCharacterToken()
		//		} else if isOperator(v.peek(0)) {
		//			v.recognizeOperatorToken()
		//		} else if isSeparator(v.peek(0)) {
		//			v.recognizeSeparatorToken()
		//		} else {
		//			v.err("Unrecognised token")
		//		}
		//fmt.Printf("No sexing defined.")
	}
}

// returns true if a comment was skipped
func (v *lexer) skipComment() bool {
	if v.skipBlockComment() {
		return true
	} else if v.skipLineComment() {
		return true
	} else {
		return false
	}
}

func (v *lexer) skipBlockComment() bool {
	pos := v.curPos
	if v.peek(0) != '/' || v.peek(1) != '\\' {
		return false
	}

	v.consume()
	v.consume()
	isDoc := v.peek(0) == '\\'

	depth := 1
	for depth > 0 {
		if isEOF(v.peek(0)) {
			v.errPos(pos, "Unterminated block comment")
		}

		if v.peek(0) == '/' && v.peek(1) == '\\' {
			v.consume()
			v.consume()
			depth++
		}

		if v.peek(0) == '\\' && v.peek(1) == '/' {
			v.consume()
			v.consume()
			depth--
		}

		v.consume()
	}

	if isDoc {
		v.pushToken(Doccomment)
	} else {
		v.discardBuffer()
	}
	fmt.Println("Recognized comment")
	return true
}

func (v *lexer) skipLineComment() bool {
	if v.peek(0) != '/' || v.peek(1) != '/' || v.peek(2) != '\\' || v.peek(3) != '\\' {
		return false
	}

	v.consume()
	v.consume()
	isDoc := v.peek(0) == '/'

	for {
		if isEOL(v.peek(0)) || isEOF(v.peek(0)) {
			if isDoc {
				v.pushToken(Doccomment)
			} else {
				v.discardBuffer()
			}
			v.consume()
			fmt.Println("Recognized line comment")
			return true
		}
		v.consume()
	}
}

func (v *lexer) skipLayoutAndComments() {
	for {
		for isLayout(v.peek(0)) {
			v.consume()
		}
		v.discardBuffer()

		if !v.skipComment() {
			break
		}

	}
	v.discardBuffer()
}

// RECOGNIZE ME
func (v *lexer) recognizeIdentifierToken() {
	v.consume()

	for isLetter(v.peek(0)) || isDecimalDigit(v.peek(0)) || v.peek(0) == '_' {
		v.consume()
	}

	v.pushToken(Identifier)
}

func (v *lexer) recognizeStringToken() {
	pos := v.curPos

	v.expect('<')
	v.expect('<')
	v.discardBuffer()

	for {
		// escape to use >
		if v.peek(0) == '\\' && v.peek(1) == '>' {
			v.consume()
			v.consume()
		} else if v.peek(0) == '>' && v.peek(1) == '>' {
			v.pushToken(String)
			v.consume()
			v.consume()
			return
		} else if isEOF(v.peek(0)) {
			v.errPos(pos, "Unterminated string literal")
		} else {
			v.consume()
		}
	}
}

func (v *lexer) recognizeSeparatorToken() {
	v.consume()
	v.pushToken(Separator)
}

// is things
func isDecimalDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isHexDigit(r rune) bool {
	return isDecimalDigit(r) || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

func isBinaryDigit(r rune) bool {
	return r == '0' || r == '1'
}

func isOctalDigit(r rune) bool {
	return r >= '0' && r <= '7'
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r)
}

func isOperator(r rune) bool {
	return strings.ContainsRune("+-*/=><!~?:|&%^#@", r)
}

func isSeparator(r rune) bool {
	return strings.ContainsRune(" ;,.`(){}[]", r)
}

func isEOL(r rune) bool {
	return r == '\n'
}

func isEOF(r rune) bool {
	return r == 0
}

func isLayout(r rune) bool {
	return (r <= ' ' || unicode.IsSpace(r)) && !isEOF(r)
}
