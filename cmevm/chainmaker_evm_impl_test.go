package cmevm

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"testing"
	"time"

	"chainmaker.org/chainmaker/common/v3/crypto"
	"chainmaker.org/chainmaker/common/v3/crypto/asym"
	"chainmaker.org/chainmaker/common/v3/evmutils"
	"chainmaker.org/chainmaker/common/v3/evmutils/abi"
	"chainmaker.org/chainmaker/logger/v3"
	acPb "chainmaker.org/chainmaker/pb-go/v3/accesscontrol"
	commonPb "chainmaker.org/chainmaker/pb-go/v3/common"
	configPb "chainmaker.org/chainmaker/pb-go/v3/config"
	storePb "chainmaker.org/chainmaker/pb-go/v3/store"
	"chainmaker.org/chainmaker/pb-go/v3/syscontract"
	vmPb "chainmaker.org/chainmaker/pb-go/v3/vm"
	"chainmaker.org/chainmaker/protocol/v3"
	"chainmaker.org/chainmaker/protocol/v3/mock"
	"chainmaker.org/chainmaker/utils/v3"
	evm "chainmaker.org/chainmaker/vm-evm/v3"
	evmGo "chainmaker.org/chainmaker/vm-evm/v3/evm-go"
	"chainmaker.org/chainmaker/vm-evm/v3/evm-go/environment"
	"chainmaker.org/chainmaker/vm-evm/v3/evm-go/storage"
	"chainmaker.org/chainmaker/vm-evm/v3/test"
	"github.com/golang/mock/gomock"
)

const (
	chainId = "chain01"
	path    = "../testdata/erc20/"
	name    = "erc20"
)

func TestInstallContract(t *testing.T) {
	parameters := make(map[string][]byte)
	contractId, txContext, byteCode := InitContextTest(commonPb.RuntimeType_EVM, t)

	runtimeInstance := &evm.RuntimeInstance{
		ChainId:      chainId,
		Log:          logger.GetLogger(logger.MODULE_VM),
		TxSimContext: txContext,
	}
	//调用合约
	abiJson, err := ioutil.ReadFile(path + name + ".abi")
	if err != nil {
		t.Errorf("Read ABI file failed, err:%v", err.Error())
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		t.Errorf("constrcut ABI obj failed, err:%v", err.Error())
	}
	args := []interface{}{}
	dataByte, err := myAbi.Pack("", args)
	if err != nil {
		t.Errorf("create ABI data failed, err:%v", err.Error())
	}
	initMethod := "init_contract"

	dataString := hex.EncodeToString(dataByte)
	byteCode, _ = hex.DecodeString(string(byteCode))
	test.BaseParam(parameters)
	parameters[protocol.ContractCreatorPkParam] = contractId.Creator.MemberInfo
	parameters[protocol.ContractSenderPkParam] = txContext.GetSender().MemberInfo
	parameters[protocol.ContractEvmParamKey] = []byte(dataString)
	contractResult, _ := runtimeInstance.Invoke(contractId, initMethod, byteCode, parameters, txContext, 0)
	t.Logf("contractResult:%v", contractResult)
}

//func (vm *ChainMakerEvmImpl) TestExecuteContract(t *testing.T) {
//	instance := &evm.RuntimeInstance{
//		Method:        "",
//		ChainId:       "",
//		Contract:      nil,
//		Log:           nil,
//		TxSimContext:  nil,
//		ContractEvent: nil,
//	}
//}

var bytes []byte
var fromAddr = []byte("0x123456789abcdef")
var privateKey, _ = asym.GenerateKeyPair(crypto.ECC_Secp256k1)
var publicKey = privateKey.PublicKey()
var ContractName = "erc20"
var ByteCodeFile = "../testdata/erc20/erc20.bin"

// InitContextTest 初始化上下文和wasm字节码
func InitContextTest(runtimeType commonPb.RuntimeType, t *testing.T) (*commonPb.Contract, *TxContextMockTest, []byte) {
	if bytes == nil {
		bytes, _ = ioutil.ReadFile(ByteCodeFile)
		fmt.Printf("byteCode file size=%d\n", len(bytes))
	}

	//addr := hex.EncodeToString(evmutils.Keccak256([]byte(ContractName)))[24:]
	addr, _ := utils.NameToAddrStr(ContractName, configPb.AddrType_ETHEREUM, 3010000)

	contractId := commonPb.Contract{
		Name:        ContractName,
		Version:     "v1",
		RuntimeType: runtimeType,
		Status:      commonPb.ContractStatus_NORMAL,
		Creator: &acPb.MemberFull{
			OrgId:      "",
			MemberType: acPb.MemberType_ADDR,
			MemberInfo: fromAddr,
		},
		Address: addr,
	}

	sender := &acPb.Member{
		OrgId:      "",
		MemberType: acPb.MemberType_ADDR,
		MemberInfo: fromAddr,
		//IsFullCert: true,
	}
	//ret0, _ := ret[0].(*common.ContractResult)
	//ret1, _ := ret[1].(protocol.ExecOrderTxType)
	//ret2, _ := ret[2].(common.TxStatusCode)

	res := commonPb.ContractResult{
		Code:    0,
		GasUsed: 1000,
		Result:  []byte("123456789abcdef"),
		Message: "OK",
	}

	ctrl := gomock.NewController(t)
	vmManager := mock.NewMockVmManager(ctrl)
	vmManager.EXPECT().RunContract(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any()).Return(&res, protocol.ExecOrderTxTypeNormal, commonPb.TxStatusCode_SUCCESS).AnyTimes()

	member := mock.NewMockMember(ctrl)
	member.EXPECT().GetUid().Return(string(fromAddr)).AnyTimes()
	member.EXPECT().GetPk().Return(publicKey).AnyTimes()
	member.EXPECT().GetMember().Return(sender, nil).AnyTimes()
	ac := mock.NewMockAccessControlProvider(ctrl)
	ac.EXPECT().NewMember(gomock.Any()).Return(member, nil).AnyTimes()

	txContext := TxContextMockTest{
		lock:      &sync.Mutex{},
		vmManager: nil,
		hisResult: make([]*callContractResult, 0),
		creator:   sender,
		sender:    sender,
		cacheMap:  make(map[string][]byte),
		ac:        ac,
	}

	data, _ := contractId.Marshal()
	key := utils.GetContractDbKey(contractId.Name)
	err := txContext.Put(syscontract.SystemContract_CONTRACT_MANAGE.String(), key, data)
	if err != nil {
		panic(err)
	}
	//versionKey := []byte(protocol.ContractVersion + ContractName)
	//runtimeTypeKey := []byte(protocol.ContractRuntimeType + ContractName)
	//versionedByteCodeKey := append([]byte(protocol.ContractByteCode+ContractName), []byte(contractId.Version)...)
	//
	//txContext.Put(syscontract.SystemContract_CONTRACT_MANAGE.String(), versionedByteCodeKey, bytes)
	//txContext.Put(syscontract.SystemContract_CONTRACT_MANAGE.String(),
	//versionKey, []byte(contractId.Version))

	//txContext.Put(syscontract.SystemContract_CONTRACT_MANAGE.String(),
	//runtimeTypeKey, []byte(strconv.Itoa(int(runtimeType))))

	return &contractId, &txContext, bytes
}

// TxContextMockTest simTxContext mock for ut
type TxContextMockTest struct {
	lock          *sync.Mutex
	vmManager     protocol.VmManager
	gasUsed       uint64 // only for callContract
	currentDepth  int
	currentResult []byte
	hisResult     []*callContractResult

	sender   *acPb.Member
	creator  *acPb.Member
	cacheMap map[string][]byte
	ac       protocol.AccessControlProvider
}

// GetSnapshot add next time
//
//	@Description:
//	@receiver s
//	@return protocol.Snapshot
func (s *TxContextMockTest) GetSnapshot() protocol.Snapshot {
	//TODO implement me
	panic("implement me")
}

// GetConsensusStateWrapper 获得共识状态
// @return protocol.ConsensusStateWrapper
func (s *TxContextMockTest) GetConsensusStateWrapper() protocol.ConsensusStateWrapper {
	//TODO implement me
	panic("implement me")
}

// SubtractGas add next time
//
//	@Description:
//	@receiver s
//	@param gasUsed
//	@return error
func (s *TxContextMockTest) SubtractGas(gasUsed uint64) error {
	//TODO implement me
	panic("implement me")
}

// GetGasRemaining add next time
//
//	@Description:
//	@receiver s
//	@return uint64
func (s *TxContextMockTest) GetGasRemaining() uint64 {
	//TODO implement me
	panic("implement me")
}

// GetBlockFingerprint returns unique id for block
func (s *TxContextMockTest) GetBlockFingerprint() string {
	//TODO implement me
	panic("implement me")
}

// GetStrAddrFromPbMember calculate string address from pb Member
func (s *TxContextMockTest) GetStrAddrFromPbMember(pbMember *acPb.Member) (string, error) {
	//TODO implement me
	panic("implement me")
}

// GetNoRecord get no record
func (s *TxContextMockTest) GetNoRecord(contractName string, key []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

// GetBlockTimestamp get block timestamp mock
func (s *TxContextMockTest) GetBlockTimestamp() int64 {
	stamp := time.Now().UnixNano()
	return stamp
}

// GetContractByName get contract by name
func (s *TxContextMockTest) GetContractByName(name string) (*commonPb.Contract, error) {
	//panic("implement me")
	return &commonPb.Contract{Name: "1007"}, nil
}

// GetContractBytecode  get contract by code
func (s *TxContextMockTest) GetContractBytecode(name string) ([]byte, error) {
	//panic("implement me")
	return []byte(""), nil
}

// PutIntoReadSet put into read set
func (s *TxContextMockTest) PutIntoReadSet(contractName string, key []byte, value []byte) {
	panic("implement me")
}

// GetBlockVersion get block version
func (s *TxContextMockTest) GetBlockVersion() uint32 {
	return 3010000
}

// SetStateKvHandle set state kv handle
func (s *TxContextMockTest) SetStateKvHandle(i int32, iterator protocol.StateIterator) {
	panic("implement me")
}

// GetStateKvHandle get state kv handle
func (s *TxContextMockTest) GetStateKvHandle(i int32) (protocol.StateIterator, bool) {
	panic("implement me")
}

// PutRecord put record
func (s *TxContextMockTest) PutRecord(contractName string, value []byte, sqlType protocol.SqlType) {
	panic("implement me")
}

// Select select
func (s *TxContextMockTest) Select(name string, startKey []byte, limit []byte) (protocol.StateIterator, error) {
	panic("implement me")
}

// GetIterHandle get iterator handle
func (s *TxContextMockTest) GetIterHandle(index int32) (interface{}, bool) {
	panic("implement me")
}

// SetIterHandle set iterator handle
func (s *TxContextMockTest) SetIterHandle(index int32, iter interface{}) {
	panic("implement me")
}

// GetHistoryIterForKey get history iterator for key
func (s *TxContextMockTest) GetHistoryIterForKey(contractName string, key []byte) (protocol.KeyHistoryIterator, error) {
	panic("implement me")
}

// GetBlockProposer get block proposer
func (s *TxContextMockTest) GetBlockProposer() *acPb.Member {
	panic("implement me")
}

// SetStateSqlHandle set state sql handle
func (s *TxContextMockTest) SetStateSqlHandle(i int32, rows protocol.SqlRows) {
	panic("implement me")
}

// GetStateSqlHandle get state sql handle
func (s *TxContextMockTest) GetStateSqlHandle(i int32) (protocol.SqlRows, bool) {
	panic("implement me")
}

type callContractResult struct {
	contractName string
	method       string
	param        map[string][]byte
	deep         int
	gasUsed      uint64
	result       []byte
}

// Get get by key
func (s *TxContextMockTest) Get(name string, key []byte) ([]byte, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	k := string(key)
	if name != "" {
		k = name + "::" + k
	}
	//println("【get】 key:" + k)
	//fms.Println("【get】 key:", k, "val:", cacheMap[k])
	return s.cacheMap[k], nil
	//return nil,nil
	//data := "hello"
	//for i := 0; i < 70; i++ {
	//	for i := 0; i < 100; i++ {//1k
	//		data += "1234567890"
	//	}
	//}
	//return []byte(data), nil
}

// Put put kv
func (s *TxContextMockTest) Put(name string, key []byte, value []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	k := string(key)
	//v := string(value)
	if name != "" {
		k = name + "::" + k
	}
	//println("【put】 key:" + k)
	//fmt.Println("【put】 key:", k, "val:", value)
	s.cacheMap[k] = value
	return nil
}

// Del del kv
func (s *TxContextMockTest) Del(name string, key []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	k := string(key)
	//v := string(value)
	if name != "" {
		k = name + "::" + k
	}
	//println("【put】 key:" + k)
	s.cacheMap[k] = nil
	return nil
}

func callback(result *evmGo.ExecuteResult, err error) {

}

func (s *TxContextMockTest) upstreamAddress(contract *commonPb.Contract, txSimContext protocol.TxSimContext,
	parameters map[string][]byte) (*evmutils.Int, *evmutils.Int, *evmutils.Int, error) {
	ac, err := txSimContext.GetAccessControl()
	if err != nil {
		return nil, nil, nil, err
	}

	creator := &acPb.Member{
		OrgId:      s.GetSender().OrgId,
		MemberType: s.GetSender().MemberType,
		MemberInfo: s.GetSender().MemberInfo,
	}

	protocolCreator, err1 := ac.NewMember(creator)
	if err1 != nil {
		return nil, nil, nil, err1
	}

	cfg, err2 := txSimContext.GetBlockchainStore().GetLastChainConfig()
	if err2 != nil {
		return nil, nil, nil, err2
	}

	creatorAddr, err3 := utils.GetIntAddrFromMember(protocolCreator, cfg.Vm.AddrType)
	if err3 != nil {
		return nil, nil, nil, err3
	}
	origin := txSimContext.GetSender()

	protocolOrigin, err4 := ac.NewMember(origin)
	if err4 != nil {
		return nil, nil, nil, err4
	}

	originAddr, err5 := utils.GetIntAddrFromMember(protocolOrigin, cfg.Vm.AddrType)
	if err5 != nil {
		return nil, nil, nil, err5
	}

	if string(parameters[syscontract.CrossParams_CALL_TYPE.String()]) == syscontract.CallType_CROSS.String() {
		senderAddr := evmutils.FromHexString(string(parameters[syscontract.CrossParams_SENDER.String()]))
		return creatorAddr, originAddr, senderAddr, nil
	}

	return creatorAddr, originAddr, originAddr, nil
}

func (s *TxContextMockTest) invoke(contract *commonPb.Contract, method string, byteCode []byte,
	parameter map[string][]byte, gasUsed uint64) (*commonPb.ContractResult, protocol.ExecOrderTxType) {

	params := string(parameter[protocol.ContractEvmParamKey])
	isDeploy := false
	if method == protocol.ContractInitMethod || method == protocol.ContractUpgradeMethod {
		isDeploy = true
	}

	if evmutils.Has0xPrefix(params) {
		params = params[2:]
	}

	if len(params)%2 == 1 {
		params = "0" + params
	}
	messageData, _ := hex.DecodeString(params)

	if isDeploy {
		messageData = append(byteCode, messageData...)
		byteCode = messageData
	}

	addr, _ := utils.NameToAddrStr(contract.Name, configPb.AddrType_ETHEREUM, s.GetBlockVersion())
	contract.Address = addr
	address := evmutils.FromHexString(contract.Address)
	codeHash := evmutils.BytesDataToEVMIntHash(byteCode)
	eContract := environment.Contract{
		Address: address,
		Code:    byteCode,
		Hash:    codeHash,
	}

	creatorAddress, originAddress, senderAddress, _ := s.upstreamAddress(contract, s, parameter)
	gasLeft := protocol.GasLimit - gasUsed
	evmTransaction := environment.Transaction{
		TxHash:   []byte(s.GetTx().Payload.TxId),
		Origin:   originAddress,
		GasPrice: evmutils.New(protocol.EvmGasPrice),
		GasLimit: evmutils.New(int64(gasLeft)),
		BaseFee:  evmutils.New(0),
	}

	externalStore := &storage.ContractStorage{
		Ctx:       s,
		OutParams: storage.NewParamsCache(),
		Contract:  contract,
	}

	evm := evmGo.New(evmGo.EVMParam{
		MaxStackDepth:  protocol.EvmMaxStackDepth,
		ExternalStore:  externalStore,
		UpperStorage:   nil,
		ResultCallback: callback, //will be called as evm.resultNotify when evm.ExecuteContract() end
		Context: &environment.Context{
			Block: environment.Block{
				Coinbase:   creatorAddress, //proposer ski
				Timestamp:  evmutils.New(7),
				Number:     evmutils.New(int64(40)), // height
				Difficulty: evmutils.New(0),
				GasLimit:   evmutils.New(protocol.GasLimit),
			},
			Contract:    eContract,
			Transaction: evmTransaction,
			Message: environment.Message{
				Caller: senderAddress,
				Value:  evmutils.New(0),
				Data:   messageData,
			},
			Parameters: parameter,
			Cfg: environment.Config{
				AddrType: 2,
				ChainId:  s.GetTx().Payload.ChainId,
			},
		},
	})

	result, _ := evm.ExecuteContract(isDeploy)
	var contractResult commonPb.ContractResult
	contractResult.Code = 0
	contractResult.GasUsed = gasLeft - result.GasLeft
	contractResult.Result = result.ResultData
	return &contractResult, protocol.ExecOrderTxTypeNormal
}

// CallContract call contract
func (s *TxContextMockTest) CallContract(caller, contract *commonPb.Contract, method string, byteCode []byte,
	parameter map[string][]byte, gasUsed uint64, refTxType commonPb.TxType) (*commonPb.ContractResult,
	protocol.ExecOrderTxType, commonPb.TxStatusCode) {
	s.gasUsed = gasUsed
	s.currentDepth = s.currentDepth + 1
	if s.currentDepth > protocol.CallContractDepth {
		contractResult := &commonPb.ContractResult{
			Code:    uint32(1),
			Result:  nil,
			Message: fmt.Sprintf("CallContract too deep %d", s.currentDepth),
		}
		return contractResult, protocol.ExecOrderTxTypeNormal, commonPb.TxStatusCode_CONTRACT_TOO_DEEP_FAILED
	}
	if s.gasUsed > protocol.GasLimit {
		contractResult := &commonPb.ContractResult{
			Code:    uint32(1),
			Result:  nil,
			Message: fmt.Sprintf("There is not enough gas, gasUsed %d GasLimit %d ", gasUsed, int64(protocol.GasLimit)),
		}
		return contractResult, protocol.ExecOrderTxTypeNormal, commonPb.TxStatusCode_CONTRACT_FAIL
	}

	if len(byteCode) == 0 {
		dbByteCode, err := s.GetContractBytecode(contract.Name)
		if err != nil {
			return nil, protocol.ExecOrderTxTypeNormal, commonPb.TxStatusCode_CONTRACT_FAIL
		}
		byteCode = dbByteCode
	}

	//r, _, code := s.vmManager.RunContract(contract, method, byteCode, parameter, s, s.gasUsed, refTxType)
	contractResult, t := s.invoke(contract, method, byteCode, parameter, s.gasUsed)
	cr, _ := contractResult.Marshal()
	result := callContractResult{
		deep:         s.currentDepth,
		gasUsed:      s.gasUsed,
		result:       cr,
		contractName: contract.Name,
		method:       method,
		param:        parameter,
	}
	s.hisResult = append(s.hisResult, &result)
	s.currentResult = cr
	s.currentDepth = s.currentDepth - 1
	return contractResult, t, commonPb.TxStatusCode_SUCCESS
}

// GetCurrentResult Get current result
func (s *TxContextMockTest) GetCurrentResult() []byte {
	return s.currentResult
}

// GetTx get tx
func (s *TxContextMockTest) GetTx() *commonPb.Transaction {
	return &commonPb.Transaction{
		Payload: &commonPb.Payload{
			ChainId:        chainId,
			TxType:         commonPb.TxType_INVOKE_CONTRACT,
			TxId:           "12345678",
			Timestamp:      0,
			ExpirationTime: 0,
		},
		Result: nil,
	}
}

// GetBlockHeight get block height
func (*TxContextMockTest) GetBlockHeight() uint64 {
	return 7
}

// GetTxResult get tx result
func (s *TxContextMockTest) GetTxResult() *commonPb.Result {
	panic("implement me")
}

// SetTxResult set tx result
func (s *TxContextMockTest) SetTxResult(txResult *commonPb.Result) {
	panic("implement me")
}

// GetTxRWSet get tx read write set
func (TxContextMockTest) GetTxRWSet(runVmSuccess bool) *commonPb.TxRWSet {
	return &commonPb.TxRWSet{
		TxId:     "txId",
		TxReads:  nil,
		TxWrites: nil,
	}
}

// GetCreator get creator
func (s *TxContextMockTest) GetCreator(namespace string) *acPb.Member {
	return s.creator
}

// GetSender get sender
func (s *TxContextMockTest) GetSender() *acPb.Member {
	return s.sender
}

// GetBlockchainStore get block chain store
func (*TxContextMockTest) GetBlockchainStore() protocol.BlockchainStore {
	//protocol.BlockchainStore
	return &mockBlockchainStore{}
}

// GetLastChainConfig returns last chain config
func (*TxContextMockTest) GetLastChainConfig() *configPb.ChainConfig {
	panic("implement me")
}

// GetAccessControl get access control
func (s *TxContextMockTest) GetAccessControl() (protocol.AccessControlProvider, error) {
	return s.ac, nil
}

// GetChainNodesInfoProvider get chain nodes info provider
func (s *TxContextMockTest) GetChainNodesInfoProvider() (protocol.ChainNodesInfoProvider, error) {
	panic("implement me")
}

// GetTxExecSeq get tx execute sequence
func (*TxContextMockTest) GetTxExecSeq() int {
	panic("implement me")
}

// SetTxExecSeq set tx execute sequence
func (*TxContextMockTest) SetTxExecSeq(i int) {
	panic("implement me")
}

// GetDepth get depth
func (s *TxContextMockTest) GetDepth() int {
	return s.currentDepth
}

// GetCrossInfo get contract call link information
func (s *TxContextMockTest) GetCrossInfo() uint64 {
	return 2
}

// GetKeys key from cache, record this operation to read set
func (s *TxContextMockTest) GetKeys(keys []*vmPb.BatchKey) ([]*vmPb.BatchKey, error) {
	panic("implement me")
}

// GetTxRWMapByContractName get the read-write map of the specified contract of the current transaction
func (s *TxContextMockTest) GetTxRWMapByContractName(contractName string) (map[string]*commonPb.TxRead,
	map[string]*commonPb.TxWrite) {
	panic("implement me")
}

// HasUsed judge whether the specified commonPb.RuntimeType has appeared in the previous depth
// in the current cross-link
func (s *TxContextMockTest) HasUsed(runtimeType commonPb.RuntimeType) bool {
	panic("implement me")
}

// RecordRuntimeTypeIntoCrossInfo record the new contract call information to the top of crossInfo
func (s *TxContextMockTest) RecordRuntimeTypeIntoCrossInfo(runtimeType commonPb.RuntimeType) {
	panic("implement me")
}

// RemoveRuntimeTypeFromCrossInfo remove the top-level information from the crossInfo
func (s *TxContextMockTest) RemoveRuntimeTypeFromCrossInfo() {
	panic("implement me")
}

type mockBlockchainStore struct {
}

func (m mockBlockchainStore) MakeSnapshot(snapshotHeight uint64) error {
	//TODO implement me
	panic("implement me")
}

func (m mockBlockchainStore) GetSnapshotStatus() uint64 {
	//TODO implement me
	panic("implement me")
}

// GetHotColdDataSeparationMaxHeight get the max height which can be used for do hot-cold data separations
func (m mockBlockchainStore) GetHotColdDataSeparationMaxHeight() (uint64, error) {
	//TODO implement me
	panic("implement me")
}

// DoHotColdDataSeparation create a new task , that move cold block data to archive file system
func (m mockBlockchainStore) DoHotColdDataSeparation(startHeight uint64, endHeight uint64) (string, error) {
	//TODO implement me
	panic("implement me")
}

// GetHotColdDataSeparationJobByID return HotColdDataSeparation job info by job ID
func (m mockBlockchainStore) GetHotColdDataSeparationJobByID(jobID string) (storePb.ArchiveJob, error) {
	//TODO implement me
	panic("implement me")
}

// GetArchiveStatus add next time
func (m mockBlockchainStore) GetArchiveStatus() (*storePb.ArchiveStatus, error) {
	//TODO implement me
	panic("implement me")
}

// GetLastHeight add next time
func (m mockBlockchainStore) GetLastHeight() (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockBlockchainStore) TxExistsInIncrementDBState(txId string, startHeight uint64) (bool, bool, error) {
	//TODO implement me
	panic("implement me")
}

// TxExistsInFullDB if tx exists in full db
func (m mockBlockchainStore) TxExistsInFullDB(txId string) (bool, uint64, error) {
	panic("implement me")
}

// TxExistsInIncrementDB if tx exists in increment db
func (m mockBlockchainStore) TxExistsInIncrementDB(txId string, startHeight uint64) (bool, error) {
	panic("implement me")
}

// GetTxWithRWSet get tx with rwset
func (m mockBlockchainStore) GetTxWithRWSet(txId string) (*commonPb.TransactionWithRWSet, error) {
	//TODO implement me
	panic("implement me")
}

// GetTxInfoWithRWSet get tx info with rwset
func (m mockBlockchainStore) GetTxInfoWithRWSet(txId string) (*commonPb.TransactionInfoWithRWSet, error) {
	//TODO implement me
	panic("implement me")
}

// GetTxWithInfo get tx with info
func (m mockBlockchainStore) GetTxWithInfo(txId string) (*commonPb.TransactionInfo, error) {
	//TODO implement me
	panic("implement me")
}

// GetTxInfoOnly get tx info only
func (m mockBlockchainStore) GetTxInfoOnly(txId string) (*commonPb.TransactionInfo, error) {
	//TODO implement me
	panic("implement me")
}

// CreateDatabase create database
func (m mockBlockchainStore) CreateDatabase(contractName string) error {
	panic("implement me")
}

// DropDatabase drop database
func (m mockBlockchainStore) DropDatabase(contractName string) error {
	panic("implement me")
}

// GetContractDbName getcontract db name
func (m mockBlockchainStore) GetContractDbName(contractName string) string {
	panic("implement me")
}

// GetMemberExtraData get member extra data
func (m mockBlockchainStore) GetMemberExtraData(member *acPb.Member) (*acPb.MemberExtraData, error) {
	panic("implement me")
}

// GetContractByName get contract by name
func (m mockBlockchainStore) GetContractByName(name string) (*commonPb.Contract, error) {
	panic("implement me")
}

// GetContractBytecode get contract byte code
func (m mockBlockchainStore) GetContractBytecode(name string) ([]byte, error) {
	panic("implement me")
}

// GetHeightByHash get height by hash
func (m mockBlockchainStore) GetHeightByHash(blockHash []byte) (uint64, error) {
	panic("implement me")
}

// GetBlockHeaderByHeight get block header by height
func (m mockBlockchainStore) GetBlockHeaderByHeight(height uint64) (*commonPb.BlockHeader, error) {
	panic("implement me")
}

// GetLastChainConfig get last chain cfg
func (m mockBlockchainStore) GetLastChainConfig() (*configPb.ChainConfig, error) {
	cc := configPb.ChainConfig{
		Vm: &configPb.Vm{
			AddrType: configPb.AddrType_ETHEREUM,
		},
	}

	return &cc, nil
}

// GetTxHeight get tx height
func (m mockBlockchainStore) GetTxHeight(txId string) (uint64, error) {
	panic("implement me")
}

// GetArchivedPivot get archived pivot
func (m mockBlockchainStore) GetArchivedPivot() uint64 {
	panic("implement me")
}

//// GetArchiveStatus get archived pivot
//func (m mockBlockchainStore) GetArchiveStatus() (*storePb.ArchiveStatus, error) {
//	panic("implement me")
//}

// ArchiveBlock get archive block
func (m mockBlockchainStore) ArchiveBlock(archiveHeight uint64) error {
	panic("implement me")
}

// RestoreBlocks restore blocks
func (m mockBlockchainStore) RestoreBlocks(serializedBlocks [][]byte) error {
	panic("implement me")
}

// QuerySingle query single
func (m mockBlockchainStore) QuerySingle(contractName, sql string, values ...interface{}) (protocol.SqlRow, error) {
	panic("implement me")
}

// QueryMulti query multi
func (m mockBlockchainStore) QueryMulti(contractName, sql string, values ...interface{}) (protocol.SqlRows, error) {
	panic("implement me")
}

// ExecDdlSql execute ddl sql
func (m mockBlockchainStore) ExecDdlSql(contractName, sql string, version string) error {
	panic("implement me")
}

// BeginDbTransaction beigin db transaction
func (m mockBlockchainStore) BeginDbTransaction(txName string) (protocol.SqlDBTransaction, error) {
	panic("implement me")
}

// GetDbTransaction get db transaction
func (m mockBlockchainStore) GetDbTransaction(txName string) (protocol.SqlDBTransaction, error) {
	panic("implement me")
}

// CommitDbTransaction commit db transaction
func (m mockBlockchainStore) CommitDbTransaction(txName string) error {
	panic("implement me")
}

// RollbackDbTransaction rollback db transaction
func (m mockBlockchainStore) RollbackDbTransaction(txName string) error {
	panic("implement me")
}

// InitGenesis init genesis
func (m mockBlockchainStore) InitGenesis(genesisBlock *storePb.BlockWithRWSet) error {
	panic("implement me")
}

// PutBlock put block
func (m mockBlockchainStore) PutBlock(block *commonPb.Block, txRWSets []*commonPb.TxRWSet) error {
	panic("implement me")
}

// SelectObject select object
func (m mockBlockchainStore) SelectObject(contractName string,
	startKey []byte, limit []byte) (protocol.StateIterator, error) {
	panic("implement me")
}

// GetHistoryForKey get history for key
func (m mockBlockchainStore) GetHistoryForKey(contractName string, key []byte) (protocol.KeyHistoryIterator, error) {
	panic("implement me")
}

// GetAccountTxHistory get account tx history
func (m mockBlockchainStore) GetAccountTxHistory(accountId []byte) (protocol.TxHistoryIterator, error) {
	panic("implement me")
}

// GetContractTxHistory get contract tx history
func (m mockBlockchainStore) GetContractTxHistory(contractName string) (protocol.TxHistoryIterator, error) {
	panic("implement me")
}

// GetBlockByHash get block by hash
func (m mockBlockchainStore) GetBlockByHash(blockHash []byte) (*commonPb.Block, error) {
	panic("implement me")
}

// BlockExists if block exists
func (m mockBlockchainStore) BlockExists(blockHash []byte) (bool, error) {
	panic("implement me")
}

// GetBlock get block
func (m mockBlockchainStore) GetBlock(height uint64) (*commonPb.Block, error) {
	panic("implement me")
}

// GetLastConfigBlock get last config block
func (m mockBlockchainStore) GetLastConfigBlock() (*commonPb.Block, error) {
	panic("implement me")
}

// GetBlockByTx get block by tx
func (m mockBlockchainStore) GetBlockByTx(txId string) (*commonPb.Block, error) {
	panic("implement me")
}

// GetBlockWithRWSets get block with rwsets
func (m mockBlockchainStore) GetBlockWithRWSets(height uint64) (*storePb.BlockWithRWSet, error) {
	panic("implement me")
}

// GetTx get tx
func (m mockBlockchainStore) GetTx(txId string) (*commonPb.Transaction, error) {
	panic("implement me")
}

// TxExists tx exists
func (m mockBlockchainStore) TxExists(txId string) (bool, error) {
	panic("implement me")
}

// GetTxConfirmedTime get tx confirmed time
func (m mockBlockchainStore) GetTxConfirmedTime(txId string) (int64, error) {
	panic("implement me")
}

// GetLastBlock get last block
func (m mockBlockchainStore) GetLastBlock() (*commonPb.Block, error) {
	return &commonPb.Block{
		Header: &commonPb.BlockHeader{
			ChainId:        "",
			BlockHeight:    0,
			PreBlockHash:   nil,
			BlockHash:      nil,
			PreConfHeight:  0,
			BlockVersion:   0,
			DagHash:        nil,
			RwSetRoot:      nil,
			TxRoot:         nil,
			BlockTimestamp: 0,
			Proposer:       nil,
			ConsensusArgs:  nil,
			TxCount:        0,
			Signature:      nil,
		},
		Dag:            nil,
		Txs:            nil,
		AdditionalData: nil,
	}, nil
}

// ReadObject read object
func (m mockBlockchainStore) ReadObject(contractName string, key []byte) ([]byte, error) {
	panic("implement me")
}

// GetTxRWSet get tx rwset
func (m mockBlockchainStore) GetTxRWSet(txId string) (*commonPb.TxRWSet, error) {
	panic("implement me")
}

// GetTxRWSetsByHeight get txrwsets by height
func (m mockBlockchainStore) GetTxRWSetsByHeight(height uint64) ([]*commonPb.TxRWSet, error) {
	panic("implement me")
}

// GetDBHandle get db handle
func (m mockBlockchainStore) GetDBHandle(dbName string) protocol.DBHandle {
	panic("implement me")
}

// Close close
func (m mockBlockchainStore) Close() error {
	panic("implement me")
}

func (m mockBlockchainStore) ReadObjects(contractName string, keys [][]byte) ([][]byte, error) {
	panic("implement me")
}
