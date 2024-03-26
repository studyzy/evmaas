package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

var TestChainConfig = &params.ChainConfig{
	ChainID:                       big.NewInt(1),
	HomesteadBlock:                big.NewInt(0),
	DAOForkBlock:                  nil,
	DAOForkSupport:                false,
	EIP150Block:                   big.NewInt(0),
	EIP155Block:                   big.NewInt(0),
	EIP158Block:                   big.NewInt(0),
	ByzantiumBlock:                big.NewInt(0),
	ConstantinopleBlock:           big.NewInt(0),
	PetersburgBlock:               big.NewInt(0),
	IstanbulBlock:                 big.NewInt(0),
	MuirGlacierBlock:              big.NewInt(0),
	BerlinBlock:                   big.NewInt(0),
	LondonBlock:                   big.NewInt(0),
	ArrowGlacierBlock:             big.NewInt(0),
	GrayGlacierBlock:              big.NewInt(0),
	MergeNetsplitBlock:            nil,
	ShanghaiTime:                  newUint64(0),
	CancunTime:                    nil,
	PragueTime:                    nil,
	VerkleTime:                    nil,
	TerminalTotalDifficulty:       nil,
	TerminalTotalDifficultyPassed: false,
	Ethash:                        new(params.EthashConfig),
	Clique:                        nil,
}

func newUint64(val uint64) *uint64 { return &val }
