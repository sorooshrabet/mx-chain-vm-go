package host

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/require"
)

var defaultVMType = []byte{0xF, 0xF}
var errAccountNotFound = errors.New("account not found")

var userAddress = []byte("userAccount.....................")

// AddressSize is the size of an account address, in bytes.
const AddressSize = 32

// SCAddressPrefix is the prefix of any smart contract address used for testing.
var SCAddressPrefix = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0f")
var parentAddress = MakeTestSCAddress("parentSC")
var childAddress = MakeTestSCAddress("childSC")

var customGasSchedule = config.GasScheduleMap(nil)

// MakeTestSCAddress generates a new smart contract address to be used for
// testing based on the given identifier.
func MakeTestSCAddress(identifier string) []byte {
	numberOfTrailingDots := AddressSize - len(SCAddressPrefix) - len(identifier)
	leftBytes := SCAddressPrefix
	rightBytes := []byte(identifier + strings.Repeat(".", numberOfTrailingDots))
	return append(leftBytes, rightBytes...)
}

// GetSCCode retrieves the bytecode of a WASM module from a file
func GetSCCode(fileName string) []byte {
	code, err := ioutil.ReadFile(filepath.Clean(fileName))
	if err != nil {
		panic(fmt.Sprintf("GetSCCode(): %s", fileName))
	}

	return code
}

// GetTestSCCode retrieves the bytecode of a WASM testing contract
func GetTestSCCode(scName string, prefixToTestSCs string) []byte {
	pathToSC := prefixToTestSCs + "test/contracts/" + scName + "/output/" + scName + ".wasm"
	return GetSCCode(pathToSC)
}

// GetTestSCCodeModule retrieves the bytecode of a WASM testing contract, given
// a specific name of the WASM module
func GetTestSCCodeModule(scName string, moduleName string, prefixToTestSCs string) []byte {
	pathToSC := prefixToTestSCs + "test/contracts/" + scName + "/output/" + moduleName + ".wasm"
	return GetSCCode(pathToSC)
}

// BuildSCModule invokes erdpy to build the contract into a WASM module
func BuildSCModule(scName string, prefixToTestSCs string) {
	pathToSCDir := prefixToTestSCs + "test/contracts/" + scName
	out, err := exec.Command("erdpy", "contract", "build", "--no-optimization", pathToSCDir).Output()
	if err != nil {
		log.Error("error building contract", "err", err, "contract", pathToSCDir)
		return
	}

	log.Info("contract built", "output", fmt.Sprintf("\n%s", out))
}

// defaultTestArwenForDeployment creates an Arwen vmHost configured for testing deployments
func defaultTestArwenForDeployment(t *testing.T, _ uint64, newAddress []byte) (*vmHost, *contextmock.BlockchainHookStub) {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		return &contextmock.StubAccount{
			Nonce: 24,
		}, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		return newAddress, nil
	}

	host := defaultTestArwen(t, stubBlockchainHook)
	return host, stubBlockchainHook
}

func defaultTestArwenForCall(tb testing.TB, code []byte, balance *big.Int) (*vmHost, *contextmock.BlockchainHookStub) {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
		if bytes.Equal(scAddress, parentAddress) {
			return &contextmock.StubAccount{
				Balance: balance,
			}, nil
		}
		return nil, errAccountNotFound
	}
	stubBlockchainHook.GetCodeCalled = func(account vmcommon.UserAccountHandler) []byte {
		return code
	}

	host := defaultTestArwen(tb, stubBlockchainHook)
	return host, stubBlockchainHook
}

func defaultTestArwenForCallWithInstanceRecorderMock(tb testing.TB, code []byte, balance *big.Int) (*vmHost, *contextmock.InstanceBuilderRecorderMock) {
	// this uses a Blockchain Hook Stub that does not cache the compiled code
	host, _ := defaultTestArwenForCall(tb, code, balance)

	instanceBuilderRecorderMock := contextmock.NewInstanceBuilderRecorderMock()
	host.Runtime().ReplaceInstanceBuilder(instanceBuilderRecorderMock)

	return host, instanceBuilderRecorderMock
}
func defaultTestArwenForCallWithInstanceMocks(tb testing.TB) (*vmHost, *worldmock.MockWorld, *contextmock.InstanceBuilderMock) {
	world := worldmock.NewMockWorld()
	host := defaultTestArwen(tb, world)

	instanceBuilderMock := contextmock.NewInstanceBuilderMock(world)
	host.Runtime().ReplaceInstanceBuilder(instanceBuilderMock)

	return host, world, instanceBuilderMock
}

func defaultTestArwenForCallWithWorldMock(tb testing.TB, code []byte, balance *big.Int) (*vmHost, *worldmock.MockWorld) {
	world := worldmock.NewMockWorld()
	host := defaultTestArwen(tb, world)

	err := world.InitBuiltinFunctions(host.GetGasScheduleMap())
	require.Nil(tb, err)

	host.protocolBuiltinFunctions = world.BuiltinFuncs.GetBuiltinFunctionNames()

	parentAccount := world.AcctMap.CreateSmartContractAccount(userAddress, parentAddress, code)
	parentAccount.Balance = balance

	return host, world
}

// defaultTestArwenForTwoSCs creates an Arwen vmHost configured for testing calls between 2 SmartContracts
func defaultTestArwenForTwoSCs(
	t *testing.T,
	parentCode []byte,
	childCode []byte,
	parentSCBalance *big.Int,
	childSCBalance *big.Int,
) (*vmHost, *contextmock.BlockchainHookStub) {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}

	if parentSCBalance == nil {
		parentSCBalance = big.NewInt(1000)
	}

	if childSCBalance == nil {
		childSCBalance = big.NewInt(1000)
	}

	stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
		if bytes.Equal(scAddress, parentAddress) {
			return &contextmock.StubAccount{
				Address: parentAddress,
				Balance: parentSCBalance,
			}, nil
		}
		if bytes.Equal(scAddress, childAddress) {
			return &contextmock.StubAccount{
				Address: childAddress,
				Balance: childSCBalance,
			}, nil
		}

		return nil, errAccountNotFound
	}
	stubBlockchainHook.GetCodeCalled = func(account vmcommon.UserAccountHandler) []byte {
		if bytes.Equal(account.AddressBytes(), parentAddress) {
			return parentCode
		}
		if bytes.Equal(account.AddressBytes(), childAddress) {
			return childCode
		}
		return nil
	}

	host := defaultTestArwen(t, stubBlockchainHook)
	return host, stubBlockchainHook
}

func defaultTestArwenForContracts(
	t *testing.T,
	contracts []*instanceTestSmartContract,
) (*vmHost, *contextmock.BlockchainHookStub) {

	stubBlockchainHook := &contextmock.BlockchainHookStub{}

	contractsMap := make(map[string]*contextmock.StubAccount)
	codeMap := make(map[string]*[]byte)

	for _, contract := range contracts {
		contractsMap[string(contract.address)] = &contextmock.StubAccount{Address: contract.address, Balance: big.NewInt(contract.balance)}
		codeMap[string(contract.address)] = &contract.code
	}

	stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
		contract, found := contractsMap[string(scAddress)]
		if found {
			return contract, nil
		}
		return nil, errAccountNotFound
	}
	stubBlockchainHook.GetCodeCalled = func(account vmcommon.UserAccountHandler) []byte {
		code, found := codeMap[string(account.AddressBytes())]
		if found {
			return *code
		}
		return nil
	}

	host := defaultTestArwen(t, stubBlockchainHook)
	return host, stubBlockchainHook
}

func defaultTestArwenWithWorldMock(tb testing.TB) (*vmHost, *worldmock.MockWorld) {
	world := worldmock.NewMockWorld()
	host := defaultTestArwen(tb, world)

	err := world.InitBuiltinFunctions(host.GetGasScheduleMap())
	require.Nil(tb, err)

	host.protocolBuiltinFunctions = world.BuiltinFuncs.GetBuiltinFunctionNames()
	return host, world
}

func defaultTestArwen(tb testing.TB, blockchain vmcommon.BlockchainHook) *vmHost {
	gasSchedule := customGasSchedule
	if gasSchedule == nil {
		gasSchedule = config.MakeGasMapForTests()
	}

	host, err := NewArwenVM(blockchain, &arwen.VMHostParameters{
		VMType:                   defaultVMType,
		BlockGasLimit:            uint64(1000),
		GasSchedule:              gasSchedule,
		ProtocolBuiltinFunctions: make(vmcommon.FunctionNames),
		ElrondProtectedKeyPrefix: []byte("ELROND"),
		UseWarmInstance:          false,
		DynGasLockEnableEpoch:    0,
	})
	require.Nil(tb, err)
	require.NotNil(tb, host)

	return host
}

// AddTestSmartContractToWorld directly deploys the provided code into the
// given MockWorld under a SC address built with the given identifier.
func AddTestSmartContractToWorld(world *worldmock.MockWorld, identifier string, code []byte) *worldmock.Account {
	address := MakeTestSCAddress(identifier)
	return world.AcctMap.CreateSmartContractAccount(userAddress, address, code)
}

// DefaultTestContractCreateInput creates a vmcommon.ContractCreateInput struct
// with default values.
func DefaultTestContractCreateInput() *vmcommon.ContractCreateInput {
	return &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: []byte("caller"),
			Arguments: [][]byte{
				[]byte("argument 1"),
				[]byte("argument 2"),
			},
			CallValue:   big.NewInt(0),
			CallType:    vmcommon.DirectCall,
			GasPrice:    0,
			GasProvided: 0,
		},
		ContractCode: []byte("contract"),
	}
}

// DefaultTestContractCallInput creates a vmcommon.ContractCallInput struct
// with default values.
func DefaultTestContractCallInput() *vmcommon.ContractCallInput {
	return &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  userAddress,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    vmcommon.DirectCall,
			GasPrice:    0,
			GasProvided: 0,
		},
		RecipientAddr: parentAddress,
		Function:      "function",
	}
}

type contractCallInputBuilder struct {
	vmcommon.ContractCallInput
}

func createTestContractCallInputBuilder() *contractCallInputBuilder {
	return &contractCallInputBuilder{
		ContractCallInput: *DefaultTestContractCallInput(),
	}
}

func (contractInput *contractCallInputBuilder) withRecipientAddr(address []byte) *contractCallInputBuilder {
	contractInput.ContractCallInput.RecipientAddr = address
	return contractInput
}

func (contractInput *contractCallInputBuilder) withGasProvided(gas uint64) *contractCallInputBuilder {
	contractInput.ContractCallInput.VMInput.GasProvided = gas
	return contractInput
}

func (contractInput *contractCallInputBuilder) withFunction(function string) *contractCallInputBuilder {
	contractInput.ContractCallInput.Function = function
	return contractInput
}

func (contractInput *contractCallInputBuilder) withArguments(arguments ...[]byte) *contractCallInputBuilder {
	contractInput.ContractCallInput.VMInput.Arguments = arguments
	return contractInput
}

func (contractInput *contractCallInputBuilder) withCurrentTxHash(txHash []byte) *contractCallInputBuilder {
	contractInput.ContractCallInput.CurrentTxHash = txHash
	return contractInput
}

func (contractInput *contractCallInputBuilder) withESDTValue(esdtValue *big.Int) *contractCallInputBuilder {
	contractInput.ContractCallInput.ESDTValue = esdtValue
	return contractInput
}

func (contractInput *contractCallInputBuilder) withESDTTokenName(esdtTokenName []byte) *contractCallInputBuilder {
	contractInput.ContractCallInput.ESDTTokenName = esdtTokenName
	return contractInput
}

func (contractInput *contractCallInputBuilder) build() *vmcommon.ContractCallInput {
	return &contractInput.ContractCallInput
}

type contractCreateInputBuilder struct {
	vmcommon.ContractCreateInput
}

func createTestContractCreateInputBuilder() *contractCreateInputBuilder {
	return &contractCreateInputBuilder{
		ContractCreateInput: *DefaultTestContractCreateInput(),
	}
}

func (contractInput *contractCreateInputBuilder) withGasProvided(gas uint64) *contractCreateInputBuilder {
	contractInput.ContractCreateInput.GasProvided = gas
	return contractInput
}

func (contractInput *contractCreateInputBuilder) withContractCode(code []byte) *contractCreateInputBuilder {
	contractInput.ContractCreateInput.ContractCode = code
	return contractInput
}

func (contractInput *contractCreateInputBuilder) withCallerAddr(address []byte) *contractCreateInputBuilder {
	contractInput.ContractCreateInput.CallerAddr = address
	return contractInput
}

func (contractInput *contractCreateInputBuilder) withCallValue(callValue int64) *contractCreateInputBuilder {
	contractInput.ContractCreateInput.CallValue = big.NewInt(callValue)
	return contractInput
}

func (contractInput *contractCreateInputBuilder) withArguments(arguments ...[]byte) *contractCreateInputBuilder {
	contractInput.ContractCreateInput.Arguments = arguments
	return contractInput
}

func (contractInput *contractCreateInputBuilder) build() *vmcommon.ContractCreateInput {
	return &contractInput.ContractCreateInput
}

// MakeVMOutput creates a vmcommon.VMOutput struct with default values
func MakeVMOutput() *vmcommon.VMOutput {
	return &vmcommon.VMOutput{
		ReturnCode:      vmcommon.Ok,
		ReturnMessage:   "",
		ReturnData:      make([][]byte, 0),
		GasRemaining:    0,
		GasRefund:       big.NewInt(0),
		DeletedAccounts: make([][]byte, 0),
		TouchedAccounts: make([][]byte, 0),
		Logs:            make([]*vmcommon.LogEntry, 0),
		OutputAccounts:  make(map[string]*vmcommon.OutputAccount),
	}
}

// MakeVMOutputError creates a vmcommon.VMOutput struct with default values
// for errors
func MakeVMOutputError() *vmcommon.VMOutput {
	return &vmcommon.VMOutput{
		ReturnCode:      vmcommon.ExecutionFailed,
		ReturnMessage:   "",
		ReturnData:      nil,
		GasRemaining:    0,
		GasRefund:       big.NewInt(0),
		DeletedAccounts: nil,
		TouchedAccounts: nil,
		Logs:            nil,
		OutputAccounts:  nil,
	}
}

// AddFinishData appends the provided []byte to the ReturnData of the given vmOutput
func AddFinishData(vmOutput *vmcommon.VMOutput, data []byte) {
	vmOutput.ReturnData = append(vmOutput.ReturnData, data)
}

// AddNewOutputAccount creates a new vmcommon.OutputAccount from the provided arguments and adds it to OutputAccounts of the provided vmOutput
func AddNewOutputAccount(vmOutput *vmcommon.VMOutput, sender []byte, address []byte, balanceDelta int64, data []byte) *vmcommon.OutputAccount {
	account := &vmcommon.OutputAccount{
		Address:        address,
		Nonce:          0,
		BalanceDelta:   big.NewInt(balanceDelta),
		Balance:        nil,
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
		Code:           nil,
	}
	if data != nil {
		account.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data:          data,
				Value:         big.NewInt(balanceDelta),
				SenderAddress: sender,
			},
		}
	}
	vmOutput.OutputAccounts[string(address)] = account
	return account
}

// SetStorageUpdate sets a storage update to the provided vmcommon.OutputAccount
func SetStorageUpdate(account *vmcommon.OutputAccount, key []byte, data []byte) {
	keyString := string(key)
	update, exists := account.StorageUpdates[keyString]
	if !exists {
		update = &vmcommon.StorageUpdate{}
		account.StorageUpdates[keyString] = update
	}
	update.Offset = key
	update.Data = data
}

// SetStorageUpdateStrings sets a storage update to the provided vmcommon.OutputAccount, from string arguments
func SetStorageUpdateStrings(account *vmcommon.OutputAccount, key string, data string) {
	SetStorageUpdate(account, []byte(key), []byte(data))
}

// OpenFile method opens the file from given path - does not close the file
func OpenFile(relativePath string) (*os.File, error) {
	path, err := filepath.Abs(relativePath)
	if err != nil {
		fmt.Printf("cannot create absolute path for the provided file: %s", err.Error())
		return nil, err
	}
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	return f, nil
}

// LoadTomlFileToMap opens and decodes a toml file as a map[string]interface{}
func LoadTomlFileToMap(relativePath string) (map[string]interface{}, error) {
	f, err := OpenFile(relativePath)
	if err != nil {
		return nil, err
	}

	fileinfo, err := f.Stat()
	if err != nil {
		fmt.Printf("cannot stat file: %s", err.Error())
		return nil, err
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = f.Read(buffer)
	if err != nil {
		fmt.Printf("cannot read from file: %s", err.Error())
		return nil, err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			fmt.Printf("cannot close file: %s", err.Error())
		}
	}()

	loadedTree, err := toml.Load(string(buffer))
	if err != nil {
		fmt.Printf("cannot interpret file contents as toml: %s", err.Error())
		return nil, err
	}

	loadedMap := loadedTree.ToMap()

	return loadedMap, nil
}

// LoadGasScheduleConfig parses and prepares a gas schedule read from file.
func LoadGasScheduleConfig(filepath string) (config.GasScheduleMap, error) {
	gasScheduleConfig, err := LoadTomlFileToMap(filepath)
	if err != nil {
		return nil, err
	}

	flattenedGasSchedule := make(config.GasScheduleMap)
	for libType, costs := range gasScheduleConfig {
		flattenedGasSchedule[libType] = make(map[string]uint64)
		costsMap := costs.(map[string]interface{})
		for operationName, cost := range costsMap {
			flattenedGasSchedule[libType][operationName] = uint64(cost.(int64))
		}
	}

	return flattenedGasSchedule, nil
}
