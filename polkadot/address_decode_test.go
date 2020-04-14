package polkadot

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

func TestAddressDecoder_AddressEncode(t *testing.T) {
	result := "HcKJ5nYJoTXuKTKv96mdLbmS7MYwLZYRbUgLTXB4D8VRFLj"
	p2pk, _ := hex.DecodeString("deb2e8c58b53dbe18cf9b103d06f8e3a83bff78443ba8750ecc1cb9a31532363")
	p2pkAddr, _ := tw.Decoder.AddressEncode(p2pk)
	fmt.Println("p2pkAddr: ", p2pkAddr, strings.EqualFold(result, p2pkAddr) )

	result = "EG3yPrrBkKoeLa9edakAM74oNDTuQMWpTDyHQHHvJMJ9ZxK"
	p2pk, _ = hex.DecodeString("4a89b8be1072ee1e71ea46430715fbfe8244c76512896980eb2adeed3518aa74")
	p2pkAddr, _ = tw.Decoder.AddressEncode(p2pk)
	fmt.Println("p2pkAddr: ", p2pkAddr, strings.EqualFold(result, p2pkAddr) )
}

func TestAddressDecoder_AddressDecode(t *testing.T) {
	result := "deb2e8c58b53dbe18cf9b103d06f8e3a83bff78443ba8750ecc1cb9a31532363"
	p2pkAddr := "HcKJ5nYJoTXuKTKv96mdLbmS7MYwLZYRbUgLTXB4D8VRFLj"
	p2pk, _ := tw.Decoder.AddressDecode(p2pkAddr)
	p2pkStr := hex.EncodeToString(p2pk)
	fmt.Println("p2pkStr: ", p2pkStr, strings.EqualFold(result, p2pkStr) )

	result = "4a89b8be1072ee1e71ea46430715fbfe8244c76512896980eb2adeed3518aa74"
	p2pkAddr = "EG3yPrrBkKoeLa9edakAM74oNDTuQMWpTDyHQHHvJMJ9ZxK"
	p2pk, _ = tw.Decoder.AddressDecode(p2pkAddr)
	p2pkStr = hex.EncodeToString(p2pk)
	fmt.Println("p2pkStr: ", p2pkStr, strings.EqualFold(result, p2pkStr) )
}

func TestAddressDecoder(t *testing.T) {
	pubkey := "deb2e8c58b53dbe18cf9b103d06f8e3a83bff78443ba8750ecc1cb9a31532363"
	p2pkAddr := "HcKJ5nYJoTXuKTKv96mdLbmS7MYwLZYRbUgLTXB4D8VRFLj"

	testP2pkAddr, _ := tw.Decoder.AddressEncode( hex.DecodeString(pubkey) )
	fmt.Println("testP2pkAddr: ", testP2pkAddr, strings.EqualFold(testP2pkAddr, p2pkAddr) )

	p2pk, _ := tw.Decoder.AddressDecode(p2pkAddr)
	p2pkStr := hex.EncodeToString(p2pk)
	fmt.Println("p2pkStr: ", p2pkStr, strings.EqualFold(pubkey, p2pkStr) )
}