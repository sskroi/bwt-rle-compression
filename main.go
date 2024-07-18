package main

import (
	"bwtrlecompr/internal/compression"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	usageText = "Usage:\n  bwt-rle-comp <file to compress>\n  bwt_rle_comp -d <file to decompress>\n\nFlags:\n  -d <file> - decompress <file>"
)

func openFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func writeResult(data []byte, filename string) error {
	if len(filename) > 14 && filename[len(filename)-14:] == ".compr.decompr" {
		filename = filename[:len(filename)-14] + ".decompr"
	}

	if _, err := os.Stat(filename); !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file already exists")
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func decompressFile(path string) error {
	data, err := openFile(path)
	if err != nil {
		return fmt.Errorf("can't open file %q: %w", path, err)
	}

	decompressedData := compression.DecompressData(data)
	return writeResult(decompressedData, path+".decompr")
}

func compressFile(path string) error {
	data, err := openFile(path)
	if err != nil {
		return fmt.Errorf("can't open file %q: %w", path, err)
	}

	compressedData := compression.CompressData(data)
	return writeResult(compressedData, path+".compr")
}

func main() {
	toDecompress := flag.String("d", "", "compressed file to decompress")

	flag.Parse()

	if len(flag.Args()) < 1 && *toDecompress == "" {
		fmt.Println(usageText)
		return
	}

	if *toDecompress != "" {
		if err := decompressFile(*toDecompress); err != nil {
			fmt.Println(err)
		}
	} else {
		toCompress := flag.Arg(0)
		if err := compressFile(toCompress); err != nil {
			fmt.Println(err)
		}
	}
}
