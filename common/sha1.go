package common

import (
	"crypto/sha1"
	"encoding/hex"
)

func Sha1(text, salt string) string {
	h := sha1.New()
	h.Write([]byte(text))
	h.Write([]byte(salt))
	return hex.EncodeToString(h.Sum(nil))
}
