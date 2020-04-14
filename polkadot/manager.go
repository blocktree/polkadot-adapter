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
	"github.com/astaxie/beego/config"
	"path/filepath"

	"github.com/blocktree/openwallet/v2/common"
	"github.com/blocktree/openwallet/v2/hdkeystore"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
)

type WalletManager struct {
	openwallet.AssetsAdapterBase

	Storage *hdkeystore.HDKeystore //秘钥存取
	ApiClient *ApiClient

	Config          *WalletConfig                 //钱包管理配置
	WalletsInSum    map[string]*openwallet.Wallet //参与汇总的钱包
	Blockscanner    *DOTBlockScanner             //区块扫描器
	Decoder         openwallet.AddressDecoderV2     //地址编码器
	TxDecoder       openwallet.TransactionDecoder //交易单编码器
	Log             *log.OWLogger                 //日志工具
	ContractDecoder *ContractDecoder              //智能合约解析器
}

func NewWalletManager() *WalletManager {
	wm := WalletManager{}
	wm.Config = NewConfig(Symbol, MasterKey, GenesisHash, SpecVersion)
	storage := hdkeystore.NewHDKeystore(wm.Config.keyDir, hdkeystore.StandardScryptN, hdkeystore.StandardScryptP)
	wm.Storage = storage
	//参与汇总的钱包
	wm.WalletsInSum = make(map[string]*openwallet.Wallet)
	//区块扫描器
	wm.Blockscanner = NewDOTBlockScanner(&wm)
	wm.Decoder = NewAddressDecoderV2(&wm)
	wm.TxDecoder = NewTransactionDecoder(&wm)
	wm.Log = log.NewOWLogger(wm.Symbol())
	wm.ContractDecoder = NewContractDecoder(&wm)

	//	wm.RPCClient = NewRpcClient("http://localhost:20336/")
	return &wm
}

//GetWalletInfo 获取钱包列表
func (wm *WalletManager) GetWalletInfo(walletID string) (*openwallet.Wallet, error) {

	wallets, err := wm.GetWallets()
	if err != nil {
		return nil, err
	}

	//获取钱包余额
	for _, w := range wallets {
		if w.WalletID == walletID {
			return w, nil
		}

	}

	return nil, errors.New("The wallet that your given name is not exist!")
}

//GetWallets 获取钱包列表
func (wm *WalletManager) GetWallets() ([]*openwallet.Wallet, error) {

	wallets, err := openwallet.GetWalletsByKeyDir(wm.Config.keyDir)
	if err != nil {
		return nil, err
	}

	for _, w := range wallets {
		w.DBFile = filepath.Join(wm.Config.dbPath, w.FileName()+".db")
	}

	return wallets, nil

}

//SendRawTransaction 广播交易
func (wm *WalletManager) SendRawTransaction(txHex string) (string, error) {

	return wm.sendRawTransactionByNode(txHex)
}

func (wm *WalletManager) sendRawTransactionByNode(txHex string) (string, error) {
	//var (
	//	txid string
	//	err error
	//)
	//if wm.Config.APIChoose == "rpc" {
	//	txid, err = wm.Client.sendTransaction(txHex)
	//} else if wm.Config.APIChoose == "ws" {
	//	txid, err = wm.WSClient.sendTransaction(txHex)
	//}else {
	//	return "",errors.New("Invalid config, check the ini file!")
	//}
	txid, err := wm.ApiClient.sendTransaction(txHex)

	if err != nil {
		return "", err
	}
	return txid, nil
}

// GetAddressNonce
func (wm *WalletManager) GetAddressNonce(wrapper openwallet.WalletDAI, account *AddrBalance) uint64 {
	var (
		key           = wm.Symbol() + "-nonce"
		nonce         uint64
		nonce_db      interface{}
		nonce_onchain uint64
	)

	//获取db记录的nonce并确认nonce值
	nonce_db, _ = wrapper.GetAddressExtParam(account.Address, key)

	//判断nonce_db是否为空,为空则说明当前nonce是0
	if nonce_db == nil {
		nonce = 0
	} else {
		nonce = common.NewString(nonce_db).UInt64()
	}

	nonce_onchain = account.Nonce

	wm.Log.Info(account.Address, " get nonce : ", nonce, ", nonce_onchain : ", nonce_onchain)

	//如果本地nonce_db > 链上nonce,采用本地nonce,否则采用链上nonce
	if nonce > nonce_onchain {
		//wm.Log.Debugf("%s nonce_db=%v > nonce_chain=%v,Use nonce_db...", address, nonce_db, nonce_onchain)
	} else {
		nonce = nonce_onchain
		//wm.Log.Debugf("%s nonce_db=%v <= nonce_chain=%v,Use nonce_chain...", address, nonce_db, nonce_onchain)
	}

	return nonce
}

// UpdateAddressNonce
func (wm *WalletManager) UpdateAddressNonce(wrapper openwallet.WalletDAI, address string, nonce uint64) {
	key := wm.Symbol() + "-nonce"
	wm.Log.Info(address, " set nonce ", nonce)
	err := wrapper.SetAddressExtParam(address, key, nonce)
	if err != nil {
		wm.Log.Errorf("WalletDAI SetAddressExtParam failed, err: %v", err)
	}
}

//LoadAssetsConfig 加载外部配置
func (wm *WalletManager) LoadAssetsConfig(c config.Configer) error {

	wm.Config.NodeAPI = c.String("nodeAPI")
	wm.Config.WSAPI = c.String("wsAPI")
	wm.Config.APIChoose = c.String("apiChoose")
	//if wm.Config.APIChoose == "rpc" {
	//	wm.Client = NewClient(wm.Config.NodeAPI, false)
	//}else if wm.Config.APIChoose == "ws" {
	//	wm.WSClient = NewWSClient(wm, wm.Config.WSAPI, 0, false)
	//}
	NewApiClient(wm)

	wm.Config.FixedFee, _ = c.Int64("fixedFee")
	wm.Config.ReserveAmount, _ = c.Int64("reserveAmount")
	wm.Config.IgnoreReserve, _ = c.Bool("ignoreReserve")
	wm.Config.LastLedgerSequenceNumber, _ = c.Int64("lastLedgerSequenceNumber")
	wm.Config.DataDir = c.String("dataDir")

	wm.Config.MemoType = c.String("memoType")
	wm.Config.MemoFormat = c.String("memoFormat")
	wm.Config.MemoScan = c.String("memoScan")
	//数据文件夹
	wm.Config.makeDataDir()

	return nil
}
