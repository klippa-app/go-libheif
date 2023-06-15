package plugin

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/klippa-app/go-libheif/library/requests"
	"github.com/klippa-app/go-libheif/library/responses"
	"github.com/klippa-app/go-libheif/library/shared"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/hashicorp/go-plugin"
	_ "github.com/strukturag/libheif/go/heif"
)

func init() {
	// Needed to serialize the image interface.
	gob.Register(&image.YCbCr{})
	gob.Register(&image.RGBA64{})
	gob.Register(&image.RGBA{})
	gob.Register(&image.Gray{})
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "libheif",
}

func StartPlugin() {
	var pluginMap = map[string]plugin.Plugin{
		"libheif": &shared.LibheifPlugin{Impl: &libHeifImplementation{}},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}

type libHeifImplementation struct{}

func (l *libHeifImplementation) Ping() (string, error) {
	return "Pong", nil
}

func (l *libHeifImplementation) DecodeImage(request *requests.DecodeImage) (*responses.DecodeImage, error) {
	decodedImage, format, err := image.Decode(bytes.NewReader(*request.Data))
	if err != nil {
		return nil, err
	}

	return &responses.DecodeImage{
		Format: format,
		Image:  decodedImage,
	}, nil
}

func (l *libHeifImplementation) DecodeConfig(request *requests.DecodeConfig) (*responses.DecodeConfig, error) {
	config, format, err := image.DecodeConfig(bytes.NewReader(*request.Data))
	if err != nil {
		return nil, err
	}
	return &responses.DecodeConfig{
		Format: format,
		Config: config,
	}, nil
}

func (l *libHeifImplementation) RenderFile(request *requests.RenderFile) (*responses.RenderFile, error) {
	decodedImage, _, err := image.Decode(bytes.NewReader(*request.Data))
	if err != nil {
		return nil, err
	}

	var imgBuf bytes.Buffer
	if request.OutputFormat == requests.RenderFileOutputFormatJPG {
		var opt jpeg.Options
		opt.Quality = 95

		for {
			err := jpeg.Encode(&imgBuf, decodedImage, &opt)
			if err != nil {
				return nil, err
			}

			if request.MaxFileSize == 0 || int64(imgBuf.Len()) < request.MaxFileSize {
				break
			}

			opt.Quality -= 10

			if opt.Quality <= 45 {
				return nil, errors.New("image would exceed maximum filesize")
			}

			imgBuf.Reset()
		}
	} else if request.OutputFormat == requests.RenderFileOutputFormatPNG {
		err := png.Encode(&imgBuf, decodedImage)
		if err != nil {
			return nil, err
		}

		if request.MaxFileSize != 0 && int64(imgBuf.Len()) > request.MaxFileSize {
			return nil, errors.New("image would exceed maximum filesize")
		}
	} else {
		return nil, errors.New("invalid output format given")
	}

	output := imgBuf.Bytes()

	return &responses.RenderFile{
		Output: &output,
	}, nil
}
