package neotransaction

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/x-contract/neo-go-sdk/neotransaction/OpCode"
	"github.com/x-contract/neo-go-sdk/neoutils"
)

// ScriptBuilder NEO 智能合约脚本构建器
type ScriptBuilder struct {
	buff bytes.Buffer
}

// Bytes 获取脚本构建器输出的二进制脚本数据
func (sb *ScriptBuilder) Bytes() []byte {
	return sb.buff.Bytes()
}

// Emit 在脚本构建器中加入一条不带参数的指令
func (sb *ScriptBuilder) Emit(op OpCode.OPCODE) {
	sb.buff.WriteByte(byte(op))
}

// EmitOpArgs 在脚本构建器中加入一条指令以及它的参数
func (sb *ScriptBuilder) EmitOpArgs(op OpCode.OPCODE, args []byte) {
	sb.buff.WriteByte(byte(op))
	sb.buff.Write(args)
}

// EmitAppCall 在脚本构建器中加入一条合约调用指令，参数为被调用的合约的脚本哈希
func (sb *ScriptBuilder) EmitAppCall(scriptHash neoutils.HASH160) {
	sb.Emit(OpCode.APPCALL)
	sb.buff.Write(neoutils.Reverse(scriptHash))
}

// EmitPushBool 在脚本构建器中加入一条压栈布尔值的指令
func (sb *ScriptBuilder) EmitPushBool(arg bool) {
	if arg {
		sb.Emit(OpCode.PUSHT)
	} else {
		sb.Emit(OpCode.PUSHF)
	}
}

// EmitPushBytes 在脚本构建器中加入一条压栈字节数组的指令
func (sb *ScriptBuilder) EmitPushBytes(arg []byte) {
	if len(arg) <= int(OpCode.PUSHBYTES75) {
		sb.buff.WriteByte(byte(len(arg)))
		sb.buff.Write(arg)
	} else if len(arg) <= 0xff {
		sb.Emit(OpCode.PUSHDATA1)
		sb.buff.WriteByte(byte(len(arg)))
		sb.buff.Write(arg)
	} else if len(arg) <= 0xffff {
		sb.Emit(OpCode.PUSHDATA2)
		neoutils.WriteUint16ToBuffer(&sb.buff, uint16(len(arg)))
		sb.buff.Write(arg)
	} else {
		sb.Emit(OpCode.PUSHDATA4)
		neoutils.WriteUint32ToBuffer(&sb.buff, uint32(len(arg)))
		sb.buff.Write(arg)
	}
}

// EmitPushNumber 在脚本构建器中加入一条压栈数字的指令
func (sb *ScriptBuilder) EmitPushNumber(arg int64) {
	if arg == -1 {
		sb.Emit(OpCode.PUSHM1)
		return
	}
	if arg == 0 {
		sb.Emit(OpCode.PUSH0)
		return
	}
	if arg > 0 && arg <= 16 {
		sb.Emit(OpCode.PUSH1 - 1 + OpCode.OPCODE(arg))
		return
	}
	bytes := neoutils.Reverse(big.NewInt(arg).Bytes())
	sb.EmitPushBytes(bytes)
}

// EmitPushString 在脚本构建器中加入一条压栈字符串的指令，压栈字符串实际上是压栈字节数组
func (sb *ScriptBuilder) EmitPushString(arg string) {
	sb.EmitPushBytes([]byte(arg))
}

// EmitPushArray 在脚本构建器中加入一条压栈数组的指令，数组的元素可以是 HASH256 HASH160 string []byte number bool
// 将数组压栈需要将数组元素按照从右至左压入栈中，然后压入数组长度，最后压栈 Pack 指令
func (sb *ScriptBuilder) EmitPushArray(arg []interface{}) error {
	for i := len(arg) - 1; i >= 0; i-- {
		arg := arg[i]
		switch v := arg.(type) {
		case neoutils.HASH256:
			sb.EmitPushBytes(neoutils.Reverse(v))
		case neoutils.HASH160:
			sb.EmitPushBytes(neoutils.Reverse(v))
		case []byte:
			sb.EmitPushBytes(v)
		case uint64:
			sb.EmitPushNumber(int64(v))
		case int64:
			sb.EmitPushNumber(int64(v))
		case uint32:
			sb.EmitPushNumber(int64(v))
		case int32:
			sb.EmitPushNumber(int64(v))
		case uint16:
			sb.EmitPushNumber(int64(v))
		case int16:
			sb.EmitPushNumber(int64(v))
		case byte:
			sb.EmitPushNumber(int64(v))
		case string:
			sb.EmitPushString(v)
		case bool:
			sb.EmitPushBool(v)
		default:
			return fmt.Errorf(`EmitPushArray Error: script params[%T] not supported`, v)
		}
	}
	sb.EmitPushNumber(int64(len(arg)))
	sb.Emit(OpCode.PACK)
	return nil
}

// BuildBasicWitnessScript 创建一个基础的鉴证人脚本,包含基础压栈脚本和基础鉴权脚本
// =============基础压栈脚本============
// Push Signature
// ================================
// =============基础鉴权脚本============
// Push PublicKey
// CheckSig
// ================================
func BuildBasicWitnessScript(keyPair *KeyPair, rawTx []byte) (*Script, error) {

	script := &Script{}

	// 对原始交易进行签名
	signature, err := keyPair.Sign(neoutils.Sha256(rawTx))
	if err != nil {
		return script, err
	}

	// 创建压栈脚本
	script.InvocationScript = make([]byte, len(signature)+1)
	script.InvocationScript[0] = byte(len(signature))
	copy(script.InvocationScript[1:], signature)
	script.InvScriptLength.Value = uint64(len(script.InvocationScript))

	// 压缩公钥数据串
	pubKey := keyPair.EncodePubkeyCompressed()

	// 创建鉴权脚本
	script.VerificationScript = make([]byte, len(pubKey)+2)
	script.VerificationScript[0] = byte(len(pubKey))
	copy(script.VerificationScript[1:], pubKey)
	script.VerificationScript[1+len(pubKey)] = 0xac
	script.VrifScriptLength.Value = uint64(len(script.VerificationScript))

	return script, nil
}

// BuildBasicVerifyScript 创建基本账户鉴权脚本
func BuildBasicVerifyScript(keyPair *KeyPair) []byte {

	// 压缩公钥数据串
	pubKey := keyPair.EncodePubkeyCompressed()

	// 创建鉴权脚本
	VerificationScript := make([]byte, len(pubKey)+2)
	VerificationScript[0] = byte(len(pubKey))
	copy(VerificationScript[1:], pubKey)
	VerificationScript[1+len(pubKey)] = 0xac

	return VerificationScript
}

// BuildCallMethodScript 生成一个合约调用脚本，合约的第一个参数必须是一个字符串
// withNonce 表示是否要在调用指令前插入一个随机数，这样可以让调用脚本以及所在交易hash值产生变化
// 如果交易中不包含utxo资产的输入输出及其它可变数据的情况下，多次调用同一个合约的同一个接口
// 构建出来的交易结构是不变的，因此hash值会冲突，加入随机数可以避免这个冲突
func BuildCallMethodScript(contractHash neoutils.HASH160, method string, args []interface{}, withNonce bool) ([]byte, error) {
	sb := ScriptBuilder{}
	if withNonce {
		rand.Seed(time.Now().UnixNano())
		sb.EmitPushNumber(int64(rand.Uint32()))
		sb.Emit(OpCode.DROP)
	}
	for i := len(args) - 1; i >= 0; i-- {
		arg := args[i]
		switch v := arg.(type) {
		case neoutils.HASH256:
			sb.EmitPushBytes(v)
		case neoutils.HASH160:
			sb.EmitPushBytes(v)
		case []byte:
			sb.EmitPushBytes(v)
		case uint64:
			sb.EmitPushNumber(int64(v))
		case int64:
			sb.EmitPushNumber(int64(v))
		case uint32:
			sb.EmitPushNumber(int64(v))
		case int32:
			sb.EmitPushNumber(int64(v))
		case uint16:
			sb.EmitPushNumber(int64(v))
		case int16:
			sb.EmitPushNumber(int64(v))
		case byte:
			sb.EmitPushNumber(int64(v))
		case string:
			sb.EmitPushString(v)
		case bool:
			sb.EmitPushBool(v)
		default:
			return nil, fmt.Errorf(`BuildCallMethodScript Error: script params[%T] not supported`, v)
		}
	}
	sb.EmitPushNumber(int64(len(args)))
	sb.Emit(OpCode.PACK)
	sb.EmitPushString(method)
	sb.EmitAppCall(contractHash)
	return sb.Bytes(), nil
}
