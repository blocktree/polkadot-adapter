package polkadot

type ApiClient struct {
	Client    *Client
	WSClient  *WSClient
	APIChoose string
}

func NewApiClient(wm *WalletManager) error {
	api := ApiClient{}

	if len(wm.Config.APIChoose) == 0 {
		wm.Config.APIChoose = "rpc" //默认采用rpc连接
	}
	api.APIChoose = wm.Config.APIChoose
	if api.APIChoose == "rpc" {
		api.Client = NewClient(wm.Config.NodeAPI, false)
	} else if api.APIChoose == "ws" {
		api.WSClient = NewWSClient(wm, wm.Config.WSAPI, 0, false)
	}

	wm.ApiClient = &api

	return nil
}

// 获取当前最高区块
func (c *ApiClient) getBlockHeight() (uint64, error) {
	var (
		currentHeight uint64
		err           error
	)
	if c.APIChoose == "rpc" {
		currentHeight, err = c.Client.getBlockHeight()
	} else if c.APIChoose == "ws" {
		currentHeight, err = c.WSClient.getBlockHeight()
	}

	return currentHeight, err
}

//获取当前最新高度
func (c *ApiClient) getMostHeightBlock() (*Block, error) {
	var (
		mostHeightBlock *Block
		err             error
	)
	if c.APIChoose == "rpc" {
		mostHeightBlock, err = c.Client.getMostHeightBlock()
	} else if c.APIChoose == "ws" {
		//mostHeightBlock, err = decoder.wm.WSClient.getBlockHeight()
	}

	return mostHeightBlock, err
}

// 获取地址余额
func (c *ApiClient) getBalance(address string, ignoreReserve bool, reserveAmount int64) (*AddrBalance, error) {
	var (
		balance *AddrBalance
		err     error
	)

	if c.APIChoose == "rpc" {
		balance, err = c.Client.getBalance(address, ignoreReserve, reserveAmount)
	} else if c.APIChoose == "ws" {
		balance, err = c.WSClient.getBalance(address, ignoreReserve, reserveAmount)
	}

	return balance, err
}

func (c *ApiClient) getBlockByHeight(height uint64) (*Block, error) {
	var (
		block *Block
		err   error
	)
	if c.APIChoose == "rpc" {
		block, err = c.Client.getBlockByHeight(height)
	} else if c.APIChoose == "ws" {
		block, err = c.WSClient.getBlockByHeight(height)
	}

	return block, err
}

func (c *ApiClient) sendTransaction(rawTx string) (string, error) {
	var (
		txid string
		err  error
	)
	if c.APIChoose == "rpc" {
		txid, err = c.Client.sendTransaction(rawTx)
	} else if c.APIChoose == "ws" {
		txid, err = c.WSClient.sendTransaction(rawTx)
	}

	return txid, err
}

func (c *ApiClient) getTxArtifacts() (*TxArtifacts, error) {
	var (
		txArtifacts *TxArtifacts
		err         error
	)
	if c.APIChoose == "rpc" {
		txArtifacts, err = c.Client.getTxArtifacts()
	} else if c.APIChoose == "ws" {
		//block, err = c.WSClient.TxArtifacts()
	}

	return txArtifacts, err
}
