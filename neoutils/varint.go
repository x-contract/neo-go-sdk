package neoutils

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// VarInt 变长整数，可以根据表达的值进行编码以节省空间。
// 最长表示64位无符号整数，通常用于表示数组长度
type VarInt struct {
	Value uint64
}

// Length returns the serialized bytes length of the var int
func (varint VarInt) Length() int {
	if varint.Value < 0xfd {
		return 1
	}
	if varint.Value <= 0xffff {
		return 3
	}
	if varint.Value <= 0xffffffff {
		return 5
	}
	return 9
}

// Bytes returns the serialized bytes of the var int
func (varint VarInt) Bytes() []byte {

	if varint.Value < 0xfd {
		ret := make([]byte, 1)
		ret[0] = byte(varint.Value)
		return ret
	}
	if varint.Value <= 0xffff {
		ret := make([]byte, 3)
		ret[0] = 0xfd
		binary.LittleEndian.PutUint16(ret[1:], uint16(varint.Value))
		return ret
	}
	if varint.Value <= 0xffffffff {
		ret := make([]byte, 5)
		ret[0] = 0xfe
		binary.LittleEndian.PutUint32(ret[1:], uint32(varint.Value))
		return ret
	}
	ret := make([]byte, 9)
	ret[0] = 0xff
	binary.LittleEndian.PutUint64(ret[1:], uint64(varint.Value))
	return ret
}

// ParseVarInt parse the serialized bytes of the var int and return VarInt
func ParseVarInt(bytes []byte) (VarInt, error) {
	ret := VarInt{}
	if len(bytes) < 1 {
		return ret, errors.New("ParseVarInt: input bytes length 0")
	}
	if bytes[0] < 0xfd {
		ret.Value = uint64(bytes[0])
		return ret, nil
	}
	if bytes[0] == 0xfd {
		if len(bytes) < 3 {
			return ret, errors.New(fmt.Sprint("ParseVarInt: input bytes starts with 0xfd but length", len(bytes)))
		}
		ret.Value = uint64(binary.LittleEndian.Uint16(bytes[1:]))
		return ret, nil
	}
	if bytes[0] == 0xfe {
		if len(bytes) < 5 {
			return ret, errors.New(fmt.Sprint("ParseVarInt: input bytes starts with 0xfe but length", len(bytes)))
		}
		ret.Value = uint64(binary.LittleEndian.Uint32(bytes[1:]))
		return ret, nil
	}

	if len(bytes) < 9 {
		return ret, errors.New(fmt.Sprint("ParseVarInt: input bytes starts with 0xfe but length", len(bytes)))
	}
	ret.Value = binary.LittleEndian.Uint64(bytes[1:])
	return ret, nil
}
