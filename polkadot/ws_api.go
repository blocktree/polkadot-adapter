package polkadot

import (
	"errors"
	"math/big"
	"time"
)

//func (ws *WSClient) Call(request map[string]interface{})(*gjson.Result, error) {
//    //flag.Parse()
//    u := url.URL{Scheme: "ws", Host: *ws.addr, Path: ""}
//    c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
//
//    defer c.Close()
//
//    done := make(chan struct{})
//
//    defer close(done)
//    tmp,_ := json.Marshal(request)
//    err = c.WriteMessage(2, tmp)
//    if err != nil {
//        if ws.debug {
//            log.Std.Info("write:", err)
//        }
//    }
//
//    _, message, err := c.ReadMessage()
//    if err != nil {
//        if ws.debug{
//            log.Std.Info("read:", err)
//        }
//        return nil, err
//    }
//
//    if ws.debug {
//        log.Std.Info("recv: %s", message)
//    }
//    ret := gjson.Parse(string(message))
//    return &ret, nil
//}

func (c *WSClient) getBlockHeight() (uint64, error) {

	request := map[string]interface{}{
		//"id" : 14,
		//"command":"ledger",
		"ledger_index": "validated",
		"accounts":     false,
		"full":         false,
		"transactions": false,
		"expand":       false,
		"owner_funds":  false,
	}

	resp, err := c.Call("ledger", request)

	if err != nil {
		return 0, err
	}

	if resp.Get("error").String() != "" {
		return 0, errors.New(resp.Get("error").String())
	}

	return resp.Get("result").Get("ledger").Get("ledger_index").Uint(), nil
}

func (c *WSClient) getBlockHash(height uint64) (string, error) {
	request := map[string]interface{}{
		//"id" : 14,
		//"command":"ledger",
		"ledger_index": height,
		"accounts":     false,
		"full":         false,
		"transactions": false,
		"expand":       false,
		"owner_funds":  false,
	}

	resp, err := c.Call("ledger", request)
	if err != nil {
		return "", err
	}

	if resp.Get("error").String() != "" {
		return "", errors.New(resp.Get("error").String())
	}

	return resp.Get("result").Get("ledger").Get("hash").String(), nil
}

func (c *WSClient) getSequence(address string) (uint32, error) {
	request := map[string]interface{}{
		//"id":           14,
		//"command":      "account_info",
		"account":      address,
		"strict":       true,
		"ledger_index": "current",
		"queue":        true,
	}
	resp, err := c.Call("account_info", request)
	if err != nil {
		return 0, err
	}

	if resp.Get("error").String() == "actNotFound" {
		return 0, nil
	}

	if resp.Get("error").String() != "" {
		return 0, errors.New(resp.Get("error").String())
	}

	return uint32(resp.Get("result").Get("account_data").Get("Sequence").Int()), nil
}

func (c *WSClient) getBalance(address string, ignoreReserve bool, reserveAmount int64) (*AddrBalance, error) {
	request := map[string]interface{}{
		//"id":           14,
		//"command":      "account_info",
		"account":      address,
		"strict":       true,
		"ledger_index": "current",
		"queue":        true,
	}

	r, err := c.Call("account_info", request)

	if err != nil {
		return nil, err
	}

	if r.Get("error").String() == "actNotFound" {
		return &AddrBalance{Address: address, Balance: big.NewInt(0), Actived: false}, nil
	}

	if r.Get("error").String() != "" {
		return nil, errors.New(r.Get("error").String())
	}

	if ignoreReserve {
		return &AddrBalance{Address: address, Balance: big.NewInt(r.Get("result").Get("account_data").Get("Balance").Int() - reserveAmount), Actived: true}, nil
	}

	return &AddrBalance{Address: address, Balance: big.NewInt(r.Get("result").Get("account_data").Get("Balance").Int()), Actived: true}, nil
}

func (c *WSClient) isActived(address string) (bool, error) {
	request := map[string]interface{}{
		//"id":           14,
		//"command":      "account_info",
		"account":      address,
		"strict":       true,
		"ledger_index": "current",
		"queue":        true,
	}

	r, err := c.Call("account_info", request)

	if err != nil {
		return false, err
	}

	if r.Get("error").String() == "actNotFound" {
		return false, nil
	}

	if r.Get("error").String() != "" {
		return false, errors.New(r.Get("error").String())
	}

	return true, nil
}

// 获取区块信息
func (c *WSClient) getBlock(hash string) (*Block, error) {
	request := map[string]interface{}{
		//"id":           14,
		//"command":      "ledger",
		"ledger_index": "validated",
		"accounts":     false,
		"full":         false,
		"transactions": true,
		"expand":       false,
		"owner_funds":  false,
	}
	resp, err := c.Call("ledger", request)

	if err != nil {
		return nil, err
	}

	if resp.Get("error").String() != "" {
		return nil, errors.New(resp.Get("error").String())
	}

	block := resp.Get("result")
	return NewBlock(&block), nil
}

func (c *WSClient) getBlockByHeight(height uint64) (*Block, error) {
	request := map[string]interface{}{
		//"id":           14,
		//"command":      "ledger",
		"ledger_index": height,
		"accounts":     false,
		"full":         false,
		"transactions": true,
		"expand":       false,
		"owner_funds":  false,
	}
	resp, err := c.Call("ledger", request)

	if err != nil {
		return nil, err
	}

	if resp.Get("error").String() != "" {
		return nil, errors.New(resp.Get("error").String())
	}

	block := resp.Get("result")

	return NewBlock(&block), nil
}

func (c *WSClient) sendTransaction(rawTx string) (string, error) {
	request := map[string]interface{}{
		//"id":      14,
		//"command": "submit",
		"tx_blob": rawTx,
	}

	resp, err := c.Call("submit", request)
	if err != nil {
		return "", err
	}

	time.Sleep(time.Duration(1) * time.Second)

	if resp.Get("result").Get("engine_result").String() != "tesSUCCESS" && resp.Get("result").Get("engine_result").String() != "terQUEUED" {
		return "", errors.New("Submit transaction with error: " + resp.String())//resp.Get("result").Get("engine_result_message").String())
	}

	return resp.Get("result").Get("tx_json").Get("hash").String(), nil
}
