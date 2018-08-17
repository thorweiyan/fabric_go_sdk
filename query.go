package fabric_go_sdk

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// Query query the chaincode to get the state of ids
func (setup *FabricSetup) Query(args []string) (string, error) {
	var inargs [][]byte
	if len(args) > 1 {
		for _, i := range args {
			inargs = append(inargs, []byte(i))
		}
		inargs = inargs[1:]
	}

	if !setup.initialize {
		err := setup.initsdk()
		if err != nil {
			return "", err
		}
		err = setup.initclient()
		if err != nil {
			return "", err
		}
		err = setup.initevent()
		if err != nil {
			return "", err
		}
	}

	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: inargs})
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}

	setup.initialize = true
	return string(response.Payload), nil
}