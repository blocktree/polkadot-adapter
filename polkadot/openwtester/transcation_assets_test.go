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
	"testing"
	"time"

	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openw"
	"github.com/blocktree/openwallet/v2/openwallet"
)

func testGetAssetsAccountBalance(tm *openw.WalletManager, walletID, accountID string) {
	balance, err := tm.GetAssetsAccountBalance(testApp, walletID, accountID)
	if err != nil {
		log.Error("GetAssetsAccountBalance failed, unexpected error:", err)
		return
	}
	log.Info("balance:", balance)
}

func testGetAssetsAccountTokenBalance(tm *openw.WalletManager, walletID, accountID string, contract openwallet.SmartContract) {
	balance, err := tm.GetAssetsAccountTokenBalance(testApp, walletID, accountID, contract)
	if err != nil {
		log.Error("GetAssetsAccountTokenBalance failed, unexpected error:", err)
		return
	}
	log.Info("token balance:", balance.Balance)
}

func testCreateTransactionStep(tm *openw.WalletManager, walletID, accountID, to, amount, feeRate string, contract *openwallet.SmartContract) (*openwallet.RawTransaction, error) {

	//err := tm.RefreshAssetsAccountBalance(testApp, accountID)
	//if err != nil {
	//	log.Error("RefreshAssetsAccountBalance failed, unexpected error:", err)
	//	return nil, err
	//}

	rawTx, err := tm.CreateTransaction(testApp, walletID, accountID, amount, to, feeRate, "test", contract)

	if err != nil {
		log.Error("CreateTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTx, nil
}

func testCreateSummaryTransactionStep(
	tm *openw.WalletManager,
	walletID, accountID, summaryAddress, minTransfer, retainedBalance, feeRate string,
	start, limit int,
	contract *openwallet.SmartContract) ([]*openwallet.RawTransactionWithError, error) {

	rawTxArray, err := tm.CreateSummaryRawTransactionWithError(testApp, walletID, accountID, summaryAddress, minTransfer,
		retainedBalance, feeRate, start, limit, contract, nil)

	if err != nil {
		log.Error("CreateSummaryTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTxArray, nil
}

func testSignTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	_, err := tm.SignTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, "12345678", rawTx)
	if err != nil {
		log.Error("SignTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Infof("rawTx: %+v", rawTx)
	return rawTx, nil
}

func testVerifyTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	//log.Info("rawTx.Signatures:", rawTx.Signatures)

	_, err := tm.VerifyTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("VerifyTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Infof("rawTx: %+v", rawTx)
	return rawTx, nil
}

func testSubmitTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	tx, err := tm.SubmitTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("SubmitTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Std.Info("tx: %+v", tx)
	log.Info("wxID:", tx.WxID)
	log.Info("txID:", rawTx.TxID)

	return rawTx, nil
}

func ClearAddressNonce(tm *openw.WalletManager, walletID string, accountID string) error{
	wrapper, err := tm.NewWalletWrapper(testApp, "")
	if err != nil {
		return err
	}

	list, err := tm.GetAddressList(testApp, walletID, accountID, 0, -1, false)
	if err != nil {
		log.Error("unexpected error:", err)
		return err
	}
	for i, w := range list {
		log.Info("address[", i, "] :", w.Address)

		key := "DOT-nonce"
		wrapper.SetAddressExtParam(w.Address, key, 0)
	}
	log.Info("address count:", len(list))

	tm.CloseDB(testApp)

	return nil
}

/*
withdraw
wallet : W5tGn1AKvsgB4gMW3fj6YRebravgjyKYGs
account : EC2qLETifz9PVQVLhHTsojCq4DxJ1FtqxvcwhvJetG7
1 address : CdxHHi8EZqkqnnqpqRKghieaeeDHm9mjtQKc3y7LZLGvQWz
*/
func TestTransfer(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "W5tGn1AKvsgB4gMW3fj6YRebravgjyKYGs"
	accountID := "EC2qLETifz9PVQVLhHTsojCq4DxJ1FtqxvcwhvJetG7"
	to := "CyVGZAGjfD9bbQcN7Ja7cFrazM4yAzawRhKcWYN1ZmMrpwf"

	testGetAssetsAccountBalance(tm, walletID, accountID)

	rawTx, err := testCreateTransactionStep(tm, walletID, accountID, to, "0.744508", "", nil)
	if err != nil {
		return
	}

	log.Std.Info("rawTx: %+v", rawTx)

	_, err = testSignTransactionStep(tm, rawTx)
	if err != nil {
		return
	}

	_, err = testVerifyTransactionStep(tm, rawTx)
	if err != nil {
		return
	}

	_, err = testSubmitTransactionStep(tm, rawTx)
	if err != nil {
		return
	}
}

/*
withdraw
wallet : W5tGn1AKvsgB4gMW3fj6YRebravgjyKYGs
account : EC2qLETifz9PVQVLhHTsojCq4DxJ1FtqxvcwhvJetG7
1 address : CdxHHi8EZqkqnnqpqRKghieaeeDHm9mjtQKc3y7LZLGvQWz
---------------
charge
wallet : VyfWKDkCqdKRj3AzeRoJ2F5Rqe8VYXnK8w
account : 5ZBY5578VF8tjwMT18rGBAvMK6vVWNefaoNUGnttLyQk
1 address : D3vsyPWgGWRYgmnph8JzLYjvbWQ9hHx8jbYKH2n21wbHxcA
2 address : DYP2dDKrxdgyAVhwN9cgdVyyrUR1234We8qMahrj5SWNFy5
3 address : DqgFAywY2jqKtKgntKnN1jNpGV6TxnkXgpUtRYo1LCjsai2
4 address : EXRio1TXjLjr7fz9LmdsjPDME5y5eQz2ffJotbcXwKyhXbC
5 address : HUwSCNJH2YjZuwoNbWmsr86C1ViPs6MSQkKUpMkEdmh98hb
*/
func TestBatchTransfer(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "W5tGn1AKvsgB4gMW3fj6YRebravgjyKYGs"
	accountID := "EC2qLETifz9PVQVLhHTsojCq4DxJ1FtqxvcwhvJetG7"
	toArr := make([]string, 0)
	toArr = append(toArr, "D3vsyPWgGWRYgmnph8JzLYjvbWQ9hHx8jbYKH2n21wbHxcA")
	toArr = append(toArr, "DYP2dDKrxdgyAVhwN9cgdVyyrUR1234We8qMahrj5SWNFy5")
	toArr = append(toArr, "DqgFAywY2jqKtKgntKnN1jNpGV6TxnkXgpUtRYo1LCjsai2")
	toArr = append(toArr, "EXRio1TXjLjr7fz9LmdsjPDME5y5eQz2ffJotbcXwKyhXbC")
	toArr = append(toArr, "HUwSCNJH2YjZuwoNbWmsr86C1ViPs6MSQkKUpMkEdmh98hb")

	for i := 0; i < len(toArr); i++{
		to := toArr[i]
		testGetAssetsAccountBalance(tm, walletID, accountID)

		rawTx, err := testCreateTransactionStep(tm, walletID, accountID, to, "0.031", "", nil)
		if err != nil {
			return
		}

		log.Std.Info("rawTx: %+v", rawTx)

		_, err = testSignTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		time.Sleep(time.Duration(5) * time.Second)
	}
}

func TestSummary(t *testing.T) {
	tm := testInitWalletManager()

	walletID := "VyfWKDkCqdKRj3AzeRoJ2F5Rqe8VYXnK8w"
	accountID := "5ZBY5578VF8tjwMT18rGBAvMK6vVWNefaoNUGnttLyQk"
	summaryAddress := "CyVGZAGjfD9bbQcN7Ja7cFrazM4yAzawRhKcWYN1ZmMrpwf"

	ClearAddressNonce(tm, walletID, accountID)

	testGetAssetsAccountBalance(tm, walletID, accountID)

	rawTxArray, err := testCreateSummaryTransactionStep(tm, walletID, accountID,
		summaryAddress, "0.01", "0.01", "",
		0, 100, nil)
	if err != nil {
		log.Errorf("CreateSummaryTransaction failed, unexpected error: %v", err)
		return
	}

	//执行汇总交易
	for _, rawTxWithErr := range rawTxArray {

		if rawTxWithErr.Error != nil {
			log.Error(rawTxWithErr.Error.Error())
			continue
		}

		_, err = testSignTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}
	}

}
