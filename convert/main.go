package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"log"
	"os"

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
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "usage: convert <in-file> <out-file> \n")
		os.Exit(1)
	}

	fin, fout := flag.Arg(0), flag.Arg(1)
	fi, err := os.Open(fin)
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	img, err := libheif.DecodeImage(fi)
	if err != nil {
		log.Fatalf("Failed to parse %s: %v\n", fin, err)
	}

	fo, err := os.OpenFile(fout, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to create output file %s: %v\n", fout, err)
	}
	defer fo.Close()

	err = jpeg.Encode(fo, img, nil)
	if err != nil {
		log.Fatalf("Failed to encode %s: %v\n", fout, err)
	}

	log.Printf("Convert %s to %s successfully\n", fin, fout)
}
