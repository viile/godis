package main

const (
	// STUnknown Unknown
	STUnknown = iota
	// STInited Inited
	STInited
	// STRunning Running
	STRunning
	// STStop Stop
	STStop
)

const  (
	// MaxDBNum 最大DB数
	MaxDBNum = 16
)

const (
	// RedisString string
	RedisString = iota
	// RedisList list
	RedisList
	// RedisSet set
	RedisSet
	// RedisZSet zset
	RedisZSet
	// RedisHash hash
	RedisHash
)

const (
	// RedisEncodingRaw raw
	RedisEncodingRaw = iota
	// RedisEncodingInt int
	RedisEncodingInt
	// RedisEncodingHt ht
	RedisEncodingHt
	// RedisEncodingZipMap zipmap
	RedisEncodingZipMap
	// RedisEncodingLinkedList linked-list
	RedisEncodingLinkedList
	// RedisEncodingZipList zip-list
	RedisEncodingZipList
	// RedisEncodingIntSet int-set
	RedisEncodingIntSet
	// RedisEncodingSkipList skip-list
	RedisEncodingSkipList
)
