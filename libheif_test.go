package libheif

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/klippa-app/go-libheif/library"
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

func TestFormatRegistered(t *testing.T) {
	err := initLib()
	if err != nil {
		t.Fatal(err)
	}

	type expectedFile struct {
		Name   string
		Format string
		Width  int
		Height int
	}
	files := []expectedFile{
		{
			Name:   "camel.heic",
			Format: "heif",
			Width:  1596,
			Height: 1064,
		},
		{
			Name:   "receipt.avif",
			Format: "avif",
			Width:  826,
			Height: 826,
		},
	}
	for _, file := range files {
		b, err := os.ReadFile(fmt.Sprintf("testdata/%s", file.Name))
		if err != nil {
			t.Fatal(err)
		}

		img, dec, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			t.Fatalf("unable to decode image %s: %s", file.Name, err)
		}

		if got, want := dec, file.Format; got != want {
			t.Errorf("unexpected decoder for %s: got %s, want %s", file.Name, got, want)
		}

		if w, h := img.Bounds().Dx(), img.Bounds().Dy(); w != file.Width || h != file.Height {
			t.Errorf("unexpected decoded image size for %s: got %dx%d, want %dx%d", file.Name, w, h, file.Width, file.Height)
		}
	}
}

func TestRenderJPEG(t *testing.T) {
	err := initLib()
	if err != nil {
		t.Fatal(err)
	}

	b, err := os.ReadFile("testdata/camel.heic")
	if err != nil {
		t.Fatal(err)
	}

	renderedFile, err := library.RenderFile(&b, library.RenderOptions{
		OutputFormat: library.RenderFileOutputFormatJPG,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := renderedFile.OriginalFormat, "heif"; got != want {
		t.Errorf("unexpected original format: got %s, want %s", got, want)
	}

	if got, want := renderedFile.NewFormat, "jpeg"; got != want {
		t.Errorf("unexpected new format: got %s, want %s", got, want)
	}

	if w, h := renderedFile.Width, renderedFile.Height; w != 1596 || h != 1064 {
		t.Errorf("unexpected rendered image size: got %dx%d, want 1596x1064", w, h)
	}

	img, dec, err := image.Decode(bytes.NewReader(*renderedFile.Output))
	if err != nil {
		t.Fatalf("unable to decode jpeg image: %s", err)
	}

	if got, want := dec, "jpeg"; got != want {
		t.Errorf("unexpected decoder: got %s, want %s", got, want)
	}

	if w, h := img.Bounds().Dx(), img.Bounds().Dy(); w != 1596 || h != 1064 {
		t.Errorf("unexpected decoded image size: got %dx%d, want 1596x1064", w, h)
	}
}

func TestRenderPNG(t *testing.T) {
	err := initLib()
	if err != nil {
		t.Fatal(err)
	}

	b, err := os.ReadFile("testdata/camel.heic")
	if err != nil {
		t.Fatal(err)
	}

	renderedFile, err := library.RenderFile(&b, library.RenderOptions{
		OutputFormat: library.RenderFileOutputFormatPNG,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := renderedFile.OriginalFormat, "heif"; got != want {
		t.Errorf("unexpected original format: got %s, want %s", got, want)
	}

	if got, want := renderedFile.NewFormat, "png"; got != want {
		t.Errorf("unexpected new format: got %s, want %s", got, want)
	}

	if w, h := renderedFile.Width, renderedFile.Height; w != 1596 || h != 1064 {
		t.Errorf("unexpected rendered image size: got %dx%d, want 1596x1064", w, h)
	}

	img, dec, err := image.Decode(bytes.NewReader(*renderedFile.Output))
	if err != nil {
		t.Fatalf("unable to decode jpeg image: %s", err)
	}

	if got, want := dec, "png"; got != want {
		t.Errorf("unexpected decoder: got %s, want %s", got, want)
	}

	if w, h := img.Bounds().Dx(), img.Bounds().Dy(); w != 1596 || h != 1064 {
		t.Errorf("unexpected decoded image size: got %dx%d, want 1596x1064", w, h)
	}
}

func Benchmark(b *testing.B) {
	err := initLib()
	if err != nil {
		b.Fatal(err)
	}

	f, err := ioutil.ReadFile("testdata/camel.heic")
	if err != nil {
		b.Fatal(err)
	}
	r := bytes.NewReader(f)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err = DecodeImage(r)
		if err != nil {
			b.Fatal(err)
		}

		r.Seek(0, io.SeekStart)
	}
}
