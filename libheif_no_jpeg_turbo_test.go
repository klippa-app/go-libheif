//go:build !go_libheif_use_turbojpeg

package libheif

import (
	"github.com/klippa-app/go-libheif/library"
	_ "image/jpeg"
	_ "image/png"
)

func initLib() error {
	err := Init(Config{LibraryConfig: library.Config{
		Command: library.Command{
			BinPath: "go",
			Args:    []string{"run", "library/worker_example/main.go"},
		},
	}})
	if err != nil {
		return err
	}
	return nil
}
