package vmcommon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"strconv"
)

const (
	ALIGN_LENGTH = 32
)

func Int16ToBytes(n int16) []byte  {
	tmp := int16(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

func Uint16ToBytes(n uint16) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint16(buf, n)
	return buf
}

func BytesToInt16(b []byte) int16  {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int16
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int16(tmp)
}

func BytesToUint16(b []byte) uint16  {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp uint16
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return uint16(tmp)
}

func Int32ToBytes(n int32) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

func Uint32ToBytes(val uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, val)
	return buf[:]
}

func BytesToUint32(b []byte) uint32 {
	b = append(make([]byte, 4-len(b)), b...)
	return binary.BigEndian.Uint32(b)
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
		b = b[len(b)-ALIGN_LENGTH:]
	}
	copy(tmp[ALIGN_LENGTH-len(b):], b)
	return tmp
}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

func Int64ToBytes(n int64) []byte {
	tmp := int64(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
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

func BytesToUint64(b []byte) uint64 {
	b = append(make([]byte, 8-len(b)), b...)
	return binary.BigEndian.Uint64(b)
}

func Float32ToBytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func BytesToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func Float64ToBytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func BytesToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}
func BoolToBytes(b bool) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, b)
	return buf.Bytes()
}

func BytesConverter(source []byte, t string) interface{} {
	switch t {
	case "int16":
		return BytesToInt16(source)
	case "uint16":
		return BytesToUint16(source)
	case "int32", "int":
		return BytesToInt32(source)
	case "uint32", "uint":
		return BytesToUint32(source)
	case "int64":
		return BytesToInt64(source)
	case "uint64":
		return BytesToUint64(source)
	case "float32":
		return BytesToFloat32(source)
	case "float64":
		return BytesToFloat64(source)
	case "string":
		return string(source[:])
	case "bool":
		return bytes.Equal(source, []byte{1})
	default:
		return source
	}
}

func StringConverter(source string, t string) ([]byte, error) {
	switch t {
	case "int16":
		dest, err := strconv.ParseInt(source, 10, 16)
		return Int16ToBytes(int16(dest)),err
	case "uint16":
		dest, err := strconv.ParseUint(source, 10, 16)
		return Uint16ToBytes(uint16(dest)),err
	case "int32", "int":
		dest, err := strconv.Atoi(source)
		return Int32ToBytes(int32(dest)), err
	case "uint32", "uint":
		dest, err := strconv.Atoi(source)
		return Uint32ToBytes(uint32(dest)), err
	case "int64":
		dest, err := strconv.ParseInt(source, 10, 64)
		return Int64ToBytes(dest), err
	case "uint64":
		dest, err := strconv.ParseUint(source, 10, 64)
		return Uint64ToBytes(dest), err
	case "float32":
		dest, err := strconv.ParseFloat(source, 32)
		return Float32ToBytes(float32(dest)), err
	case "float64":
		dest, err := strconv.ParseFloat(source, 64)
		return Float64ToBytes(dest), err
	case "bool":
		if "true" == source || "false" == source {
			return BoolToBytes("true" == source), nil
		} else {
			return []byte{}, errors.New("invalid boolean param")
		}
	default:
		return []byte(source), nil
	}

}
