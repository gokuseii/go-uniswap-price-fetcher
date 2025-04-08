package helpers

import (
	"context"
	"fmt"
	"math/big"

	"go-uniswap-price-fetcher/internal/config"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

type Token struct {
	Address  common.Address
	Symbol   string
	Decimals uint8
}

type Slot0 struct {
	SqrtPriceX96               *big.Int
	Tick                       *big.Int
	ObservationIndex           uint16
	ObservationCardinality     uint16
	ObservationCardinalityNext uint16
	FeeProtocol                uint8
	Unlocked                   bool
}

func Pack(abi abi.ABI, method string, args ...any) []byte {
	data, err := abi.Pack(method, args...)
	if err != nil {
		panic(fmt.Sprintf("failed to pack ABI method %s: %v", method, err))
	}
	return data
}

func ContractCall(client *ethclient.Client, address string, data []byte) (*[]byte, error) {
	addr := common.HexToAddress(address)

	callMsg := ethereum.CallMsg{
		To:   &addr,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call contract")
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result, nil
}

func GetTokenData(client *ethclient.Client, address string, idx int) (*Token, error) {
	if idx != 0 && idx != 1 {
		panic("token id must be 0 or 1")
	}

	callData := Pack(config.Cfg.ABI.UniPool, fmt.Sprintf("token%d", idx))
	result, err := ContractCall(client, address, callData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch token in uniswap pool")
	}

	if result == nil {
		return nil, nil
	}

	tokenAddress := common.BytesToAddress([]byte(*result)[len(*result)-20:])
	callData = Pack(config.Cfg.ABI.Erc20, "decimals")

	result, err = ContractCall(client, tokenAddress.String(), callData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token decimals")
	}

	var decimals uint8
	err = config.Cfg.ABI.Erc20.UnpackIntoInterface(&decimals, "decimals", *result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to cast token decimals")
	}

	callData = Pack(config.Cfg.ABI.Erc20, "symbol")
	result, err = ContractCall(client, tokenAddress.String(), callData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token symbol")
	}

	var symbol string
	err = config.Cfg.ABI.Erc20.UnpackIntoInterface(&symbol, "symbol", *result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to cast token symbol")
	}

	return &Token{Address: tokenAddress, Symbol: symbol, Decimals: decimals}, nil
}

func GetSlot0(client *ethclient.Client, address string) (*Slot0, error) {
	callData := Pack(config.Cfg.ABI.UniPool, "slot0")
	result, err := ContractCall(client, address, callData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch slot0")
	}

	var slot Slot0
	err = config.Cfg.ABI.UniPool.UnpackIntoInterface(&slot, "slot0", *result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to cast slot")
	}

	return &slot, nil
}
