package connector

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
)

const (
	clientIdentifier = "geth" // Client identifier to advertise over the network
)

type gethConfig struct {
	Eth      ethconfig.Config
	Node     node.Config
	Metrics  metrics.Config
	Ethstats ethstatsConfig
}

type ethstatsConfig struct {
	URL string `toml:",omitempty"`
}

func defaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = params.VersionWithCommit("", "")
	cfg.HTTPModules = append(cfg.HTTPModules, "eth")
	cfg.WSModules = append(cfg.WSModules, "eth")
	cfg.IPCPath = "geth.ipc"
	return cfg
}

func makeTxPool() *core.TxPool {
	ethConfigIns := ethconfig.Defaults
	nodeCfg := defaultNodeConfig()
	stack, err := node.New(&nodeCfg)
	if err != nil {
		panic(err)
	}
	ethereumIns, err := eth.New(stack, &ethConfigIns)
	if err != nil {
		panic(err)
	}
	return core.NewTxPool(core.TxPoolConfig{
		Locals:       nil,
		NoLocals:     true,
		Journal:      "transactions.rlp",
		Rejournal:    3600000000000,
		PriceLimit:   1,
		PriceBump:    10,
		AccountSlots: 16,
		GlobalSlots:  4096,
		AccountQueue: 64,
		GlobalQueue:  1024,
		Lifetime:     10800000000000,
	}, params.MainnetChainConfig, ethereumIns.BlockChain())
}

func makeEthereum() *eth.Ethereum {
	ethConfigIns := ethconfig.Defaults
	ethConfigIns.Miner.Etherbase = common.HexToAddress("0x773355277126cbDCf8EB80702f6bc1A3Cb843Bbb")
	nodeCfg := defaultNodeConfig()
	stack, err := node.New(&nodeCfg)
	if err != nil {
		panic(err)
	}
	ethereumIns, err := eth.New(stack, &ethConfigIns)
	if err != nil {
		panic(err)
	}
	return ethereumIns
}
