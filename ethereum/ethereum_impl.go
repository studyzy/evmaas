package ethereum

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/studyzy/evmaas"
)

type EthereumImpl struct {
}

func NewEthereumImpl() *EthereumImpl {
	return &EthereumImpl{}
}

func (v *EthereumImpl) InstallContract(tx evmaas.Transaction, stateDB evmaas.StateDB, block evmaas.Block) (
	*evmaas.ExecutionResult, error) {
	// 初始化区块链状态数据库
	statedb := NewStateDb(stateDB)
	header := &types.Header{Difficulty: big.NewInt(0), Number: big.NewInt(int64(block.Number))}
	chainConfig := params.TestChainConfig
	cfg := vm.Config{}
	author := &common.Address{}
	var (
		context = core.NewEVMBlockContext(header, nil, author)
		vmenv   = vm.NewEVM(context, vm.TxContext{}, statedb, chainConfig, cfg)
		signer  = types.MakeSigner(chainConfig, header.Number, header.Time)
	)
	etx := types.NewTransaction(0, common.Address{}, tx.Value, tx.Gas, tx.GasPrice, tx.Data)

	msg, err := core.TransactionToMessage(etx, signer, header.BaseFee)
	if err != nil {
		return nil, fmt.Errorf("could not apply tx [%v]: %w", etx.Hash().Hex(), err)
	}
	//statedb.SetTxContext(tx.Hash(), i)
	var gp = new(core.GasPool).AddGas(header.GasLimit)
	blockNumber := big.NewInt(int64(block.Number))
	blockHash := common.BytesToHash(block.BlockHash)
	usedGas := new(uint64)
	receipt, err := core.ApplyTransactionWithEVM(msg, chainConfig, gp, statedb, blockNumber, blockHash, etx, usedGas, vmenv)
	if err != nil {
		return nil, fmt.Errorf("could not apply tx [%v]: %w", etx.Hash().Hex(), err)
	}

	return convertResult(receipt)
}

func convertResult(res *types.Receipt) (*evmaas.ExecutionResult, error) {
	result := &evmaas.ExecutionResult{
		Success:      res.Status == types.ReceiptStatusSuccessful,
		ReturnData:   res.PostState,
		StateChanges: make(map[evmaas.Address]map[string][]byte),
		Balance:      make(map[evmaas.Address]*big.Int),
		Events:       make([]evmaas.EventLog, 0),
	}
	//log
	for _, log := range res.Logs {
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
	return result, nil
}
func (vm *EthereumImpl) ExecuteContract(tx evmaas.Transaction, stateDB evmaas.StateDB, block evmaas.Block) (
	*evmaas.ExecutionResult, error) {

	return nil, nil
}
