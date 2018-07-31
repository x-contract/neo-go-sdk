# neo-go-sdk
A Go Sdk for the NEO blockchain

Neo is a blockchain written in C# programing language which depends on Microsoft DotNet Core for cross-platform. More information could be found on the website https://neo.org/. GitHub Project for Neo is here https://github.com/neo-project

Go has been more and more popular in blockchain programing and decentralized application developing. The neo-go-sdk provides basic functionality to Neo block chain. It's designed to work with Neo-Cli node, not to relace it. With neo-go-sdk you can handle making kinds of transactions, create addresses, encryption and signature stuff, Neo script and so on. It makes developing Neo in Go more easy.



## Download

Just run the command
    
    go get -u github.com/x-contract/neo-go-sdk
    
The sdk depends on golang.org/x/crypto/ripemd160 to work, use the command
  
    go get -u golang.org/x/crypto/ripemd160

## How to use

###  Neo-Cli Node

First of all you need a Neo-Cli full node. The neo-go-sdk calls Neo-Cli json-rpc APIs to gain access to block chain data.

You can create a node yourself, the tools and knowledge all you need to create a Neo-Cli node could be found here http://docs.neo.org/ . However you may want to have a quick start then you can directly use the official node without need for permission, these nodes are from http://seed1.neo.org:10332 to http://seed5.neo.org:10332. You can change port 10332 to 20332 for test net.


### Algorithm, keys, addresses and script hashes etc.

#### Base58 and Hashes

    endata, checksum := neoutils.EncodeBase58WithChecksum(data)
    data, checked := neoutils.DecodeBase58WithChecksum(endata)
    hash256 := neoutils.Hash256(data)
    hash160 := neoutils.Hash160(data)

#### Get a private key from WIF string

    key, _ := neotransaction.DecodeFromWif(`KzSToRnDi9V********************************`)

#### Encode private key to WIF string

    wif := key.EncodeWif()
    
#### Create new key pair and new account address with it

    key := neotransaction.GenerateKeyPair()
    addr := key.CreateBasicAddress()
    
#### Address to ScriptHash 

    addr, _ := neotransaction.ParseAddress("ASMGHQPzZqxFB2yKmzvfv82jtKVnjhp1ES")
    scripthash := addr.ScripHash
   
#### Sign and CheckSign

    sig, _ := key.Sign(data)
    check := key.Verify(data, sig)


### Make a ContractTransaction (Transfer with common UTXO assets)

UTXO assets like NEO and GAS was transfered using a kind of transaction called ContractTransaction. It's the basic transaction in all block chains.
  
If you want to make a ContractTransaction from your code you should first get the UTXOs for the account address which you want to transfer from. Unfortunately the Neo-Cli node does not provide the ability to get that from rpc APIs. So I have to achieve that by myself using a Block-Spider (or called Block-Listener). 
  
Block-Spider fetch all the blocks in sequence and parse the transactions within those blocks to storage all the outputs data then filter all the UTXOs for any address. For now I'm using Block-Spider contributed by NEL (https://github.com/NewEconoLab/NEO_Block_API). The Go implemented Block-Spider is just on the way.
  
#### Make ContractTransaction is so easy

    utxos, _ := neoextapi.FetchUTXO(config.NEOEXTURL, addr, neotransaction.AssetNeoID)
    tx := neotransaction.CreateContractTransaction()
    tx.AppendInput(utxos[0])
    tx.AppendOutput(taddr, utxos[0].AssetID, utxos[0].Value)
    txid := tx.TXID()
    tx.AppendBasicSignWitness(key)
    result := neocliapi.SendRawTransaction(config.NEOCLIURL, tx.RawTransactionString())
    

### Make InvocationTransaction (Using Neo smart contract)

NEO smart contract is published and invoked by InvocationTransaction, with the script that push params and call to specific contract hash. In fact, InvocationTransaction just invokes a slice of compiled NeoVM script regardless of what the script means. Call to another contract or publish a new contract is some NeoVM functions just like others, there's nothing special except the cost GAS differs.
  
The result of the script is logged in the ApplicationLog, which you can get from the Neo-Cli node via rpc APIs. Remember to start the Neo-Cli node with --log args to enable the ApplicationLog.
  
Before you make a ContractTransaction you have to build your script invoked by the transaction. There's a simple ScriptBuilder in neo-go-sdk that works with NeoVM OpCodes and parameters. Complicated contract script like NEP-5 ICO contract script should be written and compiled using tools from Neo and finally transferd to bytes, then impacted into a InvocationTransaction.
  
#### Build script
	
    sb := neotransaction.ScriptBuilder{}
    // If you want to make an invocation transaction without utxo transferd
    // then you need to push a random number so that the hash(txid) could vary on each transaction
    rand.Seed(time.Now().UnixNano())
    sb.EmitPushNumber(int64(rand.Uint32()))
    sb.Emit(OpCode.DROP)
    args := []interface{}{205, addr.ScripHash}
    sb.EmitPushArray(args)
    sb.EmitPushBool(false)
    sb.EmitPushString(`name`)
    contractHash, _ := hex.DecodeString(contractHashString)
    sb.EmitAppCall(contractHash)
    
#### Make InvocationTransaction

    tx := neotransaction.CreateInvocationTransaction()
    extra := tx.ExtraData.(*neotransaction.InvocationExtraData)
    extra.Script = sb.Bytes()
    // If the transaction need additional Witness then put the ScriptHash in attributes
    tx.AppendAttribute(neotransaction.UsageScript, addr.ScripHash)
    // Perhaps the transaction need Witness
    tx.AppendBasicSignWitness(key)
    txid := tx.TXID()
    rawtx := tx.RawTransactionString()
    result := neocliapi.SendRawTransaction(neocliurl, rawtx)
    
    


### Helper function to call Neo-Cli rpc APIs

There's some helper function that making call to Neo-Cli rpc APIs more easily, just like:
  
    neocliapi.FetchBalance(config.NEOCLIURL, user.UserNeoAddress.Addr)
  
This could get an account's balance from a Neo-Cli node
  
  
  
## TODO List

The list here contains the job that I am currently struggling with. There's more work to do with the neo-go-sdk to make it more convenient to use :). 
  
  ### Block Spider
  ### Nep6 Wallet and Private Key encryption
  ### Nep5 Contract asset support (Balance, Transfer, ICO etc.)



## License
- Open-source MIT.
- Main author is @terender.
