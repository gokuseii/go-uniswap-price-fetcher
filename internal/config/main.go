package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"

	"github.com/joho/godotenv"
)

type ContractsABI struct {
	Erc20   abi.ABI
	UniPool abi.ABI
}

type Config struct {
	Port    string
	RpcURL  string
	ChainID string
	ABI     ContractsABI
}

var Cfg Config

func loadAbiContract(path string) abi.ABI {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("failed to read ABI file at path: %s", path))
	}

	contract, err := abi.JSON(strings.NewReader(string(bytes)))
	if err != nil {
		panic(fmt.Sprintf("failed to load erc20 contract ABI file at path: %s", path))
	}

	return contract
}

func Init() {
	if err := godotenv.Load(); err != nil {
		panic(errors.Wrap(err, "failed to load .env file"))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	rpcURL := os.Getenv("RPC_URL")
	if !strings.HasPrefix(rpcURL, "http") {
		panic(errors.New("invalid or empty RPC_URL"))
	}

	chainID := os.Getenv("CHAIN_ID")
	if chainID == "" {
		panic(errors.New("CHAIN_ID is required"))
	}

	erc20 := loadAbiContract(filepath.Join("internal", "contracts", "IERC20.abi"))
	uniPool := loadAbiContract(filepath.Join("internal", "contracts", "IUniswapV3Pool.abi"))

	Cfg = Config{
		Port:    port,
		RpcURL:  rpcURL,
		ChainID: chainID,
		ABI: ContractsABI{
			Erc20:   erc20,
			UniPool: uniPool,
		},
	}
}
