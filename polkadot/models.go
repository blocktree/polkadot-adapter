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
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/tidwall/gjson"
	"math/big"
	"strconv"
	"time"
)

const BATCH_CHARGE_TO_TAG  = "batch_charge"

type Block struct {
	Hash          string        `json:"block"`         // actually block signature in XRP chain
	PrevBlockHash string        `json:"previousBlock"` // actually block signature in DOT chain
	Timestamp     uint64        `json:"timestamp"`
	Height        uint64        `json:"height"`
	Transactions  []Transaction `json:"transactions"`
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
	ToArr       []string //@required 格式："地址":"数量"
	ToDecArr    []string //@required 格式："地址":"数量(带小数)"
}

func GetTransactionInBlock(json *gjson.Result) []Transaction {
	blockHash := gjson.Get(json.Raw, "hash").String()
	blockHeight := gjson.Get(json.Raw, "number").Uint()
	transactions := make([]Transaction, 0)

	blockTime := uint64(time.Now().Unix())

	for _, extrinsic := range gjson.Get(json.Raw, "extrinsics").Array() {
		//fmt.Println("extrinsic : " , extrinsic)
		method := gjson.Get(extrinsic.Raw, "method").String()
		success := gjson.Get(extrinsic.Raw, "success").Bool()
		//fmt.Println("method : ", method, ", success : ", success)
		//hasUtilityComplete := false

		if !success {
			continue
		}

		//获取这个区块的时间
		if method == "timestamp.set" {
			args := gjson.Get(extrinsic.Raw, "args")
			if len(args.Raw) >0 {
				blockTime = gjson.Get(args.Raw, "now").Uint()
			}
		}

		//解析批量转账
		if method == "utility.batch" {

			txid := gjson.Get(extrinsic.Raw, "hash").String()

			toArr := make([]string, 0)
			toAmount := uint64(0)
			batchTransaction := make([]Transaction, 0)

			args := gjson.Get(extrinsic.Raw, "args")
			if len(args.Raw) >0 {
				calls := gjson.Get(args.Raw, "calls")
				if len(calls.Raw) >0 {
					for _, callItem := range calls.Array() {
						if gjson.Get(callItem.Raw, "method").String() == "balances.transfer" {	//在交易体，发现转账方法

							callIndex := gjson.Get(callItem.Raw, "callIndex")
							callIndex0 := gjson.Get(callIndex.Raw, "0")
							if len(callIndex0.Raw)>0 {
								if callIndex0.String() != "5" {
									continue
								}
							}
							callIndex1 := gjson.Get(callIndex.Raw, "1")
							if len(callIndex1.Raw)>0 {
								if callIndex1.String() != "0" {
									continue
								}
							}

							dest := ""
							value := ""
							callArgs := gjson.Get(callItem.Raw, "args")
							if len(callArgs.Raw)>0 {
								dest = gjson.Get(callArgs.Raw, "dest").String()
								value = gjson.Get(callArgs.Raw, "value").String()
							}

							amountInt, err := strconv.ParseInt(value, 10, 64)
							if err == nil {
								amount := uint64(amountInt)
								transaction := Transaction{
									TxID:        txid,
									TimeStamp:   blockTime,
									From:        "",
									To:          dest,
									Amount:      amount,
									BlockHeight: blockHeight,
									BlockHash:   blockHash,
									Status:      "-1",
								}

								batchTransaction = append(batchTransaction, transaction)
							}
						}
					}
				}
			}

			for _, event := range gjson.Get(extrinsic.Raw, "events").Array() {
				//if gjson.Get(event.Raw, "method").String() == "utility.BatchCompleted" {	//事件是否执行完成
				//	hasUtilityComplete = true
				//}
				if gjson.Get(event.Raw, "method").String() == "balances.Transfer" {
					data := gjson.Get(event.Raw, "data").Array()
					if len(data) == 3 {
						from := data[0].String()
						to := data[1].String()
						amountStr := data[2].String()

						for batchTransactionIndex, batchTransactionItem := range batchTransaction {
							if batchTransactionItem.Status == "-1" { //未被认定
								if batchTransactionItem.To == to { //转入地址与传入参数对得上
									amountInt, err := strconv.ParseInt(amountStr, 10, 64)
									if err == nil {
										amount := uint64(amountInt)
										if batchTransactionItem.Amount == amount { //金额与传入参数对得上
											batchTransaction[batchTransactionIndex].From = from
											toArr = append(toArr, to+":"+amountStr)
											toAmount, _ = math.SafeAdd(toAmount, amount)
											batchTransaction[batchTransactionIndex].Status = "1"

											break
										}
									}
								}
							}
						}
					}
				}
			}

			//if hasUtilityComplete==false{
			//	continue
			//}

			fee := uint64(0)

			tip := uint64(gjson.Get(extrinsic.Raw, "tip").Uint())

			info := gjson.Get(extrinsic.Raw, "info")
			if info.Exists() {
				partialFee := uint64(gjson.Get(info.Raw, "partialFee").Uint())

				fee, _ = math.SafeAdd(tip, partialFee)
			}

			for _, batchTransactionItem := range batchTransaction {
				transaction := Transaction{
					TxID:        batchTransactionItem.TxID,
					Fee:         fee,
					TimeStamp:   batchTransactionItem.TimeStamp,
					From:        batchTransactionItem.From,
					To:          BATCH_CHARGE_TO_TAG,
					Amount:      toAmount,
					BlockHeight: batchTransactionItem.BlockHeight,
					BlockHash:   batchTransactionItem.BlockHash,
					Status:      batchTransactionItem.Status,
					ToArr:       toArr,
				}

				transactions = append(transactions, transaction)
				break
			}

			continue
		}

		if method != "balances.transfer" && method != "claims.attest" && method != "balances.transferKeepAlive" { //不是这个method的全部不要
			continue
		}

		argsTo := ""        //检测到的接收地址
		argsAmountStr := "" //检测到的接收金额
		from := ""          //来源地址
		to := ""            //目标地址
		amountStr := ""     //金额
		args := gjson.Get(extrinsic.Raw, "args")
		if len(args.Raw) >0 {
			argsTo = gjson.Get(args.Raw, "dest").String()
			argsAmountStr = gjson.Get(args.Raw, "value").String()
		}

		for _, event := range gjson.Get(extrinsic.Raw, "events").Array() {
			if gjson.Get(event.Raw, "method").String() == "balances.Transfer" {
				data := gjson.Get(event.Raw, "data").Array()
				if len(data) == 3 {
					from = data[0].String()
					to = data[1].String()
					amountStr = data[2].String()
				}
			}
			if gjson.Get(event.Raw, "method").String() == "claims.Claimed" {
				data := gjson.Get(event.Raw, "data").Array()
				if len(data) == 3 {
					//from = data[1].String()
					to = data[0].String()
					amountStr = data[2].String()
				}
			}
		}

		if argsTo == "" && to == "" { //没有取到值
			continue
		}
		if argsAmountStr == "" && amountStr == "" { //没有取到值
			continue
		}
		if method == "balances.transfer" && argsTo != to { //值不一样
			continue
		}
		if method == "balances.transfer" && argsAmountStr != amountStr { //值不一样
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
		//log.Debug("txid : ", txid, ",from: ", from, ",to: ", to, ",amount: ", amountStr, ",time: ", blockTime, ",fee: ", fee)

		amountInt, err := strconv.ParseInt(amountStr, 10, 64)
		if err == nil {
			amount := uint64(amountInt)

			transaction := Transaction{
				TxID:        txid,
				Fee:         fee,
				TimeStamp:   blockTime,
				From:        from,
				To:          to,
				Amount:      amount,
				BlockHeight: blockHeight,
				BlockHash:   blockHash,
				Status:      "1",
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

type TxArtifacts struct {
	Hash        string
	Height      int64
	GenesisHash string
	SpecVersion uint32
	Metadata    string
	TxVersion   uint32
	ChainName   string
}

func GetTxArtifacts(json *gjson.Result) *TxArtifacts {
	obj := &TxArtifacts{}

	obj.Hash = gjson.Get(json.Raw, "at").Get("hash").String()
	obj.Height = gjson.Get(json.Raw, "at").Get("height").Int()
	obj.GenesisHash = gjson.Get(json.Raw, "genesisHash").String()
	obj.SpecVersion = uint32(gjson.Get(json.Raw, "specVersion").Uint())
	obj.Metadata = gjson.Get(json.Raw, "metadata").String()
	obj.TxVersion = uint32(gjson.Get(json.Raw, "txVersion").Uint())
	obj.ChainName = gjson.Get(json.Raw, "chainName").String()

	return obj
}
