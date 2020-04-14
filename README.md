# polkadot-adapter

本项目适配了openwallet.AssetsAdapter接口，给应用提供了底层的区块链协议支持。

## 如何测试

openwtester包下的测试用例已经集成了openwallet钱包体系，创建conf文件，新建DOT.ini文件，编辑如下内容：

```ini

# mainnet node api url
mainnetNodeAPI = "http://127.0.0.1:10026/"
# testnet node api url
testnetNodeAPI = "http://localhost:9922/"
# Is network test?
isTestNet = false
# feeScale
feeScale = 100
# feeCharge is now a fixed value 10000000(aka 0.1VSYS)
feeCharge = 10000000
# Cache data file directory, default = "", current directory: ./data
dataDir = ""

```
