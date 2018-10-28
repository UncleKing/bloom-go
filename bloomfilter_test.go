package bloomfilter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var words = []string{"Test1", "Test2", "Test3", "Test4"}

func TestBloomfilter(t *testing.T) {
	filter := BloomFilter{}
	filter.Default()

	for i := 0; i < len(words); i++ {

		filter.Add(&DefaultFilterObject{str: words[i]})
	}

	for i := 0; i < len(words); i++ {
		assert.True(t, filter.Exists(&DefaultFilterObject{str: words[i]}),
			fmt.Sprintf("'%s' should be marked\n", words[i]))
	}

	filter.Clear()

	for i := 0; i < len(words); i++ {
		assert.False(t, filter.Exists(&DefaultFilterObject{str: words[i]}),
			fmt.Sprintf("'%s' shouldn't be marked\n", words[i]))
	}
}

func benchmarkNaiveBloomFilter_Add(b *testing.B, count int) {

	filter := BloomFilter{}
	filter.Default()
	for i := 0; i < b.N; i++ {

		for j := 0; j < count; j++ {
			filter.Add(&DefaultFilterObject{str: fmt.Sprintf("%d", j)})
		}
	}

}

func BenchmarkNaiveBloomFilter_Add10(b *testing.B) {
	benchmarkNaiveBloomFilter_Add(b, 10)
}
func BenchmarkNaiveBloomFilter_Add100(b *testing.B) {
	benchmarkNaiveBloomFilter_Add(b, 100)

}
func BenchmarkNaiveBloomFilter_Add1000(b *testing.B) {
	benchmarkNaiveBloomFilter_Add(b, 1000)
}

func BenchmarkNaiveBloomFilter_Add10000(b *testing.B) {
	benchmarkNaiveBloomFilter_Add(b, 10000)
}

func benchmarkSimpleHasher_Hash(b *testing.B, count int) {

	sh := SimpleHasher{}
	sh.prefix = "REALLY VERY VERY LONG STRING WHICH IS COMPLETELY RANDOM AS RANDOM AS IT CAN GET !!"
	dfo := DefaultFilterObject{}
	dfo.str = "SOME RANDOME STRING WHICH WILL BE COMBINED TO CREATE A KEY"

	for i := 0; i < b.N; i++ {

		for j := 0; j < count; j++ {
			dfo.str = fmt.Sprintf("SOME RANDOME STRING WHICH WILL BE COMBINED TO CREATE A KEY %d", i)
			sh.Hash(&dfo)
		}
	}
}
func BenchmarkSimpleHasher_Hash10(b *testing.B) {
	benchmarkSimpleHasher_Hash(b, 10)
}
func BenchmarkSimpleHasher_Hash100(b *testing.B) {
	benchmarkSimpleHasher_Hash(b, 100)
}

func BenchmarkSimpleHasher_Hash1000(b *testing.B) {
	benchmarkSimpleHasher_Hash(b, 1000)
}

func BenchmarkSimpleHasher_Hash10000(b *testing.B) {
	benchmarkSimpleHasher_Hash(b, 10000)
}

func BenchmarkSimpleHasher_Hash100000(b *testing.B) {
	benchmarkSimpleHasher_Hash(b, 100000)
}

func TestCBF(t *testing.T) {

	filter := CountingBloomFilter{}
	filter.Default()
	for i := 0; i < len(words); i++ {

		filter.Add(&DefaultFilterObject{str: words[i]})
	}

	for i := 0; i < len(words); i++ {
		assert.True(t, filter.Exists(&DefaultFilterObject{str: words[i]}),
			fmt.Sprintf("'%s' should be marked\n", words[i]))
	}

	filter.Clear()

	for i := 0; i < len(words); i++ {
		assert.False(t, filter.Exists(&DefaultFilterObject{str: words[i]}),
			fmt.Sprintf("'%s' shouldn't be marked\n", words[i]))
	}
}

func TestCBFDelete(t *testing.T) {

	filter := CountingBloomFilter{}
	filter.Default()
	for i := 0; i < len(words); i++ {
		filter.Add(&DefaultFilterObject{str: words[i]})
	}

	for i := 0; i < len(words); i++ {
		assert.True(t, filter.Exists(&DefaultFilterObject{str: words[i]}),
			fmt.Sprintf("'%s' should be marked\n", words[i]))
	}

	for i := 0; i < len(words); i++ {
		assert.True(t, filter.Delete(&DefaultFilterObject{str: words[i]}),
			fmt.Sprintf("'%s' should be Deleted\n", words[i]))
	}

	for i := 0; i < len(words); i++ {
		assert.False(t, filter.Exists(&DefaultFilterObject{str: words[i]}),
			fmt.Sprintf("'%s' shouldn't be marked\n", words[i]))
	}

	for i := 0; i < len(words); i++ {
		assert.False(t, filter.Delete(&DefaultFilterObject{str: words[i]}),
			fmt.Sprintf("'%s' should be Deleted\n", words[i]))
	}

}
