package polygon

import (
	"math/big"

	"github.com/0xPolygon/polygon-edge/chain"
	"github.com/0xPolygon/polygon-edge/crypto"
	"github.com/0xPolygon/polygon-edge/state/runtime"
	"github.com/0xPolygon/polygon-edge/state/runtime/evm"
	"github.com/0xPolygon/polygon-edge/types"
	"github.com/studyzy/evmaas"
)

type PolygonImpl struct {
}

func (v *PolygonImpl) InstallContract(tx evmaas.Transaction, stateDB evmaas.StateDB, block evmaas.Block) (
	*evmaas.ExecutionResult, error) {
	vm := evm.NewEVM()
	sender := types.BytesToAddress(tx.From[:])
	contractAddr := crypto.CreateAddress(sender, 0)
	contract := runtime.NewContractCreation(0, sender, sender, contractAddr, tx.Value, tx.Gas, tx.Data)
	txContext := runtime.TxContext{
		GasPrice:     types.Hash{},
		Origin:       sender,
		Coinbase:     types.Address{},
		Number:       0,
		Timestamp:    int64(block.Timestamp),
		GasLimit:     int64(tx.Gas),
		ChainID:      666,
		Difficulty:   types.Hash{},
		Tracer:       nil,
		NonPayable:   false,
		BaseFee:      nil,
		BurnContract: types.Address{},
	}
	host := NewMemHost(stateDB, txContext)
	config := chain.AllForksEnabled.At(1)
	result := vm.Run(contract, host, &config)
	return ConvertResult(result, host)
}

func ConvertResult(result *runtime.ExecutionResult, host *MemHost) (*evmaas.ExecutionResult, error) {
	if result.Err != nil {
		return nil, result.Err
	}
	evmResult := &evmaas.ExecutionResult{
		Success:      true,
		ReturnData:   result.ReturnValue,
		GasUsed:      result.GasUsed,
		Events:       make([]evmaas.EventLog, 0),
		StateChanges: make(map[evmaas.Address]map[string][]byte),
		ContractCode: make(map[evmaas.Address][]byte),
		Balance:      make(map[evmaas.Address]*big.Int),
	}
	//set logs
	evmResult.Events = host.logs
	//set state
	for contract, kv := range host.state {
		addr := evmaas.BytesToAddress(contract[:])
		evmResult.StateChanges[addr] = make(map[string][]byte)
		for k, v := range kv {
			evmResult.StateChanges[addr][string(k[:])] = v[:]
		}
	}
	//set contract code
	for contract, code := range host.contracts {
		addr := evmaas.BytesToAddress(contract[:])
		evmResult.ContractCode[addr] = code

	}
	//set account balance
	for addr, balance := range host.accountBalance {

		evmResult.Balance[evmaas.BytesToAddress(addr[:])] = balance
	}

	return evmResult, nil
}
func (v *PolygonImpl) ExecuteContract(tx evmaas.Transaction, stateDB evmaas.StateDB, block evmaas.Block) (
	*evmaas.ExecutionResult, error) {
	vm := evm.NewEVM()
	sender := types.BytesToAddress(tx.From[:])
	to := types.BytesToAddress(tx.To[:])
	contract := runtime.NewContract(0, sender, sender, to, tx.Value, tx.Gas, tx.Data)
	txContext := runtime.TxContext{
		GasPrice:     types.Hash{},
		Origin:       sender,
		Coinbase:     types.Address{},
		Number:       0,
		Timestamp:    int64(block.Timestamp),
		GasLimit:     int64(tx.Gas),
		ChainID:      666,
		Difficulty:   types.Hash{},
		Tracer:       nil,
		NonPayable:   false,
		BaseFee:      nil,
		BurnContract: types.Address{},
	}
	host := NewMemHost(stateDB, txContext)
	config := chain.AllForksEnabled.At(1)
	result := vm.Run(contract, host, &config)
	return ConvertResult(result, host)
}
