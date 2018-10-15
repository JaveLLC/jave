package lexer

// TokenType maps to tokens
type TokenType int

const (
	// Rune is a tokentype iota
	Rune TokenType = iota
	Identifier
	Separator
	Operator
	Number
	Erroneous
	String
	Doccomment
)

var tokenStrings = []string{"rune", "identifier", "separator", "operator", "number", "erroneous", "string", "doccomment"}

func (v TokenType) String() string {
	return tokenStrings[v]
}

// Token defines a token
type Token struct {
	Type     TokenType
	Contents string
	Where    Span
}

// Position positional
type Position struct {
	Filename   string
	Line, Char int
}

// Span where a token is
type Span struct {
	Filename string

	StartLine int
	StartChar int
	EndLine   int
	EndChar   int
}

// NewSpan creates a new span
func NewSpan(start, end Position) Span {
	return Span{Filename: start.Filename,
		StartLine: start.Line, StartChar: start.Char,
		EndLine: end.Line, EndChar: end.Char,
	}
}

// NewSpanFromTokens is kind of self explanatory
func NewSpanFromTokens(start, end *Token) Span {
	return Span{Filename: start.Where.Filename,
		StartLine: start.Where.StartLine, StartChar: start.Where.StartChar,
		EndLine: end.Where.EndLine, EndChar: end.Where.EndChar,
	}
}

// Start stop maybe go
func (s Span) Start() Position {
	return Position{Filename: s.Filename,
		Line: s.StartLine, Char: s.StartChar}
}

// End all life
func (s Span) End() Position {
	return Position{Filename: s.Filename,
		Line: s.EndLine, Char: s.EndChar}
}
