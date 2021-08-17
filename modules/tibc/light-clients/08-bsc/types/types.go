package types

import (
	"bytes"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func ParseValidators(validatorsBytes []byte) ([][]byte, error) {
	if len(validatorsBytes)%AddressLength != 0 {
		return nil, sdkerrors.Wrap(ErrInvalidValidatorBytes, "(validatorsBytes % AddressLength) should bz zero")
	}
	n := len(validatorsBytes) / AddressLength
	result := make([][]byte, n)
	for i := 0; i < n; i++ {
		address := make([]byte, AddressLength)
		copy(address, validatorsBytes[i*AddressLength:(i+1)*AddressLength])
		result[i] = address
	}
	return result, nil
}

func (vs ValidatorSet) Has(validator []byte) bool {
	for _, v := range vs.Validators {
		if bytes.Equal(v, validator) {
			return true
		}
	}
	return false
}

func (vs ValidatorSet) inturn(height uint64, validator []byte) bool {
	validators := vs.Validators
	offset := (height + 1) % uint64(len(validators))
	return bytes.Equal(validators[offset], validator)
}
