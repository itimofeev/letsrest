package letsrest

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func Must(err error, msg ...interface{}) {
	if err != nil {
		log.Fatal(err, msg)
	}
}

func PrintJSON(i interface{}) {
	fmt.Println("JSON: ", GetJSON(i))
}

func GetJSON(i interface{}) string {
	j, err := json.Marshal(i)
	Must(err, "Ma")
	return string(j)
}
