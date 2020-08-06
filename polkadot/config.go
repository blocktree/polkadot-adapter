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
	"fmt"
	"path/filepath"
	"strings"
	"time"

	owcrypt "github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/v2/common/file"
	"github.com/shopspring/decimal"
)

const (
	//币种
	Symbol      = "DOT"
	MasterKey   = "Polkadot seed"
	CurveType   = owcrypt.ECC_CURVE_ED25519
	GenesisHash = "b0a8d493285c2df73290dfb7e61f870f17b41801197a149ca93654499ea3dafe"
	SpecVersion = 1058
	AddrPrefix = 0x00
)

type WalletConfig struct {

	//币种
	Symbol    string
	MasterKey string
	//钥匙备份路径
	keyDir string
	//地址导出路径
	addressDir string
	//配置文件路径
	configFilePath string
	//配置文件名
	configFileName string
	//本地数据库文件路径
	dbPath string
	//备份路径
	backupDir string
	// node API
	NodeAPI string
	// websocket API
	WSAPI string
	// rpc - NodeAPI   ws - WSAPI
	APIChoose string
	//钱包安装的路径
	NodeInstallPath string
	//钱包数据文件目录
	WalletDataPath string
	//汇总阀值
	Threshold decimal.Decimal
	//汇总地址
	SumAddress string
	//汇总执行间隔时间
	CycleSeconds time.Duration
	//默认配置内容
	DefaultConfig string
	//曲线类型
	CurveType uint32
	//fixed fee in sawi
	FixedFee int64
	// reserve amount in smallest unit
	ReserveAmount int64
	// ignore reserve amount or not
	IgnoreReserve bool
	// data directory
	DataDir string

	GenesisHash string

	SpecVersion uint32

	AddrPrefix byte
	Decimal int32
}

func NewConfig(symbol string, masterKey string, GenesisHash string, SpecVersion uint32, AddrPrefix byte) *WalletConfig {

	c := WalletConfig{}

	//币种
	c.Symbol = symbol
	c.MasterKey = masterKey
	c.CurveType = CurveType
	c.GenesisHash = GenesisHash
	c.SpecVersion = SpecVersion
	c.AddrPrefix = AddrPrefix

	//钥匙备份路径
	c.keyDir = filepath.Join("data", strings.ToLower(c.Symbol), "key")
	//地址导出路径
	c.addressDir = filepath.Join("data", strings.ToLower(c.Symbol), "address")
	//配置文件路径
	c.configFilePath = filepath.Join("conf")
	//配置文件名
	c.configFileName = c.Symbol + ".ini"
	//本地数据库文件路径
	c.dbPath = filepath.Join("data", strings.ToLower(c.Symbol), "db")
	//备份路径
	c.backupDir = filepath.Join("data", strings.ToLower(c.Symbol), "backup")
	//钱包安装的路径
	c.NodeInstallPath = ""
	//钱包数据文件目录
	c.WalletDataPath = ""
	//汇总阀值
	c.Threshold = decimal.NewFromFloat(5)
	//汇总地址
	c.SumAddress = ""
	//汇总执行间隔时间
	c.CycleSeconds = time.Second * 10

	//默认配置内容
	c.DefaultConfig = `
# start node command
startNodeCMD = ""
# stop node command
stopNodeCMD = ""
# node install path
nodeInstallPath = ""
# mainnet data path
mainNetDataPath = ""
# testnet data path
testNetDataPath = ""
# RPC Server Type，0: Rest API
rpcServerType = 0
# RPC api url
serverAPI = ""
# RPC Authentication Username
rpcUser = ""
# RPC Authentication Password
rpcPassword = ""
# Is network test?
isTestNet = false
# the safe address that wallet send money to.
sumAddress = ""
# when wallet's balance is over this value, the wallet willl send money to [sumAddress]
threshold = ""
# summary task timer cycle time, sample: 1m , 30s, 3m20s etc
cycleSeconds = ""
# walletPassword use to unlock bitcoin core wallet
walletPassword = ""
`

	//创建目录
	// file.MkdirAll(c.dbPath)
	// file.MkdirAll(c.backupDir)
	// file.MkdirAll(c.keyDir)

	return &c
}

//printConfig Print config information
func (wc *WalletConfig) PrintConfig() error {

	wc.InitConfig()
	//读取配置
	absFile := filepath.Join(wc.configFilePath, wc.configFileName)
	//apiURL := c.String("apiURL")
	//walletPath := c.String("walletPath")
	//threshold := c.String("threshold")
	//minSendAmount := c.String("minSendAmount")
	//minFees := c.String("minFees")
	//sumAddress := c.String("sumAddress")
	//isTestNet, _ := c.Bool("isTestNet")

	fmt.Printf("-----------------------------------------------------------\n")
	//fmt.Printf("Wallet API URL: %s\n", apiURL)
	//fmt.Printf("Wallet Data FilePath: %s\n", walletPath)
	//fmt.Printf("Summary Address: %s\n", sumAddress)
	//fmt.Printf("Summary Threshold: %s\n", threshold)
	//fmt.Printf("Min Send Amount: %s\n", minSendAmount)
	//fmt.Printf("Transfer Fees: %s\n", minFees)
	//if isTestNet {
	//	fmt.Printf("Network: TestNet\n")
	//} else {
	//	fmt.Printf("Network: MainNet\n")
	//}
	file.PrintFile(absFile)
	fmt.Printf("-----------------------------------------------------------\n")

	return nil

}

//创建文件夹
func (wc *WalletConfig) makeDataDir() {

	if len(wc.DataDir) == 0 {
		//默认路径当前文件夹./data
		wc.DataDir = "data"
	}

	//本地数据库文件路径
	wc.dbPath = filepath.Join(wc.DataDir, strings.ToLower(wc.Symbol), "db")

	//创建目录
	file.MkdirAll(wc.dbPath)
}

//initConfig 初始化配置文件
func (wc *WalletConfig) InitConfig() {

	//读取配置
	absFile := filepath.Join(wc.configFilePath, wc.configFileName)
	if !file.Exists(absFile) {
		file.MkdirAll(wc.configFilePath)
		file.WriteFile(absFile, []byte(wc.DefaultConfig), false)
	}

}
