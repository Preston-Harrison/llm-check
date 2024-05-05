package lang

import (
	"fmt"
	"os"
	"path"
)

const GO = "go"

func Detect(dir string) (string, error) {
	goModPath := path.Join(dir, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		return GO, nil
	}

	return "", fmt.Errorf("could not detect language")
}
