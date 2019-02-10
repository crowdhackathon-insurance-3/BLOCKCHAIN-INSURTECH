package client

import (
	"os"
	"log"
	"fmt"
	"bytes"
	"errors"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	contextApi "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	org1          = "Org1"
	channelID     = "channel1"
	//chaincodeName = "mycc"
	//user 		  =	"Admin"
	//user 		  =	"fi1"
	//user 		  =	"admin2"
	chaincodeName = "cc"
	blockchainUser 		  =	"testfi1"
)

const (
	ca_server = "ca.org1.example.com"
	//peer_url = "peer1.org1.example.com"
	//orderer_url = "orderer.example.com"
	//peer_url = "192.168.9.51"
	//orderer_url = "192.168.9.51"
	peer_url          = "127.0.0.1:7051"
	event_url          = "127.0.0.1:7053"
	orderer_url       = "127.0.0.1:7050"
	adminEnrollmentId = "admin2"
)

func loadConfigBytesFromFile(filePath string) ([]byte, error) {
	// read test config file into bytes array
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to read config file. Error: %s", err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		log.Fatalf("Failed to read config file stat. Error: %s", err)
	}
	s := fi.Size()
	cBytes := make([]byte, s)
	n, err := f.Read(cBytes)
	if err != nil {
		log.Fatalf("Failed to read test config for bytes array testing. Error: %s", err)
	}
	if n == 0 {
		log.Fatalf("Failed to read test config for bytes array testing. Mock bytes array is empty")
	}
	return cBytes, err
}

func getConfigFile() (core.ConfigProvider, error) {
	filename := "./config/config.yaml"
	log.Printf("loading file: " + filename)
	configBytes, err := loadConfigBytesFromFile(filename)
	if err != nil {
		return nil, err
	}
	file := config.FromRaw(configBytes, "yaml")
	return file, nil
}
func GetClientContext(channelID string, org string, userName string) (*fabsdk.FabricSDK, error) {

	configFile := fmt.Sprintf("config/config.yaml")
	//file := config.FromFile(configFile)
	//
	configBytes, err := loadConfigBytesFromFile(configFile)
	if err != nil {
		return nil, err
	}
	configBytesWithUser := bytes.Replace(configBytes, []byte("{USERNAME}"), []byte(userName), -1)
	configBytesWithPeer := bytes.Replace(configBytesWithUser, []byte("{PEER_URL}"), []byte(peer_url), -1)
	configBytesWithEvent := bytes.Replace(configBytesWithPeer, []byte("{EVENT_URL}"), []byte(event_url), -1)
	configBytesWithOrderer := bytes.Replace(configBytesWithEvent, []byte("{ORDERER_URL}"), []byte(orderer_url), -1)
	configBytesFinal := bytes.Replace(configBytesWithOrderer, []byte("{CA_URL}"), []byte(""), -1)
	//fmt.Println(string(configBytesFinal))
	file := config.FromRaw(configBytesFinal, "yaml")

	sdk, err := fabsdk.New(file)
	return sdk, err
}
func GetOrgChannelContext(channelID string, org string, userName string) (contextApi.ChannelProvider, error) {

	configFile := fmt.Sprintf("config/config.yaml")
	//file := config.FromFile(configFile)
	//
	configBytes, err := loadConfigBytesFromFile(configFile)
	if err!=nil {
		return nil, err
	}
	configBytesWithUser := bytes.Replace(configBytes, []byte("{USERNAME}"), []byte(userName), -1)
	configBytesWithPeer := bytes.Replace(configBytesWithUser, []byte("{PEER_URL}"), []byte(peer_url), -1)
	configBytesWithEvent := bytes.Replace(configBytesWithPeer, []byte("{EVENT_URL}"), []byte(event_url), -1)
	configBytesWithOrderer := bytes.Replace(configBytesWithEvent, []byte("{ORDERER_URL}"), []byte(orderer_url), -1)
	configBytesFinal := bytes.Replace(configBytesWithOrderer, []byte("{CA_URL}"), []byte(""), -1)
	//fmt.Println(string(configBytesFinal))
	file := config.FromRaw(configBytesFinal, "yaml")

	sdk, err := fabsdk.New(file)
	if err != nil {
		log.Printf("Failed to create new SDK: %s", err)
		msg:= fmt.Sprintf("Failed to create new SDK: %+v\n", err)
		return nil, errors.New(msg)
	}
	//defer sdk.Close()

	//prepare contexts
	userOption := fabsdk.WithUser(userName)
	orgOption := fabsdk.WithOrg(org)
	org1ChannelClientContext := sdk.ChannelContext(channelID, userOption,  orgOption)

	return org1ChannelClientContext, err
}

func getOrgChannelClient(peer_url string, event_url string, orderer_url string, channelID string, org string, userName string) (*channel.Client, error) {

	org1ChannelClientContext, err := GetOrgChannelContext(channelID, org, userName)
	if err != nil {
		log.Printf("Failed to create new client: %s", err)
		msg:= fmt.Sprintf("Failed to create new client: %+v\n", err)
		return nil, errors.New(msg)
	}

	client, err := channel.New(org1ChannelClientContext)

	return client, err
}
