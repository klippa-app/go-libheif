//go:build go_libheif_use_turbojpeg

package image_jpeg

import (
	"bufio"
	"image"
	"image/jpeg"
	"io"
	"unsafe"
)

/*
#cgo pkg-config: libturbojpeg
#include <turbojpeg.h>
*/
import "C"
import "fmt"

type Sampling C.int

const (
	Sampling444  Sampling = C.TJSAMP_444
	Sampling422  Sampling = C.TJSAMP_422
	Sampling420  Sampling = C.TJSAMP_420
	SamplingGray Sampling = C.TJSAMP_GRAY
)

type PixelFormat C.int

const (
	PixelFormatRGB     PixelFormat = C.TJPF_RGB
	PixelFormatBGR     PixelFormat = C.TJPF_BGR
	PixelFormatRGBX    PixelFormat = C.TJPF_RGBX
	PixelFormatBGRX    PixelFormat = C.TJPF_BGRX
	PixelFormatXBGR    PixelFormat = C.TJPF_XBGR
	PixelFormatXRGB    PixelFormat = C.TJPF_XRGB
	PixelFormatGRAY    PixelFormat = C.TJPF_GRAY
	PixelFormatRGBA    PixelFormat = C.TJPF_RGBA
	PixelFormatBGRA    PixelFormat = C.TJPF_BGRA
	PixelFormatABGR    PixelFormat = C.TJPF_ABGR
	PixelFormatARGB    PixelFormat = C.TJPF_ARGB
	PixelFormatCMYK    PixelFormat = C.TJPF_CMYK
	PixelFormatUNKNOWN PixelFormat = C.TJPF_UNKNOWN
)

type Flags C.int

const (
	FlagAccurateDCT   Flags = C.TJFLAG_ACCURATEDCT
	FlagBottomUp      Flags = C.TJFLAG_BOTTOMUP
	FlagFastDCT       Flags = C.TJFLAG_FASTDCT
	FlagFastUpsample  Flags = C.TJFLAG_FASTUPSAMPLE
	FlagNoRealloc     Flags = C.TJFLAG_NOREALLOC
	FlagProgressive   Flags = C.TJFLAG_PROGRESSIVE
	FlagStopOnWarning Flags = C.TJFLAG_STOPONWARNING
)

func makeError(handler C.tjhandle, returnVal C.int) error {
	if returnVal == 0 {
		return nil
	}
	str := C.GoString(C.tjGetErrorStr2(handler))
	return fmt.Errorf("turbojpeg error: %v", str)
}

type Image struct {
	Width  int
	Height int
	Stride int
	Pixels []byte
}

type CompressParams struct {
	PixelFormat PixelFormat
	Sampling    Sampling
	Quality     int // 1 .. 100
	Flags       Flags
}

func MakeCompressParams(pixelFormat PixelFormat, sampling Sampling, quality int, flags Flags) CompressParams {
	return CompressParams{
		PixelFormat: pixelFormat,
		Sampling:    sampling,
		Quality:     quality,
		Flags:       flags,
	}
}

func Compress(img *Image, params CompressParams) ([]byte, error) {
	encoder := C.tjInitCompress()
	defer C.tjDestroy(encoder)

	var outBuf *C.uchar
	var outBufSize C.ulong

	// int tjCompress2(tjhandle handle, const unsigned char *srcBuf, int width, int pitch, int height, int pixelFormat,
	// unsigned char **jpegBuf, unsigned long *jpegSize, int jpegSubsamp, int jpegQual, int flags);
	res := C.tjCompress2(encoder, (*C.uchar)(&img.Pixels[0]), C.int(img.Width), C.int(img.Stride), C.int(img.Height), C.int(params.PixelFormat),
		&outBuf, &outBufSize, C.int(params.Sampling), C.int(params.Quality), C.int(params.Flags))

	var enc []byte
	err := makeError(encoder, res)
	if outBuf != nil {
		enc = C.GoBytes(unsafe.Pointer(outBuf), C.int(outBufSize))
		C.tjFree(outBuf)
	}

	if err != nil {
		return nil, err
	}
	return enc, nil
}

func Encode(w io.Writer, m image.Image, o Options) error {
	imageWriter := bufio.NewWriter(w)

	// Clip quality to [1, 100].
	quality := jpeg.DefaultQuality
	if o.Options != nil {
		quality = o.Options.Quality
		if quality < 1 {
			quality = 1
		} else if quality > 100 {
			quality = 100
		}
	}

	raw := FromImage(m, true)

	flags := Flags(0)
	if o.Progressive {
		flags |= FlagProgressive
	}

	params := MakeCompressParams(PixelFormatRGBA, Sampling420, quality, flags)
	jpg, err := Compress(raw, params)
	if err != nil {
		return err
	}

	_, err = imageWriter.Write(jpg)
	if err != nil {
		return err
	}

	err = imageWriter.Flush()
	if err != nil {
		return err
	}

	return nil
}

// Convert a Go image.Image into a turbo.Image
// If allowDeepClone is true, and the source image is type NRGBA or RGBA,
// then the resulting Image points directly to the pixel buffer of the source image.
func FromImage(src image.Image, allowDeepClone bool) *Image {
	dst := &Image{
		Width:  src.Bounds().Dx(),
		Height: src.Bounds().Dy(),
		Stride: src.Bounds().Dx() * 4,
	}
	switch v := src.(type) {
	case *image.RGBA:
		if allowDeepClone {
			dst.Pixels = v.Pix
		} else {
			dst.Pixels = make([]byte, dst.Stride*dst.Height)
			copy(dst.Pixels, v.Pix)
		}
		return dst
	case *image.NRGBA:
		if allowDeepClone {
			dst.Pixels = v.Pix
		} else {
			dst.Pixels = make([]byte, dst.Stride*dst.Height)
			copy(dst.Pixels, v.Pix)
		}
		return dst
	}

	// This must be super slow - I haven't tested
	dst.Pixels = make([]byte, dst.Stride*dst.Height)
	p := 0
	for y := 0; y < dst.Height; y++ {
		for x := 0; x < dst.Width; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			dst.Pixels[p] = byte(r >> 8)
			dst.Pixels[p+1] = byte(g >> 8)
			dst.Pixels[p+2] = byte(b >> 8)
			dst.Pixels[p+3] = byte(a >> 8)
			p += 4
		}
	}

	return dst
}
