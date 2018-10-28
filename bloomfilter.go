package bloomfilter

type Hasher interface {
	Hash(fo FilterObject) uint64
}

type FilterObject interface {
	String() string
}

type NaiveBloomFilter struct {
	bitBuffer []uint64
	size      uint64
	hash1     Hasher
	hash2     Hasher
	mask      []uint64
}

type NaiveCountingBloomFilter struct {
	byteBuffer []byte
	size       uint64
	hash1      Hasher
	hash2      Hasher
}
type SimpleHasher struct {
	prefix string
}

type DefaultFilterObject struct {
	str string
}

const defaultSize = 512 * 1024 * 1024
const randomString1 = "a;lkjfi234lkasdoiquyw34lklsbfvhweoi4u5y4i2egslkgo3i23t4eoqihsboseroret3ASDFAADSFASDt"
const randomString2 = "af;lkjlahbldkjfh'jkyu'ueyket;fsjhbkhxgzcjbvjhaGWEQWEFRGKERJHNLNV,NASDJHGAWFKWHTRLBNF"

func (dfo *DefaultFilterObject) String() string {
	return dfo.str
}

// Hash copies the Java hash approach to find a key.
func (sh *SimpleHasher) Hash(fo FilterObject) uint64 {
	str := sh.prefix + fo.String()
	var val uint64 = 1
	for i := 0; i < len(str); i++ {
		val += (val * 31) + uint64(str[i])
	}
	return val
}

func (bf *NaiveBloomFilter) Default() {
	bf.init(defaultSize)
	bf.hash1 = &SimpleHasher{prefix: randomString1}
	bf.hash2 = &SimpleHasher{prefix: randomString2}
}

func (bf *NaiveBloomFilter) New(size uint64, hasher1 Hasher, hasher2 Hasher) {

	bf.init(size)
	bf.hash1 = hasher1
	bf.hash2 = hasher2

}

func (nbf *NaiveBloomFilter) init(size uint64) {
	var bitBufferSize uint64
	bitBufferSize = 1
	if size > 64 {
		bitBufferSize = (size / 64) + 1
	}
	nbf.bitBuffer = make([]uint64, bitBufferSize)

	nbf.size = size

	// init the mask
	nbf.mask = make([]uint64, 64)
	nbf.mask[0] = 1

	// create a 64 bit mask just to speed lookup & writes.
	for i := 1; i < 64; i++ {
		nbf.mask[i] = nbf.mask[i-1] << 1
	}
}

func (nbf *NaiveBloomFilter) Add(fo FilterObject) {
	h1 := nbf.hash1.Hash(fo) % nbf.size
	h2 := nbf.hash2.Hash(fo) % nbf.size
	c1 := h1 / 64
	c2 := h2 / 64
	o1 := h1 % 64
	o2 := h2 % 64

	nbf.bitBuffer[c1] = nbf.bitBuffer[c1] | nbf.mask[o1]
	nbf.bitBuffer[c2] = nbf.bitBuffer[c2] | nbf.mask[o2]
}

func (nbf *NaiveBloomFilter) Exists(fo FilterObject) bool {
	h1 := nbf.hash1.Hash(fo) % nbf.size
	h2 := nbf.hash2.Hash(fo) % nbf.size
	c1 := h1 / 64
	c2 := h2 / 64
	o1 := h1 % 64
	o2 := h2 % 64
	return (nbf.bitBuffer[c1]&nbf.mask[o1] > 0) && (nbf.bitBuffer[c2]&nbf.mask[o2] > 0)

}

func (nbf *NaiveBloomFilter) Clear() {
	nbf.init(nbf.size)
}

func (bf *NaiveCountingBloomFilter) Default() {
	bf.init(defaultSize)
	bf.hash1 = &SimpleHasher{prefix: randomString1}
	bf.hash2 = &SimpleHasher{prefix: randomString2}
}

func (bf *NaiveCountingBloomFilter) New(size uint64, hasher1 Hasher, hasher2 Hasher) {

	bf.init(size)
	bf.hash1 = hasher1
	bf.hash2 = hasher2

}

func (nbf *NaiveCountingBloomFilter) init(size uint64) {
	nbf.size = size
	nbf.byteBuffer = make([]byte, size)
}

func (ncbf *NaiveCountingBloomFilter) Add(fo FilterObject) {
	h1 := ncbf.hash1.Hash(fo) % ncbf.size
	h2 := ncbf.hash2.Hash(fo) % ncbf.size
	ncbf.byteBuffer[h1] = ncbf.byteBuffer[h1] + 1
	ncbf.byteBuffer[h2] = ncbf.byteBuffer[h2] + 1
}

func (ncbf *NaiveCountingBloomFilter) Exists(fo FilterObject) bool {
	h1 := ncbf.hash1.Hash(fo) % ncbf.size
	h2 := ncbf.hash2.Hash(fo) % ncbf.size
	return ncbf.byteBuffer[h1] > 0 && ncbf.byteBuffer[h2] > 0
}

func (ncbf *NaiveCountingBloomFilter) Delete(fo FilterObject) bool {
	h1 := ncbf.hash1.Hash(fo) % ncbf.size
	h2 := ncbf.hash2.Hash(fo) % ncbf.size
	if ncbf.byteBuffer[h1] > 0 && ncbf.byteBuffer[h2] > 0 {
		ncbf.byteBuffer[h1] = ncbf.byteBuffer[h1] - 1
		ncbf.byteBuffer[h2] = ncbf.byteBuffer[h2] - 1
		return true
	}
	return false
}

func (ncbf *NaiveCountingBloomFilter) Clear() {
	ncbf.init(ncbf.size)
}
