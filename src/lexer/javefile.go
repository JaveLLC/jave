package lexer

import (
	"io/ioutil"
	"path/filepath"
)

// JaveFile defines incoming jave source
type JaveFile struct {
	Path     string
	Name     string
	Contents []rune
	NewLines []int
}

// NewJaveFile ingests jave source file
func NewJaveFile(path string) (*JaveFile, error) {
	// Get the filename
	filename := filepath.Base(path)

	jf := &JaveFile{Name: filename, Path: path}
	jf.NewLines = append(jf.NewLines, -1)

	contents, err := ioutil.ReadFile(jf.Path)
	if err != nil {
		return nil, err
	}

	jf.Contents = []rune(string(contents))
	return jf, nil
}
