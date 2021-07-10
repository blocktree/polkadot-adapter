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
	"github.com/blocktree/openwallet/v2/common"
	"github.com/ethereum/go-ethereum/common/math"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/blocktree/openwallet/v2/hdkeystore"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
)

type WalletManager struct {
	openwallet.AssetsAdapterBase

	Storage   *hdkeystore.HDKeystore //秘钥存取
	ApiClient *ApiClient

	Config          *WalletConfig                 //钱包管理配置
	WalletsInSum    map[string]*openwallet.Wallet //参与汇总的钱包
	Blockscanner    *DOTBlockScanner              //区块扫描器
	Decoder         openwallet.AddressDecoderV2   //地址编码器
	TxDecoder       openwallet.TransactionDecoder //交易单编码器
	Log             *log.OWLogger                 //日志工具
	ContractDecoder *ContractDecoder              //智能合约解析器
}

func NewWalletManager() *WalletManager {
	wm := WalletManager{}
	wm.Config = NewConfig(Symbol, MasterKey, AddrPrefix)
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
		//nonce = common.NewString(nonce_db).UInt64()
		nonceStr := common.NewString(nonce_db)
		nonceStrArr := strings.Split(string(nonceStr), "_")

		if len(nonceStrArr)>1 {
			saveTime := common.NewString(nonceStrArr[1]).UInt64()
			now := uint64(time.Now().Unix())
			diff, _ := math.SafeSub(now, saveTime)

			if diff > wm.Config.NonceDiff { //当前时间减去保存时间，超过1小时，就不算了
				nonce = 0
			} else {
				nonce = common.NewString(nonceStrArr[0]).UInt64()
			}
		}else{
			nonce = common.NewString(nonce_db).UInt64()
		}
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

	////临时
	//if nonce > 430 && nonce < 450 {
	//	nonce = nonce_onchain
	//}

	return nonce
}

// UpdateAddressNonce
func (wm *WalletManager) UpdateAddressNonce(wrapper openwallet.WalletDAI, address string, nonce uint64) {
	key := wm.Symbol() + "-nonce"

	nonceStr := strconv.FormatUint(nonce, 10) + "_" + strconv.FormatUint( uint64(time.Now().Unix()), 10)
	wm.Log.Info(address, " set nonce ", nonceStr)

	err := wrapper.SetAddressExtParam(address, key, nonceStr)
	if err != nil {
		wm.Log.Errorf("WalletDAI SetAddressExtParam failed, err: %v", err)
	}
}
