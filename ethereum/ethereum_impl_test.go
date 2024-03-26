package ethereum

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/studyzy/evmaas"
)

func TestInstallErc20(t *testing.T) {

	impl := NewEthereumImpl()
	var userA, _ = hex.DecodeString("c3b4d19da33e3934c9e5383d9e38bccd591e2287")
	addressA := evmaas.BytesToAddress(userA)
	var userB, _ = hex.DecodeString("ce2355fcfcb26414a254f28404c6040d0d4559c2")
	//create memStorage
	db := evmaas.NewMemStateDB()

	// deployCode is for deploy the contract.sol
	bin, _ := os.ReadFile("./testdata/erc20.bin")
	var deployCode, _ = hex.DecodeString(string(bin))
	//初始化一定的账户余额
	db.SetAccountBalance(addressA, big.NewInt(math.MaxInt64))

	fmt.Println("安装合约")

	tx := evmaas.Transaction{
		TxHash:   userA,
		From:     addressA,
		To:       evmaas.Address{},
		Value:    nil,
		Gas:      1000000,
		GasPrice: nil,
		Data:     deployCode,
	}
	block := evmaas.Block{
		BlockHash:  userB,
		Number:     100,
		Timestamp:  uint64(time.Now().Unix()),
		Difficulty: nil,
		GasLimit:   1000000,
	}
	//deploy contract

	ret, err := impl.InstallContract(tx, db, block)

	//check error
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	//store the result to ms
	printResult(ret)
	var contractAddr string
	//保存合约
	for contractAddrK, code := range ret.ContractCode {
		contractAddr = hex.EncodeToString(contractAddrK[:])
		db.SetContractCode(evmaas.NewAddress(contractAddr), code)
	}
	assert.Equal(t, "53d19b414a839c5589c59ca28a1c79d69d39efb2", contractAddr)
	//保存状态数据
	for contract, kv := range ret.StateChanges {
		for k, v := range kv {
			db.PutState(contract, []byte(k), v)
		}
	}

	//查询A的余额
	fmt.Println("查询A的余额")
	var userABalance, _ = hex.DecodeString("70a08231000000000000000000000000" + hex.EncodeToString(userA))
	tx = evmaas.Transaction{
		TxHash:   userA,
		From:     addressA,
		To:       evmaas.NewAddress(contractAddr),
		Value:    nil,
		Gas:      100000,
		GasPrice: nil,
		Data:     userABalance,
	}
	ret, err = impl.QueryContract(tx, db, block)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	fmt.Println("balance: ", hex.EncodeToString(ret.ReturnData))
	require.Equal(t, "0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000", hex.EncodeToString(ret.ReturnData))
	fmt.Println("A转账B")
	var transfer, _ = hex.DecodeString("a9059cbb000000000000000000000000ce2355fcfcb26414a254f28404c6040d0d4559c20000000000000000000000000000000000000000000000000000000000000064")
	tx = evmaas.Transaction{
		TxHash:   userA,
		From:     addressA,
		To:       evmaas.NewAddress(contractAddr),
		Value:    nil,
		Gas:      100000,
		GasPrice: nil,
		Data:     transfer,
	}
	ret, err = impl.ExecuteContract(tx, db, block)
	if err != nil {
		fmt.Println(err.Error())
	}
	//the event logs
	printResult(ret)
	//保存状态数据
	for contract, kv := range ret.StateChanges {
		for k, v := range kv {
			db.PutState(contract, []byte(k), v)
		}
	}

	fmt.Println("查询A的余额")
	tx = evmaas.Transaction{
		TxHash:   userA,
		From:     addressA,
		To:       evmaas.NewAddress(contractAddr),
		Value:    nil,
		Gas:      100000,
		GasPrice: nil,
		Data:     userABalance,
	}
	ret, err = impl.QueryContract(tx, db, block)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("balance: ", hex.EncodeToString(ret.ReturnData))
	require.Equal(t, "0000000000000000000000000000000000000000033b2e3c9fd0803ce7ffff9c", hex.EncodeToString(ret.ReturnData))

}
func logPrint(result *evmaas.ExecutionResult) {
	for _, log := range result.Events {
		fmt.Printf("Contract:%x\n", log.ContractAddress)
		for _, t := range log.Topics {
			fmt.Printf("topic:%x\n", t)
		}
		fmt.Printf("data :%x\n", log.Data)
		//fmt.Println("data as string:", string(l.Data))
	}
}

func printResult(result *evmaas.ExecutionResult) {
	fmt.Println("Success:", result.Success)
	fmt.Println("ReturnData:", hex.EncodeToString(result.ReturnData))
	fmt.Println("StateChanges:", result.StateChanges)
	fmt.Println("Balance:", result.Balance)
	fmt.Println("ContractCode:", result.ContractCode)
	fmt.Println("GasUsed:", result.GasUsed)
	logPrint(result)
}
