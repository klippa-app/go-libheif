package responses

import "image"

type DecodeImage struct {
	Format string
	Image  image.Image
}

type DecodeConfig struct {
	Format string
	Config image.Config
}

type RenderFile struct {
	Output *[]byte
}
