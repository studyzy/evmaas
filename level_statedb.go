package evmaas

import (
	"encoding/json"
	"log"
	"math/big"

	"github.com/syndtr/goleveldb/leveldb"
)

type LevelStateDB struct {
	db *leveldb.DB
}

func NewLevelStateDB(path string) *LevelStateDB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		panic(err)

	}
	return &LevelStateDB{db: db}
}

const (
	accBalancePrefix = "a"
	codePrefix       = "c"
	statePrefix      = "s"
	blockPrefix      = "b"
	receiptPrefix    = "r"
)

func (m *LevelStateDB) GetAccountBalance(address Address) *big.Int {
	key := append([]byte{}, accBalancePrefix...)
	key = append(key, address[:]...)
	data, err := m.db.Get(key, nil)
	if err != nil {
		return nil
	}
	return new(big.Int).SetBytes(data)
}

func (m *LevelStateDB) SetAccountBalance(address Address, balance *big.Int) {
	key := append([]byte{}, accBalancePrefix...)
	key = append(key, address[:]...)
	log.Printf("SetAccountBalance: addr[%x]=balance[%d]\n", address, balance)
	m.db.Put(key, balance.Bytes(), nil)
}

func (m *LevelStateDB) GetContractCode(address Address) ([]byte, error) {
	key := append([]byte{}, codePrefix...)
	key = append(key, address[:]...)
	data, err := m.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *LevelStateDB) SetContractCode(address Address, code []byte) {
	key := append([]byte{}, codePrefix...)
	key = append(key, address[:]...)
	m.db.Put(key, code, nil)
}

// DeleteContractCode delete contract code
func (m *LevelStateDB) DeleteContractCode(address Address) {
	key := append([]byte{}, codePrefix...)
	key = append(key, address[:]...)
	m.db.Delete(key, nil)

}

func (m *LevelStateDB) GetState(address Address, key []byte) ([]byte, error) {
	k := append([]byte{}, statePrefix...)
	k = append(k, address[:]...)
	k = append(k, key...)
	data, err := m.db.Get(k, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *LevelStateDB) PutState(address Address, key []byte, value []byte) {
	log.Printf("PutState: addr[%x],key[%x]=value[%x]\n", address, key, value)
	k := append([]byte{}, statePrefix...)
	k = append(k, address[:]...)
	k = append(k, key...)
	m.db.Put(k, value, nil)
}

func (m *LevelStateDB) GetBlockHash(number uint64) ([]byte, error) {
	key := append([]byte{}, blockPrefix...)
	key = append(key, []byte{byte(number)}...)
	data, err := m.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// PutBlockHash
func (m *LevelStateDB) PutBlockHash(number uint64, hash []byte) {
	key := append([]byte{}, blockPrefix...)
	key = append(key, []byte{byte(number)}...)
	m.db.Put(key, hash, nil)
}

// PutReceipt
func (m *LevelStateDB) PutReceipt(txHash []byte, receipt *Receipt) {
	key := append([]byte{}, receiptPrefix...)
	key = append(key, txHash...)
	value, _ := json.Marshal(receipt)
	log.Printf("PutReceipt: key[%x]:\n%s", key, receipt.ToString())
	m.db.Put(key, value, nil)
}

// GetReceipt
func (m *LevelStateDB) GetReceipt(txHash []byte) (*Receipt, error) {
	key := append([]byte{}, receiptPrefix...)
	key = append(key, txHash...)
	data, err := m.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	receipt := &Receipt{}
	err = json.Unmarshal(data, receipt)
	if err != nil {
		return nil, err
	}
	return receipt, nil

}

var _ StateDB = &LevelStateDB{}
