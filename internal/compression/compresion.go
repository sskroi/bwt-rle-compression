package compression

import (
	"bytes"
	"encoding/binary"
	"sort"
	"sync"
)

const (
	blockSize = 1024 * 20
)

type rotation struct {
	Data []byte
	Num  int
}

type block struct {
	Data []byte
	Num  int
}

func createBwtBlock(data []byte) ([]byte, int) {
	n := len(data)

	doubleData := make([]byte, n*2)
	copy(doubleData, data)
	copy(doubleData[n:], data)

	rotations := make([]rotation, 0, n)

	for i := 0; i < n; i++ {
		rotations = append(rotations, rotation{doubleData[i : i+n], i})
	}

	sort.Slice(rotations, func(i, j int) bool {
		return bytes.Compare(rotations[i].Data, rotations[j].Data) < 0
	})

	var resNum int
	resData := make([]byte, 0, len(data))

	for i := 0; i < n; i++ {
		resData = append(resData, rotations[i].Data[n-1])

		if rotations[i].Num == 0 {
			resNum = i
		}
	}

	return resData, resNum
}

func reverseBwtBlock(data []byte, num int) []byte {
	if len(data) == 0 {
		return []byte{}
	}

	n := len(data)

	count := make([]int, 256)

	cntEqualCurBeforeCur := make([]int, n)

	for i := 0; i < n; i++ {
		cntEqualCurBeforeCur[i] = count[data[i]]

		count[data[i]]++
	}

	prefixSums := make([]int, 256)
	prefixSums[0] = count[0]
	for i := 1; i < 256; i++ {
		prefixSums[i] = prefixSums[i-1] + count[i]
	}

	cntLessThenCur := map[byte]int{}
	for i := 1; i < 256; i++ {
		if count[i] != 0 {
			cntLessThenCur[byte(i)] = prefixSums[i-1]
		}
	}

	res := make([]byte, n)
	res[n-1] = data[num]

	prevIdx := num
	for i := n - 2; i >= 0; i-- {
		idx := cntEqualCurBeforeCur[prevIdx] + cntLessThenCur[data[prevIdx]]
		res[i] = data[idx]
		prevIdx = idx
	}

	return res
}

func createRleBlock(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	} else if len(data) == 1 {
		return []byte{1, data[0]}
	}

	res := []byte{}
	i := 0

	for i < len(data) {
		cnt := 1
		for i+1 < len(data) && data[i] == data[i+1] {
			cnt++
			i++
			if cnt == 127 {
				break
			}
		}
		if cnt > 1 {
			res = append(res, byte(cnt), data[i])
			i++
			continue
		}

		start := i
		for i+1 < len(data) && (i+1 == start || data[i] != data[i+1]) {
			cnt++
			i++
			if cnt == 128 {
				break
			}
		}
		res = append(res, byte(127+cnt))
		res = append(res, data[start:i+1]...)
		i++
	}

	return res
}

func reverseRleBlock(data []byte) []byte {
	res := []byte{}

	i := 0
	for i < len(data) {
		if data[i] < 128 {
			for j := byte(0); j < data[i]; j++ {
				res = append(res, data[i+1])
			}
			i += 2
		} else {
			n := data[i] - 127
			i++
			end := i + int(n)

			for i < end {
				res = append(res, data[i])
				i++
			}
		}
	}
	return res
}

func compressBlock(data []byte) []byte {
	res := make([]byte, 0)

	bwtProcessed, bwtNum := createBwtBlock(data)

	rleProcessed := createRleBlock(bwtProcessed)

	res = binary.BigEndian.AppendUint64(res, uint64(len(rleProcessed)))
	res = binary.BigEndian.AppendUint64(res, uint64(bwtNum))

	res = append(res, rleProcessed...)

	return res
}

// | rleBloack size | bwtNumber | rleBlock |
func CompressData(data []byte) []byte {
	if len(data) == 0 {
		res := []byte{}
		res = binary.BigEndian.AppendUint64(res, 0)
		res = binary.BigEndian.AppendUint64(res, 0)
		return res
	}

	wg := sync.WaitGroup{}
	ch := make(chan block)

	blockNum := 0
	for i := 0; i < len(data); i += blockSize {

		wg.Add(1)
		go func(blockNum int) {
			defer wg.Done()

			compressedBlock := compressBlock(data[i:min(i+blockSize, len(data))])

			ch <- block{compressedBlock, blockNum}
		}(blockNum)

		blockNum++
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	compressedBlocks := make([]block, 0)
	for blck := range ch {
		compressedBlocks = append(compressedBlocks, blck)
	}

	sort.Slice(compressedBlocks, func(i, j int) bool {
		return compressedBlocks[i].Num < compressedBlocks[j].Num
	})

	res := make([]byte, 0)
	for _, blck := range compressedBlocks {
		res = append(res, blck.Data...)
	}

	return res
}

func DecompressData(data []byte) []byte {
	res := make([]byte, 0)

	i := 0
	for i < len(data) {
		rleBlockSize := int(binary.BigEndian.Uint64(data[i : i+8]))
		bwtNum := int(binary.BigEndian.Uint64(data[i+8 : i+16]))

		extracedRle := reverseRleBlock(data[i+16 : i+rleBlockSize+16])
		extracedBwt := reverseBwtBlock(extracedRle, bwtNum)

		res = append(res, extracedBwt...)

		i += (rleBlockSize + 16)
	}

	return res
}
