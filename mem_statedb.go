package evmaas

import (
	"log"
	"math/big"
)

type MemStateDB struct {
	AccountBalance   map[Address]*big.Int
	ContractByteCode map[Address][]byte
	StateMap         map[Address]map[string][]byte
	BlockHashes      map[uint64][]byte
}

func NewMemStateDB() *MemStateDB {
	return &MemStateDB{
		AccountBalance:   make(map[Address]*big.Int),
		ContractByteCode: make(map[Address][]byte),
		StateMap:         make(map[Address]map[string][]byte),
		BlockHashes:      make(map[uint64][]byte),
	}
}

func (m *MemStateDB) GetAccountBalance(address Address) *big.Int {
	return m.AccountBalance[address]
}

func (m *MemStateDB) SetAccountBalance(address Address, balance *big.Int) {
	m.AccountBalance[address] = balance
}

func (m *MemStateDB) GetContractCode(address Address) ([]byte, error) {
	return m.ContractByteCode[address], nil
}

func (m *MemStateDB) SetContractCode(address Address, code []byte) {
	m.ContractByteCode[address] = code
}

func (m *MemStateDB) GetState(address Address, key []byte) ([]byte, error) {
	return m.StateMap[address][string(key)], nil
}

func (m *MemStateDB) PutState(address Address, key []byte, value []byte) {
	log.Printf("PutState: addr[%x],key[%x]=value[%x]\n", address, key, value)
	if m.StateMap[address] == nil {
		m.StateMap[address] = make(map[string][]byte)
	}
	m.StateMap[address][string(key)] = value
}

func (m *MemStateDB) GetBlockHash(number uint64) ([]byte, error) {
	return m.BlockHashes[number], nil
}

var _ StateDB = &MemStateDB{}
