package cmevm

import (
	"github.com/studyzy/evmaas"
)

type ChainMakerEvmImpl struct {
}

func (vm *ChainMakerEvmImpl) InstallContract(tx evmaas.Transaction, stateDB evmaas.StateDB, block evmaas.Block) (
	evmaas.ExecutionResult, error) {
	//instance := &evm.RuntimeInstance{
	//	Method:        "",
	//	ChainId:       "",
	//	Contract:      nil,
	//	Log:           nil,
	//	TxSimContext:  nil,
	//	ContractEvent: nil,
	//}
	//instance.Invoke()
	panic("implement me")
}
func (vm *ChainMakerEvmImpl) ExecuteContract(tx evmaas.Transaction, stateDB evmaas.StateDB, block evmaas.Block) (
	evmaas.ExecutionResult, error) {
	//instance := &evm.RuntimeInstance{
	//	Method:        "",
	//	ChainId:       "",
	//	Contract:      nil,
	//	Log:           nil,
	//	TxSimContext:  nil,
	//	ContractEvent: nil,
	//}
	panic("implement me")
}
