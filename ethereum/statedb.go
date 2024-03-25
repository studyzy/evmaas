package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
	"github.com/studyzy/evmaas"
)

type StateDb struct {
	innerdb    evmaas.StateDB
	state      map[common.Address]map[string][]byte
	contracts  map[common.Address][]byte
	accBalance map[common.Address]*uint256.Int
	logs       []*types.Log
}

func (s *StateDb) Finalise(deleteEmptyObjects bool) {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) IntermediateRoot(deleteEmptyObjects bool) common.Hash {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) GetLogs(hash common.Hash, blockNumber uint64, blockHash common.Hash) []*types.Log {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) TxIndex() int {
	//TODO implement me
	panic("implement me")
}

func NewStateDb(innerdb evmaas.StateDB) *StateDb {
	return &StateDb{innerdb: innerdb}
}
func (s *StateDb) CreateAccount(address common.Address) {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) SubBalance(address common.Address, u *uint256.Int, t tracing.BalanceChangeReason) {
	balance := s.accBalance[address]
	newBalance := balance.Sub(balance, u)
	s.accBalance[address] = newBalance
}

func (s *StateDb) AddBalance(address common.Address, u *uint256.Int, t tracing.BalanceChangeReason) {
	balance := s.accBalance[address]
	newBalance := balance.Add(balance, u)
	s.accBalance[address] = newBalance
}

func (s *StateDb) GetBalance(address common.Address) *uint256.Int {
	balance, ok := s.accBalance[address]
	if !ok {
		dbBalance := s.innerdb.GetAccountBalance(evmaas.BytesToAddress(address[:]))
		balance, ok = uint256.FromBig(dbBalance)
	}
	return balance
}

func (s *StateDb) GetNonce(address common.Address) uint64 {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) SetNonce(address common.Address, u uint64) {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) GetCodeHash(address common.Address) common.Hash {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) GetCode(address common.Address) []byte {
	code, ok := s.contracts[address]
	if !ok {
		code, _ = s.innerdb.GetContractCode(evmaas.BytesToAddress(address[:]))
	}
	return code
}

func (s *StateDb) SetCode(address common.Address, bytes []byte) {
	s.contracts[address] = bytes
}

func (s *StateDb) GetCodeSize(address common.Address) int {
	code := s.GetCode(address)
	return len(code)
}

func (s *StateDb) AddRefund(u uint64) {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) SubRefund(u uint64) {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) GetRefund() uint64 {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) GetCommittedState(address common.Address, hash common.Hash) common.Hash {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) GetState(address common.Address, hash common.Hash) common.Hash {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) SetState(address common.Address, hash common.Hash, hash2 common.Hash) {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) SetTransientState(addr common.Address, key, value common.Hash) {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) SelfDestruct(address common.Address) {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) HasSelfDestructed(address common.Address) bool {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) Selfdestruct6780(address common.Address) {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) Exist(address common.Address) bool {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) Empty(address common.Address) bool {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) AddressInAccessList(addr common.Address) bool {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool) {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) AddAddressToAccessList(addr common.Address) {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) RevertToSnapshot(i int) {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) Snapshot() int {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) AddLog(log *types.Log) {
	s.logs = append(s.logs, log)
}

func (s *StateDb) AddPreimage(hash common.Hash, bytes []byte) {
	//TODO implement me
	panic("implement me")
}

var _ vm.StateDB = (*StateDb)(nil)
