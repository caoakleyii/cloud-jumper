/*
Package hasher implements a simple library for handling Sha512 hashing

	The hasher library allows for salting your hash.
	Functions include returning byte arrays, strings, and base64 url safe encoded strings
*/
package hasher

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
)

// Sha512Hash Hash a value and salt
// Returns the hashed value as a byte array
func Sha512Hash(value, salt string) [64]byte {
	return sha512.Sum512([]byte(value + salt))
}

/*
	1. Hash and Encode a Password string
	Provides the ability to take a password as a string and return a Base64 encoded strings.
*/

// Sha512HashBase64Encode Hashes and Encodes a value and salt
// Returns the hash and encoded value as a string
func Sha512HashBase64Encode(value, salt string) string {
	b := Sha512Hash(value, salt)
	return base64.StdEncoding.EncodeToString(b[:])
}

// GenerateRandomString Creates a secure random string that is base64 URL Encoded
// Returns the string and nil or empty string and error
func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), err
}
