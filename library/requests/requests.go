package requests

type DecodeImage struct {
	Data *[]byte
}

type DecodeConfig struct {
	Data *[]byte
}

type RenderFileOutputFormat string // The file format to render output as.

const (
	RenderFileOutputFormatJPG RenderFileOutputFormat = "jpg" // Render the file as a JPEG file.
	RenderFileOutputFormatPNG RenderFileOutputFormat = "png" // Render the file as a PNG file.
)

type RenderFile struct {
	Data         *[]byte                // The file data.
	OutputFormat RenderFileOutputFormat // The format to output the image as
	MaxFileSize  int64                  // The maximum filesize, if jpg is chosen as output format, it will try to compress it until it fits
}
