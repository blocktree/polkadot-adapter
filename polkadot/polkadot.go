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
	"errors"
	"fmt"
	"path/filepath"

	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
)

//初始化配置流程
func (wm *WalletManager) InitConfigFlow() error {

	wm.Config.InitConfig()
	file := filepath.Join(wm.Config.configFilePath, wm.Config.configFileName)
	fmt.Printf("You can run 'vim %s' to edit wallet's Config.\n", file)
	return nil
}

//查看配置信息
func (wm *WalletManager) ShowConfig() error {
	return wm.Config.PrintConfig()
}

//InstallNode 安装节点
func (wm *WalletManager) InstallNodeFlow() error {
	return errors.New("Install node is unsupport now. ")
}

//InitNodeConfig 初始化节点配置文件
func (wm *WalletManager) InitNodeConfigFlow() error {
	return errors.New("Install node is unsupport now. ")
}

//StartNodeFlow 开启节点
func (wm *WalletManager) StartNodeFlow() error {

	//return wm.startNode()
	return nil
}

//StopNodeFlow 关闭节点
func (wm *WalletManager) StopNodeFlow() error {

	//return wm.stopNode()
	return nil
}

//RestartNodeFlow 重启节点
func (wm *WalletManager) RestartNodeFlow() error {
	return errors.New("Install node is unsupport now. ")
}

//ShowNodeInfo 显示节点信息
func (wm *WalletManager) ShowNodeInfo() error {
	return errors.New("Install node is unsupport now. ")
}

//SetConfigFlow 初始化配置流程
func (wm *WalletManager) SetConfigFlow(subCmd string) error {
	file := wm.Config.configFilePath + wm.Config.configFileName
	fmt.Printf("You can run 'vim %s' to edit %s Config.\n", file, subCmd)
	return nil
}

//ShowConfigInfo 查看配置信息
func (wm *WalletManager) ShowConfigInfo(subCmd string) error {
	wm.Config.PrintConfig()
	return nil
}

//CurveType 曲线类型
func (wm *WalletManager) CurveType() uint32 {
	return wm.Config.CurveType
}

//FullName 币种全名
func (wm *WalletManager) FullName() string {
	return "Polkadot"
}

//获取地址格式开头
func (wm *WalletManager) AddrPrefix() byte {
	return wm.Config.AddrPrefix
}

//Symbol 币种标识
func (wm *WalletManager) Symbol() string {
	return wm.Config.Symbol
}

//小数位精度
func (wm *WalletManager) Decimal() int32 {
	return wm.Config.Decimal
}

//AddressDecode 地址解析器
func (wm *WalletManager) GetAddressDecode() openwallet.AddressDecoder {
	return wm.Decoder
}

//AddressDecode 地址解析器
func (wm *WalletManager) GetAddressDecoderV2() openwallet.AddressDecoderV2 {
	return wm.Decoder
}

//TransactionDecoder 交易单解析器
func (wm *WalletManager) GetTransactionDecoder() openwallet.TransactionDecoder {
	return wm.TxDecoder
}

//GetBlockScanner 获取区块链
func (wm *WalletManager) GetBlockScanner() openwallet.BlockScanner {

	return wm.Blockscanner
}

//InitAssetsConfig 初始化默认配置
func (wm *WalletManager) InitAssetsConfig() (config.Configer, error) {
	return config.NewConfigData("ini", []byte(wm.Config.DefaultConfig))
}

//GetAssetsLogger 获取资产账户日志工具
func (wm *WalletManager) GetAssetsLogger() *log.OWLogger {
	//return nil
	return wm.Log
}

// GetSmartContractDecoder 获取智能合约解析器
func (wm *WalletManager) GetSmartContractDecoder() openwallet.SmartContractDecoder {
	return wm.ContractDecoder
}

//LoadAssetsConfig 加载外部配置
func (wm *WalletManager) LoadAssetsConfig(c config.Configer) error {

	wm.Config.NodeAPI = c.String("nodeAPI")
	wm.Config.WSAPI = c.String("wsAPI")
	wm.Config.APIChoose = c.String("apiChoose")
	NewApiClient(wm)

	wm.Config.FixedFee, _ = c.Int64("fixedFee")
	wm.Config.ReserveAmount, _ = c.Int64("reserveAmount")
	wm.Config.IgnoreReserve, _ = c.Bool("ignoreReserve")

	decimalInt, err := c.Int("decimal")
	if err!=nil {
		decimalInt = 12
	}
	wm.Config.Decimal = int32(decimalInt)

	wm.Config.DataDir = c.String("dataDir")

	//数据文件夹
	wm.Config.makeDataDir()

	return nil
}

