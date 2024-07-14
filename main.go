package main

import (
	"bwtrlecompr/internal/compression"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/spf13/pflag"
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

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func writeResult(data []byte, filename string) error {
	if len(filename) > 14 && filename[len(filename)-14:] == ".compr.decompr" {
		filename = filename[:len(filename)-14] + ".decompr"
	}

	_, err := os.Stat(filename)
	if !errors.Is(err, fs.ErrNotExist) {
		return errors.New("File alreay exists")
	}

	file, err := os.Create(filename)

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	toDecompress := pflag.StringP("decompress", "d", "", "compressed file to decompress")

	pflag.Parse()

	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println(usageText)
		return
	}

	if pflag.Lookup("decompress").Changed {
		dataToDecompress, err := openFile(*toDecompress)
		if err != nil {
			fmt.Println("Can't open file \""+*toDecompress+"\":", err)
			return
		}

		decompressedData := compression.DecompressData(dataToDecompress)

		err = writeResult(decompressedData, *toDecompress+".decompr")
		if err != nil {
			fmt.Println("Can't write result:", err)
			return
		}

		return
	} else {
		toCompress := os.Args[1]

		dataToCompress, err := openFile(toCompress)
		if err != nil {
			fmt.Println("Can't open file \""+*toDecompress+"\":", err)
			return
		}

		compressedData := compression.CompressData(dataToCompress)

		err = writeResult(compressedData, toCompress+".compr")
		if err != nil {
			fmt.Println("Can't write result:", err)
			return
		}

		return
	}
}
