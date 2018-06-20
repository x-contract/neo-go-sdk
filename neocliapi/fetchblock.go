package neocliapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// FetchBlock 获取区块
func FetchBlock(url string, height uint64) (map[string]interface{}, error) {
	reader := strings.NewReader(fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "getblock",
		"params": [%v, 1],
		"id": 1
	}`, height))
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

	return result, nil
}
