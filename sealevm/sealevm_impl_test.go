package sealevm

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/studyzy/evmaas"
)

func TestSealEvmImpl_Erc20(t *testing.T) {
	impl := NewSealEVMImpl()

	//create memStorage
	db := evmaas.NewMemStateDB()

	// deployCode is for deploy the contract.sol
	bin, _ := os.ReadFile("../testdata/erc20/erc20.bin")
	var deployCode, _ = hex.DecodeString(string(bin))

	fmt.Println("安装合约")
	var userA, _ = hex.DecodeString("ab108fc6c3850e01cee01e419d07f097186c3982")
	var userB, _ = hex.DecodeString("ce2355fcfcb26414a254f28404c6040d0d4559c2")
	tx := evmaas.Transaction{
		TxHash:   userA,
		From:     evmaas.NewAddress("ab108fc6c3850e01cee01e419d07f097186c3982"),
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
	//保存状态数据
	for contract, v := range ret.StateChanges {
		for k, v := range v {
			db.PutState(contract, []byte(k), v.([]byte))
		}
	}

	//查询A的余额
	fmt.Println("查询A的余额")
	var userABalance, _ = hex.DecodeString("70a08231000000000000000000000000ab108fc6c3850e01cee01e419d07f097186c3982")
	tx = evmaas.Transaction{
		TxHash:   userA,
		From:     evmaas.NewAddress("ab108fc6c3850e01cee01e419d07f097186c3982"),
		To:       evmaas.NewAddress(contractAddr),
		Value:    nil,
		Gas:      100000,
		GasPrice: nil,
		Data:     userABalance,
	}
	ret, err = impl.ExecuteContract(tx, db, block)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("balance: ", hex.EncodeToString(ret.ReturnData))
	require.Equal(t, "0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000", hex.EncodeToString(ret.ReturnData))
	fmt.Println("A转账B")
	var transfer, _ = hex.DecodeString("a9059cbb000000000000000000000000ce2355fcfcb26414a254f28404c6040d0d4559c20000000000000000000000000000000000000000000000000000000000000064")
	tx = evmaas.Transaction{
		TxHash:   userA,
		From:     evmaas.NewAddress("ab108fc6c3850e01cee01e419d07f097186c3982"),
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
	for contract, v := range ret.StateChanges {
		for k, v := range v {
			db.PutState(contract, []byte(k), v.([]byte))
		}
	}

	fmt.Println("查询A的余额")
	tx = evmaas.Transaction{
		TxHash:   userA,
		From:     evmaas.NewAddress("ab108fc6c3850e01cee01e419d07f097186c3982"),
		To:       evmaas.NewAddress(contractAddr),
		Value:    nil,
		Gas:      100000,
		GasPrice: nil,
		Data:     userABalance,
	}
	ret, err = impl.ExecuteContract(tx, db, block)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("balance: ", hex.EncodeToString(ret.ReturnData))
	require.Equal(t, "0000000000000000000000000000000000000000033b2e3c9fd0803ce7ffff9c", hex.EncodeToString(ret.ReturnData))

}
func logPrint(result *evmaas.ExecutionResult) {
	for _, log := range result.Events {

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

func TestSealEvmImpl_OutofGas(t *testing.T) {
	impl := NewSealEVMImpl()

	//create memStorage
	db := evmaas.NewMemStateDB()

	// deployCode is for deploy the contract.sol
	bin, _ := os.ReadFile("../testdata/erc20/erc20.bin")
	var deployCode, _ = hex.DecodeString(string(bin))

	fmt.Println("安装合约")
	var userA, _ = hex.DecodeString("ab108fc6c3850e01cee01e419d07f097186c3982")
	var userB, _ = hex.DecodeString("ce2355fcfcb26414a254f28404c6040d0d4559c2")
	tx := evmaas.Transaction{
		TxHash:   userA,
		From:     evmaas.NewAddress("ab108fc6c3850e01cee01e419d07f097186c3982"),
		To:       evmaas.Address{},
		Value:    nil,
		Gas:      10,
		GasPrice: nil,
		Data:     deployCode,
	}
	block := evmaas.Block{
		BlockHash:  userB,
		Number:     100,
		Timestamp:  uint64(time.Now().Unix()),
		Difficulty: nil,
		GasLimit:   10,
	}
	//deploy contract

	_, err := impl.InstallContract(tx, db, block)
	assert.Error(t, err)
	fmt.Println(err.Error())
}
