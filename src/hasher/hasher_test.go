package hasher

import (
	"encoding/base64"
	"testing"
)

func TestSha512Hash(t *testing.T) {
	expected := `ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==`

	password := "angryMonkey"
	salt := ""
	b := Sha512Hash(password, salt)

	if len(b) < 64 {
		t.Errorf("Sha512Hash did not return 64 length array of bytes")
	}

	h := base64.StdEncoding.EncodeToString(b[:])

	if h != expected {
		t.Errorf("Sha512Hash did not properly hash the password. Expected: %v | Returned: %v", expected, h)
	}

}

func TestSha512HashBase64Encode(t *testing.T) {
	expected := `ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==`
	password := "angryMonkey"
	salt := ""
	h := Sha512HashBase64Encode(password, salt)

	if h != expected {
		t.Errorf("Sha512HashBase64Encode did not properly hash the password. Expected: %v | Returned: %v", expected, h)
	}
}

func TestGenerateRandomString(t *testing.T) {
	s, err := GenerateRandomString(8)
	if err != nil {
		t.Errorf("GenerateRandomString errored \n\n %v", err)
	}
	if s == "" {
		t.Errorf("GenerateRandomString returned an empty string")
	}
	if len(s) < 12 || len(s) > 12 {
		t.Errorf("GenerateRandomString returned a string of the wrong length based on byte size")
	}
}
