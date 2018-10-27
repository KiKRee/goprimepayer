package goprimepayer

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

func sortByKeys(params map[string]string) []string {
	sortedKeys := make([]string, 0)
	for key := range params {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

func hash(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}
