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

package openwtester

import (
	"github.com/blocktree/openwallet/v2/common/file"
	"os"
	"path/filepath"
	"testing"

	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openw"
	"github.com/blocktree/openwallet/v2/openwallet"
)

type SubscribeTestCase struct{
	Name		string
	ToAddress 	string
	Amount		string
	Height		uint64
	Txid		string
	Ok			bool
}

////////////////////////// 测试单个扫描器 //////////////////////////

type subscriberSingle struct {
	testCaseArr	[]SubscribeTestCase
}

//BlockScanNotify 新区块扫描完成通知
func (sub *subscriberSingle) BlockScanNotify(header *openwallet.BlockHeader) error {
	//log.Notice("header:", header)
	return nil
}

//BlockTxExtractDataNotify 区块提取结果通知
func (sub *subscriberSingle) BlockExtractDataNotify(sourceKey string, data *openwallet.TxExtractData) error {
	log.Notice("account:", sourceKey)

	for i, input := range data.TxInputs {
		log.Std.Notice("data.TxInputs[%d]: %+v", i, input)
	}

	for i, output := range data.TxOutputs {
		log.Std.Notice("data.TxOutputs[%d]: %+v", i, output)
	}

	log.Std.Notice("data.Transaction: %+v", data.Transaction)

	for i := 0; i < len(sub.testCaseArr); i++{
		testCase := sub.testCaseArr[i]

		if testCase.Ok {
			continue
		}

		if testCase.Name == sourceKey {
			for _, output := range data.TxOutputs {
				isAddressRight := false
				if output.Address == testCase.ToAddress+":"+testCase.Amount{	//output 符合
					isAddressRight = true
				}

				isOutputTxIDRight := false
				if output.TxID == testCase.Txid {
					isOutputTxIDRight = true
				}

				isOutputAmountRight := false
				if output.Amount == testCase.Amount {
					isOutputAmountRight = true
				}

				isOutputHeightRight := false
				if output.BlockHeight == testCase.Height {
					isOutputHeightRight = true
				}

				if isAddressRight && isOutputTxIDRight && isOutputAmountRight && isOutputHeightRight{
					isTransactionTxIDRight := false
					if data.Transaction.TxID == testCase.Txid {
						isTransactionTxIDRight = true
					}

					isTransactionToRight := false
					for _, toItem := range data.Transaction.To {
						if toItem==testCase.ToAddress+":"+testCase.Amount{
							isTransactionToRight = true
							break
						}
					}

					isTransactionHeightRight := false
					if data.Transaction.BlockHeight == testCase.Height {
						isTransactionHeightRight = true
					}

					isTransactionStatusRight := false
					if data.Transaction.Status == "1" {
						isTransactionStatusRight = true
					}

					if isTransactionTxIDRight && isTransactionToRight && isTransactionHeightRight && isTransactionStatusRight{
						sub.testCaseArr[i].Ok = true
						log.Std.Notice("Test Case %+v is ok", testCase.Name )
					}
				}
			}
		}
	}

	okTimes := 0
	for i := 0; i < len(sub.testCaseArr); i++{
		testCase := sub.testCaseArr[i]

		if testCase.Ok {
			okTimes++
		}

		if okTimes==len(sub.testCaseArr) {
			log.Std.Info("All Test Case Finished, Over !!!")
			os.Exit(0)
		}
	}

	return nil
}

func (sub *subscriberSingle) BlockExtractSmartContractDataNotify(sourceKey string, data *openwallet.SmartContractReceipt) error {
	return nil
}

//BlockScanNotify 新区块扫描完成通知
func (sub *subscriberSingle) InitTestCases() error {
	//testCase := SubscribeTestCase{}

	testCaseArr := make([]SubscribeTestCase, 0)

	sub.testCaseArr = testCaseArr

	return nil
}

func TestSubscribeAddress(t *testing.T) {

	var (
		endRunning = make(chan bool, 1)
		symbol     = "DOT"
		addrs      = map[string]string{}
	)

	sub := subscriberSingle{}
	sub.InitTestCases()

	for i := 0; i < len(sub.testCaseArr); i++{
		testCase := sub.testCaseArr[i]

		addrs[ testCase.ToAddress ] = testCase.Name
	}

	//GetSourceKeyByAddress 获取地址对应的数据源标识
	scanTargetFunc := func(target openwallet.ScanTarget) (string, bool) {
		key, ok := addrs[target.Address]
		if !ok {
			return "", false
		}
		return key, true
	}

	assetsMgr, err := openw.GetAssetsAdapter(symbol)
	if err != nil {
		log.Error(symbol, "is not support")
		return
	}

	//读取配置
	absFile := filepath.Join(configFilePath, symbol+".ini")

	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return
	}
	assetsMgr.LoadAssetsConfig(c)

	assetsLogger := assetsMgr.GetAssetsLogger()
	if assetsLogger != nil {
		assetsLogger.SetLogFuncCall(true)
	}

	//log.Debug("already got scanner:", assetsMgr)
	scanner := assetsMgr.GetBlockScanner()

	if scanner.SupportBlockchainDAI() {
		file.MkdirAll(dbFilePath)
		dai, err := openwallet.NewBlockchainLocal(filepath.Join(dbFilePath, dbFileName), false)
		if err != nil {
			log.Error("NewBlockchainLocal err: %v", err)
			return
		}

		scanner.SetBlockchainDAI(dai)
	}

	scanner.SetRescanBlockHeight( sub.testCaseArr[0].Height )

	if scanner == nil {
		log.Error(symbol, "is not support block scan")
		return
	}

	scanner.SetBlockScanTargetFunc(scanTargetFunc)

	scanner.AddObserver(&sub)

	for i := 0; i < len(sub.testCaseArr); i++{
		testCase := sub.testCaseArr[i]

		scanner.ScanBlock(testCase.Height)
	}

	scanner.Run()

	<-endRunning
}