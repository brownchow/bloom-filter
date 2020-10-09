package bf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math/rand"
	"testing"
	"time"

	"github.com/spaolacci/murmur3"
)

func TestBasicCBF(t *testing.T) {
	var number int
	var possibility float32
	number = 100
	possibility = 0.01
	cbf := NewCountingBloomFilter(100, 0.01)
	t.Logf("with number %d, possibility %f, we cauculate m = %d, k = %d", number, possibility, cbf.m, cbf.k)

	cbf.Add([]byte("foo"))
	t.Logf("#1 cbf is %d", cbf.bfList)
	cbf.Add([]byte("bar"))
	t.Logf("#2 cbf is %d", cbf.bfList)
	cbf.Add([]byte("tom"))
	t.Logf("#3 cbf is %d", cbf.bfList)
	cbf.Add([]byte("john"))
	t.Logf("#4 cbf is %d", cbf.bfList)

	if cbf.Test([]byte("tom")) == false {
		t.Errorf("CBF: test failed")
	}

	cbf.Remove([]byte("john"))
	if cbf.Test([]byte("john")) == true {
		t.Errorf("CBF: remove failed")
	}

}

func BenchmarkTen(b *testing.B) {
	m := 100

	h1 := func(input []byte) int {
		h := fnv.New32()
		h.Write(input)
		return int(h.Sum32()+1) % m
	}

	h2 := func(input []byte) int {
		h := fnv.New32()
		h.Write(input)
		return int(h.Sum32()+2) % m
	}

	var hashFuncs [](func(data []byte) int)
	hashFuncs = append(hashFuncs, h1)
	hashFuncs = append(hashFuncs, h2)
	var inputList = [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e"), []byte("f"), []byte("g"),
		[]byte("h"), []byte("i"), []byte("j"), []byte("k"), []byte("l"), []byte("m"), []byte("n")}

	b.ResetTimer()
	bf := CreateBloomFilter(inputList, hashFuncs, m)
	b.Logf("bf is: %d", bf)

	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(len(inputList))
	// 之前把inputList里每个元素都放进去了，所以这里测试的时候如果发现某个元素不在，说明有问题
	if MembershipTest(inputList[randIndex], bf, hashFuncs) == false {
		fmt.Printf("the %dth member test failed\n", randIndex)
	}
}

func TestBasicFun(t *testing.T) {
	m := 100
	h1 := func(input []byte) int {
		h := fnv.New32()
		h.Write(input)
		return int(h.Sum32()+1) % m
	}

	h2 := func(input []byte) int {
		h := fnv.New32()
		h.Write(input)
		return int(h.Sum32()+2) % m
	}

	var hashFuncs [](func(data []byte) int)
	hashFuncs = append(hashFuncs, h1)
	hashFuncs = append(hashFuncs, h2)
	var inputList = [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e"), []byte("f"), []byte("g"),
		[]byte("h"), []byte("i"), []byte("j"), []byte("k"), []byte("l"), []byte("m"), []byte("n")}

	bf := CreateBloomFilter(inputList, hashFuncs, m)
	fmt.Println(bf)
	if len(bf) != m {
		t.Errorf("size error %d, exprect %d", len(bf), m)
	}
	// 每个元素都放进去了，所以每个元素肯定都在bf里，如果不在，说明有错
	for index, element := range inputList {
		if MembershipTest(element, bf, hashFuncs) == false {
			t.Errorf("The %dth member test failed", index)
		}
	}

	if MembershipTest([]byte("o"), bf, hashFuncs) == true {
		t.Errorf("test error, %s does't exist in the set, but bf say it exist\n", "hi")
	}
}

func TestHashRes(t *testing.T) {
	var num int = 2
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, num)

	h := fnv.New64()
	h.Write(buf.Bytes())
	t.Logf("hash 2 using fnv, the result is: %d", h.Sum64()%15)

	t.Logf("hash 2 using murmur, the result is: %d", murmur3.Sum64(buf.Bytes())%14)
}
