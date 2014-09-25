package gslang

import (
	"fmt"
	"path/filepath"
)

//Position position of source code file
type Position struct {
	FileName string //script file name
	Lines    int    //line number, starting at 1
	Column   int    //column number, starting at 1 (character count per line)
}

//ShortName get the source code file short name
func (pos Position) ShortName() string {
	return filepath.Base(pos.FileName)
}

func (pos Position) String() string {
	return fmt.Sprintf("%s(%d,%d)", pos.FileName, pos.Lines, pos.Column)
}
