package evmaas

import "math/big"

type Address [20]byte

// StateDB 接口定义了状态数据库的读取和更新方法。
type StateDB interface {
	GetAccountBalance(address Address) *big.Int
	SetAccountBalance(address Address, balance *big.Int)
	GetContractCode(address Address) []byte
	SetContractCode(address Address, code []byte)
	GetContractStorage(address Address, key []byte) []byte
	SetContractStorage(address Address, key []byte, value []byte)
	// GetBlockHash returns the hash of a block by its number. blockhash() use it.
	GetBlockHash(number uint64) []byte
}

// Transaction 结构体表示一个交易。
type Transaction struct {
	From     Address
	To       Address
	Value    *big.Int
	Gas      uint64
	GasPrice *big.Int
	Data     []byte
}

// Block 结构体表示一个区块。
type Block struct {
	Number     uint64
	Timestamp  uint64
	Difficulty *big.Int
	GasLimit   uint64
}

// ExecutionResult 结构体表示EVM执行的结果。
type ExecutionResult struct {
	Success      bool
	ReturnData   []byte
	StateChanges map[string]interface{}
	GasUsed      uint64
	Events       []Event
	Error        error
}

// Event 结构体表示合约执行过程中触发的事件。
type Event struct {
	ContractAddress Address
	Topics          [][]byte
	Data            []byte
}

type EvmInterface interface {
	ExecuteEVM(tx Transaction, stateDB StateDB, block Block) (ExecutionResult, error)
}
