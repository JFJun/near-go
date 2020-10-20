package rpc

type ChainStatus struct {
	ChainId  string        `json:"chain_id"`
	SyncInfo ChainSyncInfo `json:"sync_info"`
}

type ChainSyncInfo struct {
	LatestBlockHash   string `json:"latest_block_hash"`
	LatestBlockHeight int64  `json:"latest_block_height"`
	LatestBlockTime   string `json:"latest_block_time"`
	Syncing           bool   `json:"syncing"`
}

type NearTxStatus struct {
	Status             interface{}            `json:"status"`
	Transaction        NearTransaction        `json:"transaction"`
	TransactionOutcome NearTransactionOutcome `json:"transaction_outcome"`
}

type NearTransaction struct {
	SignerId   string        `json:"signer_id"`
	PublicKey  string        `json:"public_key"`
	Nonce      int           `json:"nonce"`
	ReceiverId string        `json:"receiver_id"`
	Actions    []interface{} `json:"actions"`
	Signature  string        `json:"signature"`
	Txid       string        `json:"hash"`
}
type NearTransactionOutcome struct {
	Proof     []interface{} `json:"proof"`
	BlockHash string        `json:"block_hash"`
	Id        string        `json:"id"`
	Outcome   NearOutcome   `json:"outcome"`
}
type NearOutcome struct {
	ReceiptIds  []string          `json:"receipt_ids"`
	GasBurnt    int64             `json:"gas_burnt"`
	TokensBurnt string            `json:"tokens_burnt"`
	ExecutorId  string            `json:"executor_id"`
	Status      map[string]string `json:"status"` // "status": {"SuccessReceiptId": "4xjmyr15T4UqRjdd5YERpfS8QRdCWTH392sMKdQWaJzM"}
}
