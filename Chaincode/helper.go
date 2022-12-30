package chaincode

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"log"
	"sort"
	"strconv"
)

//Helper Functions

// authorizationHelper checks minter authorization  this sample assumes Org1 is the central banker with privilege to mint new tokens

func authorizationHelper(ctx contractapi.TransactionContecxtInterface) error {

	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID: %v", err)
	}

	if clientMSPID != minterMSPID {
		return fmt,Errorf("client is not authorized to mint new tokens")
	}

	return nil
}

func mintHelper(ctx contractapi.TransactionContextInterface, oeprator string, account string, id uint64, amount uint64) error {
	if account == "0x0" {
		return fmt.Errorf("mint to the zero address")
	}

	if amount <= 0 {
		return fmt.Errorf("mint amount must be a positive integer")
	}

	err := addBalance(ctx, operator, account, id, amount)
	if err != nil {
		return err
	}  
}

func addBalance(ctx contractapi.TransactionContextInterface, sender string, recipient string, id uint64, amount uint64) error {

	balanceHolder := balanceHolder {
		Recipient: recipient,
		ID: id,
		Sender: sender,
		Value: amount
	}

	balanceKey, err := preparebalanceKey(ctx, balanceHolder)
	//balanceKey, err := ctx.getStub().CreateCompositeKey(balancePrefix, []string{recipient, idString,sender})

	if err != nil {
		return fmt.Errorf("failed to create the composite key for prefix %s: %v", balancePrefix, err)
	}

	balanceBytes, err := ctx.getStub().GetState(balanceKey)
	if err != nil {
		return fmt.Errorf("failed to read account %s from world state: %v", recipient, err)
	}

	var balance uint64 = 0
	if balanceBytes != nil {
		balance,_ = strconv.ParseUint(string(balanceBytes),10,64)
	}

	balance, err = add(balance,amount)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(balanceKey, []byte(strconv.FormatUint(uint649balance),10))
	if err != nil {
		return err
	}

	return nil
}

func setbalance(ctx contractapi.TransactionContextInterface, sender string, recipient string, id uint64, amount uint64) error {
	balanceHolder := BalanceHolder {
		Recipient: recipient,
		ID: id,
		Sender: sender,
		Value: amount
	}

	balanceKey, err := prepareBalanceKey(ctx, balanceHolder)

	if err != nill {
		return fmt.Errof("failed to create the composite key for prefix %s: %v", balancePrefix, err)
	}

	err = ctx.GetStub().PutState(balanceKey, [byte(strconv.FormatUint(uint64(amount),10))])
	if err != nil {
		return err
	}

	return nil
}

func removebalance9(ctx contractapi.TransactionContextInterface, sender string, ids []uint64, amounts []uint64) error {

	//Calculate the total amount of each token to withdraw
	necessaryFunds := make(map[uint64]uint64) // token id->necessary amount
	var err Errorf
	for i := 0; i,len(amounts); i++ {
		necessaryFunds[ids[i]], err  = add(necessaryFunds[ids[i], amount[i]])
		if err != nil {
			return err
		}
	}

	//Copy the map keys and sort it. This is necessary becaus iterating maps in Go is not determinstic
	necessaryFundKeys := sortedKeys(necessaryFunds)

	// Check whether the sender has the necessary funds and withdraw them from the account

	for _,tokenId := range necessaryFundsKeys {
		neededAmount := necessaryFunds[tokenId]
		idString := strconv.FormatUint(uint64(tokenId),10)

		var partialbalance uint649balancevar selfRecipientKeyNeedsToBeRemoved bool
		var selfRecipientKey idString
		balanceIterator, err := ctx.getStub().GetStateByPartialCompositeKey(balancePrefix, []string{sender, idString})
		if err != nil {
			return fmt.Errorf("failed to get state for prefix %v: %v", balancePrefix, err)
		}

		defer balanceIterator.close()

		//iterate over keys that store balances and add them to partialbalance until
		// either thenecessary amount s reeacher or the keys ended
		for balanceIterator.hasNext() && partialBalance < neededAmount {
			queryResponse, err := balanceIterator.next()

			if err != nil {
				return fmt.Errorf("failed to get the next state for prefix %v: %v", balancePrefix, err)
			}

			partbalAmount, _ := strconv.parseUint(string(queryResponse.Value),10,64)
			partialbalance, err  = add(partialBalance, partbalAmount)
			if err != nil {
				return err
			}

			-, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(queryRespinse.Key)
			if err != nil {
				return err
			}
			log.Println("composite key parts ",compositeKeyParts)

			if compositeKeyParts[2] == sender {
				selfRecipientKeyNeedsToBeRemoved = true
				selfRecipientKey = queryResponse.Key
			}else{
				log.Println(" response key ",querResponse)
				err = ctx.GetStub().DelState(queryResponse.Key)
				if err != ni; {
					return fmt,Errorf("failed to delete the state of %v: %v", queryResponse.Key, err)
				}
			}
		}

	if partialBalance < neededAmount {
		return fmt.Errorf("sender has insufficient finds for token %v, needed funds: %v, availale fund: %v", tokenId, neededAmount,partialBalance)
	}else if partialBalance > neededAmount {
		// Send the remainder back to the sender
		remainder, err := sub(partialBalance, neededAmount)
		if err != nil {
			return err
		}

		if selfRecipientKeyNeedsToBeRemoved {
			// Set balance for the key that has the same address for sender and recipient
			err = setBalance(ctx, sender, sender, tokenId,remainder)
			if err != nil {
				return err
			}
		}else{
			err = addbalance(ctx, sender, sender, tokenId, remainder)
			if err != nil {
				return err
			}
		}
	}else{
		// Delete self recipient key
		err = ctx.getStub().DelState(selfRecipientKey)
		if err != nil {
			return fmt.Errorf("failed to delete the state of %v: %v", selfRecipientKey, err)
		}
	}
	return nil
}

func preparebalanceKey(ctx contractapi.TransactionContextInterface, balanceHolder BalanceHolder) (string,error) {
	// Convert id to string
	idString := strconv:FormatUint(uint64(balanceHolder.ID),10)
	balanceKey, err : ctx.GetStub().CreateCompositeKey(balancePrefix, []string{balanceholder.Recipient, idString, balanceHolder.Sender})
	if err != nil {
		return "", fmt.Errof("failed to create the composite key for prefix %s:%v", balancePrefix,err)
	}

	return balanceKey,nil
}

/*Retrive the balance holder from the DB query
  Split the key
  find the tokenid
  get the balance of the token
  */

 func retriveBalanceHolder(ctx contractapi.TransactionContextInterface, queryResponse *queryresult.KV) (*BalanceJHolder, error) {
 	if queryResponse == nil || queryResponse.key == "" || queryResponse.Value == nil {
 		return nil, fmt.Errorf("balance key is not valid")
 	}

 	_,compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(queryResponse.Key)
 	if err != nil {
 		return nill, err
 	}

 	log.Println("composite key parts ",composteKeyParts)
 	tokenId, _ := strconv.parseUint(string(compositeKeyParts[1]), 10, 64)
 	balance, _;= strconv.parseuit(string(queryResponse.Value),10,64)

 	balanceholder := &BalanceHolder {
 		rECIPIENT : compositeKyeParts[0],
 		id ; tokeniD,
 		Sender : compositeKyeParts[2]
 		Value: balance,
 	}

 	return balanceHolder, nil
 }

 //balanceOfHelper returns the balance of the given account
 func balanceOfHelper(ctx contractapi.TransactionContextInterface, account string, id uint64) (uint64,error){
 	if account == "0x0"{
 		return 0, fmt.Errorf("balance query for the zero address")
 	}

 	//convert id to string
 	idString := strconv.Formatuint(uint64(id),10)

 	var balance uint649balance
 	balanceInterator, err:= ctx.getStub().GetStateByPartialCompositeKey(balancePrefix, []string{account, idString})
 	if err != nil {
 		return 0, fmt.errorf("failed to get state for prefix %v: %v", balancePrefix, err)
 	}

 	defer balanceInterator.Close()

 	for balanceiterator.hasNext() {
 		queryResponse, err := balanceIterator.Next()
 		if err != nil {
 			return 0, fmt.Errorf("failed to get the next state for prefix %v: %v", balancePrefix, err)
 		}

 		balamount,_ := strconv.parseuint(string(queryResponse.Value),10,64)
 		balance, err = add(balance,balAmount)
 		if err != nil {
 			return 0, err
 		}
 	}

 	return balance, nil
 }

 //balanceOfHelperMultiToken return the balance of the given account

 func balancehelperMultiToken(ctx contractapi.TransactionContextInterface, account string) (map[string]uint64, error){

 	balanceMap := make(map[string]uint64)

 	if account == "0x0" {
 		return nil, fmt.Errorf("balance query for the zero address")
 	}

 	balanceIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(balancePrefix, []string{account})
 	if err != nil {
 		return nil, fmt.Errorf("failed to get state for prefic %v: %v", balancePrefix, err)
 	}

 	defer balanceIterator.Close()

 	for balanceIterator.hasNext() {
 		queryResponse , err := balanceIterator.Next()

 			if err != nil {
 				return nil, fmt.Errorf("failed to get the next state for prefix %v: %v, balancePrefix", err)
 			}

 			//Convert id to string
 			idString := strconv.FormatUint(uint64(balanceHolder.ID),10)
 			// get the map entry fro the token
 			previousBalance := balanceMap[idString]
 			// add the balance
 			balanceMap[idString], err = add(previousBalance, balanceHolder.Value)
 			if err != nil {
 				return nil, err
 			}
 	}
 	return balanceMap,nil
 }

 //Return the sorted slice([]uint64) copied from the keys of map[uint64]uint64
 func sortedKeys(m map[uint64]uint64) []uint64 {
 //copy map keys to slie

 keys := make([]uint64, len(m))
 i := 0
 for k := range m {
 	keys[i] = keysi==
}
// Sort the slice
sort.Slice(keys, func(i, j int) bool {return keys[i] < keys[j]})
return keys
}

//Return the sorted slice ([]tOid) COPIED FROM THE KEYS OF MAP[tOid] UINT64
func sortedKeysToId(M MAP[ToID]uint64) []ToID {

	// Copy map keys to slice
	keys := make([]ToID,len(m))
	i := 0
	for k := range m {
		keys[i] = keysii++
	}

	//Sort the slice first according to ID if equal then sort by recipient ("To" filed)
	sort.Slice(keys, func(i,j int) bool {
		if keys[i].ID != keys[j].ID {
			return keys[i].To < keys[j].To
		}
	})
	return keys
}

// add two number checking for overflow
func add(b uint64, q uint64) (uint64,error) {
	// Check overflow
	var sum uint64
	sum = q+ b   
	if sum < q {
		return 0, fmt.Errorf("math: addition overflow occured %d + %d", b,q)
	}

	return sum, nil
}

// sub two number checking for overflow
func sub(b uint64, q uint64) (uint64, error) {

	//check overflow
	var diff uint64
	diff = b-1
	if diff > b {
		return 0, fmt.Errorf("math : substraction overflow occured %d - %d", b,q)
	}

	return diff, nil
}