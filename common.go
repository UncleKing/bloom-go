package bloomfilter

type Hasher interface {
	Hash(fo FilterObject) uint64
}

type FilterObject interface {
	String() string
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

