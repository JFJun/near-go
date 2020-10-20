package transaction

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/JFJun/near-go/serialize"
	"github.com/btcsuite/btcutil/base58"
	"strings"
)

type Transaction struct {
	SignerId   serialize.String
	PublicKey  serialize.PublicKey
	Nonce      serialize.U64
	ReceiverId serialize.String
	BlockHash  serialize.BlockHash
	Actions    []serialize.IAction
}
type ActionTransfer struct {
	Transfer Transfer `json:"transfer"`
}

type Transfer struct {
	Deposit string `json:"deposit"` //amount
}

/*
blockHash: latest block hash
nonce: access_key["nonce"]+1
PublicKey: hex or base58
*/
func CreateTransaction(from, to, publicKey, blockHash string, nonce int64) (*Transaction, error) {
	var err error
	tx := new(Transaction)
	tx.SignerId = serialize.String{Value: from}
	tx.Nonce = serialize.U64{Value: uint64(nonce)}
	tx.ReceiverId = serialize.String{Value: to}
	bh := base58.Decode(blockHash)
	if len(bh) == 0 {
		return nil, fmt.Errorf("base58  decode blockhash error ,BlockHash=%s", blockHash)
	}
	tx.BlockHash = serialize.BlockHash{
		Value: bh,
	}
	publicKey = strings.TrimPrefix(publicKey, "ed25519:")
	var pk []byte
	if len(publicKey) == 64 { //is hex
		pk, err = hex.DecodeString(publicKey)
		if err != nil {
			return nil, fmt.Errorf("decode public key error,Err=%v", err)
		}
	} else {
		pk = base58.Decode(publicKey)
		if len(pk) == 0 {
			return nil, fmt.Errorf("base58 decode public key error,Public Key=%s", publicKey)
		}
	}
	tx.PublicKey = serialize.PublicKey{
		KeyType: 0,
		Value:   pk,
	}
	return tx, nil
}

func (tx *Transaction) SetAction(action ...serialize.IAction) {
	tx.Actions = append(tx.Actions, action...)
}

func (tx *Transaction) Serialize() ([]byte, error) {
	var (
		data []byte
	)
	ss, err := tx.SignerId.Serialize()
	if err != nil {
		return nil, fmt.Errorf("tx serialize: signerId error,Err=%v", err)
	}
	data = append(data, ss...)
	ps, err := tx.PublicKey.Serialize()
	if err != nil {
		return nil, fmt.Errorf("tx serialize: publickey error,Err=%v", err)
	}
	data = append(data, ps...)
	ns, err := tx.Nonce.Serialize()
	if err != nil {
		return nil, fmt.Errorf("tx serialize: nonce error,Err=%v", err)
	}
	data = append(data, ns...)
	rs, err := tx.ReceiverId.Serialize()
	if err != nil {
		return nil, fmt.Errorf("tx serialize: ReceiverId error,Err=%v", err)
	}
	data = append(data, rs...)
	bs, err := tx.BlockHash.Serialize()
	if err != nil {
		return nil, fmt.Errorf("tx serialize: blockhash error,Err=%v", err)
	}
	data = append(data, bs...)
	//序列化action
	al := len(tx.Actions)
	uAL := serialize.U32{
		Value: uint32(al),
	}
	uALData, err := uAL.Serialize()
	if err != nil {
		return nil, fmt.Errorf("tx serialize: action length error,Err=%v", err)
	}
	data = append(data, uALData...)
	for _, action := range tx.Actions {
		as, err := action.Serialize()
		if err != nil {
			return nil, fmt.Errorf("tx serialize: action error,Err=%v", err)
		}
		data = append(data, as...)
	}
	return data, nil
}

func SignTransaction(tx_hex string, privateKey string) (string, error) {
	priv, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("decode private key error,Err=%v", err)
	}
	if len(priv) != 32 {
		return "", fmt.Errorf("private key kength is not equal 32,Len=%d", len(priv))
	}
	data, err := hex.DecodeString(tx_hex)
	if err != nil {
		return "", fmt.Errorf("decode tx hex error,Err=%v", err)
	}
	preSigData := sha256.Sum256(data)
	p := ed25519.NewKeyFromSeed(priv)
	sig := ed25519.Sign(p, preSigData[:])
	if len(sig) != 64 {
		return "", fmt.Errorf("sign error,length is not equal 64,length=%d", len(sig))
	}
	return hex.EncodeToString(sig), nil
}

type SignatureTransaction struct {
	Sig serialize.Signature
	Tx  *Transaction
}

func CreateSignatureTransaction(tx *Transaction, sig string) (*SignatureTransaction, error) {
	var signature []byte
	var err error
	if len(sig) == 128 {
		signature, err = hex.DecodeString(sig)
		if err != nil {
			return nil, err
		}
	} else {
		signature = base58.Decode(sig)
		if len(signature) == 0 {
			return nil, fmt.Errorf("b58 decode sig error,sig=%s", sig)
		}
	}
	stx := new(SignatureTransaction)
	stx.Tx = tx
	stx.Sig = serialize.Signature{
		KeyType: tx.PublicKey.KeyType,
		Value:   signature,
	}
	return stx, nil
}

func (stx *SignatureTransaction) Serialize() ([]byte, error) {
	data, err := stx.Tx.Serialize()
	if err != nil {
		return nil, fmt.Errorf("sign serialize: tx serialize error,Err=%v", err)
	}
	ss, err := stx.Sig.Serialize()
	if err != nil {
		return nil, fmt.Errorf("sign serialize: sig serialize error,Err=%v", err)
	}
	data = append(data, ss...)
	return data, nil
}
