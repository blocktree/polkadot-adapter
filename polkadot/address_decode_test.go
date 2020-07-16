package polkadot

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

func TestAddressDecoder_AddressEncode(t *testing.T) {
	result := "162zn6hjYDi5bCeQ75LisY4v994xpyJW3iNR76Ea8VwWrnYJ"
	p2pk, _ := hex.DecodeString("deb2e8c58b53dbe18cf9b103d06f8e3a83bff78443ba8750ecc1cb9a31532363")
	p2pkAddr, _ := tw.Decoder.AddressEncode(p2pk)
	fmt.Println("p2pkAddr: ", p2pkAddr, strings.EqualFold(result, p2pkAddr) )

	result = "12gjTQn3RAaMLDmDqZphQYaDWPvso36USa7i42zgzbAKadp9"
	p2pk, _ = hex.DecodeString("4a89b8be1072ee1e71ea46430715fbfe8244c76512896980eb2adeed3518aa74")
	p2pkAddr, _ = tw.Decoder.AddressEncode(p2pk)
	fmt.Println("p2pkAddr: ", p2pkAddr, strings.EqualFold(result, p2pkAddr) )
}

func TestAddressDecoder_AddressDecode(t *testing.T) {
	result := "deb2e8c58b53dbe18cf9b103d06f8e3a83bff78443ba8750ecc1cb9a31532363"
	p2pkAddr := "162zn6hjYDi5bCeQ75LisY4v994xpyJW3iNR76Ea8VwWrnYJ"
	p2pk, _ := tw.Decoder.AddressDecode(p2pkAddr)
	p2pkStr := hex.EncodeToString(p2pk)
	fmt.Println("p2pkStr: ", p2pkStr, strings.EqualFold(result, p2pkStr) )

	result = "4a89b8be1072ee1e71ea46430715fbfe8244c76512896980eb2adeed3518aa74"
	p2pkAddr = "12gjTQn3RAaMLDmDqZphQYaDWPvso36USa7i42zgzbAKadp9"
	p2pk, _ = tw.Decoder.AddressDecode(p2pkAddr)
	p2pkStr = hex.EncodeToString(p2pk)
	fmt.Println("p2pkStr: ", p2pkStr, strings.EqualFold(result, p2pkStr) )
}

func TestAddressDecoder(t *testing.T) {
	pubkey := "deb2e8c58b53dbe18cf9b103d06f8e3a83bff78443ba8750ecc1cb9a31532363"
	p2pkAddr := "162zn6hjYDi5bCeQ75LisY4v994xpyJW3iNR76Ea8VwWrnYJ"

	testP2pkAddr, _ := tw.Decoder.AddressEncode( hex.DecodeString(pubkey) )
	fmt.Println("testP2pkAddr: ", testP2pkAddr, strings.EqualFold(testP2pkAddr, p2pkAddr) )

	p2pk, _ := tw.Decoder.AddressDecode(p2pkAddr)
	p2pkStr := hex.EncodeToString(p2pk)
	fmt.Println("p2pkStr: ", p2pkStr, strings.EqualFold(pubkey, p2pkStr) )
}

func Test_ed25519_AddressVerify_Valid(t *testing.T) {
	addressArr := make([]string, 0)
	addressArr = append(addressArr, "162zn6hjYDi5bCeQ75LisY4v994xpyJW3iNR76Ea8VwWrnYJ")	//正确
	addressArr = append(addressArr, "162zn6hjYDi5bCeQ75LisY4v994xpyJW3iNR76Ea8VwWrnYa")	//改了最后一位
	addressArr = append(addressArr, "462zn6hjYDi5bCeQ75LisY4v994xpyJW3iNR76Ea8VwWrnYJ")	//改了第一位
	addressArr = append(addressArr, "162zn6hjYDi5bCeQ75LisY4v994xpybW3iNR76Ea8VwWrnYJ")	//改了中间

	for i := 0; i < len(addressArr); i++ {
		address := addressArr[i]
		valid := tw.Decoder.AddressVerify(address)

		fmt.Println(address, " isvalid : ", valid)
	}
}

func Test_sr25519_AddressVerify_Valid(t *testing.T) {
	addressArr := make([]string, 0)
	addressArr = append(addressArr, "14iKEitKAGicKHJprSmFjfveM6FkoVNoSEiL1bC2tykLUGac")	//正确
	addressArr = append(addressArr, "14iKEitKAGicKHJprSmFjfveM6FkoVNoSEiL1bC2tykLUGa6")	//改了最后一位
	addressArr = append(addressArr, "h4iKEitKAGicKHJprSmFjfveM6FkoVNoSEiL1bC2tykLUGac")	//改了第一位
	addressArr = append(addressArr, "14iKEitKAGicKHJprSmFjfveM6FkoVNoSdiL1bC2tykLUGac")	//改了中间

	for i := 0; i < len(addressArr); i++ {
		address := addressArr[i]
		valid := tw.Decoder.AddressVerify(address)

		fmt.Println(address, " isvalid : ", valid)
	}
}