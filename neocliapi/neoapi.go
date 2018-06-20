package neocliapi

import (
	"log"
	"time"
)

// CurrBlockHeight 当前已经抓取到的区块高度
var CurrBlockHeight = uint64(0)

// NEOCLIURL ...
var NEOCLIURL = ``

// NewBlockChan 新区块的Channel
var NewBlockChan = make(chan map[string]interface{})

// StartSpider 开始监听NEO节点
func StartSpider(cliurl string, fromHeight uint64) chan map[string]interface{} {
	NEOCLIURL = cliurl
	CurrBlockHeight = fromHeight
	log.Printf("neoapi: fetch init block height[%v]\n", CurrBlockHeight)
	go func() {
		for {
			fetchBlock()
			time.Sleep(3 * time.Second)
		}
	}()
	return NewBlockChan
}

func fetchBlock() {
	height, err := FetchBlockHeight(NEOCLIURL)
	if err != nil {
		log.Println(`neoapi: fetch block height error`, err)
		return
	}
	if height <= CurrBlockHeight {
		return
	}

	for i := CurrBlockHeight + 1; i <= height; i++ {
		block, err := FetchBlock(NEOCLIURL, i)
		if err != nil {
			log.Printf("neoapi: fetch block[%v] %v\n", i, err)
			return
		}
		NewBlockChan <- block
		CurrBlockHeight = i
	}
	CurrBlockHeight = height
}
