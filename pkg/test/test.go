package test

import (
	"fmt"
	"image"
	_ "image/png"
	"math"
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

	assert.Len(t, resultPixels, len(masterPixels))
	if len(resultPixels) == 0 || len(masterPixels) == 0 {
		return
	}
	assert.Equal(t, masterBound, resultBound)
	if masterBound != resultBound {
		return
	}

	totalError := compare(masterPixels, resultPixels, 0.05)
	assert.Equal(t, float64(0), totalError)

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

func compare(img1, img2 []uint8, authorizedComponentError float32) float64 {
	overallChange := int64(0)

	for i := 0; i < len(img1); i++ {
		overallChange += diffUInt8(img1[i], img2[i], authorizedComponentError)
	}

	return math.Sqrt(float64(overallChange)) / float64(len(img1))
}

func diffUInt8(x, y uint8, authorizedComponentError float32) int64 {
	delta := float32(x - y)
	pow := delta * delta
	percentage := pow / float32(255*255)
	if percentage > authorizedComponentError {
		return int64(pow)
	}
	return 0
}
