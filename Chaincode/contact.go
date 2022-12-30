package chaincode

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const uniKey = "uri"

const balancePrefix = "accont~tokenId~sender"
const approvalPrefix = "account~operator"

const minterMSPID = "Org1MSP"

//Define key names for options
const nameKey = "name"
const symbolKey = "symbol"

//SmartContact provides functions for transferring token between accouts
type SmartContract struct {
	contractapi.Contract
}

// To represents recipent address
// ID represents token ID

type ToID struct {
	To string
	ID uint64
}

type BalanceHolder struct {
	Recipent string `json:"recipient`
	Sender   string `json:"sender`
	ToID	 uint64 `json:"id`
	Value	 uint64 `json:"value`
}

// Mint creates amount token of token type id and assigns them to account.
// This function emits a TransferSingle event

func (s * SmartContract) Mint(ctx contractapi.TransactionContextInterface, account string, id uint64, amount uint64) error {

	//check id contract has been initialized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return fmt.Errorf(Contract options need to be se before calling any function, call Initialize() to initialize cor"")
	}

	// check minter authorization - this sample assumes Org1 is the central banker with privilefe to mint new  tokens
	err = authorizationHelper(ctx)
	if err!= nil {
		return err
	}

	// Get ID of submitting client identity

	operator, err := ctx.GetCkientIdentity().GetID()
	fmt.Printf("Minting %d tokens of type %d for account %s with operator %s\n", amount, id, account, operator)
	if err != nil {
		return fmt.Errorf("failed to get client id: %v",err)
	}

	// Mint tokens
	err - mintHelper(ctx,operator,account,id,amount)
	if err != nil {
		return err
	}

	// EMit TransferSingle eveent
	transferSingleEvent := TransferSingle{operator, "0x0", account, id, amount}
	return emitTransferSingle(ctx, transferSingleEvent)
}

// MintBatch creates amount tokens for each token type id and assigns them to account
// This function emits a TransferBatch event.

funct (s * SmartContract) MintBatch()ctx xontractapi,TransactionContextInterface, account string, ids []uint64, amounts []uint64) error {

	//check if contract has been initialized first
	initialized, err := checkInititalized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return fmt.Errorf("Contact options need to be set before calling any function, call Initialize() to initialize")

	}

	if len(ids) != len(amounts){
		return fmt.Errorf("ids and amounts must have the same lenghts")
	}

	// check minter authorization - this sample assumes Org1 is the central banker with privilege to mint new tokens
	err = authorizationHelper(ctx)
	if err != nil {
		return err
	}

	// Get ID of submitting client identity
	operator, err := ctx.GetclientIdentity().GetID()
	if  err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	// Group amount by token is because we can only send token to a recipient only one time in a block. This prevent key config
	amountToSend := nake(map[uint64]uint64) // token id => amount

	for i :=0; i < len(amounts); i++ {
		amount ToSend[ids[i]], err = add(amountToSend[idss[i]], amounts[i])
		if err != nil {
			return err
		}
	}

	// copy the map keys and sort it. This is necessary because iterating maos in Go is not deterministic
	amountToSendKeys := sortedKeys(amountToSend)

	// Mint tokens
	for _,id := range amountToSendKeys {
		amount := amountToSend[id]
		err = mintHelper(ctx, operator, account, id, amount)
		if err != nil {
			return err
		}
	}

	// EMit TransferBatch event
	transferBatchEvent := TransferBatch{operator, "0x0", account, id, amounts}
	return emitTransferBatch(ctx, transferBatchEvent)
}

// Burn destroys amount tokens of token type id from account
// This function triggers a TransferSingle event
fun (s *SmartContract) Burn(ctx contractapi.TransactionContextInterface, account string, id uint64, amount uint64) error {

	//check if contract has been initialized firs
	initialized, err := checkInitialized(ctx)
	if err != nil{
		return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize")
	}

	if account == "0x0" {
		return fmt.Errorf("burn to the zero address")
	}

	//check minter authorization - this sample assumes Org1 is the central banker with privilage to burn new tokens
	err = authorizationHelper(ctx)
	if err != nil {
			return err
	}

	// Get ID of submitting client identity
	operator, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	// Burn tokens
	err = removeBalance(ctx, account, []uint64{id}, []uint64{amount})
	if err != nil {
		return err
	}

	transferSingleEvent := TransferSingle{operator, account, "oxo", id,amount}
	return emitTransferSingle(ctx, transferSingleEvent)
}

// BurnBatch destroys amount tokens of for each token type id from account
// This function emits a TransferBtach event.

fun (s *SmartContract) BurnBatch(ctx contractapi.TransactionContextInterface, account string, ids[]uint64, amount []uint64) error {
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}

	if !initialized {
		return fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize ")
	}

	if account == "0x0" {
		return fmt.Errorf("burn to the zero address")
	}

	if len(ids) != len(amounts){
		return fmt.Errorf("ids and amounts must have the same lenghts")
	}

	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to burn new tokens
	err = authorizationHelper(ctx)
	if err != nil {
		return err
	}

	// Ger ID of submitting client identity
	operator, err := ctx.GetCkientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}
}

// TransferFrom transfers token from sender account to recipient account
// recipient account must be a valid clientID as returned by the ClientID() function
// This function triggers a TransferSingle event
fun (s *SmartContract) TransferFrom(ctx contractapi.TransactionContextInterface, sender string, recipient string, id uint64, amount){
	 //check if contract has been initialized first
	 initialized, erro := checkInitialized(ctx)
	 if err != nil {
	 	return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	 }
	 if !initialized {
	 	return fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize ")
	 }

	 if sender == recipient {
	 	log.Println("error : transfer to self")
	 	reutrn fmt.Errorf("transfer to self")
	 }

	 // Get ID of submitting client identity
	 operator, err := ctx.GetclientIdentity().GetID()
	 if err != nil {
	 	log.printf("error : failed to get client id: %v\n", err)
	 	return fmt.Errorf("failed to get client id: %v", err)
	 }

	 //check whhether operator is owner or apprved
	 if operator != sender {
	 	approved, err := _isApprovedForAll(ctx, sender, operator)
	 	if err != nil {
			return err
		}
		if !approved {
			log.Println("error : caller is not owner nor is approved")
			return fmt.Errorf("caller is not owner nor is approved")
		}
	 }

	 //Withdraw the funds from the sender address
	 err = removeBalance(ctx, sender, []uint64{id}, []uint64{amount})
	 if err != nil {
		return err
	 }

	 if recipient == "0x0" {
	 	return fmt.Errorf("transfer tp the zero address")
	 }

	 //Deposit the fnd to the recipient address
	 err = adBalance(ctx,sender, recipient, id, amount)
	 if err != nil {
		return err
	 }

	 // Emit TranferSingle event
	 transferSingleEvent := TransferSingle{operator, sender, recipient, id, amount}
	 return emitTransferSingle(ctx, transferSingleEvent)
} 

// BatchTransferForm transfers multiple tokens from sender account to recipient account
// recipient account must be a valid clientID as returned by the ClientID() funtio
// This function triggers a TransferBatch event

fun (s * SmartContract) BatchTransferFrom(ctx contractapi.TransactionContextInterface, sender string, recipient string, ids []uint64){

	//check if contract has been initialized first
	initialized, err := checkInitialized(ctx)
	if err !=nil {
		return fmt.Errorf("failed to check if contact ia already initialized: %v",err)
	}

	if !initialized {
		return fmt.Errorf("contact options need to be set before calling any function, call Initialize() to initialize config")
	}

	if sender == recipient {
		return fmt.Errorf("transfer to self")
	}

	if len(ids) != len(amounts) {
		return fmt.Errorf("ids and amounts must have the same lenght")
	}

	// Get ID of submitting client identity
	operator, err := ctx.GetCkientIdentity().GetID()
	if err != nill {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	// check whether opertor is owner or approved
	if operator != sender {
		approved, err := _isApprovedForAll(ctx,sender, operator)
		if err != nil {
			return err
		}
		if !approved {
			return fmt.Errorf("caller is not owner nor is approved")
		}

	}

	// withdraw the funds fro the sender address
	ERR = removeBalance(ctx,sender,ids,amounts)
	if err!= nil {
		return err
	}

	if recipient == "0x0" {
		return fmt.Errorf("transfer to the zero address")
	}

	//Group amount by token id because we can only send token to a recipient only one time in a block.
	amountToSend := make(map[uint64]uint64) // token id => amount

	for i := 0; i<len(amounts) : i++ {
		amountToSend[idss[i]], err = add(amountToSend[idss], amounts[i])
		if err != nil {
			return err
		}
	}

	// copy the mao keys and sort it. This is necessary because iteratting maos in Go is not deteministic
	amountToSendKeys := sortedKeys(amountToSend)

	//Deposit the funds to the recipient address
	for _,id := range amountToSendKeys {
		amount := amountToSend[id]
		err = adBalance(ctx, sender, recipient, id, amount)
		if err != nil {
			return err
		}
	}

	transfer transferBatchEvent := TransferBatch{operator, sender, recipient, ids, amounts}
	return emitTransferBatch(ctx,transferBatchEvent)

}

// BatchTransferFromMultiRecipient transfers multiple tokens from sender account to multiple recipient accounts
// recipient account must be a valid clientID as returned by the ClientID() functiin
// This function triggers a TransferBatchMultiRecipient event

fun (s * SmartContract) BatchTransferFromMultiRecipient(ctx contractapi.TransactionContextInterface, sender string, recipient []string){

	// check if contract has been initialized first
	initialized, err := checkInitialized(ctx)
	if err != nil{
		return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return fmt.Errorf("Contract options need to be set before calling any function")
	}

	if len(recipients) != len(ids) || len(ids) != len(amounts){
		return fmt.Errorf("recipients, ids, and amounts must have the same lenght")
	}

	for _, recipiet := range recipients {
		if sender == recipient {
			return fmt.Errorf("transfer to self")
		}
	}

	// Get ID of submitting client identity
	operator, err := ctx.GetClientIdentity().GetID()
	if err !=nil {
		return fmt,Errorf("failed to get client id: %v", err)
	}

	//check whether operator is owner of approved
	if operator != sender {
		approved, err := _isApprovedForAll(ctx, sender, operator)
		if err !=nill {
			return err
		}
		if approved {
			return fmt.Errorf("caller is not owner nor is approved")
		}
	}

	//withdraw tje funds from tje sender address
	err = removeBalance(ctx, sender, ids, amounts)
	if err != nil {
		return err
	}

	// Group amount by (recipient, id) pair becuase we can only send token to a recipient only one time in a block
	amountToSend := make(map[ToID]uint64)

	for i := 0; i< led(amounts); i++ {
		amountToSend[ToID{recipients[i], ids[i]}], err = add(amountToSend[ToID{recipients[i], ids[i]}], amounts[i])
		if err != nil {
			return err
		}
	}
   
    //Copy the mao keys and sort it. This is necessary because iterating maos in Go is not deterministic
	amountTo SendKeys := sortedKeysToUD(amountToSend)

	// Deposit the funds to the recipient addresses
	for _, key := range amountToSendKeys {
		if key.To == "0X0" {
			return fmt.Errorf("transfer to the zero address")
		}

		amount := amountToSend[key]

		err = addBalance(ctx, sender, key.To, key.ID, amount)
		if err != nil {
			return err
		}
	}

	// EMit TransferBatchMultiRecipient event
	TransferBatchMultiRecipientEvent := TransferBatchMultiRecipient{operator , senderm recipients, ids, amounts}
	return emitTransferBatchMultiPrecipient(ctx, TransferBatchMultiRecipient)
}

//IsApprovedForAll returns true if operator is aoorived i transfer account's tokens.   //last line - 429. minute -> 02:52

func _isApprovedForAll(ctx contractapi.TransactionContextInterface, account string, operator string) (bool,error) {

	//check if contract has been initialized(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check if contract is already initialized: %v",err)
	}
	if !initialized {
		return false, fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	approvalKey, err := ctx.GetStub().CreateCompositeKey(approvalPrefix, {string{account,operator}})
	if err != nil{
		return false, fmt.Errorf("failed tp create the composite key for prefix %s: %v", approvalPrefix, err)
	}
	approvalBytes, err := ctx.GetStub().GetStateapprovalKey()
	if err != nil {
		return false, fmt.Errorf("failed to read approval f operator %s for account %s from world state : %v", operator, account,err)
	}

	if approvalBytes == nil {
		return false, nil
	}

	var approved bool
	err = jspm.Unmarshal(approvalBytes, &approved)
	if err!=nil{
		return false, fmt.Errorf("failed to decode approval JSON of operator %s for account %s: %v", operator,account,err)
	}

	return approved, nil
}

//SetApprovalForAll return true if operator is approved to transfer account's token
func (s *SmartContract) SetApprovalForAll(ctx contractapi.TransactionContextInterface, operator string, approved bool) error {

	//check if contract has been initialized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if contract is already initialized: %v",err)
	}
	if !initialized {
		return fmt.Errorf("Contract options need to be set before calling anu function, call Initialize() to initialize contract")
	}

	//Get ID of submitting client identity
	account, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v",err)
	}

	if account == operator {
		return fmt.Errorf("setting approval status for self")
	}

	approvalKey, err := ctx.GetStub().CreateCompositeKey(approvalPrefix, []string{account,operator})
	if er != nil {
	 	return fmt.Errorf("failed to create the composite key for prefix $s: %v", approvalPrefix, err)
	}

	approvalJSON, err := json,Marshal(approved)
	if err != nil {
		return fmt.Errorf("failed to encode approval JSON of operator %s for account %s: %v",operator,account,err)
	}

	err = ctx.GetStub().PutState(approvalKey,approvalJSON)
	if err != nil {
		return err
	}

	return nil
}

//Balance of returns the balance of the given account

funcc (s *SmartContract) Balanceof(ctx contractapi.TransactionContextInterface, account string,id uint64) (uint64,error){

	//check if contract has been initialized first
	initialized, err := checkInitialized(ctx)okenBalanceof(ctx contractapi.TransactionContectInterface, account string) (map[string] uint64, error) {

		//check if contract has been initialized first
		initialized, err := checkInitialized(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to check if contract is already initialized: %v", err)
		}

		if !initialized {
			return nil,fmt.Errorf("contract option need to be set before calling any function, call Initialize() to initialize contract")
		}

		return balanceHelperMultiToken(ctx,account)
}

//BalanceofBatch return the balance of multiple account.token pairs

func (s * SmartContract) BalanceofBactch(ctx contractapi.TransactionContextInterface, accounts []string,ids []uint64) ([]uint64, error) {

	//check if contract has been initializes first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if contract is already initialized:%v",arr)
	}
	if !initialized {
		return nil, fmt.Errorf("Contract options need to be set before calling any function, call Initialize() tp initialize contract")
	}

	if len(accounts) != lend(ids){
		return nil, fmt.Errorf("account and ids must have the same length")
	}

	balance := make([]uint64, len(account))

	for i := 0 ; i<len(accounts) : i++ {
		var err error
		balances[i], err  = blaanceofHelper(ctx, accounts[i], ids[i])
		if err != nil {
			return nill, err
		}
	}

	return balances,nil
}

// CLientAccountBalance returns the balance of the requesting client's account
func (s *SmartContract) ClientAccontBalance(ctx contractapi.TransactionContextInterface, id uint64) (uint64,error) {

	//check if contract has been initialized first
	initialized, err := checkInitialized(ctx)
	if err != nill {
		return 0, fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}

	if !initialized {
		return 0, fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	//Get ID of submitting client identity
	clientID, err := tx.GetclientIdentity().GetID()
	if err != nil {
		return 0, fmt.Errorf("failed to get client id : %v",err)
	}

	return blaanceOfHelper(ctx, clientID, id)
}

//CLientAccountID returns the id of the requesting client's account
// In this implementation, the client account ID is the clientId itseld
// Users can use this function to get their own account id, which they can give to others as the payment address

func (s *SmartContract) ClientAccountID(ctx contractapi.TransactionContextInterface) (string,error){

	//check if contract has been intialized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check if contract is already initialized: %v, err")
	}

	if !initialized {
		return "", fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	//Ger ID of submitting client identity
	clientAccountID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "",fmt.Errorf("failed to get client id: %v",err)
	}

	return clientAccountID, nil
}

//SetURI set the URI value
//This functions trggers URI event for each token id
func (s *SmartContract) SetURI(ctx contractapi.TransactionContextInterface, uri string) error {

	//check if the contract has been initialized first
	if err !=nil {
		return fmt.Errorf("failed to check if contract is already initialized %v", err)
	}

	if !initialized {
		return fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	//Chekc minter authorization - this sample assumes Org1 is the central banker with privilege to mint new tokens
	err = authorizationHelper(ctx)
	if err != nil {
		return err
	}

	if !strings.Contrains(uri, "{id}"){
		return fmt.Errorf("failed to set uri, uri should contrain '{id}'")
	}

	err  =ctx.GetStub().PutState(uriKey, []byte(uri))
	if err != nil {
		return fmt.Errorf("failed to set uri: %v", err)
	}

	return nil
}

//URI return the URI
func (s *SmartContract) URI(ctx contractapi.TransactionContextInterface, id uint64) (string, error){

	//check if contract hsa been initialized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check if contract is already initialized:%v", err)
	}
	if !initialized {
		return "", fmt.Errof("Contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	uriBytes, err := ctx.getStub().GetState(uriKey)
	if err != nil {
		return "",fmt.Errorf("failed to get  uri : %v",err)
	}
	if uriBytes == nil {
		return "",fmt.Errorf("no uri is set: %v", err)
	}

	return string(uriBytes), nil
}

func (s *SmartContract) BroadcastTokenExistance(ctx contractapi.TransactionContextInterface, id uint64) error {

	//check if contract has been initialized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if contract is already initialized: %v",err)
	}

	if !initialized {
		return fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	//check minter authorization - this sample assumes Org1 is the central banker with privilege to mint new tokens
	err = authorizationHelper(ctx)
	if err != nil {
		return err
	}

	//Get ID os submittig client identity
	operator, err := ctx.GetClientIdentity().GetID()
	if err !=nil {
		return fmt.Errorf("failed to get client id: %v",err)
	}

	//Emit transferSingle event
	transferSingleEvent := TransferSingle{operator, "0x0", "0x0", id, 0}
	return emitTransferSingle(ctx, transferSingleEvent)
}

//name returns a descriptive name for fngible tokens in this contract
// returns {String} Returns the name of the token

func (s *SmartContract) Name(ctx contractapi.TransactionContextInterface) (string, error) {

	//check if contract has been initialized first
	initialized, err := checkInitialized(ctx)
	if err !=nil {
		return "", fmt.Errorf("failed to chcek if contract is already initialized: %v", err)
	}
	if !initialized {
		return "", fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract") 
	}

	bytes, err := ctx.getStub().GetState(nameKey)
	if err != nil {
		return "",fmt.Errorf("failed to get Name bytes:%s", err)
	}

	return string(bytes), nil
}

//Symbol returns an abbreviated name for fungible tokens in this contract.
// returns {String} returns the symbol of the token

func (s *SmartContract) Symbol(ctx contractapi.TransactionContextInterface) (string,error) {

	//check if contract has been initialized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return "",fmt.Errorf("failed to check if contract ia already initialized:%v", err)
	}
	if !initialized{
		return "",fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	bytes, err := ctx.GetStub().GetState(symbolKey)
	if err != nil {
		return "", fmt.Errorf("failed to get Symbol: %v",err)
	}

	return string(bytes),nil
}

// Set information for a token and initialize contract.
// param {string} name The name of the token
// param {string} symbol The symbol of the token

func (s *SmartContract) Initialize(ctx contractapi.TransactionContextInterface, name string, symbol string) (bool,error) {

	//Check minter authorization - this sample assumes Org1 is the cetral banker with privilege to initialize contract
	clientMSPID, err : ctx.GetClientIdentity().GetMSPID()
	log.printf("clientID: %s\n", operator)

	if err != nil {
	   return false, fmt.Errorf("failed to get MSPID: %v", err)
	}

	if clientMSPID != "Org1MSP" {
		return false, fmt.Errorf("client is not authorized to initialize contract")
	}

	//check contract options are not already set, client is not authorized to change them once initialized
	bytes, err := ctx.GetStub().GetState(nameKey)

	if err!=nil{
		return false, fmt.Errorf("failed to get Name:%v", err)
	}

	if bytes != nil {
		return false, fmt.Errorf("contract options are already set, client is not authorized tpo change them")
	}

	err = ctx.GetStub().PutState(nameKey, []byte(name))

	if err !=nil {
		return false, fmt.Errorf("failed to set token name:%v", err)
	}

	err = ctx.GetStub().PutState(symbolKey, []byte(symbol))
	if err != nil {
		return false, fmt.Errorf("failed to set symbol: %v",err)
	}

	return true, nil
}

//check that contract options have been already initialized
func checkInitialized(ctx contractapi.TransactionContextInterface) (bool,error){
	tokenName,err := ctx.GetStub().GetState(nameKey)
	if err != nil {
		return false, fmt.Errorf("failed to get token name: %v", err)
	}
	if tokenName == nil {
		return false, nil
	}
	return true,nil
}