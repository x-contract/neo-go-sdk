package main

import (
	"log"

	"github.com/x-contract/neo-go-sdk/neocliapi"
	"github.com/x-contract/neo-go-sdk/neoextapi"
	"github.com/x-contract/neo-go-sdk/neotransaction"
)

var (
	neoextapiurl = `https://api.nel.group/api/testnet`
	neocliurl    = `http://seed4.neo.org:20332`
)

// TransferUTXO UTXO资产转账交易
func TransferUTXO() {

	// The private-keys of 2 basic address
	// addr1 will transfer GAS to addr2
	// and addr2 will transfer NEO to addr1 within the same transaction
	// Put your own private-keys WIF string here to test
	key1, _ := neotransaction.DecodeFromWif("WIF of Private-Key 1")
	key2, _ := neotransaction.DecodeFromWif("WIF of Private-Key 2")

	addr1 := key1.CreateBasicAddress()
	addr2 := key2.CreateBasicAddress()

	log.Printf(addr1.Addr)
	log.Printf(addr2.Addr)

	utxos1, _ := neoextapi.FetchUTXO(neoextapiurl, addr1, neotransaction.AssetGasID)
	utxos2, _ := neoextapi.FetchUTXO(neoextapiurl, addr2, neotransaction.AssetNeoID)

	//log.Printf(utxos1)
	//log.Printf(utxos2)

	tx := neotransaction.CreateContractTransaction()

	//var utxo1 *UTXO
	value1 := int64(0)
	for _, utxo := range utxos1 {
		value1 += utxo.Value
		tx.AppendInput(utxo)
	}

	value2 := int64(0)
	for _, utxo := range utxos2 {
		value2 += utxo.Value
		tx.AppendInput(utxo)
	}

	wantNeoValue := (value1 / neotransaction.TxOutputValueBase) / 4 * neotransaction.TxOutputValueBase
	payGasValue := wantNeoValue * 4

	changeGasBack := value1 - payGasValue
	changeNeoBack := int64(0)

	if value2 < wantNeoValue {
		wantNeoValue = value2
		payGasValue = wantNeoValue * 4
		changeGasBack = value1 - payGasValue
	} else if value2 > wantNeoValue {
		changeNeoBack = value2 - wantNeoValue
	}

	log.Printf("账户1 支出GAS[%v]", payGasValue)
	log.Printf("账户1 获得NEO[%v]", wantNeoValue)
	log.Printf("账户1 找零GAS[%v]", changeGasBack)
	log.Printf("账户2 找零NEO[%v]", changeNeoBack)

	tx.AppendOutput(addr1, utxos2[0].AssetID, wantNeoValue)
	tx.AppendOutput(addr2, utxos1[0].AssetID, payGasValue)

	if changeGasBack > 0 {
		tx.AppendOutput(addr1, utxos1[0].AssetID, changeGasBack)
	}
	if changeNeoBack > 0 {
		tx.AppendOutput(addr2, utxos2[0].AssetID, changeNeoBack)
	}

	tx.AppendBasicSignWitness(key1)
	tx.AppendBasicSignWitness(key2)

	log.Printf(`Generate contract transaction[%s]`, tx.TXID())
	log.Printf(`  transaction content:`)
	rawtx := tx.RawTransactionString()
	log.Printf(rawtx)

	result := neocliapi.SendRawTransaction(neocliurl, rawtx)
	log.Printf(`Send transaction to neo-cli node result[%v]`, result)
}
