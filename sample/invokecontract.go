package main

import (
	"encoding/hex"
	"log"
	"math/rand"
	"time"

	"github.com/x-contract/neo-go-sdk/neocliapi"

	"github.com/x-contract/neo-go-sdk/neotransaction"
	"github.com/x-contract/neo-go-sdk/neotransaction/OpCode"
)

var (
	contractHashString = `dc675afc61a7c0f7b3d2682bf6e1d8ed865a0e5f`
)

// InvokeContract 调用一个智能合约的 BalanceOf 接口
func InvokeContract() {

	contractHash, _ := hex.DecodeString(contractHashString)

	// The key used to sign the transaction if needed
	//key, _ := neotransaction.DecodeFromWif("Your Private Key's WIF")
	//addr := key.CreateBasicAddress()

	// 创建一个 Invocation 交易
	tx := neotransaction.CreateInvocationTransaction()

	extra := tx.ExtraData.(*neotransaction.InvocationExtraData)
	sb := neotransaction.ScriptBuilder{}

	// If you want to make an invocation transaction without utxo transferd
	// then you need to push a random number so that the hash(txid) could vary on each transaction
	rand.Seed(time.Now().UnixNano())
	sb.EmitPushNumber(int64(rand.Uint32()))
	sb.Emit(OpCode.DROP)

	//args := []interface{}{205, addr.ScripHash}
	//sb.EmitPushArray(args)
	sb.EmitPushBool(false)
	sb.EmitPushString(`name`)
	sb.EmitAppCall(contractHash)

	extra.Script = sb.Bytes()

	// If the transaction need additional Witness then put the ScriptHash in attributes
	//tx.AppendAttribute(neotransaction.UsageScript, addr.ScripHash)

	// Perhaps the transaction need Witness
	//tx.AppendBasicSignWitness(key)

	log.Printf(`Generate invocation transaction[%s]`, tx.TXID())
	log.Printf(`  transaction content:`)
	rawtx := tx.RawTransactionString()
	log.Printf(rawtx)

	result := neocliapi.SendRawTransaction(neocliurl, rawtx)
	log.Printf(`Send transaction to neo-cli node result[%v]`, result)
}
