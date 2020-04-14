package polkadot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

const (
	testNodeAPI = "http://127.0.0.1:8080"
)

func PrintJsonLog(t *testing.T, logCont string){
	if strings.HasPrefix(logCont, "{") {
		var str bytes.Buffer
		_ = json.Indent(&str, []byte(logCont), "", "    ")
		t.Logf("Get Call Result return: \n\t%+v\n", str.String())
	}else{
		t.Logf("Get Call Result return: \n\t%+v\n", logCont)
	}
}

func TestGetCall(t *testing.T) {
	tw := NewClient(testNodeAPI, true)

	if r, err := tw.GetCall("/metadata/" ); err != nil {
		t.Errorf("Get Call Result failed: %v\n", err)
	} else {
		PrintJsonLog(t, r.String())
	}
}

func Test_getBlockHeight(t *testing.T) {

	c := NewClient(testNodeAPI, true)

	r, err := c.getBlockHeight()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("height:", r)
	}

}

func Test_getBlockByHeight(t *testing.T) {
	c := NewClient(testNodeAPI, true)
	r, err := c.getBlockByHeight(1830393)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}

func Test_getBalance(t *testing.T) {

	c := NewClient(testNodeAPI, true)

	address := "CyVGZAGjfD9bbQcN7Ja7cFrazM4yAzawRhKcWYN1ZmMrpwf"

	r, err := c.getBalance(address, true, 20000000)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}

	address = "J5cTdWxoMRZyQHyvvDMxB6dp7YitNpzEkj3ZrJFsGmARcC2"

	r, err = c.getBalance(address, true, 20000000)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}

func Test_sendTransaction(t *testing.T) {
	c := NewClient(testNodeAPI, true)
	r, err := c.sendTransaction("0x390284f38450ff27a19377c88aa42031c2e4b658dd2df9e839f620f21d97de4848705100ebf250f774085bf37d3b52067fafa5c3a3362b84e6274592cce6b89d3a73705f2c8478028a11007cb8b183471fb4aa43a1f56136b0c7efcb191c7e03b563e60585000800040097dbc27785f579866e9f11067dd3659b45cfd95a9c32a6a3a3e0ea43a50e173f0700e8764817")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}