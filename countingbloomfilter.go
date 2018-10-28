package bloomfilter



type CountingBloomFilter struct {
	byteBuffer []byte
	size       uint64
	hash1      Hasher
	hash2      Hasher
}


func (cbf *CountingBloomFilter) Default() {
	cbf.init(defaultSize)
	cbf.hash1 = &SimpleHasher{prefix: randomString1}
	cbf.hash2 = &SimpleHasher{prefix: randomString2}
}

func (cbf *CountingBloomFilter) New(size uint64, hasher1 Hasher, hasher2 Hasher) {

	cbf.init(size)
	cbf.hash1 = hasher1
	cbf.hash2 = hasher2

}

func (cbf *CountingBloomFilter) init(size uint64) {
	cbf.size = size
	cbf.byteBuffer = make([]byte, size)
}

func (cbf *CountingBloomFilter) Add(fo FilterObject) {
	h1 := cbf.hash1.Hash(fo) % cbf.size
	h2 := cbf.hash2.Hash(fo) % cbf.size
	cbf.byteBuffer[h1] = cbf.byteBuffer[h1] + 1
	cbf.byteBuffer[h2] = cbf.byteBuffer[h2] + 1
}

func (cbf *CountingBloomFilter) Exists(fo FilterObject) bool {
	h1 := cbf.hash1.Hash(fo) % cbf.size
	h2 := cbf.hash2.Hash(fo) % cbf.size
	return cbf.byteBuffer[h1] > 0 && cbf.byteBuffer[h2] > 0
}

func (cbf *CountingBloomFilter) Delete(fo FilterObject) bool {
	h1 := cbf.hash1.Hash(fo) % cbf.size
	h2 := cbf.hash2.Hash(fo) % cbf.size
	if cbf.byteBuffer[h1] > 0 && cbf.byteBuffer[h2] > 0 {
		cbf.byteBuffer[h1] = cbf.byteBuffer[h1] - 1
		cbf.byteBuffer[h2] = cbf.byteBuffer[h2] - 1
		return true
	}
	return false
}

func (cbf *CountingBloomFilter) Clear() {
	cbf.init(cbf.size)
}
