package sealevm

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SealSC/SealEVM"
	"github.com/SealSC/SealEVM/common"
	"github.com/SealSC/SealEVM/crypto/hashes"
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/storage"
)

func logPrinter(logCache *storage.LogCache) {
	for _, logs := range *logCache {
		for _, l := range logs {
			for _, t := range l.Topics {
				fmt.Printf("topic:%x\n", t)
			}
			fmt.Printf("data :%x\n", l.Data)
			//fmt.Println("data as string:", string(l.Data))
		}
	}
}

// store result to memStorage
func storeResult(result *SealEVM.ExecuteResult, storage *memStorage) {
	for addr, cache := range result.StorageCache.CachedData {
		for key, v := range cache {
			log.Printf("save statedb: addr[%x],key[%x]=value[%x]\n", addr, key, v.Bytes())

			storage.putState([]byte(addr), []byte(key), v.Bytes())
		}
	}
}

// create a new evm
func newEvm(contractAddr, contractCode []byte, callData []byte, caller []byte, ms *memStorage) *SealEVM.EVM {

	hashInt := evmInt256.New(0)
	hashInt.SetBytes(contractAddr)

	//same contract code has same address in this example
	cNamespace := hashInt
	contract := environment.Contract{
		Namespace: cNamespace,
		Code:      contractCode,
		Hash:      hashInt,
	}

	var callHash [32]byte
	copy(callHash[12:], caller)
	callerInt, _ := common.HashBytesToEVMInt(callHash)
	evm := SealEVM.New(SealEVM.EVMParam{
		MaxStackDepth:  0,
		ExternalStore:  ms,
		ResultCallback: nil,
		Context: &environment.Context{
			Block: environment.Block{
				ChainID:    evmInt256.New(666),
				Coinbase:   evmInt256.New(0),
				Timestamp:  evmInt256.New(int64(time.Now().Second())),
				Number:     evmInt256.New(0),
				Difficulty: evmInt256.New(0),
				GasLimit:   evmInt256.New(10000000),
				Hash:       evmInt256.New(0),
			},
			Contract: contract,
			Transaction: environment.Transaction{
				Origin:   callerInt,
				GasPrice: evmInt256.New(1),
				GasLimit: evmInt256.New(10000000),
			},
			Message: environment.Message{
				Caller: callerInt,
				Value:  evmInt256.New(0),
				Data:   callData,
			},
		},
	})

	return evm
}

func Case1() {
	//load SealEVM module
	SealEVM.Load()

	// deployCode is for deploy the contract.sol
	var deployCode, _ = hex.DecodeString("608060405234801561001057600080fd5b506102fc806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806361bc221a1461003b578063a792da5914610059575b600080fd5b610043610089565b604051610050919061010a565b60405180910390f35b610073600480360381019061006e9190610194565b61008f565b604051610080919061010a565b60405180910390f35b60005481565b600060016000808282546100a39190610210565b925050819055506000547f84fa11cd0353da7cf3201711842e07f8fdf6a488011edfc5b5d996318e339d5584846040516100de9291906102a2565b60405180910390a2600054905092915050565b6000819050919050565b610104816100f1565b82525050565b600060208201905061011f60008301846100fb565b92915050565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f8401126101545761015361012f565b5b8235905067ffffffffffffffff81111561017157610170610134565b5b60208301915083600182028301111561018d5761018c610139565b5b9250929050565b600080602083850312156101ab576101aa610125565b5b600083013567ffffffffffffffff8111156101c9576101c861012a565b5b6101d58582860161013e565b92509250509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061021b826100f1565b9150610226836100f1565b925082820190508082111561023e5761023d6101e1565b5b92915050565b600082825260208201905092915050565b82818337600083830152505050565b6000601f19601f8301169050919050565b60006102818385610244565b935061028e838584610255565b61029783610264565b840190509392505050565b600060208201905081810360008301526102bd818486610275565b9050939250505056fea2646970667358221220978e7b7b85089cdb4ce014907b80c3cb6683084e5591e9347e724cc2a65814be64736f6c63430008120033")

	// call data of increaseFor("example")
	var callIncreaseFor, _ = hex.DecodeString("a792da59000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000076578616d706c6500000000000000000000000000000000000000000000000000")

	// call data of Counter()
	var callCounter, _ = hex.DecodeString("61bc221a")

	// example caller
	var caller, _ = hex.DecodeString("0b0b")

	//create memStorage
	ms := &memStorage{}
	ms.storage = make(map[string][]byte)
	ms.contracts = make(map[string][]byte)
	contractHash := hashes.Keccak256(deployCode)
	//deploy contract
	evm := newEvm(contractHash, deployCode, nil, caller, ms)
	ret, err := evm.ExecuteContract(false)

	//check error
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	//result data of ret is the deployed code of example contract
	contractCode := ret.ResultData

	//call Counter() to get current counter's value
	evm = newEvm(contractHash, contractCode, callCounter, caller, ms)
	ret, _ = evm.ExecuteContract(false)

	//result of Counter()
	fmt.Println("counter: ", hex.EncodeToString(ret.ResultData))

	//call increaseFor("example")
	evm = newEvm(contractHash, contractCode, callIncreaseFor, caller, ms)
	ret, _ = evm.ExecuteContract(false)

	//store the result to ms
	storeResult(&ret, ms)

	//the event logs
	logPrinter(&ret.StorageCache.Logs)

	//call Counter to get counter's value after increase
	evm = newEvm(contractHash, contractCode, callCounter, caller, ms)
	ret, err = evm.ExecuteContract(false)

	//result of Counter()
	fmt.Println("counter: ", hex.EncodeToString(ret.ResultData))
}
