package rpc

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/JFJun/near-go/account"
	"github.com/btcsuite/btcutil/base58"
	"strings"
)

// 签名能用到的接口

//获取最新的block hash

func (c *Client) GetLatestBlockHash() (string, error) {
	var chainStatus ChainStatus
	err := c.Post("status", &chainStatus, []interface{}{})
	if err != nil {
		return "", err
	}
	if chainStatus.SyncInfo.LatestBlockHash == "" {
		return "", errors.New("latest block hash is null")
	}
	return chainStatus.SyncInfo.LatestBlockHash, nil
}

/*

 */
func (c *Client) GetNonce(address, publicKey, finality string) (int64, error) {
	var err error
	var pkStr string
	if !strings.HasPrefix(publicKey, "ed25519:") {
		var pk []byte
		if len(publicKey) == 64 {
			pk, err = hex.DecodeString(publicKey)
			if err != nil {
				return -1, err
			}
			pkStr = account.PublicKeyToString(pk)
		} else {
			pk = base58.Decode(publicKey)
			if len(pk) == 0 {
				return -1, fmt.Errorf("b58 decode public key error, %s", publicKey)
			}
			pkStr = "ed25519:" + publicKey
		}
	} else {
		pkStr = publicKey
	}
	params := make(map[string]interface{})
	params["request_type"] = "view_access_key"
	params["account_id"] = address
	params["public_key"] = pkStr
	params["finality"] = finality
	if finality == "" {
		params["finality"] = "optimistic"
	}
	var resp map[string]interface{}
	err = c.Post("query", &resp, params)
	if err != nil {
		return -1, err
	}
	if resp["nonce"] == nil {
		return -1, fmt.Errorf("resp nonce is null,resp=%v", resp)
	}
	nonce := int64(resp["nonce"].(float64))
	// fork: github.com/near/near-api-py/near_api/account.py    page: 32
	return nonce + 1, nil
}

func (c *Client)GetAccountBalance(address string)(string,string,error){
	params := make(map[string]interface{})
	params["request_type"] = "view_account"
	params["account_id"] = address
	params["finality"] = "final"
	var resp map[string]interface{}
	err := c.Post("query", &resp, params)
	if err != nil {
		return "", "", err
	}
	if resp["amount"]==nil {
		return "","",errors.New("amount is null")
	}
	if resp["locked"]==nil  {
		return "","",errors.New("locked amount is null")
	}
	return resp["acount"].(string),resp["locked"].(string),nil
}

func (c *Client) BroadcastTransaction(stxBase64 string) (string, error) {
	var resp NearTxStatus
	err := c.Post("broadcast_tx_commit", &resp, []interface{}{stxBase64})
	if err != nil {
		return "", fmt.Errorf("broadcast tx error,Err=%v", err)
	}
	//只关心txid，不关心是否发送成功
	txid := resp.Transaction.Txid
	return txid, nil
}
