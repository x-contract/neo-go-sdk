package neotransaction

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/x-contract/neo-go-sdk/neoutils"
)

// KeyPair 包含一对ECC加密算法的公私钥对
// 继承了ecdsa.PrivateKey
type KeyPair struct {
	ecdsa.PrivateKey
}

// HasPrivKey 判断一个KeyPair是否包含私钥，可以用于签名
func (key *KeyPair) HasPrivKey() bool {
	return key != nil && key.D.Cmp(big.NewInt(0)) != 0
}

// DecodeFromWif 从WIF字符串解码得到公私钥对
func DecodeFromWif(wif string) (*KeyPair, error) {
	buff, ok := neoutils.DecodeBase58WithChecksum([]byte(wif))
	if !ok {
		return nil, errors.New("DecodeFromWif checksum failed wif string " + wif)
	}

	if len(buff) != 34 || buff[0] != 0x80 || buff[33] != 0x01 {
		return nil, errors.New("DecodeFromWif invalid wif string " + wif)
	}

	priv := new(KeyPair)
	priv.PublicKey.Curve = elliptic.P256()
	priv.D = new(big.Int).SetBytes(buff[1:33])
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(priv.D.Bytes())
	return priv, nil
}

// DecodeFromPubkey 根据公钥二进制串解析出KeyPair,只包含公钥，仅能用于验签
func DecodeFromPubkey(pub []byte) (*KeyPair, error) {
	if len(pub) == 0 {
		return nil, errors.New("DecodeFromPubkey nil input")
	}

	priv := new(KeyPair)
	priv.PublicKey.Curve = elliptic.P256()
	priv.D = new(big.Int).SetInt64(0)

	switch pub[0] {
	case 0x02:
		fallthrough
	case 0x03:
		if len(pub) != 33 {
			return nil, errors.New("DecodeFromPubkey input length invalid")
		}
		priv.X = new(big.Int).SetBytes(pub[1:])
		y1, y2, err := CalcYOnEccCurve(priv.Curve.Params(), priv.X)
		if err != nil {
			return nil, err
		}
		if y1.Bit(0) == uint(pub[0]&1) {
			priv.Y = y1
		} else {
			priv.Y = y2
		}
	default:
		return nil, errors.New("DecodeFromPubkey type 0x" + hex.EncodeToString(pub[0:1]) + " not supported")
	}

	return priv, nil
}

// EncodeWif 将私钥编码成WIF字符串
func (key *KeyPair) EncodeWif() string {
	buff := make([]byte, 34)
	buff[0] = 0x80
	d := key.D.Bytes()
	copy(buff[1:33], d)
	buff[33] = 0x01
	buff, _ = neoutils.EncodeBase58WithChecksum(buff)
	return string(buff)
}

// EncodePubkeyCompressed 将公钥输出为压缩格式
func (key *KeyPair) EncodePubkeyCompressed() []byte {
	//return elliptic.Marshal(key.Curve, key.X, key.Y)
	data := make([]byte, 33)
	if key.Y.Bit(0) == 0 {
		data[0] = 0x02
	} else {
		data[0] = 0x03
	}
	copy(data[1:], key.X.Bytes())
	return data
}

// Sign 使用私钥对数据data进行签名
func (key *KeyPair) Sign(data []byte) ([]byte, error) {
	if key == nil || !key.HasPrivKey() {
		return nil, errors.New("The KeyPair does not contain private key")
	}
	r, s, _ := ecdsa.Sign(strings.NewReader("Nidces Ecc Signer"), &key.PrivateKey, data)
	ret := r.Bytes()
	ret = append(ret, s.Bytes()...)
	return ret, nil
}

// Verify 使用公私钥对中的公钥对数据验签
func (key *KeyPair) Verify(data []byte, sig []byte) bool {

	if len(sig) != 64 {
		return false
	}
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:])
	return ecdsa.Verify(&key.PublicKey, data, r, s)
}

// CalcYOnEccCurve 根据x值得到Ecc曲线上的点的Y值
func CalcYOnEccCurve(curve *elliptic.CurveParams, x *big.Int) (*big.Int, *big.Int, error) {

	x3 := new(big.Int).Mul(x, x)
	//x3.Mod(x3, curve.P)
	x3.Sub(x3, new(big.Int).SetInt64(3))

	x3.Mul(x3, x)
	x3.Add(x3, curve.B)

	if x3.ModSqrt(x3, curve.P) == nil {
		return nil, nil, errors.New("CalcYOnEccCurve no valid Y on curve")
	}
	return x3.Mod(x3, curve.P), new(big.Int).Sub(curve.P, x3), nil
}

// CreateBasicAddress 使用秘钥对的公钥创建基础鉴权账户
func (key *KeyPair) CreateBasicAddress() *Address {
	script := BuildBasicVerifyScript(key)
	addr, _ := CreateAddressByScript(script)
	return addr
}

// GenerateKeyPair 生成一个新的秘钥对
func GenerateKeyPair() *KeyPair {
	p, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return &KeyPair{PrivateKey: *p}
}
