package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
	"strings"
)

// ParseClassTrace parses a string with the ibc prefix (class trace) and the base class
// into a ClassTrace type.
//
// Examples:
//
// 	- "nft/A/B/dog" => ClassTrace{Path: "nft/A/B", BaseClass: "dog"}
// 	- "dog" => ClassTrace{Path: "", BaseClass: "dog"}
func ParseClassTrace(rawClass string) ClassTrace {
	classSplit := strings.Split(rawClass, "/")

	if classSplit[0] == rawClass {
		return ClassTrace{
			Path:      "",
			BaseClass: rawClass,
		}
	}

	return ClassTrace{
		Path:      strings.Join(classSplit[:len(classSplit)-1], "/"),
		BaseClass: classSplit[len(classSplit)-1],
	}
}


// Hash returns the hex bytes of the SHA256 hash of the ClassTrace fields using the following formula:
//
// hash = sha256(tracePath + "/" + baseClass)
func (ct ClassTrace) Hash() tmbytes.HexBytes {
	hash := sha256.Sum256([]byte(ct.GetFullClassPath()))
	return hash[:]
}

// GetFullClassPath returns the full class according to the ICS20 specification:
// tracePath + "/" + baseClass
// If there exists no trace then the base class is returned.
func (ct ClassTrace) GetFullClassPath() string {
	if ct.Path == "" {
		return ct.BaseClass
	}
	return ct.GetPrefix() + ct.BaseClass
}

// GetPrefix returns the receiving class prefix composed by the trace info and a separator.
func (ct ClassTrace) GetPrefix() string {
	return ct.Path + "/"
}

// IBCClass a nft class for an TICS30 fungible token in the format
// 'tibc-{hash(tracePath + baseClass)}'. If the trace is empty, it will return the base denomination.
func (ct ClassTrace) IBCClass() string {
	if ct.Path != "" {
		return fmt.Sprintf("%s-%s", ClassPrefix, ct.Hash())
	}
	return ct.BaseClass
}

// ParseHexHash parses a hex hash in string format to bytes and validates its correctness.
func ParseHexHash(hexHash string) (tmbytes.HexBytes, error) {
	hash, err := hex.DecodeString(hexHash)
	if err != nil {
		return nil, err
	}

	if err := tmtypes.ValidateHash(hash); err != nil {
		return nil, err
	}

	return hash, nil
}
