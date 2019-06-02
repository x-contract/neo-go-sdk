package neoutils

import (
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

// HASH256 the type of hash value in neo. It's a unit256 big number that length 32 byte
// PublicKey, BlockHash, TransactionHash and other hashes are in HASH256 format
type HASH256 []byte

// IsValid 判断一个 HASH256 长度是否有效
func (hash HASH256) IsValid() bool {
	return len(hash) == 32
}

// Copy 复制一份 HASH256 数据
func (hash HASH256) Copy() HASH256 {
	ret := make(HASH256, len(hash))
	copy(ret, hash)
	return ret
}

// HASH160 is the type of script hash value in neo. Usually used as an account address.
// The address in string format is Base58 coded script hash with some checksum bytes.
type HASH160 []byte

// IsValid 判断一个 HASH160 长度是否有效
func (hash HASH160) IsValid() bool {
	return len(hash) == 20
}

// Copy 复制一份 HASH160 数据
func (hash HASH160) Copy() HASH160 {
	ret := make(HASH160, len(hash))
	copy(ret, hash)
	return ret
}

//var sha, ripemd hash.Hash

func init() {

}

// Sha256 get the SHA-256 hash value of b
func Sha256(b []byte) []byte {
	sha := sha256.New()
	sha.Write(b)
	return sha.Sum(nil)
}

// Hash256 get the twice SHA-256 hash value of ba
func Hash256(ba []byte) []byte {
	sha := sha256.New()
	sha.Write(ba)
	ba = sha.Sum(nil)
	sha.Reset()
	sha.Write(ba)
	return sha.Sum(nil)
}

// Hash160 first calculate SHA-256 hash result of b, then RIPEMD-160 hash of the result
func Hash160(b []byte) []byte {
	sha := sha256.New()
	sha.Write(b)
	b = sha.Sum(nil)
	ripemd := ripemd160.New()
	ripemd.Write(b)
	return ripemd.Sum(nil)
}
