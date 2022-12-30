package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
	"strconv"
)

const escroBalancePrefix = "escrow~agent~depositor~tokenId~beneficiary~reference"

//EscrowSingle MUST emit when a single token is escro
//The operator arument MUST be msg.sender.
// The depositor argument MUST be the address of the holder whose balance is decreased.
//The agent argument MUST be the address who will hold the escrow token.
// The beneficiary argument MUST be the address of the recipient whose balance will be increased after escro release.
// The id argument MUST be the token type beign transferred
//The reference argument MUST be uniqye for all escro amounf depositor, agent and beneficiary.
//The value argument MUST be the number of tokens the holder balance is decreased
// by and match what the agent balance is increased by as escrow.

type EscrowSingle struct (
	Operator  		string `json:"operator`
	Depositor 		string `json:"depositor`
	Agent	  		string `json:"agent`
	Beneficiary 	string `json:"beneficiary`
	ID              uint64 `json:"id`
	Reference		string `json:"reference`
	Value			uint64 `json:"value`
)

/*EscrowDetails will hold the escrow detail information*/

type EscrowDetails struct {
	Depositor 	string `json:"depositor`
	Agent		string	`json:"agent"`
	Beneficiary	string	`json:beneficiary`
	ID          uint64	`json:"id"`
	Reference	string	`json:"reference`
	Value		uint64	`json:"value"`
}

// Escrow transfers tokens from sender account to admin account
// admin and recipient account must be a valid clientID as returned by the CLientID() funstin
// This function triggers a TransferSingle event

func (s *SmartContract) EscrowFrom(ctx contractapi.TransactionContecxtInterface, depositor string, agent string, beneficiary string, id uint64){ //something left at right

	//check if contract has been initialized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if contract ia already initialized: %v", err)
	}

	if !initialized {
		return fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	if depositor == beneficiary || depositor == agent {
		return fmt.Errorf("transfer to self")
	}

	//Get ID of submitting client identity
	operator, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v",err)
	}

	//check whether operator is owner or approved
	if operator != depositor {
		approved, err := _isApprovedForAll(ctx,depositor,operator)
		if err != nill {
			return err
		}

		if !approved {
			return fmt.Errorf("caller is not owner not is approved")
		}
	}

	if beneficiary == "0x0" {
		return fmt.Errorf("transfer to the zero address")
	}

	//withdraw the funds from the depositor address
	err = removeBalance(ctx,depositor, []uint64{id}, []uint64{amount})

	if err != nil {
		return err
	}

	//Deposit the fund to the agent address as escrow
	err = addEscrowBalance(ctx,depositor,agent,beneficiary,id,reference, amount)
	if err != nil {
		return err
	}

	// Emit escrowSingleEvent event
	escrowSingleEvent := EscrowSingle{ operator, depositor, agent, beneficiary, id, reference, amount }
	return emitEscrowSingle(ctx.escrowSingleEvent)
}

func addEscrowBalance(ctx contractapi.TransactionContextInterface, depositor string, agent string, beneficiary string, id uint64, reference string) { //sometext were cuttoff from screen

	//Convert id to string
	idString := strconv.FormatUnit(uint64(id),10)
	escrowBalanceKey, er := ctx.GetStub().CreateCompositeKey(escrowBalancePrefix, []string {agent, depositor, idString, beneficiary,reference})
	if err != nil {
		return fmt.Errorf("failed to read account %s from world state: %v", depositor, err)
	}

	balanceBytes, err := ctx.GetStub().GetState(escrowBalanceKey)   //line number 111
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(escrowBalanceKey, []byte(strconv,FormatUint(uint64(balance),10)))
	if err != nil {
		return err
	}

	return nil
}

//ReleaseEscroDrom transfers token from sender account to admin account
//admin and reciient account must be a valid clientID as returned by the ClientID() function
//this finction trigger a TransferSingle event

func (s *SmartContract) ReleaseEscroFrom(ctx contractapi.TransactionContextInterface, depositor string, agent string, beneficiary string, id uint64){ //some texts were cutt off

	//check if contract has been initialized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}

	if !initialized {
		return fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	if receiver != beneficiary && recover != depositor {
		return fmt.Errorf("receiver address %s not the depositor %s nor the beneficiary %s", receiver, depositor, beneficiary)
	}

	//Get ID of submitting client identity
	operator, err := ctx.GetClientIdentity().GetID()
	if err!=nil {
		return fmt.Errof("failed to get client id:%v", err)
	}

	//check whether operator is owner or approved
	if operator != agent {
		approved, err := _isApprovedForAll(ctx,agent,operator)
		if err != nil {
			return err
		}

		if !approved {
			return fmt.Errorf("caller is not owner nor is approved")
		}
	}

	if rec == "0x0" {
		return fmt.Errorf("transfer to zero address")
	}

	//check escrow balance
	balance, err := s.GetEscrowBalance(ctx,depositor,agent,beneficiary,id,reference)
	if err != nil {
		return err
	}

	// withdraw the escrow funds from the agent address
	err = remoeEscrowbalance(ctx,depositor,agent,beneficiary,idString,reference)
	if err != nil {
		return err
	}

	// Deposit the fund to the recipient address
	err = addBalance(ctx, depositor,receiver,idString,balance)
	if err!= nil {
		return err
	}

	// Emit TranferSingle event
	transferSingleEvent := TransferSingle(operator,depositor,receiver,id,balance)
	return emitTransferSingle(ctx,transferSingleEvent)
}

func removeEscrowBalance(ctx contractapi.TransactionContextInterface, depositor string, agent string, beneficiary string, id uint64, reference string){ // some text were cutt off

	//Convert id to string
	idString := strconv.FormatUint(uint64(id),10)

	escrowBalanceKey, err := ctx.GetStub().CreateCompositeKey(escrowBalancePrefix, []string(agent,depositor,idString,beneficiary,reference))
	if err != nil {
		return fmt.Errorf("failed to create the composite key for prefix %s: %v, reference,err")
	}

	if balanceBytes != nil {
		err = ctx.GetStub().DelState(escrowBalanceKey)
		if err != nil {
			return fmt.Errorf("failed to delete the state if %v: %v",escrowBalanceKey,err)
		}
	}

	return nil
}

func (s *SmartContract) GetEscrowBalance(ctx contractapi.TransactionContextInterface, depositor string, agent string, beneficiary string, id uint64){ //some text were cutt off

	//Convert if to string
	idString := strconv.FormatUint(uint64(id),10)

	balanceKey, err := ctx.GetStub().CreateCompositeKey(escrowBalancePrefix, []string(agent, depositor, idString, beneficiary, reference))
	if err != nil {
		return 0, fmt.Errorf("failed to create the composite key for prefix %s: %v", escrowBalancePrefix, err)
	}

	balanceBytes, err := ctx.GetStub().GetState(balanceKey)
	if err != nil {
		return 0, fmt.Errorf("failed to read escrow with reference %s from world state: %v", reference, err)
	}

	if balanceBytes != nil {
		balanceBytes, _ := strconv.ParseUint(string(balanceBytes),10,64)
		return balance, nil
	}

	return 0, fmt.Errorf("failed to read escrow with reference %s from world state: %v", reference, err)
}

/*Get All the escrow for a depositor*/

func (s *SmartContract) GetEscrowBalanceOfDepositor(ctx contractapi.TransactionContextInterface, depositor string, agent string ({}EscrowDetails)){  //some text were cutt off

	//Array to hold all escrow for the user

	var escrowDetails []EscrowDetails

	if depositor == "0x0" {
		return escrowDetails, fmt.Errorf("escrow query for the zero address")
	}

	escrowBalanceIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(escrowBalancePrefix, []string(agent,depositor))
	if err != nil {
		return escrowDetails, fmt.Errorf("failed to get escrow for user %v: %v",depositor,err)
	} 

	defer escrowBalanceIterator.close()

	for escrowBalanceIterator.HasNext() {
		queryResponse, err := escrowBalanceIterator.Next()
		if err != nil {
			return escroDetails, fmt.Errorf("failed to get the next state for prefix %v: %v", escrowBalancePrefix, err)
		}
		_, compositeKeyParts, err := ctx.GetStub().SplitCompositKey(queryResponse.Key)
		if err != nil {
			return escrowDetails, err
		}

		log.Println("composite key parts ",compositeKeyParts)
		// parse the cinposite key
		// "escrow~agent~depositor~tokenId~beneficiary~reference"
		if queryResponse.Value != nil {
			balance, _ := strconv.ParseUint(string(queryResponse, Value),10,64)
			tokenId,_ := strconv.ParseUint(string(compositeKeyParts[2],10,64))

			escrowDetail := EscrowDetails {
				Agent:		compositeKyeParts[0],
				Depositor;	compositeKeyParts[1],
				ID:			tokenId,
				Beneficiary: compositeKyeParts[3],
				Reference : compositeKyeParts[4],
				Value: balance
			}
			escrowDetails = append(escroDetails, escrowDetail)
		}
	}
	log.Println("escrowDetails ",escrowDetails)
	return escrowDetails, nil
}


/*emit an escrow event durin an escrow operator*/

func emitEscrowSingle(ctx contractapi.TransactionContextInterface, escrowSingleEvent EscrowSingle) error {
	escrowSingleEventJSON, err := json.Marshal(escrowSingleEvent)
	if err != nil {
		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}

	err = ctx.GetStub().SetEvent("EscrowSingle", escrowSingleEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}
	return nil
}