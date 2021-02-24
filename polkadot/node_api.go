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
	"github.com/blocktree/openwallet/v2/log"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type ClientInterface interface {
	Call(path string, request []interface{}) (*gjson.Result, error)
}

// A Client is a Elastos RPC client. It performs RPCs over HTTP using JSON
// request and responses. A Client must be configured with a secret token
// to authenticate with other Cores on the network.
type Client struct {
	BaseURL     string
	AccessToken string
	Debug       bool
	client      *req.Req
	Symbol      string
}

type Response struct {
	Code    int         `json:"code,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Message string      `json:"message,omitempty"`
	Id      string      `json:"id,omitempty"`
}

func NewClient(url string /*token string,*/, debug bool, symbol string) *Client {
	c := Client{
		BaseURL: url,
		//	AccessToken: token,
		Debug: debug,
	}

	log.Debug("BaseURL : ", url)

	api := req.New()
	//trans, _ := api.Client().Transport.(*http.Transport)
	//trans.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c.client = api
	c.Symbol = symbol

	return &c
}

// 用get方法获取内容
func (c *Client) PostCall(path string, v map[string]interface{}) (*gjson.Result, error) {
	if c.Debug {
		log.Debug("Start Request API...")
	}

	r, err := req.Post(c.BaseURL+path, req.BodyJSON(&v))

	if c.Debug {
		log.Std.Info("Request API Completed")
	}

	if c.Debug {
		log.Debugf("%+v\n", r)
	}

	if err != nil {
		return nil, err
	}

	resp := gjson.ParseBytes(r.Bytes())

	result := resp

	return &result, nil
}

// 用get方法获取内容
func (c *Client) GetCall(path string) (*gjson.Result, error) {

	if c.Debug {
		log.Debug("Start Request API...")
	}

	r, err := req.Get(c.BaseURL + path)

	if c.Debug {
		log.Std.Info("Request API Completed")
	}

	if c.Debug {
		log.Debugf("%+v\n", r)
	}

	if err != nil {
		return nil, err
	}

	resp := gjson.ParseBytes(r.Bytes())

	result := resp

	return &result, nil
}

// 获取当前最高区块
func (c *Client) getBlockHeight() (uint64, error) {
	mostHeightBlock, err := c.getMostHeightBlock()
	if err != nil {
		return 0, err
	}
	return mostHeightBlock.Height, nil
}

// 获取当前最高区块
func (c *Client) getTxArtifacts() (*TxArtifacts, error) {
	//resp, err := c.GetCall("/tx/artifacts/")
	resp, err := c.GetCall("/transaction/material")

	if err != nil {
		return nil, err
	}
	return GetTxArtifacts(resp), nil
}

//获取当前最新高度
func (c *Client) getMostHeightBlock() (*Block, error) {
	resp, err := c.GetCall("/blocks/head")

	if err != nil {
		return nil, err
	}
	return NewBlock(resp, c.Symbol), nil
}

// 获取地址余额
func (c *Client) getBalance(address string) (*AddrBalance, error) {
	r, err := c.GetCall("/accounts/"+address+"/balance-info")

	if err != nil {
		return nil, err
	}

	if r.Get("error").String() == "actNotFound" {
		return &AddrBalance{Address: address, Balance: big.NewInt(0), Actived: false, Nonce: uint64(0)}, nil
	}

	free := big.NewInt(r.Get("free").Int())
	feeFrozen := big.NewInt(r.Get("feeFrozen").Int())
	nonce := uint64(r.Get("nonce").Uint())
	balance := new(big.Int)
	balance = balance.Sub(free, feeFrozen)
	return &AddrBalance{Address: address, Balance: balance, Freeze: feeFrozen, Free: free, Actived: true, Nonce: nonce}, nil
}

func (c *Client) getBlockByHeight(height uint64) (*Block, error) {
	resp, err := c.GetCall("/blocks/" + strconv.FormatUint(height, 10))

	if err != nil {
		return nil, err
	}
	return NewBlock(resp, c.Symbol), nil
}

func (c *Client) sendTransaction(rawTx string) (string, error) {
	body := map[string]interface{}{
		"tx": rawTx,
	}

	//resp, err := c.PostCall("/tx", body)
	resp, err := c.PostCall("/transaction", body)
	if err != nil {
		return "", err
	}

	time.Sleep(time.Duration(1) * time.Second)

	log.Debug("sendTransaction result : ", resp)

	if resp.Get("error").String() != "" && resp.Get("cause").String() != "" {
		return "", errors.New("Submit transaction with error: " + resp.Get("error").String() + "," + resp.Get("cause").String())
	}

	return resp.Get("hash").String(), nil
}

func RemoveOxToAddress(addr string) string {
	if strings.Index(addr, "0x") == 0 {
		return addr[2:]
	}
	return addr
}
