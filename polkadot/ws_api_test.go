package polkadot

import (
	"fmt"
	"testing"
)

//var origin = "http://"
var wsurl = ":"

//
//func Test_ws(t *testing.T){
//
//	ws, err := websocket.Dial(wsurl, "", origin)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	defer func() {
//		err = ws.Close()
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Printf("WebSocket close\n")
//	}()
//	request := []interface{}{
//		map[string]interface{}{
//			"id" : 14,
//			"command":"ledger",
//			"ledger_index": "validated",
//			"accounts":     false,
//			"full":         false,
//			"transactions": false,
//			"expand":       false,
//			"owner_funds":  false,
//		},
//	}
//
//
//	//message := []byte("{\"id\": 14,\"command\": \"ledger\",\"ledger_index\": \"validated\",\"full\": false,\"accounts\": false,\"transactions\": false,\"expand\": false,\"owner_funds\": false}")
//message,_ := json.Marshal(request[0])
//	fmt.Println(string(message))
//	n, err := ws.Write(message)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Send(%d): %s\n", n, message)
//
//	var msg = make([]byte, 1024)
//	n, err = ws.Read(msg)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Receive(%d): %s\n", n, msg)
//}


func Test_ws_getBlockHeight(t *testing.T){

	height, err := tw.ApiClient.WSClient.getBlockHeight()
	if err != nil {
		t.Error(err)
	}

	fmt.Println("height : ", height)
}

func Test_ws_gwtBlockHash(t *testing.T) {
	//c := NewWSClient(wsurl, 0, true)
	height := uint64(48551264)
	hash, err := tw.ApiClient.WSClient.getBlockHash(height)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("hash : ", hash)
}

func Test_ws_getSequence(t *testing.T){
	c := tw.ApiClient.WSClient
	addr := "rMzax7NdBeVe5dqwo87VQepccSh9AWyP1m"

	sequence, err := c.getSequence(addr)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("sequence : ", sequence)

	addr = "rUTEn2jLLv4ESmrUqQmhZfEfDN3LorhgvZ"

	sequence, err = c.getSequence(addr)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("sequence : ", sequence)
}

func Test_ws_getBalance(t *testing.T){
	c := tw.ApiClient.WSClient
	addr := "rMzax7NdBeVe5dqwo87VQepccSh9AWyP1m"

	balance, err := c.getBalance(addr, true, 20000000)

	if err != nil {
		t.Error(err)
	}else{
		fmt.Println("balance : ", balance)
	}


	addr = "rUTEn2jLLv4ESmrUqQmhZfEfDN3LorhgvZ"
	balance, err = c.getBalance(addr, true, 20000000)

	if err != nil {
		t.Error(err)
	}else{
		fmt.Println("balance : ", balance)
	}

	addr = "rUTEn2jLLv4ESmrUqQmhZfEfDN3LorhgvZ"
	balance, err = c.getBalance(addr, true, 20000000)

	if err != nil {
		t.Error(err)
	}else{
		fmt.Println("balance : ", balance)
	}

	addr = "rUTEn2jLLv4ESmrUqQmhZfEfDN3LorhgvZ"
	balance, err = c.getBalance(addr, true, 20000000)

	if err != nil {
		t.Error(err)
	}else{
		fmt.Println("balance : ", balance)
	}

	addr = "rUTEn2jLLv4ESmrUqQmhZfEfDN3LorhgvZ"
	balance, err = c.getBalance(addr, true, 20000000)

	if err != nil {
		t.Error(err)
	}else{
		fmt.Println("balance : ", balance)
	}
}


func Test_ws_isActived(t *testing.T){
	c := tw.ApiClient.WSClient
	addr := "rMzax7NdBeVe5dqwo87VQepccSh9AWyP1m"

	isActived, err := c.isActived(addr)

	if err != nil {
		t.Error(err)
	}else{
		fmt.Println("isActived : ", isActived)
	}


	addr = "rUTEn2jLLv4ESmrUqQmhZfEfDN3LorhgvZ"
	isActived, err = c.isActived(addr)

	if err != nil {
		t.Error(err)
	}else{
		fmt.Println("isActived : ", isActived)
	}
}

func Test_ws_getBlockByHeight(t *testing.T) {
	c := tw.ApiClient.WSClient
	r, err := c.getBlockByHeight(48554232)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}



