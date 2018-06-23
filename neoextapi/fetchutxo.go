package neoextapi

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/x-contract/neo-go-sdk/neoutils"

	"github.com/x-contract/neo-go-sdk/neotransaction"
)

///////////////////////////////////////////////////////////////////////////
/// The neo-cli does not implement the api to get utxo datas
/// Using NEL extension block api (https://github.com/NewEconoLab) for now
/// Implementation of golang is on the way :)
///////////////////////////////////////////////////////////////////////////

// FetchUTXO 从neo-cli扩展节点的api接口获得一个账户的utxo数据
func FetchUTXO(url string, address *neotransaction.Address, assetFilter string) ([]*neotransaction.UTXO, error) {
	client := &http.Client{}
	client.Timeout = 60 * time.Second
	response, err := client.Post(url, "application/json", strings.NewReader(`{
		"jsonrpc": "2.0",
		"method": "getutxo",
		"params": ["`+address.Addr+`"],
		"id": 1
	}`))

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	buff, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]interface{})
	if err = json.Unmarshal(buff, &ret); err != nil {
		return nil, err
	}

	result, ok := ret["result"].([]interface{})
	if !ok {
		return nil, errors.New("fetchUTXO result missing")
	}

	utxos := make([]*neotransaction.UTXO, 0, len(result))
	for _, obj := range result {
		u, ok := obj.(map[string]interface{})
		if !ok {
			continue
		}
		utxo := &neotransaction.UTXO{}
		txid := u["txid"].(string)
		utxo.TxHash, _ = hex.DecodeString(strings.TrimPrefix(txid, "0x"))
		utxo.TxHash = neoutils.Reverse(utxo.TxHash)

		utxo.Index = uint16(u["n"].(float64))

		assetid := u["asset"].(string)
		assetid = strings.TrimPrefix(assetid, "0x")
		if assetFilter != "*" && assetFilter != assetid {
			continue
		}
		utxo.AssetID, _ = hex.DecodeString(assetid)
		utxo.AssetID = neoutils.Reverse(utxo.AssetID)

		value := u["value"].(string)
		v, _ := strconv.ParseFloat(value, 64)
		utxo.Value = int64(v * float64(neotransaction.TxOutputValueBase))

		utxo.ScriptHash = address.ScripHash

		utxos = append(utxos, utxo)
	}

	return utxos, nil
}
