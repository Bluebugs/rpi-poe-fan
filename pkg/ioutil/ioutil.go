package ioutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ReadStringFromFile(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(raw)), nil
}

func ReadIntFromFile(path string) (int, error) {
	text, err := ReadStringFromFile(path)
	if err != nil {
		return 0, err
	}

	r, err := strconv.ParseInt(text, 10, 0)
	return int(r), err
}

func WriteIntToFile(path string, value int) error {
	evaluatedPath, err := filepath.EvalSymlinks(path)
	if len(evaluatedPath) > 0 && err == nil {
		path = evaluatedPath
	}
	valueAsString := fmt.Sprintf("%d", value)
	return os.WriteFile(path, []byte(valueAsString), 0644)
}
