package util

import (
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

// GetKeyFromValue to get key from value
func GetKeyFromValue(m map[uint8]map[string]string, value string) (uint8, bool) {
	for k, v := range m {
		for _, vv := range v {
			if value == vv {
				return k, true
			}
		}
	}
	return 0, false
}

// GetUUID to generate a uuid
func GetUUID() string {
	id, _ := uuid.NewV4()
	return id.String()
}

// GenerateGUID to generate guid without hyphen
func GenerateGUID() string {
	id, _ := uuid.NewV4()
	return strings.ReplaceAll(id.String(), "-", "")
}

// HomeDir to get the home directory of the computer
func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// GetUint8MapKeys to get all keys from a map which keys is uint8
func GetUint8MapKeys(m map[uint8]map[string]string) []uint8 {
	keys := []uint8{}
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// ByValue means sort by value
type ByValue []uint8

func (a ByValue) Len() int           { return len(a) }
func (a ByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByValue) Less(i, j int) bool { return a[i] < a[j] }

// GetSortUint8Slice to sort a uint8 slice
func GetSortUint8Slice(s []uint8) []uint8 {
	sort.Sort(ByValue(s))
	return s
}

// GetSortedUint8MapKeys to get all sorted keys from a map which keys is uint8
func GetSortedUint8MapKeys(m map[uint8]map[string]string) []uint8 {
	keys := GetUint8MapKeys(m)
	return GetSortUint8Slice(keys)
}

// StringInSlice to check if a string in a slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// RandNumbers generate a random string, each character is a digit
// n is the length of the string
func RandNumbers(n int8) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// RandStringRunes generate a random string, each character is a alphabet
// n is the length of the string
func RandStringRunes(n int8) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// RandLowerStringRunes generate a random string, each character is a lower-case alphabet
// n is the length of the string
func RandLowerStringRunes(n int8) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// CheckErrorIsDuplicateEntryInDB to check whether the error is duplicate entry error
func CheckErrorIsDuplicateEntryInDB(err error) bool {
	return strings.HasPrefix(err.Error(), "Error 1062: Duplicate entry")
}
