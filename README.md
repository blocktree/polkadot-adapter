# polkadot-adapter

本项目适配了openwallet.AssetsAdapter接口，给应用提供了底层的区块链协议支持。

## 如何测试

openwtester包下的测试用例已经集成了openwallet钱包体系，创建conf文件，新建DOT.ini文件，编辑如下内容：

```ini

# node api url
nodeAPI = "http://127.0.0.1:8080"

# ws api url
# wsAPI = "ws://127.0.0.1:12345"

# fixed Fee in smallest unit
fixedFee = 10000000000
# reserve amount in smallest unit
reserveAmount = 0
# ignore reserve amount
ignoreReserve = true
# register fee in smallest unit
registerFee = 10000
# last ledger sequence number
lastLedgerSequenceNumber = 20
# memo type
memoType = "withdraw"
# memo format
memoFormat = "text/plain"
# which feild of memo to scan
memoScan = "MemoData"
# Cache data file directory, default = "", current directory: ./data
dataDir = ""
APIChoose = "rpc"

```
