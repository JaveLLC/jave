package lexer

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// JaveFile defines incoming jave source
type JaveFile struct {
	Path     string
	Name     string
	Contents []rune
	NewLines []int
	Tokens   []*Token
}

// NewJaveFile ingests jave source file
func NewJaveFile(path string) (*JaveFile, error) {
	// Get the filename
	filename := filepath.Base(path)

	jf := &JaveFile{Name: filename, Path: path}
	jf.NewLines = append(jf.NewLines, -1)
	jf.NewLines = append(jf.NewLines, -1)

	contents, err := ioutil.ReadFile(jf.Path)
	if err != nil {
		return nil, err
	}

	jf.Contents = []rune(string(contents))
	return jf, nil
}

// GetLine gets a line of coke
func (s *JaveFile) GetLine(line int) string {
	return string(s.Contents[s.NewLines[line]+1 : s.NewLines[line+1]])
}

// TabWidth idk
//const TabWidth = 4

// MarkPos because mark is a pos
func (s *JaveFile) MarkPos(pos Position) string {
	buf := new(bytes.Buffer)

	lineString := s.GetLine(pos.Line)
	lineStringRunes := []rune(lineString)
	pad := pos.Char - 1

	buf.WriteString(strings.Replace(strings.Replace(lineString, "%", "%%", -1), "\t", "    ", -1))
	buf.WriteRune('\n')
	for i := 0; i < pad; i++ {
		spaces := 1

		if lineStringRunes[i] == ' ' {
			spaces = 1
		}

		for t := 0; t < spaces; t++ {
			buf.WriteRune(' ')
		}
	}
	buf.WriteString("^")
	buf.WriteRune('\n')

	return buf.String()

}

// MarkSpan ohai mark
func (s *JaveFile) MarkSpan(span Span) string {
	// if the span is just one character, use MarkPos instead
	spanEnd := span.End()
	spanEnd.Char--
	if span.Start() == spanEnd {
		return s.MarkPos(span.Start())
	}

	// mark the span
	buf := new(bytes.Buffer)

	for line := span.StartLine; line <= span.EndLine; line++ {
		lineString := s.GetLine(line)
		lineStringRunes := []rune(lineString)

		var pad int
		if line == span.StartLine {
			pad = span.StartChar - 1
		} else {
			pad = 0
		}

		var length int
		if line == span.EndLine {
			length = span.EndChar - span.StartChar
		} else {
			length = len(lineStringRunes)
		}

		buf.WriteString(strings.Replace(strings.Replace(lineString, "%", "%%", -1), "\t", "    ", -1))
		buf.WriteRune('\n')

		for i := 0; i < pad; i++ {
			spaces := 1

			if lineStringRunes[i] == ' ' {
				spaces = 1
			}

			for t := 0; t < spaces; t++ {
				buf.WriteRune(' ')
			}
		}

		buf.WriteString("[DERP]")
		for i := 0; i < length; i++ {
			// there must be a less repetitive way to do this but oh well
			spaces := 1

			if lineStringRunes[i+pad] == ' ' {
				spaces = 1
			}

			for t := 0; t < spaces; t++ {
				buf.WriteRune('~')
			}
		}
		//buf.WriteString(util.TEXT_RESET)
		buf.WriteRune('\n')
	}

	return buf.String()
}
