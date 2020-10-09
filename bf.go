package bf

import (
	"hash"
	"hash/fnv"
	"math"
)

//CBF :  counting Bloom Filter
type CBF struct {
	m      int         // bit数组的长度
	k      int         // 不同hash函数的个数
	h      hash.Hash32 // 哈希函数，使用fnv系列的哈希函数
	bfList []int       // bit数组中，每个位置的值
}

// NewCountingBloomFilter 初始化
func NewCountingBloomFilter(totalNumber uint32, flaseDetectRate float64) *CBF {
	b := &CBF{h: fnv.New32()} // 使用fnv系列的哈希函数
	b.estimateMK(totalNumber, flaseDetectRate)
	b.bfList = make([]int, b.m)
	return b
}

// Add add an element into the cbf structure
func (b *CBF) Add(element []byte) {
	b.bfSet(element)
}

// Test 测试数据是否存在，因为 bloomFilter只能告诉你一个元素肯定不存在，或者可能存在
func (b *CBF) Test(element []byte) bool {
	return b.bfTest(element)
}

// Remove 从 cbf 数据结构中删除一个元素
func (b *CBF) Remove(element []byte) {
	if !b.bfTest(element) {
		return
	}

	for i := 0; i < b.k; i++ {
		listIndex := b.hashFuncs(i, element)
		if b.bfList[listIndex] != 0 {
			b.bfList[listIndex]--
		}
	}

}

// bfSet: 用每一个hash函数计算hash值，然后转成整数，得到bit数组中的位置idx，返回这个idx集合
func (b *CBF) bfSet(data []byte) {
	for i := 0; i < b.k; i++ {
		// 用每个hash函数计算得到的int值
		listIndex := b.hashFuncs(i, data)
		// +1 之后变成2也没关系吗？
		b.bfList[listIndex]++
	}
}

// hashFuns，给定整数，求出这个整数hash后的int值，然后模上bit数组的长度m
func (b *CBF) hashFuncs(indexFn int, data []byte) int {
	b.h.Reset()
	b.h.Write(data)
	hashData := b.h.Sum32()
	hashInt := int(hashData)
	// 这里仅仅用（hash结果 + 哈希函数的idx）表示不同的哈希函数
	return (hashInt + indexFn) % b.m
}

// estimateMK: 根据传入容量（可能插入数据集的容量大小）和误判率，估算bit数组的长度m 和hash函数的个数 k
func (b *CBF) estimateMK(number uint32, possibility float64) {
	// m = -1 *(n*lnP)/(ln2)^2 完全是数学运算
	nFloat := float64(number)
	ln2 := math.Log(2)
	b.m = int(-1 * (nFloat * math.Log(possibility)) / math.Pow(ln2, 2))

	// k = m/n *ln2
	b.k = int(math.Ceil(float64(b.m) / nFloat * ln2))
}

// bfTest 判断给定的数据是否有可能存在CBF数据结构中
func (b *CBF) bfTest(data []byte) bool {
	for i := 0; i < b.k; i++ {
		listIndex := b.hashFuncs(i, data)
		if b.bfList[listIndex] == 0 {
			return false
		}
	}
	return true
}

// CreateBloomFilter 通过自己指定Hash函数和bit数组的长度，以及初始数据，创建BloomFilter。setA: 自己提供的元素  m:bit数组的长度，hashFuncs 哈希函数列表
// 返回结果是bfList（初始化了的BloomFilter）表示每个bit是否置位，是一个数组
func CreateBloomFilter(setA [][]byte, hashFuncs [](func(data []byte) int), m int) []int {
	bfList := make([]int, m)
	for _, a := range setA {
		for _, hashFn := range hashFuncs {
			filterIndex := hashFn(a)
			bfList[filterIndex] = 1
		}
	}
	return bfList
}

// MemberShipTest 判断一个元素是否在一个集合中
// element: 需要判断的元素, filterSet: 将集合转成bloomFilter结构， hashFuncs 哈希函数列表
// 返回：元素是否存在集合中，将element 哈希后，每个bit都set了，说明有可能存在，若某个bit unset，说明肯定不存在集合中
func MembershipTest(element []byte, filterSet []int, hashFuncs [](func(data []byte) int)) bool {
	for _, hashFunc := range hashFuncs {
		index := hashFunc(element)
		if filterSet[index] == 0 {
			return false
		}
	}
	return true
}
