package compression

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestCompression(t *testing.T) {
	dirWithTestFile := "test_files"

	dir, err := os.Open(dirWithTestFile)
	if err != nil {
		t.Fatal("os.Open:", err)
	}
	defer dir.Close()

	files, err := dir.ReadDir(-1)
	if err != nil {
		t.Fatal("ReadDir:", err)
	}

	for _, file := range files {
		if file.IsDir() {
			t.Fatal("There should be no directories in the directory with text files")
		}

		fmt.Println("\""+file.Name()+"\"", "- test file")

		testFile, err := os.Open(filepath.Join(dirWithTestFile, file.Name()))
		if err != nil {
			t.Fatal("os.Open:", err)
		}

		origData, err := io.ReadAll(testFile)
		if err != nil {
			t.Fatal("io.ReadAll:", err)
		}

		sizeBeforeCompress := getPrettySizeOutput(len(origData))
		fmt.Println(sizeBeforeCompress, "- size before compress")

		startTime := time.Now().UnixMilli()
		compressedData := CompressData(origData)
		fmt.Println(time.Now().UnixMilli()-startTime, "ms", "- time to compress")

		sizeAfterCompess := getPrettySizeOutput(len(compressedData))
		fmt.Println(sizeAfterCompess, "- size after compress")

		startTime = time.Now().UnixMilli()
		decompressData := DecompressData(compressedData)
		fmt.Print(time.Now().UnixMilli()-startTime, " ms", " - time to decompress\n\n")

		if bytes.Compare(origData, decompressData) != 0 {
			t.Fatal("Incorrect result for file:", testFile.Name())
		}
	}
}

func getPrettySizeOutput(size int) string {
	floatSize := float64(size)

	kib := 1024
	mib := kib * 1024

	var s string

	if size < kib {
		s = strconv.FormatInt(int64(size), 10) + " bytes"
	} else if size < mib {
		s = strconv.FormatFloat(floatSize/float64(kib), 'f', 2, 64) + " KiB"
	} else {
		s = strconv.FormatFloat(floatSize/float64(mib), 'f', 2, 64) + " MiB"
	}

	return s
}

func TestBwt(t *testing.T) {
	tests := [][]byte{
		{1},
		{0},
		{0, 1},
		{255, 255},
		{0, 0},
	}

	for _, test := range tests {
		postBwt, bwtNum := createBwtBlock(test)

		res := reverseBwtBlock(postBwt, bwtNum)

		if bytes.Compare(res, test) != 0 {
			t.Fatal("\nExpected:", test, "\nGiven:", res)
		}
	}
}

func TestRle(t *testing.T) {
	tests := [][]byte{
		{1},
		{0},
		{0, 1},
		{255, 255},
		{0, 0},
	}

	for _, test := range tests {
		postRle := createRleBlock(test)

		res := reverseRleBlock(postRle)

		if bytes.Compare(res, test) != 0 {
			t.Fatal("\nExpected:", test, "\nGiven:", res)
		}
	}
}
