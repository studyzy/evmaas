package sealevm

import (
	"fmt"
	"log"

	"github.com/SealSC/SealEVM"
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/storage"
	"github.com/studyzy/evmaas"
)

type SealEVMImpl struct {
	statedb evmaas.StateDB
}

func NewSealEVMImpl() *SealEVMImpl {
	SealEVM.Load()
	sdb := evmaas.NewMemStateDB()

	return &SealEVMImpl{statedb: sdb}
}

// create a new evm
func (vm *SealEVMImpl) newEVMParam(tx evmaas.Transaction, ms storage.IExternalStorage, block evmaas.Block) SealEVM.EVMParam {
	callerInt := evmaas.AddressToInt(tx.From)
	blockHashInt := evmaas.BytesToInt(block.BlockHash)
	p := SealEVM.EVMParam{
		MaxStackDepth:  0,
		ExternalStore:  ms,
		ResultCallback: nil,
		Context: &environment.Context{
			Block: environment.Block{
				ChainID:    evmInt256.New(666),
				Coinbase:   evmInt256.New(0),
				Timestamp:  evmInt256.New(int64(block.Timestamp)),
				Number:     evmInt256.New(int64(block.Number)),
				Difficulty: evmInt256.New(0),
				GasLimit:   evmInt256.New(10000000),
				Hash:       blockHashInt,
			},
			//Contract: contract,
			Transaction: environment.Transaction{
				Origin:   callerInt,
				GasPrice: evmInt256.New(1),
				GasLimit: evmInt256.New(10000000),
			},
			Message: environment.Message{
				Caller: callerInt,
				Value:  evmInt256.New(0),
				Data:   tx.Data,
			},
		},
	}

	return p
}
func (vm *SealEVMImpl) calcContractAddress(tx evmaas.Transaction) evmaas.Address {
	//新合约的地址由安装合约的交易生成
	contractAddr := tx.TxHash[12:32]
	return evmaas.Address(contractAddr)
}
func (vm *SealEVMImpl) InstallContract(tx evmaas.Transaction, stateDB evmaas.StateDB, block evmaas.Block) (
	*evmaas.ExecutionResult, error) {
	contractAddr := vm.calcContractAddress(tx)
	hashInt := evmInt256.New(0)
	hashInt.SetBytes(contractAddr[:])

	//same contract code has same address in this example
	cNamespace := hashInt
	contract := environment.Contract{
		Namespace: cNamespace,
		Code:      tx.Data,
		Hash:      hashInt,
	}
	ms := newMemStorage(stateDB)

	evmP := vm.newEVMParam(tx, ms, block)
	evmP.Context.Contract = contract
	evm := SealEVM.New(evmP)
	ret, err := evm.ExecuteContract(false)

	//check error
	if err != nil {
		return nil, err
	}
	contractCode := ret.ResultData

	result := evmaas.NewExecutionResult()

	//保存合约代码
	result.ContractCode[contractAddr] = contractCode
	//保存状态更新
	for addr, cache := range ret.StorageCache.CachedData {
		for key, v := range cache {
			log.Printf("save statedb: addr[%x],key[%x]=value[%x]\n", addr, key, v.Bytes())
			result.PutState(evmaas.Address([]byte(addr)), []byte(key), v.Bytes())
		}
	}
	//保存Balance
	//TODO
	//保存Logs
	for _, logs := range ret.StorageCache.Logs {
		for _, l := range logs {

			eventLog := evmaas.EventLog{
				ContractAddress: evmaas.Address(l.Context.Contract.Namespace.Bytes()),
				Topics:          l.Topics,
				Data:            l.Data,
			}
			result.Events = append(result.Events, eventLog)
			for _, t := range l.Topics {
				fmt.Printf("topic:%x\n", t)
			}
			fmt.Printf("data :%x\n", l.Data)
		}
	}

	//返回执行结果
	return &evmaas.ExecutionResult{
		Success:      true,
		ReturnData:   ret.ResultData,
		StateChanges: nil,
		GasUsed:      ret.GasLeft,
		Events:       nil,
	}, nil

}
func (vm *SealEVMImpl) ExecuteContract(tx evmaas.Transaction, stateDB evmaas.StateDB, block evmaas.Block) (
	*evmaas.ExecutionResult, error) {
	contractAddr := tx.To
	contractCode, err := vm.statedb.GetContractCode(contractAddr)
	if err != nil {
		return nil, err

	}
	hashInt := evmInt256.New(0)
	hashInt.SetBytes(contractAddr[:])

	//same contract code has same address in this example
	cNamespace := hashInt
	contract := environment.Contract{
		Namespace: cNamespace,
		Code:      contractCode,
		Hash:      hashInt,
	}
	ms := newMemStorage(stateDB)

	evmP := vm.newEVMParam(tx, ms, block)
	evmP.Context.Contract = contract
	evm := SealEVM.New(evmP)
	ret, err := evm.ExecuteContract(false)

	//check error
	if err != nil {
		return nil, err
	}

	result := evmaas.NewExecutionResult()

	//保存合约代码
	for addr, code := range ms.contracts {
		result.ContractCode[evmaas.NewAddress(addr)] = code
	}
	//保存状态更新
	for addr, cache := range ret.StorageCache.CachedData {
		for key, v := range cache {
			log.Printf("save statedb: addr[%x],key[%x]=value[%x]\n", addr, key, v.Bytes())
			result.PutState(evmaas.Address([]byte(addr)), []byte(key), v.Bytes())
		}
	}
	//保存Balance
	//TODO
	//保存Logs
	for _, logs := range ret.StorageCache.Logs {
		for _, l := range logs {

			eventLog := evmaas.EventLog{
				ContractAddress: evmaas.Address(l.Context.Contract.Namespace.Bytes()),
				Topics:          l.Topics,
				Data:            l.Data,
			}
			result.Events = append(result.Events, eventLog)
			for _, t := range l.Topics {
				fmt.Printf("topic:%x\n", t)
			}
			fmt.Printf("data :%x\n", l.Data)
		}
	}

	//返回执行结果
	return &evmaas.ExecutionResult{
		Success:      true,
		ReturnData:   ret.ResultData,
		StateChanges: nil,
		GasUsed:      ret.GasLeft,
		Events:       nil,
	}, nil

}
