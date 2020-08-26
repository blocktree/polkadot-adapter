module github.com/blocktree/polkadot-adapter

go 1.12

require (
	github.com/asdine/storm v2.1.2+incompatible
	github.com/astaxie/beego v1.12.0
	github.com/blocktree/go-owcdrivers v1.2.9
	github.com/blocktree/go-owcrypt v1.1.4
	github.com/blocktree/openwallet/v2 v2.0.2
	github.com/bwmarrin/snowflake v0.3.0
	github.com/ethereum/go-ethereum v1.9.9
	github.com/gorilla/websocket v1.4.1
	github.com/imroc/req v0.2.4
	github.com/pborman/uuid v1.2.0
	github.com/prometheus/common v0.6.0
	github.com/shopspring/decimal v0.0.0-20200105231215-408a2507e114
	github.com/tidwall/gjson v1.3.5
)

//replace github.com/blocktree/go-owcdrivers => ../go-owcdrivers
