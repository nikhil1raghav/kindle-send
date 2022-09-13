package util

import (
	"fmt"
	"strconv"
)

func GetHash(name string) string{
	hash:=murmurHash64B([]byte(name), 0)
	hashStr:=strconv.Itoa(int(hash))

	//TODO: this works but is functionally incorrect, assign mimetype based on content
	return fmt.Sprintf("img%s.png",hashStr)
}
func murmurHash64B(key []byte, seed uint64) (hash uint64) {
	const m uint32 = 0x5bd1e995
	const r = 24

	var l int = len(key)
	var h1 uint32 = uint32(seed) ^ uint32(l)
	var h2 uint32 = uint32(seed) >> 32

	var data []byte = key

	var k1, k2 uint32

	for l >= 8 {
		k1 = uint32(data[0]) + uint32(data[1]) << 8 + uint32(data[2]) << 16 + uint32(data[3]) << 24
		k1 *= m; k1 ^= k1 >> r; k1 *=m
		h1 *= m; h1 ^= k1
		data = data[4:]
		l -= 4

		k2 = uint32(data[0]) + uint32(data[1]) << 8 + uint32(data[2]) << 16 + uint32(data[3]) << 24
		k2 *= m; k2 ^= k2 >> r; k2 *= m
		h2 *= m; h2 ^= k2
		data = data[4:]
		l -= 4
	}

	if l >= 4 {
		k1 = uint32(data[0]) + uint32(data[1]) << 8 + uint32(data[2]) << 16 + uint32(data[3]) << 24
		k1 *= m; k1 ^= k1 >> r; k1 *= m
		h1 *= m; h1 ^= k1
		data = data[4:]
		l -= 4
	}

	switch l {
	case 3:
		h2 ^= uint32(data[2]) << 16
		fallthrough
	case 2:
		h2 ^= uint32(data[1]) << 8
		fallthrough
	case 1:
		h2 ^= uint32(data[0])
		h2 *= m
	}

	h1 ^= h2 >> 18; h1 *= m
	h2 ^= h1 >> 22; h2 *= m
	h1 ^= h2 >> 17; h1 *= m
	h2 ^= h1 >> 19; h2 *= m

	var h uint64 = uint64(h1)

	h = (h << 32) | uint64(h2)

	return h

}
