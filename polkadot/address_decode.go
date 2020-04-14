package polkadot

import (
	"encoding/hex"
	"github.com/blocktree/go-owcdrivers/addressEncoder"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/v2/openwallet"
)

var (
	alphabet = addressEncoder.BTCAlphabet
	ssPrefix = []byte{0x53, 0x53, 0x35, 0x38, 0x50, 0x52, 0x45}
	prefix = []byte{0x02}
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
	data := addressEncoder.CatData(prefix, hash)
	input := addressEncoder.CatData(ssPrefix, data)
	checkSum := owcrypt.Hash(input, 64, owcrypt.HASH_ALG_BLAKE2B)[:2]
	result := addressEncoder.EncodeData( addressEncoder.CatData(data, checkSum), encodeType, alphabet)
	return result, nil
}

// AddressVerify 地址校验
func (dec *AddressDecoderV2) AddressVerify(address string, opts ...interface{}) bool {
	pub, err := hex.DecodeString(address)
	if err != nil {
		return false
	}

	if len(pub) != 33 {
		return false
	}
	return true
}
