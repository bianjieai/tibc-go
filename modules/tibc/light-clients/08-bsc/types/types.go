package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// Fixed number of extra-data prefix bytes reserved for signer vanity
	extraVanity = 32

	// Fixed number of extra-data suffix bytes reserved for signer seal
	extraSeal = 65

	// AddressLength is the expected length of the address
	AddressLength = 20
)

func ParseValidators(validatorsBytes []byte) ([][]byte, error) {
	if len(validatorsBytes)%AddressLength != 0 {
		return nil, sdkerrors.Wrap(ErrInvalidValidatorBytes, "(validatorsBytes % AddressLength) == 0")
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
