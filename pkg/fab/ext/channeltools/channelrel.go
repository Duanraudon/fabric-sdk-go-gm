package channeltools

import (
	"fmt"

	"github.com/Duanraudon/fabric-sdk-go-gm/internal/github.com/hyperledger/fabric/common/channelconfig"
	"github.com/Duanraudon/fabric-sdk-go-gm/internal/github.com/hyperledger/fabric/protoutil"
	"github.com/Duanraudon/fabric-sdk-go-gm/internal/github.com/hyperledger/fabric/sdkinternal/configtxgen/encoder"
	"github.com/Duanraudon/fabric-sdk-go-gm/internal/github.com/hyperledger/fabric/sdkinternal/configtxgen/genesisconfig"
	"github.com/Duanraudon/fabric-sdk-go-gm/internal/github.com/hyperledger/fabric/sdkinternal/configtxlator/update"
	"github.com/Duanraudon/fabric-sdk-go-gm/internal/github.com/hyperledger/fabric/sdkinternal/pkg/identity"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/common/logging"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/fab/ext/utils"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/pkg/errors"

	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
// org1 = "Org1"
// org2             = "Org2"
// ordererAdminUser = "Admin"
// ordererOrgName   = "OrdererOrg"
// org1AdminUser    = "Admin"
// org2AdminUser    = "Admin"
// org1User         = "User1"
// org2User         = "User1"
// channelID        = "orgchannel"
// ccPath = "github.com/example_cc"
)

var logger = logging.NewLogger("fabsdk/fab")

func GetProfileFromConf(channelID string, profile string, configPath string, configName string) (*genesisconfig.Profile, error) {
	if channelID == "" {
		logger.Fatal("channelID cannot be null")
	}
	config := viper.New()
	if configPath != "" {
		config.AddConfigPath(configPath)
		config.SetConfigType("yaml")
		config.SetConfigName(configName)
	} else {
		utils.InitViper(config, "configtx")
	}

	//if len(configPaths) > 0 {
	//	for _, p := range configPaths {
	//		config.SetConfigName("configtx")
	//		config.AddConfigPath(p)
	//		config.SetConfigFile(p)
	//		config.SetConfigType("yaml")
	//	}
	//	//config.AddConfigPath(p)
	//	//config.SetConfigType("yaml")
	//	//config.SetConfigName("configtx")
	//
	//} else {
	//	InitViper(config, "configtx")
	//}

	// For environment variables
	config.SetEnvPrefix(utils.Prefix)
	config.AutomaticEnv()

	replacer := strings.NewReplacer(strings.ToUpper(fmt.Sprintf("profiles.%s.", profile)), "", ".", "_")
	config.SetEnvKeyReplacer(replacer)

	err := config.ReadInConfig()
	if err != nil {
		logger.Errorf("Error reading configuration: %s", err)
	}
	fmt.Printf("mjprofile:%s", config.Get("Profiles.SampleSingleMSPSolo.Orderer.Policies"))
	//fmt.Printf("mjprofile:%s",config.Get("Profiles.TwoOrgsOrdererGenesis.Orderer.Policies"))
	logger.Debugf("Using config file: %s", config.ConfigFileUsed())

	var conf genesisconfig.TopLevel

	err = utils.EnhancedExactUnmarshal(config, &conf)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling config into struct: %s", err)
	}

	result, ok := conf.Profiles[strings.ToLower(profile)] // 这里的profile都转成小写了
	if !ok {
		logger.Errorf("Could not find profile: %s", profile)
	}

	CompleteProfileInitialization(result, filepath.Dir(config.ConfigFileUsed()))

	fixProfileUpper(result)

	if result.Policies == nil || len(result.Policies) == 0 {
		result.Policies = sampleAppPolicies()
	}
	if &result.Application.Policies == nil || len(result.Application.Policies) == 0 {
		result.Application.Policies = sampleAppPolicies()
	}

	if result.Application.Organizations != nil && len(result.Application.Organizations) > 0 {
		for _, org := range result.Application.Organizations {
			if org.Policies == nil || len(org.Policies) == 0 {
				org.Policies = sampleOrgPolicies()
			}
		}
	}

	logger.Infof("Loaded configuration: %s", config.ConfigFileUsed())

	return result, nil
}

func fixPolicyUpper(policies map[string]*genesisconfig.Policy) map[string]*genesisconfig.Policy {
	if policies == nil || len(policies) == 0 {
		return nil
	}
	result := make(map[string]*genesisconfig.Policy)
	for k, v := range policies {
		result[utils.StrFirstToUpper(k)] = v
	}
	return result
}

func fixProfileUpper(profile *genesisconfig.Profile) {
	if profile.Application.Capabilities != nil && len(profile.Application.Capabilities) > 0 {
		capabilities := make(map[string]bool)
		for k, v := range profile.Application.Capabilities {
			capabilities[strings.ToUpper(k)] = v
		}
		profile.Application.Capabilities = capabilities
	}
	if profile.Capabilities != nil && len(profile.Capabilities) > 0 {
		capabilities := make(map[string]bool)
		for k, v := range profile.Capabilities {
			capabilities[strings.ToUpper(k)] = v
		}
		profile.Capabilities = capabilities
	}

	if profile.Policies != nil && len(profile.Policies) > 0 {
		profile.Policies = fixPolicyUpper(profile.Policies)
	}

	if profile.Application.Policies != nil && len(profile.Application.Policies) > 0 {
		profile.Application.Policies = fixPolicyUpper(profile.Application.Policies)
	}
	if profile.Application.Organizations != nil && len(profile.Application.Organizations) > 0 {
		for _, org := range profile.Application.Organizations {
			if org.Policies != nil {
				org.Policies = fixPolicyUpper(org.Policies)
			}
		}
	}
}

func GetChannelTxByProfile(channelID string, signer identity.SignerSerializer, conf *genesisconfig.Profile) (*cb.Envelope, error) {

	var configtx *cb.Envelope
	// resource.CreateChannelCreateTx(conf,nil,channelID)
	configtx, err := encoder.MakeChannelCreationTransaction(channelID, nil, conf)
	return configtx, err
}

func GetAnchorPeersUpdate(profile *genesisconfig.Profile, channelID string, asOrg string) (*common.Envelope, error) {
	logger.Info("Generating anchor peer update")
	if asOrg == "" {
		return nil, fmt.Errorf("must specify an organization to update the anchor peer for")
	}

	if profile.Application == nil {
		return nil, fmt.Errorf("cannot update anchor peers without an application section")
	}

	original, err := encoder.NewChannelGroup(profile)
	if err != nil {
		return nil, errors.WithMessage(err, "error parsing profile as channel group")
	}
	original.Groups[channelconfig.ApplicationGroupKey].Version = 1

	updated := proto.Clone(original).(*cb.ConfigGroup)

	originalOrg, ok := original.Groups[channelconfig.ApplicationGroupKey].Groups[asOrg]
	if !ok {
		return nil, errors.Errorf("org with name '%s' does not exist in config", asOrg)
	}

	if _, ok = originalOrg.Values[channelconfig.AnchorPeersKey]; !ok {
		return nil, errors.Errorf("org '%s' does not have any anchor peers defined", asOrg)
	}

	delete(originalOrg.Values, channelconfig.AnchorPeersKey)

	updt, err := update.Compute(&cb.Config{ChannelGroup: original}, &cb.Config{ChannelGroup: updated})
	if err != nil {
		return nil, errors.WithMessage(err, "could not compute update")
	}
	updt.ChannelId = channelID

	newConfigUpdateEnv := &cb.ConfigUpdateEnvelope{
		ConfigUpdate: protoutil.MarshalOrPanic(updt),
	}

	updateTx, err := protoutil.CreateSignedEnvelope(cb.HeaderType_CONFIG_UPDATE, channelID, nil, newConfigUpdateEnv, 0, 0)
	if err != nil {
		return nil, errors.WithMessage(err, "could not create signed envelope")
	}
	return updateTx, err
}

// func createChannel(sdk *fabsdk.FabricSDK, channelID string, channelConfigPath string) {
// 	org1MspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(org1))
// 	if err != nil {
// 		logger.Fatal(err)
// 	}

// 	org1AdminUser, err := org1MspClient.GetSigningIdentity(org1AdminUser)
// 	if err != nil {
// 		logger.Fatalf("failed to get org1AdminUser, err : %s", err)
// 	}
// 	//org2MspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(org2))
// 	//if err != nil {
// 	//	logger.Fatal(err)
// 	//}
// 	//org2AdminUser, err := org2MspClient.GetSigningIdentity(org2AdminUser)
// 	//if err != nil {
// 	//	logger.Fatalf("failed to get org2AdminUser, err : %s", err)
// 	//}
// 	ordererClientContext := sdk.Context(fabsdk.WithUser(ordererAdminUser), fabsdk.WithOrg(ordererOrgName))

// 	req := resmgmt.SaveChannelRequest{ChannelID: channelID,
// 		ChannelConfigPath: channelConfigPath,
// 		SigningIdentities: []msp.SigningIdentity{org1AdminUser}}
// 	chMgmtClient, err := resmgmt.New(ordererClientContext)
// 	txID, err := chMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.example.com"))
// 	logger.Infof("create channel success: txID:%s", txID)
// }

func sampleOrgPolicies() map[string]*genesisconfig.Policy {
	return map[string]*genesisconfig.Policy{
		"Readers": {
			Type: "Signature",
			Rule: "OR('Org1MSP.admin', 'Org1MSP.peer', 'Org1MSP.client')",
		},
		"Writers": {
			Type: "Signature",
			Rule: "OR('Org1MSP.admin', 'Org1MSP.client')",
		},
		"Admins": {
			Type: "Signature",
			Rule: "OR('Org1MSP.admin')",
		},
		//"Endorsement": {
		//	Type: "Signature",
		//	Rule: "OR('SampleOrg.member')",
		//},
	}
}

func sampleAppPolicies() map[string]*genesisconfig.Policy {
	return map[string]*genesisconfig.Policy{
		//"LifecycleEndorsement": {
		//	Type: "ImplicitMeta",
		//	Rule: "MAJORITY Endorsement",
		//},
		//"Endorsement": {
		//	Type: "ImplicitMeta",
		//	Rule: "MAJORITY Endorsement",
		//},
		"Readers": {
			Type: "ImplicitMeta",
			Rule: "ANY Readers",
		},
		"Writers": {
			Type: "ImplicitMeta",
			Rule: "ANY Writers",
		},
		"Admins": {
			Type: "ImplicitMeta",
			Rule: "MAJORITY Admins",
		},
	}
}
