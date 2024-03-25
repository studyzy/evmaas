package polygon

import (
	"errors"
	"math/big"

	"github.com/0xPolygon/polygon-edge/chain"
	"github.com/0xPolygon/polygon-edge/state/runtime"
	"github.com/0xPolygon/polygon-edge/state/runtime/tracer/calltracer"
	"github.com/0xPolygon/polygon-edge/types"
	"github.com/studyzy/evmaas"
	"golang.org/x/crypto/sha3"
)

type MemHost struct {
	innerdb        evmaas.StateDB
	accountBalance map[types.Address]*big.Int
	contracts      map[types.Address][]byte
	state          map[types.Address]map[types.Hash]types.Hash
	logs           []evmaas.EventLog
	txContext      runtime.TxContext
}

func NewMemHost(db evmaas.StateDB, txContext runtime.TxContext) *MemHost {
	return &MemHost{
		innerdb:        db,
		accountBalance: make(map[types.Address]*big.Int),
		contracts:      make(map[types.Address][]byte),
		state:          make(map[types.Address]map[types.Hash]types.Hash),
		txContext:      txContext,
	}
}

func (m *MemHost) AccountExists(addr types.Address) bool {
	balance := m.innerdb.GetAccountBalance(evmaas.BytesToAddress(addr[:]))
	return balance != nil
}

func (m *MemHost) GetStorage(addr types.Address, key types.Hash) types.Hash {
	kv, ok := m.state[addr]
	if !ok {
		value, err := m.innerdb.GetState(evmaas.BytesToAddress(addr[:]), key[:])
		if err != nil || len(value) == 0 {
			return types.Hash{}

		}
		return types.Hash(value)
	}
	value, ok := kv[key]
	if !ok {
		return types.Hash{}
	}
	return value
}

func (m *MemHost) SetStorage(addr types.Address, key types.Hash, value types.Hash, config *chain.ForksInTime) runtime.StorageStatus {
	kv, ok := m.state[addr]
	if !ok {
		m.state[addr] = map[types.Hash]types.Hash{
			key: value,
		}
		return runtime.StorageAdded
	}
	_, ok = kv[key]
	if !ok {
		kv[key] = value
		return runtime.StorageAdded
	}
	kv[key] = value
	return runtime.StorageModified
}

func (m *MemHost) SetState(addr types.Address, key types.Hash, value types.Hash) {
	kv, ok := m.state[addr]
	if !ok {
		m.state[addr] = map[types.Hash]types.Hash{
			key: value,
		}
		return

	}
	kv[key] = value
}

func (m *MemHost) SetNonPayable(nonPayable bool) {
	//TODO implement me
	panic("implement me")
}

func (m *MemHost) GetBalance(addr types.Address) *big.Int {
	balance, ok := m.accountBalance[addr]
	if !ok {
		balance = m.innerdb.GetAccountBalance(evmaas.BytesToAddress(addr[:]))
	}
	return balance
}

func (m *MemHost) GetCodeSize(addr types.Address) int {
	code := m.GetCode(addr)
	return len(code)
}

func (m *MemHost) GetCodeHash(addr types.Address) types.Hash {
	code := m.GetCode(addr)
	//sha3 code
	h := sha3.NewLegacyKeccak256().Sum(code)
	return types.Hash(h)
}

func (m *MemHost) GetCode(addr types.Address) []byte {
	code, ok := m.contracts[addr]
	if !ok {
		var err error
		code, err = m.innerdb.GetContractCode(evmaas.BytesToAddress(addr[:]))
		if err != nil {
			return nil
		}
	}
	return code
}

func (m *MemHost) Selfdestruct(addr types.Address, beneficiary types.Address) {
	//TODO implement me
	panic("implement me")
}

func (m *MemHost) GetTxContext() runtime.TxContext {
	return m.txContext
}

func (m *MemHost) GetBlockHash(number int64) types.Hash {
	//TODO implement me
	panic("implement me")
}

func (m *MemHost) EmitLog(addr types.Address, topics []types.Hash, data []byte) {

	ts := make([][]byte, len(topics))
	for i := 0; i < len(topics); i++ {
		ts[i] = topics[i][:]
	}
	log := evmaas.EventLog{
		ContractAddress: evmaas.BytesToAddress(addr[:]),
		Topics:          ts,
		Data:            data,
	}
	m.logs = append(m.logs, log)
}

func (m *MemHost) Callx(contract *runtime.Contract, host runtime.Host) *runtime.ExecutionResult {
	//TODO implement me
	panic("implement me")
}

func (m *MemHost) Empty(addr types.Address) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MemHost) GetNonce(addr types.Address) uint64 {
	return 0
}

func (m *MemHost) Transfer(from types.Address, to types.Address, amount *big.Int) error {
	fromBalance := m.GetBalance(from)
	toBalance := m.GetBalance(to)
	if fromBalance.Cmp(amount) < 0 {
		return errors.New("insufficient balance")
	}
	fromBalance.Sub(fromBalance, amount)
	toBalance.Add(toBalance, amount)
	m.accountBalance[from] = fromBalance
	m.accountBalance[to] = toBalance
	return nil
}

func (m *MemHost) GetTracer() runtime.VMTracer {
	return &calltracer.CallTracer{}
}

func (m *MemHost) GetRefund() uint64 {
	//TODO implement me
	panic("implement me")
}
