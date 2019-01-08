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

const (
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

// Key
const (
	DEL       = "DEL"
	DUMP      = "DUMP"
	EXISTS    = "EXISTS"
	EXPIRE    = "EXPIRE"
	EXPIREAT  = "EXPIREAT"
	KEYS      = "KEYS"
	MIGRATE   = "MIGRATE"
	MOVE      = "MOVE"
	OBJECT    = "OBJECT"
	PERSIST   = "PERSIST"
	PEXPIRE   = "PEXPIRE"
	PEXPIREAT = "PEXPIREAT"
	PTTL      = "PTTL"
	RANDOMKEY = "RANDOMKEY"
	RENAME    = "RENAME"
	RENAMENX  = "RENAMENX"
	RESTORE   = "RESTORE"
	SORT      = "SORT"
	TTL       = "TTL"
	TYPE      = "TYPE"
	SCAN      = "SCAN"
)

// String
const (
	APPEND      = "APPEND"
	BITCOUNT    = "BITCOUNT"
	BITOP       = "BITOP"
	DECR        = "DECR"
	DECRBY      = "DECRBY"
	GET         = "GET"
	GETBIT      = "GETBIT"
	GETRANGE    = "GETRANGE"
	GETSET      = "GETSET"
	INCR        = "INCR"
	INCRBY      = "INCRBY"
	INCRBYFLOAT = "INCRBYFLOAT"
	MGET        = "MGET"
	MSET        = "MSET"
	MSETNX      = "MSETNX"
	PSETEX      = "PSETEX"
	SET         = "SET"
	SETBIT      = "SETBIT"
	SETEX       = "SETEX"
	SETNX       = "SETNX"
	SETRANGE    = "SETRANGE"
	STRLEN      = "STRLEN"
)

// Hash
const (
	HDEL         = "HDEL"
	HEXISTS      = "HEXISTS"
	HGET         = "HGET"
	HGETALL      = "HGETALL"
	HINCRBY      = "HINCRBY"
	HINCRBYFLOAT = "HINCRBYFLOAT"
	HKEYS        = "HKEYS"
	HLEN         = "HLEN"
	HMGET        = "HMGET"
	HMSET        = "HMSET"
	HSET         = "HSET"
	HSETNX       = "HSETNX"
	HVALS        = "HVALS"
	HSCAN        = "HSCAN"
)

// List
const (
	BLPOP      = "BLPOP"
	BRPOP      = "BRPOP"
	BRPOPLPUSH = "BRPOPLPUSH"
	LINDEX     = "LINDEX"
	LINSERT    = "LINSERT"
	LLEN       = "LLEN"
	LPOP       = "LPOP"
	LPUSH      = "LPUSH"
	LPUSHX     = "LPUSHX"
	LRANGE     = "LRANGE"
	LREM       = "LREM"
	LSET       = "LSET"
	LTRIM      = "LTRIM"
	RPOP       = "RPOP"
	RPOPLPUSH  = "RPOPLPUSH"
	RPUSH      = "RPUSH"
	RPUSHX     = "RPUSHX"
)

// Set
const (
	SADD        = "SADD"
	SCARD       = "SCARD"
	SDIFF       = "SDIFF"
	SDIFFSTORE  = "SDIFFSTORE"
	SINTER      = "SINTER"
	SINTERSTORE = "SINTERSTORE"
	SISMEMBER   = "SISMEMBER"
	SMEMBERS    = "SMEMBERS"
	SMOVE       = "SMOVE"
	SPOP        = "SPOP"
	SRANDMEMBER = "SRANDMEMBER"
	SREM        = "SREM"
	SUNION      = "SUNION"
	SUNIONSTORE = "SUNIONSTORE"
	SSCAN       = "SSCAN"
)

// SortedSet
const (
	ZADD             = "ZADD"
	ZCARD            = "ZCARD"
	ZCOUNT           = "ZCOUNT"
	ZINCRBY          = "ZINCRBY"
	ZRANGE           = "ZRANGE"
	ZRANGEBYSCORE    = "ZRANGEBYSCORE"
	ZRANK            = "ZRANK"
	ZREM             = "ZREM"
	ZREMRANGEBYRANK  = "ZREMRANGEBYRANK"
	ZREMRANGEBYSCORE = "ZREMRANGEBYSCORE"
	ZREVRANGE        = "ZREVRANGE"
	ZREVRANGEBYSCORE = "ZREVRANGEBYSCORE"
	ZREVRANK         = "ZREVRANK"
	ZSCORE           = "ZSCORE"
	ZUNIONSTORE      = "ZUNIONSTORE"
	ZINTERSTORE      = "ZINTERSTORE"
	ZSCAN            = "ZSCAN"
)

// Connection
const (
	AUTH   = "AUTH"
	ECHO   = "ECHO"
	PING   = "PING"
	QUIT   = "QUIT"
	SELECT = "SELECT"
)

// Server
const (
	FLUSHALL = "FLUSHALL"
	FLUSHDB = "FLUSHDB"
	DBSIZE = "DBSIZE"
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
