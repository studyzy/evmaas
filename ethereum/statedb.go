package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
	"github.com/studyzy/evmaas"
	"golang.org/x/crypto/sha3"
)

type StateDb struct {
	innerdb    evmaas.StateDB
	state      map[common.Address]map[string][]byte
	contracts  map[common.Address][]byte
	accBalance map[common.Address]*uint256.Int
	logs       []*types.Log
}

func NewStateDb(innerdb evmaas.StateDB) *StateDb {
	return &StateDb{
		innerdb:    innerdb,
		accBalance: make(map[common.Address]*uint256.Int),
		contracts:  make(map[common.Address][]byte),
		state:      make(map[common.Address]map[string][]byte),
		logs:       make([]*types.Log, 0),
	}
}
func (s *StateDb) Finalise(deleteEmptyObjects bool) {
	return
}

func (s *StateDb) IntermediateRoot(deleteEmptyObjects bool) common.Hash {
	return common.Hash{}
}

func (s *StateDb) GetLogs(hash common.Hash, blockNumber uint64, blockHash common.Hash) []*types.Log {
	return s.logs
}

func (s *StateDb) TxIndex() int {
	return 0
}

func (s *StateDb) CreateAccount(address common.Address) {
	s.accBalance[address] = uint256.NewInt(0)
}

func (s *StateDb) SubBalance(address common.Address, u *uint256.Int, t tracing.BalanceChangeReason) {
	balance := s.GetBalance(address)
	newBalance := balance.Sub(balance, u)
	s.accBalance[address] = newBalance
}

func (s *StateDb) AddBalance(address common.Address, u *uint256.Int, t tracing.BalanceChangeReason) {
	balance := s.GetBalance(address)
	newBalance := balance.Add(balance, u)
	s.accBalance[address] = newBalance
}

func (s *StateDb) GetBalance(address common.Address) *uint256.Int {
	balance, ok := s.accBalance[address]
	if !ok {
		dbBalance := s.innerdb.GetAccountBalance(evmaas.BytesToAddress(address[:]))
		balance, ok = uint256.FromBig(dbBalance)
	}
	if balance == nil {
		balance = uint256.NewInt(0)
	}
	return balance
}

func (s *StateDb) GetNonce(address common.Address) uint64 {
	return 0
}

func (s *StateDb) SetNonce(address common.Address, u uint64) {
	return
}

func (s *StateDb) GetCodeHash(address common.Address) common.Hash {
	code := s.GetCode(address)
	h := sha3.NewLegacyKeccak256().Sum(code)
	return common.BytesToHash(h)
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
	return
}

func (s *StateDb) SubRefund(u uint64) {
	return
}

func (s *StateDb) GetRefund() uint64 {
	return 0
}

func (s *StateDb) GetCommittedState(address common.Address, hash common.Hash) common.Hash {
	return s.GetState(address, hash)
}

func (s *StateDb) GetState(address common.Address, hash common.Hash) common.Hash {
	kv, ok := s.state[address]
	if !ok {
		kaddr := evmaas.BytesToAddress(address[:])
		v, err := s.innerdb.GetState(kaddr, hash[:])
		if err != nil {
			return common.Hash{}
		}
		return common.BytesToHash(v)
	}
	v, ok := kv[string(hash[:])]
	if ok {
		return common.BytesToHash(v)

	}
	return common.Hash{}
}

func (s *StateDb) SetState(address common.Address, hash common.Hash, hash2 common.Hash) {
	kv, ok := s.state[address]
	if !ok {
		kv = make(map[string][]byte)
	}
	kv[string(hash[:])] = hash2[:]
	s.state[address] = kv
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
	return false
}

func (s *StateDb) Selfdestruct6780(address common.Address) {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) Exist(address common.Address) bool {
	return true
}

func (s *StateDb) Empty(address common.Address) bool {
	//TODO implement me
	panic("implement me")
}

func (s *StateDb) AddressInAccessList(addr common.Address) bool {
	return true
}

func (s *StateDb) SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool) {
	return true, true
}

func (s *StateDb) AddAddressToAccessList(addr common.Address) {
	return
}

func (s *StateDb) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	return
}

func (s *StateDb) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	return
}

func (s *StateDb) RevertToSnapshot(i int) {
	return
}

func (s *StateDb) Snapshot() int {
	return 0
}

func (s *StateDb) AddLog(log *types.Log) {
	s.logs = append(s.logs, log)
}

func (s *StateDb) AddPreimage(hash common.Hash, bytes []byte) {
	//TODO implement me
	panic("implement me")
}

var _ vm.StateDB = (*StateDb)(nil)
