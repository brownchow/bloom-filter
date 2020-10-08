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
	h      hash.Hash32 // 哈希函数
	bfList []int       // hash之后，在BF中所有索引的集合
}

// NewCountingBloomFilter 初始化
func NewCountingBloomFilter(totalNumber uint32, flaseDetectRate float64) *CBF {
	b := &CBF{h: fnv.New32()}
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
		listIndex := b.hashFuns(i, element)
		if b.bfList[listIndex] != 0 {
			b.bfList[listIndex]--
		}
	}

}

// bfSet： 用每一个hash函数计算hash值，然后转成整数，得到bit数组中的位置idx，返回这个idx集合
func (b *CBF) bfSet(data []byte) {
	for i := 0; i < b.k; i++ {
		// 用每个hash函数计算得到的int值
		listIndex := b.hashFuns(i, data)
		// +1 之后变成2也没关系吗？
		b.bfList[listIndex]++
	}
}

// hashFuns，给定整数，求出这个整数hash后的int值（模上 bit数组的长度m）
func (b *CBF) hashFuns(indexFn int, data []byte) int {
	b.h.Reset()
	b.h.Write(data)
	hashData := b.h.Sum32()
	hashInt := int(hashData)
	return (hashInt + indexFn) % b.m
}

// estimateMK: 根据传入容量和误判率，估算bit数组的长度m 和hash函数的个数 k
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
		listIndex := b.hashFuns(i, data)
		if b.bfList[listIndex] == 0 {
			return false
		}
	}
	return true
}
