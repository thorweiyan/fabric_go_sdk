package fabric_go_sdk

import (
	"testing"
	"os"
	"fmt"
)

//Just a example, need environment
func TestFabricSetup_Initialize(t *testing.T) {
	// Definition of the Fabric SDK properties
	fSetup := FabricSetup{
		// Network parameters
		OrdererID: "orderer.fudan.edu.cn",
		OrgID: "org1.fudan.edu.cn",

		// Channel parameters
		ChannelID:     "fudanfabric",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/thorweiyan/fabric_go_sdk/fixtures/artifacts/fudanfabric.channel.tx",

		// Chaincode parameters
		ChainCodeID:     "fudancc",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/thorweiyan/fabric_go_sdk/chaincode/test/",
		ChaincodeVersion: "0",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",

		// User parameters
		UserName: "User1",
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}
	// Close SDK
	defer fSetup.CloseSDK()
}

func TestFabricSetup_InstallAndInstantiateCC(t *testing.T) {
	// Definition of the Fabric SDK properties
	fSetup := FabricSetup{
		// Network parameters
		OrdererID: "orderer.fudan.edu.cn",
		OrgID: "org1.fudan.edu.cn",

		// Channel parameters
		ChannelID:     "fudanfabric",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/thorweiyan/fabric_go_sdk/fixtures/artifacts/fudanfabric.channel.tx",

		// Chaincode parameters
		ChainCodeID:     "fudancc9",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/thorweiyan/fabric_go_sdk/chaincode/testcall/",
		ChaincodeVersion: "1",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",

		// User parameters
		UserName: "User1",
	}
	// Install and instantiate the chaincode
	err := fSetup.InstallAndInstantiateCC([]string{"init"})
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}
}

func TestFabricSetup_Invoke(t *testing.T) {
	// Definition of the Fabric SDK properties
	fSetup := FabricSetup{
		// Network parameters
		OrdererID: "orderer.fudan.edu.cn",
		OrgID: "org1.fudan.edu.cn",

		// Channel parameters
		ChannelID:     "fudanfabric",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/thorweiyan/fabric_go_sdk/fixtures/artifacts/fudanfabric.channel.tx",

		// Chaincode parameters
		ChainCodeID:     "fudancc8",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/thorweiyan/fabric_go_sdk/chaincode/testcall/",
		ChaincodeVersion: "1",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",

		// User parameters
		UserName: "User1",
	}

	//trcid, err := fSetup.Invoke([]string{"call","fudancc6","call","fudancc7","write","hello","world"})
	trcid, err := fSetup.Invoke([]string{"call","fudancc9","call","fudancc8","write","hello","world"})
	if err != nil {
		fmt.Println("invoke error!", err)
	}
	fmt.Println("payload:"+trcid)
}

func TestFabricSetup_Query(t *testing.T) {
	// Definition of the Fabric SDK properties
	fSetup := FabricSetup{
		// Network parameters
		OrdererID: "orderer.fudan.edu.cn",
		OrgID: "org1.fudan.edu.cn",

		// Channel parameters
		ChannelID:     "fudanfabric",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/thorweiyan/fabric_go_sdk/fixtures/artifacts/fudanfabric.channel.tx",

		// Chaincode parameters
		ChainCodeID:     "fudancc1",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/thorweiyan/fabric_go_sdk/chaincode/test/",
		ChaincodeVersion: "0",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",

		// User parameters
		UserName: "User1",
	}

	payload, err := fSetup.Query([]string{"invoke", "query", "hello"})
	if err != nil {
		fmt.Println("query error!", err)
	}
	fmt.Println(payload)
}
