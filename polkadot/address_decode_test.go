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

func Test_ed25519_AddressVerify_Valid(t *testing.T) {
	addressArr := make([]string, 0)
	addressArr = append(addressArr, "HcKJ5nYJoTXuKTKv96mdLbmS7MYwLZYRbUgLTXB4D8VRFLj")	//正确
	addressArr = append(addressArr, "HcKJ5nYJoTXuKTKv96mdLbmS7MYwLZYRbUgLTXB4D8VRFLi")	//改了最后一位
	addressArr = append(addressArr, "AcKJ5nYJoTXuKTKv96mdLbmS7MYwLZYRbUgLTXB4D8VRFLj")	//改了第一位
	addressArr = append(addressArr, "HcKJ5nYJoTXuKTKv96mdLbmS7MYwLZ26SUgLTXB4D8VRFLj")	//改了中间

	for i := 0; i < len(addressArr); i++ {
		address := addressArr[i]
		valid := tw.Decoder.AddressVerify(address)

		fmt.Println(address, " isvalid : ", valid)
	}
}

func Test_sr25519_AddressVerify_Valid(t *testing.T) {
	addressArr := make([]string, 0)
	addressArr = append(addressArr, "CreEvbSZbf7bJ6ymX6UKNWmwBxmyXK6j29QWQGVvtr7TLsn")	//正确
	addressArr = append(addressArr, "CreEvbSZbf7bJ6ymX6UKNWmwBxmyXK6j29QWQGVvtr7TLs7")	//改了最后一位
	addressArr = append(addressArr, "WreEvbSZbf7bJ6ymX6UKNWmwBxmyXK6j29QWQGVvtr7TLsn")	//改了第一位
	addressArr = append(addressArr, "WreEvbSZbf7bJ6ymX6UKNWm1516yXK6j29QWQGVvtr7TLsn")	//改了中间

	for i := 0; i < len(addressArr); i++ {
		address := addressArr[i]
		valid := tw.Decoder.AddressVerify(address)

		fmt.Println(address, " isvalid : ", valid)
	}
}