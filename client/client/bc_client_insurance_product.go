package client

import (
	"log"
	"fmt"
	"strconv"
	"errors"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"blockchain_at_insurtech/client/model"
)

func CreateInsuranceProduct(insuranceProduct model.InsuranceProductType, withUser string) (string, error) {
	if len(withUser)==0 {
		withUser = blockchainUser
	}
	clientContext, err := GetOrgChannelContext(channelID, org1, withUser)
	if err != nil {
		log.Printf("Failed to create new channel client for insuranceProduct: %s", err)
		return "", errors.New("failed to create new channel client for insuranceProduct")
	}

	// get client
	client, err := channel.New(clientContext)
	if err != nil {
		log.Printf("Failed to create new channel client for insuranceProduct: %s", err)
		return "", errors.New("failed to create new channel client for insuranceProduct")
	}

	args:=[][]byte{
		[]byte(insuranceProduct.Id),
		[]byte(insuranceProduct.Name),
		[]byte(insuranceProduct.StartDate),
		[]byte(insuranceProduct.EndDate),
		[]byte(strconv.FormatInt(insuranceProduct.MinAmount, 10)),
		[]byte(strconv.FormatInt(insuranceProduct.MaxAmount, 10)),
		[]byte(insuranceProduct.Terms),
		[]byte(strconv.FormatInt(insuranceProduct.DRate, 10)),
	}
	// create insurance product
	res, err := client.Execute(channel.Request{ChaincodeID: chaincodeName, Fcn: "createInsuranceProduct", Args: args})
	if err != nil {
		msg:= fmt.Sprintf("failed to create insuranceProduct: %+v\n", err)
		return "", errors.New(msg)
	}
	resPayload := res.Payload

	return string(resPayload), nil
}

func RetrieveInsuranceProduct(id string, withUser string) (string, error) {
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
	res, err := client.Execute(channel.Request{ChaincodeID: chaincodeName, Fcn: "retrieveInsuranceProduct", Args: args})
	if err != nil {
		msg:= fmt.Sprintf("failed to get insurance product: %+v\n", err)
		fmt.Printf(msg)
		return "", err
	}
	resPayload := string(res.Payload)

	return resPayload, nil
}

func RetrieveInsuranceProducts(withUser string) (string, error) {
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
	// retrieve insurance products
	res, err := client.Execute(channel.Request{ChaincodeID: chaincodeName, Fcn: "retrieveInsuranceProducts", Args: args})
	if err != nil {
		msg:= fmt.Sprintf("failed to get insurance products: %+v\n", err)
		fmt.Printf(msg)
		return "", err
	}
	resPayload := string(res.Payload)

	return resPayload, nil
}
