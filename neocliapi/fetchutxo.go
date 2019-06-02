package neocliapi

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/x-contract/neo-go-sdk/neoutils"

	"github.com/x-contract/neo-go-sdk/neotransaction"
)

///////////////////////////////////////////////////////////////////////////
/// The neo-cli 2.10.2 add a new json rpc api getunspents
/// Just use it :)
///////////////////////////////////////////////////////////////////////////

// FetchUTXO 从neo-cli扩展节点的api接口获得一个账户的utxo数据
func FetchUTXO(url string, address *neotransaction.Address, assetFilter string) ([]*neotransaction.UTXO, error) {
	client := &http.Client{}
	client.Timeout = 60 * time.Second
	response, err := client.Post(url, "application/json", strings.NewReader(`{
		"jsonrpc": "2.0",
		"method": "getunspents",
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

	ret := &struct {
		JsonRPC string
		ID      int
		Result  *struct {
			Address string
			Balance []struct {
				AssetHash string `json:"asset_hash"`
				Asset     string
				Symbol    string `json:"asset_symbol"`
				Amount    float64
				Unspent   []struct {
					TXID  string
					N     uint16
					Value float64
				}
			}
		}
	}{}

	if err = json.Unmarshal(buff, &ret); err != nil {
		return nil, err
	}

	if ret.Result == nil {
		return nil, fmt.Errorf(`FetchUTXO for address[%v] failed %s`, address.Addr, buff)
	}

	utxos := make([]*neotransaction.UTXO, 0)
	for _, asset := range ret.Result.Balance {
		for _, txout := range asset.Unspent {
			utxo := &neotransaction.UTXO{}
			utxo.TxHash, _ = hex.DecodeString(strings.TrimPrefix(txout.TXID, "0x"))
			utxo.TxHash = neoutils.Reverse(utxo.TxHash)
			utxo.Index = txout.N
			utxo.AssetID, _ = hex.DecodeString(asset.AssetHash)
			utxo.AssetID = neoutils.Reverse(utxo.AssetID)
			utxo.Value = int64(txout.Value * float64(neotransaction.TxOutputValueBase))
			utxo.ScriptHash = address.ScripHash
			utxos = append(utxos, utxo)
		}
	}

	return utxos, nil
}
