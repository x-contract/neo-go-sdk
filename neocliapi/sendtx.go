package neocliapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// SendRawTransaction 向一个neo-cli节点发送原始交易字符串
func SendRawTransaction(url string, rawtx string) bool {
	client := &http.Client{}
	client.Timeout = 60 * time.Second
	response, err := client.Post(url, "application/json", strings.NewReader(`{
		"jsonrpc": "2.0",
		"method": "sendrawtransaction",
		"params": ["`+rawtx+`"],
		"id": 1
	}`))

	if err != nil {
		log.Println("Try sendRawTransaction", err)
		return false
	}
	defer response.Body.Close()

	buff, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Try sendRawTransaction", err)
		return false
	}

	ret := make(map[string]interface{})
	err = json.Unmarshal(buff, &ret)
	if err != nil {
		log.Println("Try sendRawTransaction", err)
		return false
	}

	result, ok := ret[`result`].(bool)
	if !ok {
		//log.Println("Try sendRawTransaction result not a bool")
		return false
	}

	//log.Println("Try sendRawTransaction", string(buff))
	return result
}
