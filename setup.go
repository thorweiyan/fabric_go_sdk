package fabric_go_sdk

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)

// FabricSetup implementation
type FabricSetup struct {
	initialize			bool
	ConfigFile			string
	OrgID           	string
	OrdererID       	string
	ChannelID       	string
	ChainCodeID     	string
	ChannelConfig   	string
	ChaincodeGoPath 	string
	ChaincodePath   	string
	ChaincodeVersion 	string
	OrgAdmin        	string
	OrgName         	string
	UserName        	string
	client          	*channel.Client
	admin           	*resmgmt.Client
	sdk             	*fabsdk.FabricSDK
	event           	*event.Client
}

// 根据config.yaml初始化并启动client, chain 和event hub
func (setup *FabricSetup) Initialize() error {
	if setup.initialize {
		return errors.New("sdk already initialized")
	}
	err := setup.initsdk()
	if err != nil {
		return err
	}
	err = setup.initadmin()
	if err != nil {
		return err
	}

	// The MSP client allow us to retrieve user information from their identity, like its signing identity which we will need to save the channel
	mspClient, err := mspclient.New(setup.sdk.Context(), mspclient.WithOrg(setup.OrgName))
	if err != nil {
		return errors.WithMessage(err, "failed to create MSP client")
	}
	adminIdentity, err := mspClient.GetSigningIdentity(setup.OrgAdmin)
	if err != nil {
		return errors.WithMessage(err, "failed to get admin signing identity")
	}
	req := resmgmt.SaveChannelRequest{ChannelID: setup.ChannelID, ChannelConfigPath: setup.ChannelConfig, SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	txID, err := setup.admin.SaveChannel(req, resmgmt.WithOrdererEndpoint(setup.OrdererID))
	if err != nil || txID.TransactionID == "" {
		return errors.WithMessage(err, "failed to save channel")
	}
	fmt.Println("Channel created")

	// Make admin user join the previously created channel
	if err = setup.admin.JoinChannel(setup.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(setup.OrdererID)); err != nil {
		return errors.WithMessage(err, "failed to make admin join channel")
	}
	fmt.Println("Channel joined")

	setup.initialize = true
	fmt.Println("Initialization Successful")
	return nil
}

// 安装并初始实例化cc
func (setup *FabricSetup) InstallAndInstantiateCC(args []string) error {
	if !setup.initialize {
		err := setup.initsdk()
		if err != nil {
			return err
		}
		err = setup.initadmin()
		if err != nil {
			return err
		}
	}
	var inargs [][]byte

	for _, i := range args {
		inargs = append(inargs, []byte(i))
	}

	// Create the chaincode package that will be sent to the peers
	ccPkg, err := packager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return errors.WithMessage(err, "failed to create chaincode package")
	}
	fmt.Println("ccPkg created")

	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodePath, Version: setup.ChaincodeVersion, Package: ccPkg}
	_, err = setup.admin.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return errors.WithMessage(err, "failed to install chaincode")
	}
	fmt.Println("Chaincode installed")

	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{setup.OrgID})

	resp, err := setup.admin.InstantiateCC(setup.ChannelID, resmgmt.InstantiateCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodeGoPath, Version: setup.ChaincodeVersion, Args: inargs, Policy: ccPolicy})
	if err != nil || resp.TransactionID == "" {
		return errors.WithMessage(err, "failed to instantiate the chaincode")
	}
	fmt.Println("Chaincode instantiated")

	err = setup.initclient()
	if err != nil {
		return err
	}
	err = setup.initevent()
	if err != nil {
		return err
	}

	fmt.Println("Chaincode Installation & Instantiation Successful")
	setup.initialize = true
	return nil
}
// 实例化cc
func (setup *FabricSetup) InstantiateCC(args []string) error {
	if !setup.initialize {
		err := setup.initsdk()
		if err != nil {
			return err
		}
		err = setup.initadmin()
		if err != nil {
			return err
		}
	}
	var inargs [][]byte

	for _, i := range args {
		inargs = append(inargs, []byte(i))
	}
	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{setup.OrgID})

	resp, err := setup.admin.InstantiateCC(setup.ChannelID, resmgmt.InstantiateCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodeGoPath, Version: setup.ChaincodeVersion, Args: inargs, Policy: ccPolicy})
	if err != nil || resp.TransactionID == "" {
		return errors.WithMessage(err, "failed to instantiate the chaincode")
	}
	fmt.Println("Chaincode instantiated")

	err = setup.initclient()
	if err != nil {
		return err
	}
	err = setup.initevent()
	if err != nil {
		return err
	}

	fmt.Println("ChaincodeInstantiation Successful")
	setup.initialize = true
	return nil
}

func (setup *FabricSetup) CloseSDK() {
	setup.sdk.Close()
}

func (setup *FabricSetup) UpdateCC() {
	//TODO: ADD METHOD TO UPDATE CHAINCODE
	return
}

func (setup *FabricSetup) initsdk() error{
	// Initialize the SDK with the configuration file
	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		return errors.WithMessage(err, "failed to create SDK")
	}
	setup.sdk = sdk
	fmt.Println("SDK created")
	return nil
}

//need sdk
func (setup *FabricSetup) initadmin() error{
	// The resource management client is responsible for managing channels (create/update channel)
	resourceManagerClientContext := setup.sdk.Context(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName))
	resMgmtClient, err := resmgmt.New(resourceManagerClientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create channel management client from Admin identity")
	}
	setup.admin = resMgmtClient
	fmt.Println("Ressource management client created")
	return nil
}

//need sdk
func (setup *FabricSetup) initclient() error {
	var err error
	// Channel client is used to query and execute transactions
	clientContext := setup.sdk.ChannelContext(setup.ChannelID, fabsdk.WithUser(setup.UserName))
	setup.client, err = channel.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new channel client")
	}
	fmt.Println("Channel client created")
	return nil
}

//need sdk
func (setup *FabricSetup) initevent() error {
	var err error
	// Channel client is used to query and execute transactions
	clientContext := setup.sdk.ChannelContext(setup.ChannelID, fabsdk.WithUser(setup.UserName))
	// Creation of the client which will enables access to our channel events
	setup.event, err = event.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new event client")
	}
	fmt.Println("Event client created")

	return nil
}