package neocliapi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// FetchBlockHeight 获取区块高度
func FetchBlockHeight(url string) (uint64, error) {
	reader := strings.NewReader(`{
		"jsonrpc": "2.0",
		"method": "getblockcount",
		"params": [""],
		"id": 1
	}`)
	client := http.Client{Timeout: 30 * time.Second}
	response, err := client.Post(url, "application/json", reader)

	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	buff, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	ret := make(map[string]interface{})
	if err = json.Unmarshal(buff, &ret); err != nil {
		return 0, err
	}

	result, ok := ret[`result`].(float64)
	if !ok {
		return 0, errors.New(`no result`)
	}

	height := uint64(result) - 1
	return height, nil
}
