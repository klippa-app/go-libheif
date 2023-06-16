# go-libheif - A go gettable decoder/converter for HEIC/HEIF/AVIF based on libheif

This package allows you to handle the following file formats using the Go image package:

- image/heic:         HEIF file using h265 compression
- image/heif:           HEIF file using any other compression
- image/heic-sequence:  HEIF image sequence using h265 compression
- image/heif-sequence:  HEIF image sequence using any other compression
- image/avif:           AVIF image
- image/avif-sequence:  AVIF sequence

It processes the images in a subprocess to allow for a safe way to handle the images.

The package also contains a simple method (`library.RenderFile`) to go from binary data of one of the file formats above
to a JPEG or PNG.
The reason for this helper method is speed, since sending raw images over RPC between the subprocess and the main
process can be quite slow, this method can be useful if you just want to convert something.

## Install dependencies

You can install libheif, libede265 and libaom from any source, but please remember that package managers might contain
outdated versions.
The Go bindings needs a pretty recent version, and this package references v1.14.2 (latest PPA version) currently, but
it might work with older versions.

You can use the strukturag PPA to get very recent versions:

```bash
sudo add-apt-repository ppa:strukturag/libde265 
sudo add-apt-repository ppa:strukturag/libheif 
sudo apt install libheif-dev
```

## Install

```go get github.com/klippa-app/go-libheif```

- Code Sample

First make a worker package/binary in the golibheif_plugin directory.

```
package main

import "github.com/klippa-app/go-libheif/library/plugin"

func main() {
	plugin.StartPlugin()
}
```

Then use the worker file/binary in your program.
If you want to make it run through go, use the example below.

You can also go build `golibheif_plugin/main.go` and then reference it in BinPath, this is the advices run method for
deployments.

```
package main

import (
	"log"

	"github.com/klippa-app/go-libheif"
	"github.com/klippa-app/go-libheif/library"
)

func init() {
	err := libheif.Init(libheif.Config{LibraryConfig: library.Config{
		Command: library.Command{
			BinPath: "go",
			Args:    []string{"run", "library/worker_example/main.go"},
		},
	}})
	if err != nil {
		log.Fatalf("could not start libheif worker: %s", err.Error())
	}
}

func main() {
	// Load a file from somewhere and put it in reader `r`.
	img, err := libheif.DecodeImage(r)
	if err != nil {
		// Handle error.
	}

	// Do something with `img` here.

	// Once you have imported `github.com/klippa-app/go-libheif`, you can also use the image package to decode.
}
```

You can also check the convert package for a full example.

## What is done

- Includes libheif using pkg-config and a simple golang binding

- Processes the images in a subprocess to prevent crashing the main application on segfaults

- A utility `convert` to illustrate the usage.

- Registers an `image` handler so that the Go `image` package can handle it

## License

- libheif and the libraries that it includes are in their own licenses (LGPL)

- go-libheif is in MIT license

## Credits

- libheif (https://github.com/strukturag/libheif)


