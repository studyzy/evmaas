package evmaas

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

type Address [20]byte

func NewAddress(addr string) Address {
	if strings.HasPrefix(addr, "0x") {
		addr = addr[2:]
	}
	addrBytes, _ := hex.DecodeString(addr)
	var address Address
	copy(address[:], addrBytes)
	return address
}
func BytesToAddress(b []byte) Address {
	var a Address

	size := len(b)
	min := min(size, 20)

	copy(a[20-min:], b[len(b)-min:])

	return a

}

// StateDB 接口定义了状态数据库的读取和更新方法。
type StateDB interface {
	GetAccountBalance(address Address) *big.Int
	//SetAccountBalance(address Address, balance *big.Int)
	GetContractCode(address Address) ([]byte, error)
	//SetContractCode(address Address, code []byte)
	GetState(address Address, key []byte) ([]byte, error)
	//SetContractStorage(address Address, key []byte, value []byte)
	// GetBlockHash returns the hash of a block by its number. blockhash() use it.
	GetBlockHash(number uint64) ([]byte, error)
}

// Transaction 结构体表示一个交易。
type Transaction struct {
	TxHash   []byte
	From     Address
	To       Address
	Value    *big.Int
	Gas      uint64
	GasPrice *big.Int
	Data     []byte
}

// Block 结构体表示一个区块。
type Block struct {
	BlockHash  []byte
	Number     uint64
	Timestamp  uint64
	Difficulty *big.Int
	GasLimit   uint64
}

// ExecutionResult 结构体表示EVM执行的结果。
type ExecutionResult struct {
	Success      bool
	ReturnData   []byte
	StateChanges map[Address]map[string][]byte
	Balance      map[Address]*big.Int
	ContractCode map[Address][]byte
	GasUsed      uint64
	Events       []EventLog
}

func (result *ExecutionResult) Receipt() *Receipt {
	return &Receipt{
		Success:    result.Success,
		ReturnData: result.ReturnData,
		GasUsed:    result.GasUsed,
		Events:     result.Events,
	}
}

type Receipt struct {
	Success    bool
	ReturnData []byte
	GasUsed    uint64
	Events     []EventLog
}

func (result *ExecutionResult) PutState(address Address, key []byte, value []byte) {
	if result.StateChanges[address] == nil {
		result.StateChanges[address] = make(map[string][]byte)
	}
	result.StateChanges[address][string(key)] = value

}

func NewExecutionResult() *ExecutionResult {
	return &ExecutionResult{
		StateChanges: make(map[Address]map[string][]byte),
		Balance:      make(map[Address]*big.Int),
		ContractCode: make(map[Address][]byte),
	}
}

// EventLog 结构体表示合约执行过程中触发的事件。
type EventLog struct {
	ContractAddress Address
	Topics          [][]byte
	Data            []byte
}

func (r *Receipt) ToString() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Success:%v\n", r.Success))
	sb.WriteString(fmt.Sprintf("ReturnData:%x\n", r.ReturnData))
	sb.WriteString(fmt.Sprintf("GasUsed:%v\n", r.GasUsed))
	sb.WriteString("Events:\n")
	for i, e := range r.Events {
		sb.WriteString(fmt.Sprintf("Event[%d]:\n", i))
		sb.WriteString(fmt.Sprintf("\tContractAddress:%x\n", e.ContractAddress))
		sb.WriteString(fmt.Sprintf("\tTopics:%x\n", e.Topics))
		sb.WriteString(fmt.Sprintf("\tData:%x\n", e.Data))
	}
	return sb.String()
}

type EvmInterface interface {
	InstallContract(tx Transaction, stateDB StateDB, block Block) (*ExecutionResult, error)
	ExecuteContract(tx Transaction, stateDB StateDB, block Block) (*ExecutionResult, error)
	QueryContract(tx Transaction, stateDB StateDB, block Block) (*ExecutionResult, error)
}
