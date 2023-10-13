//go:build !go_libheif_use_turbojpeg

package image_jpeg

import (
	"image"
	"image/jpeg"
	"io"
)

func Encode(w io.Writer, m image.Image, o Options) error {
	return jpeg.Encode(w, m, o.Options)
}
