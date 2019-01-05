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

const (
	// MaxReadSize read size
	MaxReadSize = 1024
)

const  (
	// MaxDBNum 最大DB数
	MaxDBNum = 16
)

const (
	// InitParser .
	InitParser = iota
	// ArrayParser .
	ArrayParser
	// BulkLengthParser .
	BulkLengthParser
	// BulkStringParser .
	BulkStringParser
	// IntParser .
	IntParser
	// StatusParser .
	StatusParser
	// ErrorParser .
	ErrorParser
)

// RespReply .
type RespReply byte

const (
	// ErrorReply .
	ErrorReply = RespReply('-')
	// StatusReply .
	StatusReply = RespReply('+')
	// IntReply .
	IntReply = RespReply(':')
	// BulkLengthReply .
	BulkLengthReply = RespReply('$')
	// ArrayReply .
	ArrayReply = RespReply('*')
)

const (
	KeyNotExists = -2
)

const (
	DEL = "DEL"
	SET = "SET"
	EX = "EX"
	PX = "PX"
	NX = "NX"
	XX = "XX"
	GET = "GET"
	TTL = "TTL"
	PTTL = "PTTL"
	EXPIRE = "EXPIRE"
	EXPIREAT = "EXPIREAT"
	PEXPIREAT = "PEXPIREAT"
	PEXPIRE = "PEXPIRE"
	TYPE = "TYPE"
	EXISTS = "EXISTS"
	PERSIST = "PERSIST"
)

const (
	// TypeRedisString string
	TypeRedisString = iota
	// TypeRedisList list
	TypeRedisList
	// TypeRedisSet set
	TypeRedisSet
	// TypeRedisZSet zset
	TypeRedisZSet
	// TypeRedisHash hash
	TypeRedisHash
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
