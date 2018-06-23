package neocliapi

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/x-contract/neo-go-sdk/neotransaction"
	"github.com/x-contract/neo-go-sdk/neoutils"
)

// Argument 调用智能合约的参数和智能合约返回值的参数类型
type Argument struct {
	Type  string
	Value string
}

func serializeParamString(params []interface{}) (string, error) {
	paramString := ``
	for _, param := range params {
		if param == nil {
			continue
		}
		if len(paramString) > 0 {
			paramString += `,`
		}

		v := reflect.ValueOf(param)
		if v.Kind() == reflect.Slice && v.Type().Elem().Kind() != reflect.Uint8 {
			sub := make([]interface{}, v.Len())
			for i := 0; i < v.Len(); i++ {
				sub[i] = v.Index(i).Interface()
			}
			substring, err := serializeParamString(sub)
			if err != nil {
				return ``, err
			}
			paramString += `{ "type": "Array", "value": [` + substring + `]}`
			continue
		}

		switch v := param.(type) {
		case string:
			paramString += `{ "type": "String", "value": "` + v + `"}`
		case bool:
			paramString += `{ "type": "Boolean", "value": "` + strconv.FormatBool(v) + `"}`
		case int32, int64, uint32, uint64:
			paramString += `{ "type": "Integer", "value": "` + fmt.Sprint(v) + `"}`
		case []byte:
			paramString += `{ "type": "ByteArray", "value": "` + hex.EncodeToString(v) + `"}`
		case neoutils.HASH256:
			paramString += `{ "type": "Hash256", "value": "` + hex.EncodeToString(v) + `"}`
		case neoutils.HASH160:
			paramString += `{ "type": "Hash160", "value": "` + hex.EncodeToString(v) + `"}`
		default:
			return ``, fmt.Errorf(`InvokeScript error: params contains invalid type[%T]`, v)
		}
	}
	return paramString, nil
}

// Invoke 向一个neo-cli节点调用一个已发布的智能合约
// 注意：这种调用方法只能调用一个查询类接口，不修改智能合约内存储数据，结果也不会上链，只是在本地节点上模拟运行
// 如果需要上链的调用，需要拼一个 InvocationTransaction 然后调用 InvokeScript 接口广播此交易
func Invoke(url string, scriptHashString string, params []interface{}) ([]Argument, int64, error) {

	paramString, err := serializeParamString(params)
	//log.Println(paramString)
	if err != nil {
		return nil, 0, err
	}

	client := &http.Client{}
	client.Timeout = 60 * time.Second
	response, err := client.Post(url, "application/json", strings.NewReader(`{
		"jsonrpc": "2.0",
		"method": "invoke",
		"params": ["`+scriptHashString+`", [`+paramString+`]],
		"id": 1
	}`))

	if err != nil {
		return nil, 0, fmt.Errorf(`InvokeScript error: %v`, err)
	}
	defer response.Body.Close()

	buff, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, 0, fmt.Errorf(`InvokeScript error: %v`, err)
	}

	ret := make(map[string]interface{})
	err = json.Unmarshal(buff, &ret)
	if err != nil {
		return nil, 0, fmt.Errorf(`InvokeScript error: %v`, err)
	}

	result, ok := ret[`result`].(map[string]interface{})
	if !ok {
		return nil, 0, fmt.Errorf(`InvokeScript error: return value contains no result or result is not a object`)
	}

	state, ok := result[`state`].(string)
	if !ok {
		return nil, 0, fmt.Errorf(`InvokeScript error: no state return`)
	}

	if state != `HALT, BREAK` {
		return nil, 0, fmt.Errorf(`InvokeScript error: eval state "%s"`, string(buff))
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
		return nil, gas, fmt.Errorf(`InvokeScript error: no stack in eval result`)
	}

	args := make([]Argument, 0, len(stack))

	for _, v := range stack {
		arg, ok := v.(map[string]interface{})
		if !ok {
			return nil, gas, fmt.Errorf(`InvokeScript error: stack element[%v] parse failed`, v)
		}
		a := Argument{}
		a.Type, ok = arg[`type`].(string)
		a.Value, ok = arg[`value`].(string)
		args = append(args, a)
	}

	return args, gas, nil
}

// InvokeScript ...
func InvokeScript(url string, script []byte) ([]Argument, int64, error) {
	client := &http.Client{}
	client.Timeout = 60 * time.Second
	response, err := client.Post(url, "application/json", strings.NewReader(`{
		"jsonrpc": "2.0",
		"method": "invokescript",
		"params": ["`+hex.EncodeToString(script)+`"],
		"id": 1
	}`))

	if err != nil {
		return nil, 0, fmt.Errorf(`InvokeScript error: %v`, err)
	}
	defer response.Body.Close()

	buff, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, 0, fmt.Errorf(`InvokeScript error: %v`, err)
	}

	ret := make(map[string]interface{})
	err = json.Unmarshal(buff, &ret)
	if err != nil {
		return nil, 0, fmt.Errorf(`InvokeScript error: %v`, err)
	}

	result, ok := ret[`result`].(map[string]interface{})
	if !ok {
		return nil, 0, fmt.Errorf(`InvokeScript error: return value contains no result or result is not a object`)
	}

	state, ok := result[`state`].(string)
	if !ok {
		return nil, 0, fmt.Errorf(`InvokeScript error: no state return`)
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
		return nil, gas, fmt.Errorf(`InvokeScript error: no stack in eval result`)
	}

	args := make([]Argument, 0, len(stack))

	for _, v := range stack {
		arg, ok := v.(map[string]interface{})
		if !ok {
			return nil, gas, fmt.Errorf(`InvokeScript error: stack element[%v] parse failed`, v)
		}
		a := Argument{}
		a.Type, ok = arg[`type`].(string)
		a.Value, ok = arg[`value`].(string)
		args = append(args, a)
	}

	if state != `HALT, BREAK` {
		return args, gas, fmt.Errorf(`InvokeScript error: eval state "%s"`, state)
	}

	return args, gas, nil
}
