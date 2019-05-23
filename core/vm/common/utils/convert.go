package utils

import (
	"bytes"
	"encoding/binary"
)

const (
	ALIGN_LENGTH = 32
)

func Int64ToBytes(i int64) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, &i)
	return buf.Bytes()
}

func BytesToInt64(b []byte) int64 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int64
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int64(tmp)
}

func Uint64ToBytes(n uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, n)
	return buf
}

func Int32ToBytes(n int32) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

func BytesToInt32(b []byte) int32 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int32(tmp)
}

func Align32Bytes(b []byte) []byte {
	tmp := make([]byte, ALIGN_LENGTH)
	if len(b) > ALIGN_LENGTH {
		b = b[len(b) - ALIGN_LENGTH:]
	}
	copy(tmp[ALIGN_LENGTH-len(b):], b)
	return tmp
}