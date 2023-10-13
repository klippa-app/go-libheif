//go:build !go_libheif_use_turbojpeg

package image_jpeg

import (
	"bytes"
	"image"
	"image/jpeg"
	"testing"
)

func TestEncode(t *testing.T) {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{100, 100}})
	testWriter := bytes.NewBuffer(nil)
	err := Encode(testWriter, img, Options{})
	if err != nil {
		t.Fatalf("Encode resulted in error: %s", err.Error())
	}
	if testWriter.Len() != 789 {
		t.Fatalf("Encode resulted in wrong byte result, got %d, want %d", testWriter.Len(), 789)
	}
}

func TestEncodeQuality(t *testing.T) {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{100, 100}})
	testWriter := bytes.NewBuffer(nil)
	err := Encode(testWriter, img, Options{
		Options: &jpeg.Options{
			Quality: 100,
		},
	})
	if err != nil {
		t.Fatalf("Encode resulted in error: %s", err.Error())
	}
	if testWriter.Len() != 791 {
		t.Fatalf("Encode resulted in wrong byte result, got %d, want %d", testWriter.Len(), 791)
	}
}
