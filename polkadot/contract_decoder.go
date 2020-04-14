package polkadot

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/shopspring/decimal"
)

func convertFlostStringToBigInt(amount string) (*big.Int, error) {
	vDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		log.Error("convert from string to decimal failed, err=", err)
		return nil, err
	}

	decimalInt := big.NewInt(1)
	for i := 0; i < 9; i++ {
		decimalInt.Mul(decimalInt, big.NewInt(10))
	}
	d, _ := decimal.NewFromString(decimalInt.String())
	vDecimal = vDecimal.Mul(d)
	rst := new(big.Int)
	if _, valid := rst.SetString(vDecimal.String(), 10); !valid {
		log.Error("conver to big.int failed")
		return nil, errors.New("conver to big.int failed")
	}
	return rst, nil
}

func convertBigIntToFloatDecimal(amount string) (decimal.Decimal, error) {
	d, err := decimal.NewFromString(amount)
	if err != nil {
		log.Error("convert string to deciaml failed, err=", err)
		return d, err
	}

	decimalInt := big.NewInt(1)
	for i := 0; i < 9; i++ {
		decimalInt.Mul(decimalInt, big.NewInt(10))
	}

	w, _ := decimal.NewFromString(decimalInt.String())
	d = d.Div(w)
	return d, nil
}

func convertIntStringToBigInt(amount string) (*big.Int, error) {
	vInt64, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		log.Error("convert from string to int failed, err=", err)
		return nil, err
	}

	return big.NewInt(vInt64), nil
}

type ContractDecoder struct {
	*openwallet.SmartContractDecoderBase
	wm *WalletManager
}

//NewContractDecoder 智能合约解析器
func NewContractDecoder(wm *WalletManager) *ContractDecoder {
	decoder := ContractDecoder{}
	decoder.wm = wm
	return &decoder
}

func (decoder *ContractDecoder) GetTokenBalanceByAddress(contract openwallet.SmartContract, address ...string) ([]*openwallet.TokenBalance, error) {
	return nil, nil
}
