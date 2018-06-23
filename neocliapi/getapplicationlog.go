package neocliapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/x-contract/neo-go-sdk/neotransaction"
)

// GetApplicationLog 向一个neo-cli节点获取一次合约调用交易(InvocationTransaction)的执行结果
func GetApplicationLog(url string, txid string) ([]Argument, int64, error) {

	client := &http.Client{}
	client.Timeout = 60 * time.Second
	response, err := client.Post(url, "application/json", strings.NewReader(`{
		"jsonrpc": "2.0",
		"method": "getapplicationlog",
		"params": ["`+txid+`"],
		"id": 1
	}`))

	if err != nil {
		return nil, 0, fmt.Errorf(`GetApplicationLog error: %v`, err)
	}
	defer response.Body.Close()

	buff, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, 0, fmt.Errorf(`GetApplicationLog error: %v`, err)
	}

	ret := make(map[string]interface{})
	err = json.Unmarshal(buff, &ret)
	if err != nil {
		return nil, 0, fmt.Errorf(`GetApplicationLog error: %v`, err)
	}

	result, ok := ret[`result`].(map[string]interface{})
	if !ok {
		return nil, 0, fmt.Errorf(`GetApplicationLog error: return value contains no result or result is not a object`)
	}

	state, ok := result[`vmstate`].(string)
	if !ok {
		return nil, 0, fmt.Errorf(`GetApplicationLog error: no state return`)
	}

	if state != `HALT, BREAK` {
		return nil, 0, fmt.Errorf(`GetApplicationLog error: eval state "%s"`, string(buff))
	}

	gas := int64(0)

	gasconsumed, ok := result[`gas_consumed`].(string)
	if !ok {
		v, err := strconv.ParseFloat(gasconsumed, 10)
		if err == nil {
			gas = int64(v * float64(neotransaction.TxOutputValueBase))
		}
	}

	stack, ok := result[`stack`].([]interface{})
	if !ok {
		return nil, gas, fmt.Errorf(`GetApplicationLog error: no stack in eval result`)
	}

	args := make([]Argument, 0, len(stack))

	for _, v := range stack {
		arg, ok := v.(map[string]interface{})
		if !ok {
			return nil, gas, fmt.Errorf(`GetApplicationLog error: stack element[%v] parse failed`, v)
		}
		a := Argument{}
		a.Type, ok = arg[`type`].(string)
		a.Value, ok = arg[`value`].(string)
		args = append(args, a)
	}

	return args, gas, nil
}
