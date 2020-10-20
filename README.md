# near offline sign go sdk

## create address
    priv,pub,err:=account.GenerateKey()
    if err != nil {
    	panic(err)
    }
    fmt.Println("Priv: ",hex.EncodeToString(priv))
    fmt.Println("Address: ",account.PublicKeyToAddress(pub))
    
## transfer
    client,err:=rpc.NewRpcClient("https://rpc.mainnet.near.org","","")
    if err != nil {
    	panic(err)
    }
    blockhash,err:=client.GetLatestBlockHash()
    if err != nil {
    	panic(err)
    }
    fmt.Println("Blockhash: ",blockhash)
    pub,_:=hex.DecodeString("f0cb2082b845259526fba5953897760cdb06ccf090093e94298b829b564392f7")
    pubKey:=account.PublicKeyToString(pub)
    fmt.Println("PublicKey: ",pubKey)
    nonce,err:=client.GetNonce("f0cb2082b845259526fba5953897760cdb06ccf090093e94298b829b564392f7",pubKey,"")
    if err != nil {
    	panic(err)
    }
    fmt.Println("Nonce: ",nonce)
    tx,err:=transaction.CreateTransaction(
    	"f0cb2082b845259526fba5953897760cdb06ccf090093e94298b829b564392f7",
    	"37f453aa256e430216a5b7b01f92386846876fa4fef5ecc29d12acb179380fd6",
    	pubKey,
    	blockhash,
    	nonce,
    )
    if err != nil {
    	panic(err)
    }
    amount:=decimal.NewFromFloat(0.1).Shift(24)
    fmt.Println(amount.String())
    ta,err:=serialize.CreateTransfer(amount.String())
    if err != nil {
    	panic(err)
    }
    tx.SetAction(ta)
    txData,err:=tx.Serialize()
    if err != nil {
    	panic(err)
    }
    tx_hex:=hex.EncodeToString(txData)
    sig,err:=transaction.SignTransaction(tx_hex,"private key")
    if err != nil {
    	panic(err)
    }
    fmt.Println("Sig: ",sig)
    stx,err:=transaction.CreateSignatureTransaction(tx,sig)
    if err != nil {
    	panic(err)
    }
    stxData,err:=stx.Serialize()
    if err != nil {
    	panic(err)
    }
    b64Data:=base64.StdEncoding.EncodeToString(stxData)
    txid,err:=client.BroadcastTransaction(b64Data)
    if err != nil {
    	panic(err)
    }
    fmt.Println("Txid: ",txid)