package test

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	err := os.MkdirAll(filepath.Join("testdata", "failed"), 0755)
	if err != nil {
		panic(err)
	}
}

func VerifyImage(t *testing.T, resultPath string) {
	wd, err := os.Getwd()
	assert.NoError(t, err)

	filename := filepath.Base(resultPath)

	masterPath := filepath.Join(wd, "testdata", filename)

	masterBound, masterPixels := get(t, masterPath)
	resultBound, resultPixels := get(t, resultPath)

	if len(resultPixels) == 0 || len(masterPixels) == 0 {
		return
	}

	assert.Len(t, resultPixels, len(masterPixels))
	assert.Equal(t, masterBound, resultBound)
	assert.Equal(t, masterPixels, resultPixels)

	if !t.Failed() {
		os.Remove(resultPath)
	}
}

func get(t *testing.T, path string) (image.Rectangle, []uint8) {
	f, err := os.Open(path)
	assert.NoError(t, err)
	if err != nil {
		return image.Rectangle{}, nil
	}
	defer f.Close()

	raw, _, err := image.Decode(f)
	assert.NoError(t, err)

	pixels, err := pixels(raw)
	assert.NoError(t, err)

	if raw == nil {
		return image.Rectangle{}, nil
	}

	return raw.Bounds(), pixels
}

func pixels(i image.Image) ([]uint8, error) {
	switch i := i.(type) {
	case *image.RGBA:
		return i.Pix, nil
	case *image.NRGBA:
		return i.Pix, nil
	default:
		return nil, fmt.Errorf("unsupported image type %T", i)
	}
}
