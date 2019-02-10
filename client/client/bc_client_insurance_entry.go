package client

import (
	"log"
	"fmt"
	"strconv"
	"errors"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"blockchain_at_insurtech/client/model"
)

func CreateInsuranceEntry(insuranceEntry model.InsuranceEntryType, withUser string) (string, error) {
	if len(withUser)==0 {
		withUser = blockchainUser
	}
	clientContext, err := GetOrgChannelContext(channelID, org1, withUser)
	if err != nil {
		log.Printf("Failed to create new channel client for user: %s", err)
		return "", nil
	}

	// get client
	client, err := channel.New(clientContext)
	if err != nil {
		log.Printf("Failed to create new channel client for user: %s", err)
		return "", nil
	}

	args := [][]byte{
		[]byte(insuranceEntry.Id),
		[]byte(insuranceEntry.Item),
		[]byte(insuranceEntry.StartDate),
		[]byte(strconv.FormatInt(insuranceEntry.NumOfDays, 10)),
		[]byte(strconv.FormatInt(insuranceEntry.Amount, 10)),
		[]byte(insuranceEntry.ProductId),
	}

	// create insurance entity
	response, err := client.Execute(channel.Request{ChaincodeID: chaincodeName, Fcn: "createInsuranceEntry", Args: args})
	if err != nil {
		msg := fmt.Sprintf("failed to create insurance entry: %+v\n", err)
		return "", errors.New(msg)
	}
	resPayload := response.Payload

	return string(resPayload), nil
}

func RetrieveInsuranceEntry(id string, withUser string) (string, error) {
	if len(withUser)==0 {
		withUser = blockchainUser
	}
	client, err := getOrgChannelClient(peer_url, event_url, orderer_url, channelID, org1, withUser)
	if err != nil {
		msg := fmt.Sprintf("Failed to create new channel client for user: %s", err)
		log.Printf(msg)
		return "", err
	}

	args:=[][]byte{
		[]byte(id),
	}
	// create entity
	res, err := client.Execute(channel.Request{ChaincodeID: chaincodeName, Fcn: "retrieveInsuranceEntry", Args: args})
	if err != nil {
		msg:= fmt.Sprintf("failed to get insurance entry: %+v\n", err)
		fmt.Printf(msg)
		return "", err
	}
	resPayload := string(res.Payload)

	return resPayload, nil
}

func RetrieveInsuranceEntries(withUser string) (string, error) {
	if len(withUser)==0 {
		withUser = blockchainUser
	}
	client, err := getOrgChannelClient(peer_url, event_url, orderer_url, channelID, org1, withUser)
	if err != nil {
		msg := fmt.Sprintf("Failed to create new channel client for user: %s", err)
		log.Printf(msg)
		return "", err
	}

	args:=[][]byte{	}
	// retrieve insurance entries
	res, err := client.Execute(channel.Request{ChaincodeID: chaincodeName, Fcn: "retrieveInsuranceEntries", Args: args})
	if err != nil {
		msg:= fmt.Sprintf("failed to get insurance entries: %+v\n", err)
		fmt.Printf(msg)
		return "", err
	}
	resPayload := string(res.Payload)

	return resPayload, nil
}
