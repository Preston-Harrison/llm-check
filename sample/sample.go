package sample

import (
	"io/fs"
	"llm-check/sample/sample2"
	"os"
)

// llm-check(caller): make sure to close this file after use.
func openFileForWriting() (*os.File, error) {
    sample2.MyExportedFunc()
    return os.OpenFile("/tmp/some_file.txt", os.O_CREATE, fs.ModeAppend)
}
