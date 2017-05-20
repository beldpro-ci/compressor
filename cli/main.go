package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/beldpro-ci/compressor"
)

var out = flag.String("out", "", "Destination for the compressed file")
var in = flag.String("in", "", "Input to be compressed")

func main() {
	flag.Parse()
	fmt.Println("")
	fmt.Println("    Starting compressor")
	fmt.Println("        IN: ", *in)
	fmt.Println("        OUT: ", *out)
	fmt.Println("")

	if *in == "" || *out == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	inputAbs, err := filepath.Abs(*in)
	if err != nil {
		log.Panicf("Can't get absolute path of %s", *in)
	}

	_, err = os.Stat(inputAbs)
	if err != nil {
		if os.IsNotExist(err) {
			log.Panicf("Input at %s does not exist",
				inputAbs)
		}

		log.Panicf("Unexpected error looking for file %s",
			inputAbs)
	}

	outputAbs, err := filepath.Abs(*out)
	if err != nil {
		log.Panicf("Can't get absolute path of %s", *out)
	}

	_, err = os.Stat(outputAbs)
	if err == nil {
		log.Panicf("A file or directory already exists at %s", outputAbs)
	}

	outputFile, err := os.Create(outputAbs)
	if err != nil {
		log.Panicf("Couldn't create output file %s", outputAbs)
	}
	defer outputFile.Close()

	err = compressor.TarGz.MakeBytes([]string{
		inputAbs,
	}, outputFile, nil)
	if err != nil {
		log.Panicf("Errored compressing %s", inputAbs)
	}

	fmt.Println("Done!")
}
