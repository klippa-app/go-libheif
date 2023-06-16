package libheif

import (
	"errors"
	"image"
	"io"
	"sync"

	"github.com/klippa-app/go-libheif/library"
)

func DecodeImage(r io.Reader) (image.Image, error) {
	if !isInitialized {
		return nil, NotInitializedError
	}

	return library.DecodeImage(r)
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	var config image.Config

	if !isInitialized {
		return config, NotInitializedError
	}

	return library.DecodeConfig(r)
}

var NotInitializedError = errors.New("goheif was not initialized, you must call the Init() method")
var isInitialized = false
var initLock = sync.Mutex{}

type Config struct {
	LibraryConfig library.Config
}

func Init(config Config) error {
	initLock.Lock()
	defer initLock.Unlock()
	if isInitialized {
		return nil
	}
	err := library.Init(config.LibraryConfig)
	if err != nil {
		return err
	}
	isInitialized = true

	return nil
}

func DeInit() {
	initLock.Lock()
	defer initLock.Unlock()

	if !isInitialized {
		return
	}

	library.DeInit()
	isInitialized = true
}

func init() {
	image.RegisterFormat("heif", "????ftypheic", DecodeImage, DecodeConfig)
	image.RegisterFormat("heif", "????ftypheim", DecodeImage, DecodeConfig)
	image.RegisterFormat("heif", "????ftypheis", DecodeImage, DecodeConfig)
	image.RegisterFormat("heif", "????ftypheix", DecodeImage, DecodeConfig)
	image.RegisterFormat("heif", "????ftyphevc", DecodeImage, DecodeConfig)
	image.RegisterFormat("heif", "????ftyphevm", DecodeImage, DecodeConfig)
	image.RegisterFormat("heif", "????ftyphevs", DecodeImage, DecodeConfig)
	image.RegisterFormat("heif", "????ftypmif1", DecodeImage, DecodeConfig)
	image.RegisterFormat("avif", "????ftypavif", DecodeImage, DecodeConfig)
	image.RegisterFormat("avif", "????ftypavis", DecodeImage, DecodeConfig)
}
