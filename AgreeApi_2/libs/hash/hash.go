package im_hash

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

type (
	Hash []byte
)

func (h Hash) String() string { return hex.EncodeToString(h) }

func HashSha256(src uuid.UUID, hash bool) (Hash, error) {
	if !hash && len(src) != 32 {
		return Hash{}, fmt.Errorf("Hash is %d bytes in length.", len(src))
	}

	b, _ := src.MarshalText()

	var dst Hash = make(Hash, 32)
	if hash {
		sha := sha256.Sum256(b)
		copy(dst, sha[:])
	} else {
		copy(dst, src)
	}
	return dst, nil
}

func ValidateHashedAdId(hash Hash) (Hash, error) {
	hexBytes, _ := hex.DecodeString(hex.EncodeToString(hash)) // Ensure hash is actually in hex.
	if len(hexBytes) != 32 || !bytes.Equal(hash, hexBytes) {  // bytes.Compare(hash, hexBytes) != 0 { // Ensure hash is correct length.
		return nil, fmt.Errorf("Hash incorrect length")
	}
	return hexBytes, nil
}
