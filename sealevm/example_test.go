package sealevm

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/SealSC/SealEVM"
	"github.com/SealSC/SealEVM/crypto/hashes"
	"github.com/studyzy/evmaas"
)

func TestCase1(t *testing.T) {
	Case1()
}

func TestDirectErc20(t *testing.T) {
	//load SealEVM module
	SealEVM.Load()

	//create memStorage
	ms := newMemStorage(evmaas.NewMemStateDB())

	// deployCode is for deploy the contract.sol
	bin, _ := os.ReadFile("../testdata/erc20/erc20.bin")
	var deployCode, _ = hex.DecodeString(string(bin))

	//// call data of increaseFor("example")
	//var callIncreaseFor, _ = hex.DecodeString("a792da59000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000076578616d706c6500000000000000000000000000000000000000000000000000")
	//
	//// call data of Counter()
	//var callCounter, _ = hex.DecodeString("61bc221a")

	fmt.Println("安装合约")
	var userA, _ = hex.DecodeString("ab108fc6c3850e01cee01e419d07f097186c3982")
	//var userB, _ = hex.DecodeString("ce2355fcfcb26414a254f28404c6040d0d4559c2")
	//deploy contract
	contractHash := hashes.Keccak256(deployCode)
	evm := newEvm(contractHash, deployCode, nil, userA, ms)
	ret, err := evm.ExecuteContract(false)

	//check error
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	//store the result to ms
	storeResult(&ret, ms)
	logPrinter(&ret.StorageCache.Logs)
	//查询A的余额
	fmt.Println("查询A的余额")
	var userABalance, _ = hex.DecodeString("70a08231000000000000000000000000ab108fc6c3850e01cee01e419d07f097186c3982")

	//result data of ret is the deployed code of example contract
	contractCode := ret.ResultData
	evm = newEvm(contractHash, contractCode, userABalance, userA, ms)

	ret, err = evm.ExecuteContract(false)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	//result of Counter()
	fmt.Println("balance: ", hex.EncodeToString(ret.ResultData))
	fmt.Println("A转账B")
	var transfer, _ = hex.DecodeString("a9059cbb000000000000000000000000ce2355fcfcb26414a254f28404c6040d0d4559c20000000000000000000000000000000000000000000000000000000000000064")

	//call increaseFor("example")
	evm = newEvm(contractHash, contractCode, transfer, userA, ms)
	ret, _ = evm.ExecuteContract(false)

	//store the result to ms
	storeResult(&ret, ms)

	//the event logs
	logPrinter(&ret.StorageCache.Logs)
	fmt.Println("查询A的余额")
	//call
	evm = newEvm(contractHash, contractCode, userABalance, userA, ms)
	ret, err = evm.ExecuteContract(false)

	//result of Counter()
	fmt.Println("balance: ", hex.EncodeToString(ret.ResultData))
}
