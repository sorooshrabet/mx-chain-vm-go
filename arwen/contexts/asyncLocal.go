package contexts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/math"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func (context *asyncContext) executeAsyncLocalCalls() error {
	for {
		call := context.getNextLocalAsyncCall()
		if check.IfNil(call) {
			break
		}

		err := context.executeAsyncLocalCall(call)
		if err != nil {
			return err
		}

		context.closeCompletedAsyncCalls()
	}

	return nil
}

func (context *asyncContext) executeCompletedGroupCallbacks() {
	for _, group := range context.asyncCallGroups {
		if group.IsComplete() {
			context.executeCallGroupCallback(group)
		}
	}
}

func (context *asyncContext) getNextLocalAsyncCall() *arwen.AsyncCall {
	for _, group := range context.asyncCallGroups {
		for _, call := range group.AsyncCalls {
			if call.IsLocal() {
				return call
			}
		}
	}

	return nil
}

func (context *asyncContext) executeAsyncLocalCall(asyncCall *arwen.AsyncCall) error {
	if asyncCall.ExecutionMode == arwen.ESDTTransferOnCallBack {
		context.executeESDTTransferOnCallback(asyncCall)
		return nil
	}

	destinationCallInput, err := context.createContractCallInput(asyncCall)
	if err != nil {
		return err
	}

	// Briefly restore the AsyncCall GasLimit, after it was consumed in its
	// entirety by addAsyncCall(); this is required, because ExecuteOnDestContext()
	// must also consume the GasLimit in its entirety, before starting execution,
	// but will restore any GasRemaining to the current instance.
	metering := context.host.Metering()
	metering.RestoreGas(asyncCall.GetGasLimit())

	txHashOfNewCall := destinationCallInput.CurrentTxHash
	vmOutput, err := context.host.ExecuteOnDestContext(destinationCallInput)
	if vmOutput == nil {
		return arwen.ErrNilDestinationCallVMOutput
	}

	// The vmOutput instance returned by host.ExecuteOnDestContext() is never nil,
	// by design. Using it without checking for err is safe here.
	asyncCall.UpdateStatus(vmOutput.ReturnCode)

	receiverAddr := destinationCallInput.RecipientAddr
	deserializedAsyncCtx, _ := NewSerializedAsyncContextFromStore(context.host.Storage(), receiverAddr, txHashOfNewCall)

	if asyncCall.HasCallback() && (deserializedAsyncCtx == nil || !deserializedAsyncCtx.HasPendingCallGroups()) {
		callbackVMOutput, callbackErr := context.executeSyncCallback(asyncCall, vmOutput, err)
		context.finishAsyncLocalExecution(callbackVMOutput, callbackErr)
		context.gasAccumulated = 0
		if callbackVMOutput != nil {
			context.accumulateGas(callbackVMOutput.GasRemaining)
		}
	} else {
		context.gasAccumulated = 0
		if vmOutput != nil {
			context.accumulateGas(vmOutput.GasRemaining)
		}
	}

	return nil
}

func (context *asyncContext) executeESDTTransferOnCallback(asyncCall *arwen.AsyncCall) {
	context.host.Output().PrependFinish(asyncCall.Data)

	// The contract has already paid the gas for GasLimit and
	// GasLocked, as if the call were destined to another contract. Both
	// GasLimit and GasLocked are restored in the case of
	// ESDTTransferOnCallBack because:
	// * GasLocked isn't needed, since no callback will be called
	// * GasLimit cannot be paid here, because it's the *destination*
	// contract that ends up paying the gas for the ESDTTransfer
	context.host.Metering().RestoreGas(asyncCall.GasLimit)
	context.host.Metering().RestoreGas(asyncCall.GasLocked)
	asyncCall.UpdateStatus(vmcommon.Ok)
}

func (context *asyncContext) executeSyncCallback(
	asyncCall *arwen.AsyncCall,
	destinationVMOutput *vmcommon.VMOutput,
	destinationErr error,
) (*vmcommon.VMOutput, error) {

	callbackInput, err := context.createCallbackInput(asyncCall, destinationVMOutput, destinationErr)
	if err != nil {
		return nil, err
	}

	// Restore gas locked while still on the caller instance; otherwise, the
	// locked gas will appear to have been used twice by the caller instance.
	context.host.Metering().RestoreGas(asyncCall.GetGasLocked())
	callbackVMOutput, callBackErr := context.host.ExecuteOnDestContext(callbackInput)

	return callbackVMOutput, callBackErr
}

// executeCallGroupCallback synchronously executes the designated callback of
// the AsyncCallGroup, as it was set with SetGroupCallback().
//
// Gas for the execution has been already paid for when SetGroupCallback() was
// set. The remaining gas is refunded to context.callerAddr, which initiated
// the call and paid for the gas in the first place.
func (context *asyncContext) executeCallGroupCallback(group *arwen.AsyncCallGroup) {
	if !group.HasCallback() {
		return
	}

	input := context.createGroupCallbackInput(group)
	context.gasAccumulated = 0
	vmOutput, err := context.host.ExecuteOnDestContext(input)
	context.finishAsyncLocalExecution(vmOutput, err)
	logAsync.Trace("gas remaining after group callback", "group", group.Identifier, "gas", vmOutput.GasRemaining)
	context.accumulateGas(vmOutput.GasRemaining)
}

// executeSyncHalfOfBuiltinFunction will synchronously call the requested
// built-in function. This is required for all cross-shard calls to built-in
// functions, because they will handle cross-shard calls themselves, by
// generating entries in vmOutput.OutputAccounts, and they need to be executed
// synchronously to do that. As a consequence, it is not necessary to call
// sendAsyncCallCrossShard(). The vmOutput produced by the built-in function,
// containing the cross-shard call, has ALREADY been merged into the main
// output by the inner call to host.ExecuteOnDestContext(). Moreover, the
// status of the AsyncCall is not updated here - it will be updated by
// PostprocessCrossShardCallback(), when the cross-shard call returns.
func (context *asyncContext) executeSyncHalfOfBuiltinFunction(asyncCall *arwen.AsyncCall) error {
	destinationCallInput, err := context.createContractCallInput(asyncCall)
	if err != nil {
		return err
	}

	// Briefly restore the AsyncCall GasLimit, after it was consumed in its
	// entirety by addAsyncCall(); this is required, because ExecuteOnDestContext()
	// must also consume the GasLimit in its entirety, before starting execution,
	// but will restore any GasRemaining to the current instance.
	metering := context.host.Metering()
	metering.RestoreGas(asyncCall.GetGasLimit())

	vmOutput, err := context.host.ExecuteOnDestContext(destinationCallInput)
	if err != nil {
		return err
	}

	// If the in-shard half of the built-in function call has failed, go no
	// further and execute the error callback of this AsyncCall.
	if vmOutput.ReturnCode != vmcommon.Ok {
		asyncCall.Reject()
		callbackVMOutput, callbackErr := context.executeSyncCallback(asyncCall, vmOutput, err)
		context.finishAsyncLocalExecution(callbackVMOutput, callbackErr)
	}

	// The gas that remains after executing the in-shard half of the built-in
	// function is provided to the cross-shard half.
	asyncCall.GasLimit = vmOutput.GasRemaining

	return nil
}

// TODO return values are never used by code that calls finishAsyncLocalExecution
func (context *asyncContext) finishAsyncLocalExecution(vmOutput *vmcommon.VMOutput, err error) {
	if err == nil {
		return
	}

	runtime := context.host.Runtime()
	output := context.host.Output()

	runtime.GetVMInput().GasProvided = 0

	if vmOutput == nil {
		vmOutput = output.CreateVMOutputInCaseOfError(err)
	}

	output.SetReturnCode(vmOutput.ReturnCode)

	output.SetReturnMessage(vmOutput.ReturnMessage)
	output.Finish([]byte(vmOutput.ReturnCode.String()))
	output.Finish(runtime.GetCurrentTxHash())
}

func (context *asyncContext) createContractCallInput(asyncCall *arwen.AsyncCall) (*vmcommon.ContractCallInput, error) {
	host := context.host
	runtime := host.Runtime()
	sender := runtime.GetSCAddress()

	function, arguments, err := context.callArgsParser.ParseData(string(asyncCall.GetData()))
	if err != nil {
		return nil, err
	}

	gasLimit := asyncCall.GetGasLimit()
	gasToUse := host.Metering().GasSchedule().ElrondAPICost.AsyncCallStep
	if gasLimit <= gasToUse {
		return nil, arwen.ErrNotEnoughGas
	}
	gasLimit -= gasToUse

	// TODO matei-p factor out to use also in tests
	// newTxHash := NewTxHashFromExisting([]byte(asyncCall.Identifier), runtime.GetCurrentTxHash(), runtime.GetPrevTxHash())
	newTxHash := append([]byte(asyncCall.Identifier), runtime.GetCurrentTxHash()...)
	newTxHash = append(newTxHash, runtime.GetPrevTxHash()...)
	newTxHash, _ = host.Crypto().Sha256(newTxHash)

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     sender,
			Arguments:      arguments,
			CallValue:      big.NewInt(0).SetBytes(asyncCall.GetValue()),
			CallType:       vm.AsynchronousCall,
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    gasLimit,
			GasLocked:      asyncCall.GetGasLocked(),
			CurrentTxHash:  newTxHash,
			OriginalTxHash: runtime.GetOriginalTxHash(),
			PrevTxHash:     runtime.GetCurrentTxHash(),
		},
		RecipientAddr: asyncCall.GetDestination(),
		Function:      function,
	}

	return contractCallInput, nil
}

// TODO function too large; refactor needed
func (context *asyncContext) createCallbackInput(
	asyncCall *arwen.AsyncCall,
	vmOutput *vmcommon.VMOutput,
	destinationErr error,
) (*vmcommon.ContractCallInput, error) {
	metering := context.host.Metering()
	runtime := context.host.Runtime()

	// always provide return code as the first argument to callback function
	arguments := [][]byte{
		big.NewInt(int64(vmOutput.ReturnCode)).Bytes(),
	}

	esdtFunction := ""
	isESDTOnCallBack := false
	esdtArgs := make([][]byte, 0)
	returnWithError := false
	if destinationErr == nil && vmOutput.ReturnCode == vmcommon.Ok {
		// when execution went Ok, callBack arguments are:
		// [0, result1, result2, ....]
		isESDTOnCallBack, esdtFunction, esdtArgs = context.isESDTTransferOnReturnDataWithNoAdditionalData(
			asyncCall.Destination,
			runtime.GetSCAddress(),
			vmOutput)
		arguments = append(arguments, vmOutput.ReturnData...)
	} else {
		// when execution returned error, callBack arguments are:
		// [error code, error message]
		arguments = append(arguments, []byte(vmOutput.ReturnMessage))
		returnWithError = true
	}

	callbackFunction := asyncCall.GetCallbackName()

	gasLimit := math.AddUint64(vmOutput.GasRemaining, asyncCall.GetGasLocked())
	gasLimit = math.AddUint64(gasLimit, context.gasAccumulated)
	dataLength := computeDataLengthFromArguments(callbackFunction, arguments)

	gasToUse := metering.GasSchedule().ElrondAPICost.AsyncCallStep
	copyPerByte := metering.GasSchedule().BaseOperationCost.DataCopyPerByte
	gas := math.MulUint64(copyPerByte, uint64(dataLength))
	gasToUse = math.AddUint64(gasToUse, gas)
	if gasLimit <= gasToUse {
		return nil, arwen.ErrNotEnoughGas
	}
	gasLimit -= gasToUse

	// Return to the sender SC, calling its specified callback method.
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:           asyncCall.Destination,
			Arguments:            arguments,
			CallValue:            context.computeCallValueFromVMOutput(vmOutput),
			CallType:             vm.AsynchronousCallBack,
			GasPrice:             runtime.GetVMInput().GasPrice,
			GasProvided:          gasLimit,
			GasLocked:            0,
			CurrentTxHash:        runtime.GetCurrentTxHash(),
			OriginalTxHash:       runtime.GetOriginalTxHash(),
			PrevTxHash:           runtime.GetPrevTxHash(),
			ReturnCallAfterError: returnWithError,
		},
		RecipientAddr: runtime.GetSCAddress(),
		Function:      callbackFunction,
	}

	if isESDTOnCallBack {
		contractCallInput.Function = esdtFunction
		contractCallInput.Arguments = make([][]byte, 0, len(arguments))
		contractCallInput.Arguments = append(contractCallInput.Arguments, esdtArgs...)
		contractCallInput.Arguments = append(contractCallInput.Arguments, []byte(callbackFunction))
		contractCallInput.Arguments = append(contractCallInput.Arguments, big.NewInt(int64(vmOutput.ReturnCode)).Bytes())
		if len(vmOutput.ReturnData) > 1 {
			contractCallInput.Arguments = append(contractCallInput.Arguments, vmOutput.ReturnData[1:]...)
		}
	}

	return contractCallInput, nil
}

func (context *asyncContext) createGroupCallbackInput(group *arwen.AsyncCallGroup) *vmcommon.ContractCallInput {
	runtime := context.host.Runtime()

	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallType:       vm.AsynchronousCallBack,
			CallerAddr:     context.callerAddr,
			Arguments:      [][]byte{group.CallbackData},
			CallValue:      big.NewInt(0),
			GasPrice:       context.gasPrice,
			GasProvided:    group.GasLocked + context.gasAccumulated,
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
			PrevTxHash:     runtime.GetPrevTxHash(),
		},
		RecipientAddr: runtime.GetSCAddress(),
		Function:      group.Callback,
	}

	logAsync.Trace("created group callback input", "group", group.Identifier, "function", input.Function)
	logAsync.Trace("created group callback input gas", "provided", input.GasProvided, "locked", group.GasLocked, "accumulated", context.gasAccumulated)
	return input
}

func (context *asyncContext) createContextCallbackInput() *vmcommon.ContractCallInput {
	host := context.host
	runtime := host.Runtime()

	arguments := [][]byte{context.callbackData}

	// TODO ensure a new value for VMInput.CurrentTxHash
	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     context.callerAddr,
			Arguments:      arguments,
			CallValue:      runtime.GetVMInput().CallValue,
			CallType:       vm.AsynchronousCallBack,
			GasPrice:       runtime.GetVMInput().GasPrice,
			GasProvided:    context.gasAccumulated,
			CurrentTxHash:  runtime.GetCurrentTxHash(),
			OriginalTxHash: runtime.GetOriginalTxHash(),
			PrevTxHash:     runtime.GetPrevTxHash(),
		},
		RecipientAddr: runtime.GetSCAddress(),
		Function:      context.callback,
	}

	logAsync.Trace("created context callback input", "sc", runtime.GetSCAddress(), "function", input.Function)
	logAsync.Trace("created context callback input gas", "provided", input.GasProvided, "accumulated", context.gasAccumulated)
	return input
}

func (context *asyncContext) isESDTTransferOnReturnDataWithNoAdditionalData(
	sndAddr, dstAddr []byte,
	destinationVMOutput *vmcommon.VMOutput,
) (bool, string, [][]byte) {
	if len(destinationVMOutput.ReturnData) == 0 {
		return false, "", nil
	}

	functionName, args, err := context.callArgsParser.ParseData(string(destinationVMOutput.ReturnData[0]))
	if err != nil {
		return false, "", nil
	}

	return context.isESDTTransferOnReturnDataFromFunctionAndArgs(sndAddr, dstAddr, functionName, args)
}

func (context *asyncContext) isESDTTransferOnReturnDataFromFunctionAndArgs(
	sndAddr, dstAddr []byte,
	functionName string,
	args [][]byte,
) (bool, string, [][]byte) {
	parsedTransfer, err := context.esdtTransferParser.ParseESDTTransfers(sndAddr, dstAddr, functionName, args)
	if err != nil {
		return false, functionName, args
	}

	isNoCallAfter := len(parsedTransfer.CallFunction) == 0
	return isNoCallAfter, functionName, args
}

func (context *asyncContext) computeCallValueFromVMOutput(destinationVMOutput *vmcommon.VMOutput) *big.Int {
	if len(destinationVMOutput.ReturnData) > 0 {
		return big.NewInt(0)
	}

	returnTransfer := big.NewInt(0)
	callBackReceiver := context.host.Runtime().GetSCAddress()
	outAcc, ok := destinationVMOutput.OutputAccounts[string(callBackReceiver)]
	if !ok {
		return returnTransfer
	}

	if len(outAcc.OutputTransfers) == 0 {
		return returnTransfer
	}

	lastOutTransfer := outAcc.OutputTransfers[len(outAcc.OutputTransfers)-1]
	if len(lastOutTransfer.Data) == 0 {
		returnTransfer.Set(lastOutTransfer.Value)
	}

	return returnTransfer
}
