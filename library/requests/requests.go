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
	Data          *[]byte                // The file data.
	OutputFormat  RenderFileOutputFormat // The format to output the image as
	MaxFileSize   int64                  // Only used when OutputFormat RenderFileOutputFormatJPG. The maximum filesize, if jpg is chosen as output format, it will try to lower the quality it until it fits.
	OutputQuality int                    // Only used when OutputFormat RenderFileOutputFormatJPG. Ranges from 1 to 100 inclusive, higher is better. The default is 95.
	Progressive   bool                   // Only used when OutputFormat RenderFileOutputFormatJPG and with build tag go_libheif_use_turbojpeg. Will render a progressive jpeg.
}
