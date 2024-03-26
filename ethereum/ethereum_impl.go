package ethereum

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	cmath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/studyzy/evmaas"
)

var baseFee *big.Int = big.NewInt(1000000000) //1 GWei

type EthereumImpl struct {
}

func NewEthereumImpl() *EthereumImpl {
	return &EthereumImpl{}
}

func (v *EthereumImpl) InstallContract(tx evmaas.Transaction, stateDB evmaas.StateDB, block evmaas.Block) (
	*evmaas.ExecutionResult, error) {
	// 初始化区块链状态数据库
	statedb := NewStateDb(stateDB)
	header := &types.Header{
		Difficulty: big.NewInt(0),
		Number:     big.NewInt(int64(block.Number)),
		BaseFee:    baseFee,
		Time:       uint64(time.Now().Unix()),
	}
	chainConfig := TestChainConfig
	cfg := vm.Config{}
	author := &common.Address{}
	var (
		context = core.NewEVMBlockContext(header, nil, author)
		vmenv   = vm.NewEVM(context, vm.TxContext{}, statedb, chainConfig, cfg)
	)
	etx := types.NewTransaction(0, common.Address{}, tx.Value, tx.Gas, tx.GasPrice, tx.Data)

	msg, err := TransactionToMessage(etx, tx.From, nil, header.BaseFee)
	if err != nil {
		return nil, fmt.Errorf("could not apply tx [%v]: %w", etx.Hash().Hex(), err)
	}
	//statedb.SetTxContext(tx.Hash(), i)
	var gp = new(core.GasPool).AddGas(block.GasLimit)
	//blockNumber := big.NewInt(int64(block.Number))
	//blockHash := common.BytesToHash(block.BlockHash)
	//usedGas := new(uint64)

	// Create a new context to be used in the EVM environment.
	txContext := core.NewEVMTxContext(msg)
	vmenv.Reset(txContext, statedb)

	// Apply the transaction to the current state (included in the env).
	result, err := core.ApplyMessage(vmenv, msg, gp)
	if err != nil {
		return nil, err
	}

	return convertReceipt2Result(result, statedb)
}

func convertReceipt2Result(res *core.ExecutionResult, statedb *StateDb) (*evmaas.ExecutionResult, error) {
	if res.Err != nil {
		return nil, res.Err
	}

	result := &evmaas.ExecutionResult{
		Success:      true,
		ReturnData:   res.ReturnData,
		StateChanges: make(map[evmaas.Address]map[string][]byte),
		Balance:      make(map[evmaas.Address]*big.Int),
		Events:       make([]evmaas.EventLog, 0),
		ContractCode: make(map[evmaas.Address][]byte),
		GasUsed:      res.UsedGas,
	}
	//log
	for _, log := range statedb.logs {
		ev := evmaas.EventLog{
			ContractAddress: evmaas.BytesToAddress(log.Address[:]),
			Topics:          make([][]byte, 0),
			Data:            log.Data,
		}
		for _, topic := range log.Topics {
			ev.Topics = append(ev.Topics, topic[:])
		}
		result.Events = append(result.Events, ev)
	}
	//state
	for addr, kv := range statedb.state {
		a := evmaas.BytesToAddress(addr[:])
		for k, v := range kv {
			if _, ok := result.StateChanges[a]; !ok {
				result.StateChanges[a] = make(map[string][]byte)
			}
			result.StateChanges[a][k] = v
		}
	}
	//account balance
	for addr, balance := range statedb.accBalance {
		a := evmaas.BytesToAddress(addr[:])
		result.Balance[a] = balance.ToBig()
	}
	//contract code
	for addr, code := range statedb.contracts {
		a := evmaas.BytesToAddress(addr[:])
		result.ContractCode[a] = code
	}

	return result, nil
}
func (v *EthereumImpl) ExecuteContract(tx evmaas.Transaction, stateDB evmaas.StateDB, block evmaas.Block) (
	*evmaas.ExecutionResult, error) {
	// 初始化区块链状态数据库
	statedb := NewStateDb(stateDB)
	header := &types.Header{
		Difficulty: big.NewInt(0),
		Number:     big.NewInt(int64(block.Number)),
		BaseFee:    baseFee,
		Time:       uint64(time.Now().Unix()),
	}
	chainConfig := TestChainConfig
	cfg := vm.Config{}
	author := &common.Address{}
	var (
		context = core.NewEVMBlockContext(header, nil, author)
		vmenv   = vm.NewEVM(context, vm.TxContext{}, statedb, chainConfig, cfg)
	)
	etx := types.NewTransaction(0, common.Address{}, tx.Value, tx.Gas, tx.GasPrice, tx.Data)
	to := common.BytesToAddress(tx.To[:])
	msg, err := TransactionToMessage(etx, tx.From, &to, header.BaseFee)
	if err != nil {
		return nil, fmt.Errorf("could not apply tx [%v]: %w", etx.Hash().Hex(), err)
	}
	//statedb.SetTxContext(tx.Hash(), i)
	var gp = new(core.GasPool).AddGas(block.GasLimit)
	//blockNumber := big.NewInt(int64(block.Number))
	//blockHash := common.BytesToHash(block.BlockHash)
	//usedGas := new(uint64)
	//receipt, err := core.ApplyTransactionWithEVM(msg, chainConfig, gp, statedb, blockNumber, blockHash, etx, usedGas, vmenv)
	//if err != nil {
	//	return nil, fmt.Errorf("could not apply tx [%v]: %w", etx.Hash().Hex(), err)
	//}
	txContext := core.NewEVMTxContext(msg)
	vmenv.Reset(txContext, statedb)

	// Apply the transaction to the current state (included in the env).
	result, err := core.ApplyMessage(vmenv, msg, gp)
	if err != nil {
		return nil, err
	}
	return convertReceipt2Result(result, statedb)
}

func (v *EthereumImpl) QueryContract(tx evmaas.Transaction, stateDB evmaas.StateDB, block evmaas.Block) (
	*evmaas.ExecutionResult, error) {
	// 初始化区块链状态数据库
	statedb := NewStateDb(stateDB)
	header := &types.Header{
		Difficulty: big.NewInt(0),
		Number:     big.NewInt(int64(block.Number)),
		BaseFee:    baseFee,
		Time:       uint64(time.Now().Unix()),
	}
	chainConfig := TestChainConfig
	cfg := vm.Config{}
	author := &common.Address{}
	var (
		context = core.NewEVMBlockContext(header, nil, author)
		vmenv   = vm.NewEVM(context, vm.TxContext{}, statedb, chainConfig, cfg)
	)
	etx := types.NewTransaction(0, common.Address{}, tx.Value, tx.Gas, tx.GasPrice, tx.Data)
	to := common.BytesToAddress(tx.To[:])
	msg, err := TransactionToMessage(etx, tx.From, &to, header.BaseFee)
	if err != nil {
		return nil, fmt.Errorf("could not apply tx [%v]: %w", etx.Hash().Hex(), err)
	}
	//statedb.SetTxContext(tx.Hash(), i)
	var gp = new(core.GasPool).AddGas(block.GasLimit)
	result, err := core.ApplyMessage(vmenv, msg, gp)
	if err != nil {
		return nil, fmt.Errorf("could not apply tx [%v]: %w", etx.Hash().Hex(), err)
	}

	return convertResult2Result(result)
}

func convertResult2Result(result *core.ExecutionResult) (*evmaas.ExecutionResult, error) {
	if result.Err != nil {
		return nil, result.Err

	}
	return &evmaas.ExecutionResult{
		ReturnData: result.ReturnData,
		GasUsed:    result.UsedGas,
	}, nil
}
func TransactionToMessage(tx *types.Transaction, from evmaas.Address, to *common.Address, baseFee *big.Int) (*core.Message, error) {
	msg := &core.Message{
		Nonce:             tx.Nonce(),
		GasLimit:          tx.Gas(),
		GasPrice:          new(big.Int).Set(tx.GasPrice()),
		GasFeeCap:         baseFee,
		GasTipCap:         new(big.Int).Set(tx.GasTipCap()),
		To:                to,
		Value:             tx.Value(),
		Data:              tx.Data(),
		AccessList:        tx.AccessList(),
		SkipAccountChecks: true,
		BlobHashes:        tx.BlobHashes(),
		BlobGasFeeCap:     tx.BlobGasFeeCap(),
	}
	// If baseFee provided, set gasPrice to effectiveGasPrice.
	if baseFee != nil {
		msg.GasPrice = cmath.BigMin(msg.GasPrice.Add(msg.GasTipCap, baseFee), msg.GasFeeCap)
	}
	var err error
	msg.From = common.BytesToAddress(from[:])
	return msg, err
}
