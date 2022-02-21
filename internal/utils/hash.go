package utils

import (
	"crypto/sha1"
	"encoding/base32"
	"strings"
)

func HashString(s string, length int) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	return strings.ToLower(base32.HexEncoding.EncodeToString(bs))[:length]
}
