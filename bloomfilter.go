package bloomfilter



type BloomFilter struct {
	bitBuffer []uint64
	size      uint64
	hash1     Hasher
	hash2     Hasher
	mask      []uint64
}



func (bf *BloomFilter) Default() {
	bf.init(defaultSize)
	bf.hash1 = &SimpleHasher{prefix: randomString1}
	bf.hash2 = &SimpleHasher{prefix: randomString2}
}

func (bf *BloomFilter) New(size uint64, hasher1 Hasher, hasher2 Hasher) {

	bf.init(size)
	bf.hash1 = hasher1
	bf.hash2 = hasher2

}

func (bf *BloomFilter) init(size uint64) {
	var bitBufferSize uint64
	bitBufferSize = 1
	if size > 64 {
		bitBufferSize = (size / 64) + 1
	}
	bf.bitBuffer = make([]uint64, bitBufferSize)

	bf.size = size

	// init the mask
	bf.mask = make([]uint64, 64)
	bf.mask[0] = 1

	// create a 64 bit mask just to speed lookup & writes.
	for i := 1; i < 64; i++ {
		bf.mask[i] = bf.mask[i-1] << 1
	}
}

func (bf *BloomFilter) Add(fo FilterObject) {
	h1 := bf.hash1.Hash(fo) % bf.size
	h2 := bf.hash2.Hash(fo) % bf.size
	c1 := h1 / 64
	c2 := h2 / 64
	o1 := h1 % 64
	o2 := h2 % 64

	bf.bitBuffer[c1] = bf.bitBuffer[c1] | bf.mask[o1]
	bf.bitBuffer[c2] = bf.bitBuffer[c2] | bf.mask[o2]
}

func (bf *BloomFilter) Exists(fo FilterObject) bool {
	h1 := bf.hash1.Hash(fo) % bf.size
	h2 := bf.hash2.Hash(fo) % bf.size
	c1 := h1 / 64
	c2 := h2 / 64
	o1 := h1 % 64
	o2 := h2 % 64
	return (bf.bitBuffer[c1]&bf.mask[o1] > 0) && (bf.bitBuffer[c2]&bf.mask[o2] > 0)

}

func (bf *BloomFilter) Clear() {
	bf.init(bf.size)
}

