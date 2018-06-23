package neocliapi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/x-contract/neo-go-sdk/neotransaction"
)

// NeoBalance 用户NEO账户余额
type NeoBalance map[string]float64

// FetchBalance 获取账户余额
func FetchBalance(url string, addr string) (NeoBalance, error) {
	reader := strings.NewReader(`{
		"jsonrpc": "2.0",
		"method": "getaccountstate",
		"params": ["` + addr + `"],
		"id": 1
	}`)
	client := http.Client{Timeout: 30 * time.Second}
	response, err := client.Post(url, "application/json", reader)

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

	result, ok := ret[`result`].(map[string]interface{})
	if !ok {
		return nil, errors.New(`no result`)
	}

	balances, ok := result[`balances`].([]interface{})
	if !ok {
		return nil, errors.New(`no balances`)
	}

	neobalance := make(map[string]float64)
	neobalance[`NEO`] = 0
	neobalance[`GAS`] = 0

	for _, i := range balances {
		balance, ok := i.(map[string]interface{})
		if !ok {
			continue
		}
		asset, ok := balance[`asset`].(string)
		if !ok {
			continue
		}
		value, ok := balance[`value`].(string)
		if !ok {
			continue
		}
		v, err := strconv.ParseFloat(value, 10)
		if err != nil {
			continue
		}
		asset = strings.TrimPrefix(asset, `0x`)
		switch asset {
		case neotransaction.AssetNeoID:
			neobalance[`NEO`] = v
		case neotransaction.AssetGasID:
			neobalance[`GAS`] = v
		}
	}

	return neobalance, nil
}
