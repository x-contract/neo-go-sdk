package neoutils

import "bytes"

// WriteUint16ToBuffer 向bytes.Buffer中以小端序写入一个uint16
func WriteUint16ToBuffer(b *bytes.Buffer, v uint16) {
	b.WriteByte(byte(v))
	b.WriteByte(byte(v >> 8))
}

// WriteUint32ToBuffer 向bytes.Buffer中以小端序写入一个uint32
func WriteUint32ToBuffer(b *bytes.Buffer, v uint32) {
	b.WriteByte(byte(v))
	b.WriteByte(byte(v >> 8))
	b.WriteByte(byte(v >> 16))
	b.WriteByte(byte(v >> 24))
}

// WriteUint64ToBuffer 向bytes.Buffer中以小端序写入一个uint64
func WriteUint64ToBuffer(b *bytes.Buffer, v uint64) {
	b.WriteByte(byte(v))
	b.WriteByte(byte(v >> 8))
	b.WriteByte(byte(v >> 16))
	b.WriteByte(byte(v >> 24))
	b.WriteByte(byte(v >> 32))
	b.WriteByte(byte(v >> 40))
	b.WriteByte(byte(v >> 48))
	b.WriteByte(byte(v >> 56))
}

// Reverse 将s反序
func Reverse(s []byte) []byte {
	k := make([]byte, len(s))
	copy(k, s)
	for i, j := 0, len(k)-1; i < j; i, j = i+1, j-1 {
		k[i], k[j] = k[j], k[i]
	}
	return k
}
