package fabric_go_sdk

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	//"time"
)

// Invoke,传入args包括func和参数
func (setup *FabricSetup) Invoke(args []string) (string, error) {
	var inargs [][]byte
	if len(args) > 1 {
		for _, i := range args {
			inargs = append(inargs, []byte(i))
		}
		inargs = inargs[1:]
	}
	if !setup.initialize{
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

	//eventID := "eventInvoke"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in invoke")
	//
	//reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
	//if err != nil {
	//	return "", err
	//}
	//defer setup.event.Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: inargs, TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	// Wait for the result of the submission
	//select {
	//case ccEvent := <-notifier:
	//	fmt.Printf("Received CC event: %v\n", ccEvent)
	//case <-time.After(time.Second * 20):
	//	return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	//}

	setup.initialize = true
	return string(response.TransactionID) +"  "+ string(response.Payload), nil
}