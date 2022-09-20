package wasmer

import (
	"fmt"

	"github.com/ElrondNetwork/wasm-vm/executorinterface"
)

const noArity = -1

// getSignature returns the signature for the given functionName
func (instance *Instance) getSignature(functionName string) (*ExportedFunctionSignature, bool) {
	signature, ok := instance.Signatures[functionName]
	return signature, ok
}

func (instance *Instance) verifyVoidFunction(functionName string) error {
	inArity, err := instance.getInputArity(functionName)
	if err != nil {
		return err
	}

	outArity, err := instance.getOutputArity(functionName)
	if err != nil {
		return err
	}

	isVoid := inArity == 0 && outArity == 0
	if !isVoid {
		return fmt.Errorf("%w: %s", executorinterface.ErrFunctionNonvoidSignature, functionName)
	}
	return nil
}

func (instance *Instance) getInputArity(functionName string) (int, error) {
	signature, ok := instance.getSignature(functionName)
	if !ok {
		return noArity, fmt.Errorf("%w: %s", executorinterface.ErrFuncNotFound, functionName)
	}
	return signature.InputArity, nil
}

func (instance *Instance) getOutputArity(functionName string) (int, error) {
	signature, ok := instance.getSignature(functionName)
	if !ok {
		return noArity, fmt.Errorf("%w: %s", executorinterface.ErrFuncNotFound, functionName)
	}
	return signature.OutputArity, nil
}
