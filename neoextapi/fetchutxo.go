package neoextapi

import (
	"github.com/x-contract/neo-go-sdk/neocliapi"
	"github.com/x-contract/neo-go-sdk/neotransaction"
)

///////////////////////////////////////////////////////////////////////////
/// The neo-cli 2.10.2 add a new json rpc api getunspents
/// Just use it :)
///////////////////////////////////////////////////////////////////////////

// FetchUTXO 从neo-cli扩展节点的api接口获得一个账户的utxo数据
func FetchUTXO(url string, address *neotransaction.Address, assetFilter string) ([]*neotransaction.UTXO, error) {
	return neocliapi.FetchUTXO(url, address, assetFilter)
}
