package evmaas

import (
	"math/big"

	"github.com/0xPolygon/polygon-edge/types"
	"github.com/SealSC/SealEVM/evmInt256"
	"golang.org/x/crypto/sha3"
)

func GetCodeHash(code []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(code)
	return hash.Sum(nil)
}
func AddressToInt(address Address) *evmInt256.Int {
	var i big.Int
	i.SetBytes(address[:])
	bigInt := evmInt256.FromBigInt(&i)
	return bigInt
}
func BytesToInt(data []byte) *evmInt256.Int {
	var i big.Int
	i.SetBytes(data)
	bigInt := evmInt256.FromBigInt(&i)
	return bigInt
}
func GetContractAddr(txHash []byte) Address {
	return Address(txHash[12:32])
}

func DecodeTx(txData []byte) (*types.Transaction, error) {
	tx := types.Transaction{}
	err := tx.UnmarshalRLP(txData)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}
