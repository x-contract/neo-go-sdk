package neoutils

import (
	"math/big"
	"strings"
)

const base58 = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// EncodeBase58 ...
func EncodeBase58(ba []byte) []byte {
	if len(ba) == 0 {
		return nil
	}
	//Expected size increase from base58 conversion 25 approximately 137%,use 138% to be safe
	ri := len(ba) * 138 / 100
	ra := make([]byte, ri+1)

	x := new(big.Int).SetBytes(ba) // ba is big-endian
	x.Abs(x)
	y := big.NewInt(58)
	m := new(big.Int)

	for x.Sign() > 0 {
		x, m = x.DivMod(x, y, m)
		ra[ri] = base58[int32(m.Int64())]
		ri--
	}

	//Leading zeros encoded as base58 zeros
	for i := 0; i < len(ba); i++ {
		if ba[i] != 0 {
			break
		}
		ra[ri] = '1'
		ri--
	}
	return ra[ri+1:]
}

// DecodeBase58 ...
func DecodeBase58(ba []byte) []byte {
	if len(ba) == 0 {
		return nil
	}

	x := new(big.Int)
	y := big.NewInt(58)
	z := new(big.Int)
	for _, b := range ba {
		v := strings.IndexRune(base58, rune(b))
		z.SetInt64(int64(v))
		x.Mul(x, y)
		x.Add(x, z)
	}
	xa := x.Bytes()

	// Restore leading zeros
	i := 0
	for i < len(ba) && ba[i] == '1' {
		i++
	}
	ra := make([]byte, i+len(xa))
	copy(ra[i:], xa)
	return ra
}

// EncodeBase58WithChecksum 在数据ba后面添加4 bytes的HASH256的校验数据，然后整体进行Base58编码
// 同时返回添加4 byte校验数据之后的数据
func EncodeBase58WithChecksum(ba []byte) ([]byte, []byte) {
	//add 4-byte hash check to the end
	hash := Hash256(ba)
	ba = append(ba, hash[:4]...)
	return EncodeBase58(ba), ba
}

// DecodeBase58WithChecksum 将数据ba进行Base58解码，并且校验后4 bytes是否是前面数据的HASH256
// 返回的数据移除了最后4 bytes 的 Checksum
func DecodeBase58WithChecksum(ba []byte) ([]byte, bool) {
	ba = DecodeBase58(ba)
	if len(ba) < 4 || ba == nil {
		return nil, false
	}

	k := len(ba) - 4
	hash := Hash256(ba[:k])
	for i := 0; i < 4; i++ {
		if hash[i] != ba[k+i] {
			return nil, false
		}
	}
	return ba[:k], true
}
