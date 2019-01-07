package main

import "errors"

var (
	ErrDBNotFound          = errors.New("DB not found ")
	ErrRedisProtocolParser = errors.New("ErrRedisProtocolParser")
	ErrRedisLength = errors.New("ErrRedisLength")
	ErrKeyNotFound = errors.New("key not found")
	ErrTypeNotMatch = errors.New("type not match")
	ErrDontSupportThisCommand = errors.New("don't support this command  !!!!!!")
	ErrCommandArgsWrongNumber = errors.New("ERR wrong number of arguments for '%s' command")
	ErrCommand                = errors.New("ERR command")
	ErrValueType         = errors.New("ERR value is not an integer or out of range")
	ErrWrongKeyType         = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	ErrIndexOutOfRange        = errors.New("ERR index out of range")
)