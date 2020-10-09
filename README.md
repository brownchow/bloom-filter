# Go 语言实现 Bloom Filter

原理：https://llimllib.github.io/bloomfilter-tutorial/zh_CN/

BloomFilter 是一个数据结构，可以判断某个元素是否在集合内，具有**运行速度快**，**占用内存小**的特点。

高效查询的代价是，BloomFilter 是一个**基于概率的数据结构**：它只能告诉你**一个元素绝对不在集合内或者可能在集合内**。

BloomFilter 的基础数据结构是一个**比特向量**

Bloom Filter 里的哈希函数是需要**彼此独立**且**均匀分布**，同时，也要尽可能快

目前代码里使用的并不是 bitarray实现的，第一版是用 int array 实现的，而且使用不同的哈希函数，也只是用函数列表的index做区分

在所有编程语言里都不能直接生明比特，只能声明字节，其实问题不大，可以用字节数字，也可以用整形数组，最好是用bit数组。



~~每个样本占一个bit，100亿样本，10^10 bit = 1.16GB，~~bloomfilter 可以高效判断一个数据是否在100亿的集合里，而且运行速度很快，内存占用也小。

如果一个个和数组中的元素判断，时间复杂度是O(N)，这在数据量很大的情况下，效率很低。而用 BloomFilter虽然有一定的误判率，但效率高很多了。



参考：

https://www.bilibili.com/video/BV16g4y1B7Mp

https://www.bilibili.com/video/BV1v54y1m7hC



TODO:

1、goroutine 写入bit数组，好像更慢了（可能是太多的协程创建和销毁更耗费时间了）

https://www.bilibili.com/video/BV1bJ411t7A8

2、把之前的数据存到Redis，加快启动速度？就是一个序列化的过程





