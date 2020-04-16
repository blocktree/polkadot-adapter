/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */
package polkadot

import (
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/tidwall/gjson"
	"math/big"
	"strconv"
	"time"
)

type Block struct {
	Hash                  string		`json:"block"` // actually block signature in XRP chain
	PrevBlockHash         string 		`json:"previousBlock"` // actually block signature in DOT chain
	Timestamp             uint64		`json:"timestamp"`
	Height                uint64   		`json:"height"`
	Transactions          []Transaction`json:"transactions"`
}

type Transaction struct {
	TxID        string
	Fee         uint64
	TimeStamp   uint64
	From        string
	To          string
	Amount      uint64
	BlockHeight uint64
	BlockHash   string
	Status      string
}

func GetTransactionInBlock(json *gjson.Result) []Transaction{
	blockHash := gjson.Get(json.Raw, "hash").String()
	blockHeight := gjson.Get(json.Raw, "number").Uint()
	transactions := make([]Transaction, 0)

	blockTime := uint64(time.Now().Unix())

	for _, extrinsic := range gjson.Get(json.Raw, "extrinsics").Array(){
		//fmt.Println("extrinsic : " , extrinsic)
		method := gjson.Get(extrinsic.Raw, "method").String()
		success := gjson.Get(extrinsic.Raw, "success").Bool()
		//fmt.Println("method : ", method, ", success : ", success)

		if !success {
			continue
		}

		//获取这个区块的时间
		if method=="timestamp.set" {
			args := gjson.Get(extrinsic.Raw, "args").Array()
			if len(args)==2 {
				blockTime = args[0].Uint()
			}
		}

		if method!="balances.transfer" {	//不是这个method的全部不要
			continue
		}

		argsTo := ""			//检测到的接收地址
		argsAmountStr := ""		//检测到的接收金额
		from := ""				//来源地址
		to := ""				//目标地址
		amountStr := ""			//金额
		args := gjson.Get(extrinsic.Raw, "args").Array()
		if len(args)==2 {
			argsTo = args[0].String()
			argsAmountStr = args[1].String()
		}

		for _, event := range gjson.Get(extrinsic.Raw, "events").Array(){
			if gjson.Get(event.Raw, "method").String() == "balances.Transfer" {
				data := gjson.Get(event.Raw, "data").Array()
				if len(data)==3{
					from = data[0].String()
					to = data[1].String()
					amountStr = data[2].String()
				}
			}
		}

		if argsTo=="" && to=="" {	//没有取到值
			continue
		}
		if argsAmountStr=="" && amountStr=="" {	//没有取到值
			continue
		}
		if argsTo!=to{	//值不一样
			continue
		}
		if argsAmountStr!=amountStr{	//值不一样
			continue
		}

		txid := gjson.Get(extrinsic.Raw, "hash").String()

		fee := uint64(0)

		tip := uint64(gjson.Get(extrinsic.Raw, "tip").Uint())

		info := gjson.Get(extrinsic.Raw, "info")
		if info.Exists() {
			partialFee := uint64(gjson.Get(info.Raw, "partialFee").Uint())

			fee, _ = math.SafeAdd(tip, partialFee)
		}

		//fmt.Println("txid : ", txid, ",from: ", from, ",to: ", to, ",amount: ", amountStr, ",time: " ,blockTime, ",fee: ", fee)
		log.Debug("txid : ", txid, ",from: ", from, ",to: ", to, ",amount: ", amountStr, ",time: " ,blockTime, ",fee: ", fee)

		amountInt, err := strconv.ParseInt(amountStr, 10, 64)
		if err == nil{
			amount := uint64(amountInt)

			transaction := Transaction{
				TxID        : txid,
				Fee         : fee,
				TimeStamp   : blockTime,
				From        : from,
				To          : to,
				Amount      : amount,
				BlockHeight : blockHeight,
				BlockHash   : blockHash,
				Status      : "1",
			}

			transactions = append(transactions, transaction)
		}
	}

	return transactions
}

func NewBlock(json *gjson.Result) *Block {
	obj := &Block{}
	// 解析
	obj.Hash = gjson.Get(json.Raw, "hash").String()
	obj.PrevBlockHash = gjson.Get(json.Raw, "parentHash").String()
	obj.Height = gjson.Get(json.Raw, "number").Uint()
	obj.Transactions = GetTransactionInBlock(json)

	if obj.Hash == "" {
		time.Sleep(5 * time.Second)
	}
	return obj
}

//BlockHeader 区块链头
func (b *Block) BlockHeader() *openwallet.BlockHeader {

	obj := openwallet.BlockHeader{}
	//解析json
	obj.Hash = b.Hash
	//obj.Confirmations = b.Confirmations
	obj.Previousblockhash = b.PrevBlockHash
	obj.Height = b.Height
	//obj.Symbol = Symbol

	return &obj
}

type AddrBalance struct {
	Address string
	Balance *big.Int
	Free    *big.Int
	Freeze  *big.Int
	Nonce   uint64
	index   int
	Actived bool
}
