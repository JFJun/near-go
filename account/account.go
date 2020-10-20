package account

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcutil/base58"
)

func GenerateKey() ([]byte, []byte, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return priv.Seed(), pub, nil
}
func GeneratePubKeyBySeed(seed []byte) ([]byte, error) {
	priv := ed25519.NewKeyFromSeed(seed)

	pub := priv[32:]
	if len(pub) != 32 {
		return nil, errors.New("public key length is not equal 32")
	}
	return pub, nil
}

func GeneratePubKeyByBase58(b58Key string) ([]byte, error) {
	seed := base58.Decode(b58Key)
	if len(seed) == 0 {
		return nil, errors.New("base 58 decode error")
	}
	if len(seed) != 32 {
		return nil, errors.New("seed length is not equal 32")
	}
	return GeneratePubKeyBySeed(seed)
}

func PublicKeyToString(pub []byte) string {
	publicKey := base58.Encode(pub)
	return "ed25519:" + publicKey
}

func PublicKeyToAddress(pub []byte) string {
	return hex.EncodeToString(pub)
}
