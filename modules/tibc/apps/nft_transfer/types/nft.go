package types

import "fmt"

// GetClassPrefix returns the receiving class prefix
func GetClassPrefix(sourceChain, destChain string) string {
	return fmt.Sprintf("%s/%s/", sourceChain, destChain)
}
