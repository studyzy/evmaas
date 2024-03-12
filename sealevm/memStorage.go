package sealevm

import (
	"errors"
	"log"
	"math/big"
	"time"

	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/studyzy/evmaas"
)

// external storage for example
type memStorage struct {
	statedb evmaas.StateDB
	//storage   key: contractAddr+"#"+key, value: value
	storage   map[string][]byte
	contracts map[string][]byte
}

func newMemStorage(statedb evmaas.StateDB) *memStorage {
	return &memStorage{statedb: statedb, storage: make(map[string][]byte), contracts: make(map[string][]byte)}
}

func (r *memStorage) getAccountBalance(address evmaas.Address) *big.Int {
	return big.NewInt(1000000000000000000)
}

func (r *memStorage) getContractCode(address evmaas.Address) ([]byte, error) {
	code, exit := r.contracts[string(address[:])]
	if exit {
		return code, nil
	}
	return r.statedb.GetContractCode(address)
}

func (r *memStorage) getState(address, key []byte) ([]byte, error) {
	v, exit := r.storage[string(address[:])+"#"+string(key)]
	if exit {
		return v, nil

	}
	return r.statedb.GetState(evmaas.Address(address), key)
}
func (r *memStorage) putState(address, key []byte, value []byte) error {
	k := string(address[:]) + "#" + string(key)
	r.storage[k] = value
	return nil
}

func (r *memStorage) getBlockHash(number uint64) ([]byte, error) {
	return r.statedb.GetBlockHash(number)
}

func (r *memStorage) GetBalance(address *evmInt256.Int) (*evmInt256.Int, error) {
	b := r.getAccountBalance(evmaas.NewAddress(address.AsStringKey()))
	return evmInt256.FromBigInt(b), nil
}

func (r *memStorage) CanTransfer(from, to, val *evmInt256.Int) bool {
	return true
}

func (r *memStorage) GetCode(address *evmInt256.Int) ([]byte, error) {
	addr := evmaas.NewAddress(address.AsStringKey())
	return r.getContractCode(addr)
}

func (r *memStorage) GetCodeSize(address *evmInt256.Int) (*evmInt256.Int, error) {
	code, err := r.GetCode(address)
	if err != nil {
		return nil, errors.New("no code for: 0x" + address.Text(16))
	}
	return evmInt256.New(int64(len(code))), nil
}

func (r *memStorage) GetCodeHash(address *evmInt256.Int) (*evmInt256.Int, error) {
	code, err := r.GetCode(address)
	if err != nil {
		return nil, errors.New("no code for: 0x" + address.Text(16))
	}
	h := evmaas.GetCodeHash(code)
	var i big.Int
	i.SetBytes(h)
	return evmInt256.FromBigInt(&i), nil
}

func (r *memStorage) GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error) {
	height := block.Uint64()
	bhash, err := r.getBlockHash(height)
	if err != nil {
		return nil, err

	}
	var i big.Int
	i.SetBytes(bhash)
	return evmInt256.FromBigInt(&i), nil
}

func (r *memStorage) GetChainID() (*evmInt256.Int, error) {
	return evmInt256.New(666), nil
}

func (r *memStorage) CreateAddress(caller *evmInt256.Int, tx environment.Transaction) *evmInt256.Int {
	return evmInt256.New(time.Now().UnixNano())
}

func (r *memStorage) CreateFixedAddress(caller *evmInt256.Int, salt *evmInt256.Int, tx environment.Transaction) *evmInt256.Int {
	return evmInt256.New(time.Now().UnixNano())
}

func (r *memStorage) Load(n string, k string) (*evmInt256.Int, error) {
	ret := evmInt256.New(0)
	//addr := evmaas.NewAddress(n)
	val, err := r.getState([]byte(n), []byte(k))
	if err != nil || len(val) == 0 {
		log.Printf("query statedb: %x[%x] but not found\n", n, k)
	}

	ret.SetBytes(val)
	log.Printf("query statedb: %x[%x] = %x\n", n, k, ret.Bytes())

	return ret, nil
}

func (r *memStorage) NewContract(n string, code []byte) error {
	r.contracts[n] = code
	return nil
}
