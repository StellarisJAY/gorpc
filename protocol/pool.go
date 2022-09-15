package protocol

import "sync"

// messagePool Message资源池，避免重复地分配header等数组内存
var messagePool = &sync.Pool{New: func() interface{} {
	header := Header([10]byte{})
	header[0] = magicNumber
	return &Message{Header: &header}
}}
