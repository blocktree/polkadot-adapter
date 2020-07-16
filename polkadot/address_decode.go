package polkadot

import (
	"github.com/blocktree/go-owcdrivers/addressEncoder"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/v2/openwallet"
)

var (
	alphabet = addressEncoder.BTCAlphabet
	ssPrefix = []byte{0x53, 0x53, 0x35, 0x38, 0x50, 0x52, 0x45}
	encodeType = "base58"
)

var (
	Default = AddressDecoderV2{}
)

//AddressDecoderV2
type AddressDecoderV2 struct {
	*openwallet.AddressDecoderV2Base
	wm *WalletManager
}

//NewAddressDecoder 地址解析器
func NewAddressDecoderV2(wm *WalletManager) *AddressDecoderV2 {
	decoder := AddressDecoderV2{}
	decoder.wm = wm
	return &decoder
}

//AddressDecode 地址解析
func (dec *AddressDecoderV2) AddressDecode(addr string, opts ...interface{}) ([]byte, error) {
	data, err := addressEncoder.Base58Decode(addr, addressEncoder.NewBase58Alphabet(alphabet))
	if err != nil {
		return nil, err
	}
	pubkey := data[1: len(data)-2 ]
	return pubkey, nil
}

//AddressEncode 地址编码
func (dec *AddressDecoderV2) AddressEncode(hash []byte, opts ...interface{}) (string, error) {
	if len(hash) != 32 {
		hash, _= owcrypt.CURVE25519_convert_Ed_to_X(hash)
	}
	prefix := []byte{ dec.wm.AddrPrefix() }
	data := addressEncoder.CatData(prefix, hash)
	input := addressEncoder.CatData(ssPrefix, data)
	checkSum := owcrypt.Hash(input, 64, owcrypt.HASH_ALG_BLAKE2B)[:2]
	result := addressEncoder.EncodeData( addressEncoder.CatData(data, checkSum), encodeType, alphabet)
	return result, nil
}

// AddressVerify 地址校验
func (dec *AddressDecoderV2) AddressVerify(address string, opts ...interface{}) bool {
	P2PKHPrefix := byte( dec.wm.AddrPrefix() )
	decodeBytes, err := addressEncoder.Base58Decode(address, addressEncoder.NewBase58Alphabet(alphabet) )
	if err != nil || len(decodeBytes) != 35 {
		return false
	}
	if decodeBytes[0] != P2PKHPrefix {
		return false
	}
	pub := decodeBytes[1: len(decodeBytes)-2 ]
	prefix := []byte{ dec.wm.AddrPrefix() }
	data := append(prefix, pub...)
	input := append(ssPrefix, data...)
	checkSum := owcrypt.Hash(input, 64, owcrypt.HASH_ALG_BLAKE2B)[:2]

	for i := 0; i < 2; i ++ {
		if checkSum[i] != decodeBytes[33 + i] {
			return false
		}
	}
	if len(pub) != 32 {
		return false
	}
	return true
}
