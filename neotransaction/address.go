package neotransaction

import (
	"errors"

	"github.com/x-contract/neo-go-sdk/neoutils"
)

// AddressVersion 地址版本号
var AddressVersion = byte(23)

// 地址类型
const (
	AddresTypeUnKown int = 0
)

// Address represents an account address in neo
// An account address is a base58 string of raw address
// raw address is the script-hash with some checksum bits
// the script-hash is HASH160 of the script
type Address struct {
	Script    []byte
	ScripHash neoutils.HASH160
	Addr      string
	RawAddr   []byte
}

// Version 获取地址的版本号，版本号是RawAddr的第一个byte，目前默认都是23
// 配置见Neo节点客户端的protocol.json配置文件
func (addr *Address) Version() (byte, bool) {
	if addr.RawAddr == nil {
		return 0, false
	}
	return addr.RawAddr[0], true
}

// GetAddrString 获取地址字符串
func (addr *Address) GetAddrString() string {
	return addr.Addr
}

// HaveScript 判断是否拥有地址的鉴权脚本，如果没有鉴权脚本的话不能使用此地址中的资产
func (addr *Address) HaveScript() bool {
	return addr.Script != nil
}

// ParseAddress 根据一个地址字符串解析为地址结构体
// 注意：这样解析出来的地址结构体是不带鉴权脚本的
// 并且只有此地址的ScriptHash，没有此地址的公钥
func ParseAddress(addr string) (*Address, error) {
	rawAddr, ok := neoutils.DecodeBase58WithChecksum([]byte(addr))
	if !ok {
		return nil, errors.New("Address checksum failed")
	}
	ret := &Address{Addr: addr, RawAddr: rawAddr, ScripHash: rawAddr[1:21]}
	return ret, nil
}

// ParseAddressHash 根据ScriptHash来解析创建地址结构体，使用当前配置版本号来生成地址字符串
// 注意：这样解析出来的地址结构体是不带鉴权脚本的
func ParseAddressHash(scriptHash neoutils.HASH160) (*Address, error) {
	if !scriptHash.IsValid() {
		return nil, errors.New("Input script hash invalid")
	}
	ret := &Address{}
	data := make([]byte, 21)
	data[0] = AddressVersion
	copy(data[1:], scriptHash)
	var addr []byte
	addr, ret.RawAddr = neoutils.EncodeBase58WithChecksum(data)
	ret.Addr = string(addr[:])
	ret.ScripHash = ret.RawAddr[1:21]
	return ret, nil
}

// CreateAddressByScript 根据鉴权脚本来创建地址结构体，使用当前配置版本号来生成地址字符串
func CreateAddressByScript(script []byte) (*Address, error) {
	ret, err := ParseAddressHash(neoutils.Hash160(script))
	if err != nil {
		return nil, err
	}
	ret.Script = script
	return ret, nil
}
