package neotransaction

import (
	"bytes"
	"encoding/hex"
	"errors"

	"github.com/x-contract/neo-go-sdk/neoutils"
)

// The AssetId of some neo token
const (
	AssetNeoID = "c56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b"
	AssetGasID = "602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7"
)

// GetAssetID 根据货币符号获取资产ID
func GetAssetID(assetName string) string {
	switch assetName {
	case `NEO`:
		return AssetNeoID
	case `GAS`:
		return AssetGasID
	default:
		return ``
	}
}

// GetAssetSymbol 根据资产ID获取货币符号
func GetAssetSymbol(assetID string) string {
	switch assetID {
	case AssetNeoID:
		return `NEO`
	case AssetGasID:
		return `GAS`
	default:
		return ``
	}
}

// The type of NEO Transaction
// Length = 1 byte
const (
	MinerTranscation     byte = 0x00 // 用于分配字节费的交易
	IssueTransaction     byte = 0x01 // 用于分发资产的交易
	ClaimTransaction     byte = 0x02 // 用于分配 NeoGas 的交易
	ContractTransaction  byte = 0x80 // 合约交易，这是最常用的一种交易
	InvocationTransacton byte = 0xd1 // 调用智能合约的特殊交易
)

// The type of NEO transaction attribute usage
// Length = 1 byte
const (
	UsageContractHash   byte = 0x00 //  外部合同的散列值
	UsageECDH02         byte = 0x02
	UsageECDH03         byte = 0x03 //  用于 ECDH 密钥交换的公钥
	UsageScript         byte = 0x20 //  额外鉴证人的ScriptHash
	UsageVote           byte = 0x30 //	用于投票选出记账人
	UsageCertURL        byte = 0x80 //	证书地址
	UsageDescriptionURL byte = 0x81 //	外部介绍信息地址
	UsageDescription    byte = 0x90 //	简短的介绍信息
	UsageHash1          byte = 0xa1 //	用于存放自定义的散列值
	//-0xaf  -Hash15
	UsageRemark byte = 0xf0 // 备注
	//-0xff	Remark-Remark15
)

// TxOutputValueBase Neo交易金额的基数（定点数的小数位数)
const TxOutputValueBase = int64(100000000)

// Attribute attribute of a NeoTransaction
// 对于 ContractHash，ECDH 系列，Vote，Hash 系列，数据长度固定为 32 字节，length 字段省略
// 对于 Script 固定为20个字节，并且是 Big-Endian 字节序
// 对于 CertUrl，DescriptionUrl，Description，Remark 系列，必须明确给出数据长度，且长度不能超过 255
type Attribute struct {
	Usage byte   // 用途
	Data  []byte // 特定于用途的外部数据
}

// TxInput input struct of a NeoTransaction
type TxInput struct {
	PrevHash  neoutils.HASH256 // 引用交易的散列值
	PrevIndex uint16           // 引用交易输出的索引
}

// TxOutput output struct of a NeoTransaction
type TxOutput struct {
	AssetID    neoutils.HASH256 // 资产编号
	Value      int64            // 金额 金额固定为 10e-8 单位
	ScriptHash neoutils.HASH160 // 收款地址
}

// Script is the script part of a NeoTransaction
type Script struct {
	InvScriptLength  neoutils.VarInt
	InvocationScript []byte // StackScript 栈脚本代码

	VrifScriptLength   neoutils.VarInt
	VerificationScript []byte // RedeemScript 合约脚本代码
}

type extraData interface {
	Bytes() []byte
}

// InvocationExtraData 调用交易的额外数据
type InvocationExtraData struct {
	ScriptLength neoutils.VarInt
	Script       []byte
	GasConsumed  int64
}

// Bytes ...
func (extra *InvocationExtraData) Bytes() []byte {
	buff := new(bytes.Buffer)
	extra.ScriptLength.Value = uint64(len(extra.Script))
	buff.Write(extra.ScriptLength.Bytes())
	buff.Write(extra.Script)
	//neoutils.WriteUint64ToBuffer(buff, uint64(extra.GasConsumed))
	return buff.Bytes()
}

// NeoTransaction struct
type NeoTransaction struct {
	Type    byte
	Version byte

	ExtraData extraData

	AttributeCount neoutils.VarInt
	Attributes     []Attribute

	InputsCount neoutils.VarInt
	Inputs      []TxInput

	OutputsCount neoutils.VarInt
	Outputs      []TxOutput

	ScriptsCount neoutils.VarInt
	Scripts      []Script

	dirty bool
	//txid        string
	unsingedraw []byte
	witness     []byte
}

// AppendAttribute 向交易添加一条Attribute
func (tx *NeoTransaction) AppendAttribute(usage byte, data []byte) {
	tx.Attributes = append(tx.Attributes, Attribute{
		Usage: usage,
		Data:  data,
	})
	tx.dirty = true
}

// AppendInput 向交易添加一笔UTXO作为输入
func (tx *NeoTransaction) AppendInput(utxo *UTXO) {
	tx.Inputs = append(tx.Inputs, TxInput{utxo.TxHash, utxo.Index})
	tx.dirty = true
}

// AppendInputByTxHash 向交易添加一笔输入
func (tx *NeoTransaction) AppendInputByTxHash(transactionID string, index uint16) error {
	var hash neoutils.HASH256
	hash, _ = hex.DecodeString(transactionID)
	if !hash.IsValid() {
		return errors.New("NeoTransaction.AppendInput invalid transactionID")
	}
	hash = neoutils.Reverse(hash)
	tx.Inputs = append(tx.Inputs, TxInput{PrevHash: hash, PrevIndex: index})
	tx.dirty = true
	return nil
}

// AppendOutputByAddrString 向交易添加一笔输出
func (tx *NeoTransaction) AppendOutputByAddrString(address string, assetID string, count int64) error {

	addr, err := ParseAddress(address)
	if err != nil {
		return err
	}

	var assetHash neoutils.HASH256
	assetHash, err = hex.DecodeString(assetID)
	if err != nil {
		return err
	}
	if !assetHash.IsValid() {
		return errors.New("NeoTransaction.AppendOutputByAddString invalid assetID")
	}
	assetHash = neoutils.Reverse(assetHash)

	tx.Outputs = append(tx.Outputs, TxOutput{AssetID: assetHash, ScriptHash: addr.ScripHash, Value: count})
	tx.dirty = true
	return nil
}

// AppendOutput 向交易添加一笔输出
func (tx *NeoTransaction) AppendOutput(addr *Address, assetHash neoutils.HASH256, count int64) error {

	if !assetHash.IsValid() {
		return errors.New("NeoTransaction.AppendOutput invalid assetHash")
	}
	tx.Outputs = append(tx.Outputs, TxOutput{AssetID: assetHash, ScriptHash: addr.ScripHash, Value: count})
	tx.dirty = true
	return nil
}

// UnsignedRawTransaction 返回不包含脚本的原始交易
func (tx *NeoTransaction) UnsignedRawTransaction() []byte {
	if !tx.dirty {
		return tx.unsingedraw
	}
	tx.dirty = false
	buff := new(bytes.Buffer)
	buff.WriteByte(tx.Type)
	buff.WriteByte(tx.Version)
	if tx.ExtraData != nil {
		buff.Write(tx.ExtraData.Bytes())
	}

	tx.AttributeCount.Value = uint64(len(tx.Attributes))
	buff.Write(tx.AttributeCount.Bytes())
	for _, attr := range tx.Attributes {
		buff.WriteByte(attr.Usage)
		if attr.Usage == UsageCertURL || attr.Usage == UsageDescription || attr.Usage == UsageDescriptionURL || attr.Usage >= UsageRemark {
			buff.WriteByte(byte(len(attr.Data)))
		}
		buff.Write(attr.Data)
	}

	tx.InputsCount.Value = uint64(len(tx.Inputs))
	buff.Write(tx.InputsCount.Bytes())
	for i := 0; i < len(tx.Inputs); i++ {
		buff.Write(tx.Inputs[i].PrevHash)
		neoutils.WriteUint16ToBuffer(buff, tx.Inputs[i].PrevIndex)
	}

	tx.OutputsCount.Value = uint64(len(tx.Outputs))
	buff.Write(tx.OutputsCount.Bytes())
	for i := 0; i < len(tx.Outputs); i++ {
		buff.Write(tx.Outputs[i].AssetID)
		neoutils.WriteUint64ToBuffer(buff, uint64(tx.Outputs[i].Value))
		buff.Write(tx.Outputs[i].ScriptHash)
	}
	tx.unsingedraw = buff.Bytes()
	return tx.unsingedraw
}

// TXID 返回一个交易的交易ID
func (tx *NeoTransaction) TXID() string {
	// if len(tx.txid) != 0 {
	// 	return tx.txid
	// }
	txid := neoutils.Hash256(tx.UnsignedRawTransaction())
	//tx.txid = hex.EncodeToString(neoutils.Reverse(txid))
	return hex.EncodeToString(neoutils.Reverse(txid))
}

// AppendWitness 向交易添加一个鉴证人。一个独立的鉴证人有一个鉴证人脚本，包括一个压栈脚本和一个鉴权脚本
func (tx *NeoTransaction) AppendWitness(witness *Script) {
	tx.Scripts = append(tx.Scripts, *witness)
}

// AppendBasicSignWitness 向交易添加一个基本签名账户的鉴证人脚本，压栈脚本为一条将签名压栈的指令，
// 鉴权脚本就是基本账户鉴权脚本，将公钥压栈然后调用验签
func (tx *NeoTransaction) AppendBasicSignWitness(key *KeyPair) {
	script, _ := BuildBasicWitnessScript(key, tx.UnsignedRawTransaction())
	tx.AppendWitness(script)
}

// Witnesses 返回鉴证人脚本的二进制数据块
func (tx *NeoTransaction) Witnesses() []byte {
	if tx.witness != nil {
		return tx.witness
	}
	buff := new(bytes.Buffer)
	tx.ScriptsCount.Value = uint64(len(tx.Scripts))
	buff.Write(tx.ScriptsCount.Bytes())
	for i := 0; i < len(tx.Scripts); i++ {
		buff.Write(tx.Scripts[i].InvScriptLength.Bytes())
		buff.Write(tx.Scripts[i].InvocationScript)

		buff.Write(tx.Scripts[i].VrifScriptLength.Bytes())
		buff.Write(tx.Scripts[i].VerificationScript)
	}
	tx.witness = buff.Bytes()
	return tx.witness
}

// RawTransaction 返回签名后的完整二进制交易
func (tx *NeoTransaction) RawTransaction() []byte {
	return append(tx.UnsignedRawTransaction(), tx.Witnesses()...)
}

// RawTransactionString 返回完整交易的二进制字符串
func (tx *NeoTransaction) RawTransactionString() string {
	return hex.EncodeToString(tx.RawTransaction())
}

// CreateContractTransaction 创建一个合约交易（utxo转账交易）
func CreateContractTransaction() *NeoTransaction {
	tx := &NeoTransaction{
		Type:  ContractTransaction,
		dirty: true,
	}
	return tx
}

// CreateInvocationTransaction 创建一个调用交易（调用只能合约）
func CreateInvocationTransaction() *NeoTransaction {
	tx := &NeoTransaction{
		Type: InvocationTransacton,
		//Version:   1,
		ExtraData: &InvocationExtraData{},
		dirty:     true,
	}
	return tx
}
