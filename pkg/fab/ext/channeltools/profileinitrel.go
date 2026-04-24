package channeltools

import (
	"time"

	"github.com/Duanraudon/fabric-sdk-go-gm/internal/github.com/hyperledger/fabric/sdkinternal/configtxgen/genesisconfig"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/fab/ext/utils"
	"github.com/hyperledger/fabric-protos-go/orderer/etcdraft"
)

const (
	AdminRoleAdminPrincipal = "Role.ADMIN"
	EtcdRaft                = "etcdraft"
)

var genesisDefaults = genesisconfig.TopLevel{
	Orderer: &genesisconfig.Orderer{
		OrdererType:  "solo",
		BatchTimeout: 2 * time.Second,
		BatchSize: genesisconfig.BatchSize{
			MaxMessageCount:   500,
			AbsoluteMaxBytes:  10 * 1024 * 1024,
			PreferredMaxBytes: 2 * 1024 * 1024,
		},
		Kafka: genesisconfig.Kafka{
			Brokers: []string{"127.0.0.1:9092"},
		},
		EtcdRaft: &etcdraft.ConfigMetadata{
			Options: &etcdraft.Options{
				TickInterval:         "500ms",
				ElectionTick:         10,
				HeartbeatTick:        1,
				MaxInflightBlocks:    5,
				SnapshotIntervalSize: 16 * 1024 * 1024, // 16 MB
			},
		},
	},
}

func CompleteProfileInitialization(p *genesisconfig.Profile, configDir string) {
	if p.Application != nil {
		for _, org := range p.Application.Organizations {
			completeOrgInitialization(org, configDir)
		}
	}

	if p.Consortiums != nil {
		for _, consortium := range p.Consortiums {
			for _, org := range consortium.Organizations {
				completeOrgInitialization(org, configDir)
			}
		}
	}

	if p.Orderer != nil {
		for _, org := range p.Orderer.Organizations {
			completeOrgInitialization(org, configDir)
		}
		// Some profiles will not define orderer parameters
		completeOrderInitialization(p.Orderer, configDir)
	}

}

func completeOrgInitialization(org *genesisconfig.Organization, configDir string) {
	// set the MSP type; if none is specified we assume BCCSP
	if org.MSPType == "" {
		org.MSPType = "bccsp"
	}

	if org.AdminPrincipal == "" {
		org.AdminPrincipal = AdminRoleAdminPrincipal
	}
	var p *string = &org.MSPDir
	utils.TranslatePaths(configDir, p)
}

func completeOrderInitialization(ord *genesisconfig.Orderer, configDir string) {
loop:
	for {
		switch {
		case ord.OrdererType == "":
			logger.Infof("Orderer.OrdererType unset, setting to %v", genesisDefaults.Orderer.OrdererType)
			ord.OrdererType = genesisDefaults.Orderer.OrdererType
		case ord.BatchTimeout == 0:
			logger.Infof("Orderer.BatchTimeout unset, setting to %s", genesisDefaults.Orderer.BatchTimeout)
			ord.BatchTimeout = genesisDefaults.Orderer.BatchTimeout
		case ord.BatchSize.MaxMessageCount == 0:
			logger.Infof("Orderer.BatchSize.MaxMessageCount unset, setting to %v", genesisDefaults.Orderer.BatchSize.MaxMessageCount)
			ord.BatchSize.MaxMessageCount = genesisDefaults.Orderer.BatchSize.MaxMessageCount
		case ord.BatchSize.AbsoluteMaxBytes == 0:
			logger.Infof("Orderer.BatchSize.AbsoluteMaxBytes unset, setting to %v", genesisDefaults.Orderer.BatchSize.AbsoluteMaxBytes)
			ord.BatchSize.AbsoluteMaxBytes = genesisDefaults.Orderer.BatchSize.AbsoluteMaxBytes
		case ord.BatchSize.PreferredMaxBytes == 0:
			logger.Infof("Orderer.BatchSize.PreferredMaxBytes unset, setting to %v", genesisDefaults.Orderer.BatchSize.PreferredMaxBytes)
			ord.BatchSize.PreferredMaxBytes = genesisDefaults.Orderer.BatchSize.PreferredMaxBytes
		default:
			break loop
		}
	}

	logger.Infof("orderer type: %s", ord.OrdererType)
	// Additional, consensus type-dependent initialization goes here
	// Also using this to panic on unknown orderer type.
	switch ord.OrdererType {
	case "solo":
		// nothing to be done here
	case "kafka":
		if ord.Kafka.Brokers == nil {
			logger.Infof("Orderer.Kafka unset, setting to %v", genesisDefaults.Orderer.Kafka.Brokers)
			ord.Kafka.Brokers = genesisDefaults.Orderer.Kafka.Brokers
		}
	case EtcdRaft:
		if ord.EtcdRaft == nil {
			logger.Errorf("%s configuration missing", EtcdRaft)
		}
		if ord.EtcdRaft.Options == nil {
			logger.Infof("Orderer.EtcdRaft.Options unset, setting to %v", genesisDefaults.Orderer.EtcdRaft.Options)
			ord.EtcdRaft.Options = genesisDefaults.Orderer.EtcdRaft.Options
		}
	second_loop:
		for {
			switch {
			case ord.EtcdRaft.Options.TickInterval == "":
				logger.Infof("Orderer.EtcdRaft.Options.TickInterval unset, setting to %v", genesisDefaults.Orderer.EtcdRaft.Options.TickInterval)
				ord.EtcdRaft.Options.TickInterval = genesisDefaults.Orderer.EtcdRaft.Options.TickInterval

			case ord.EtcdRaft.Options.ElectionTick == 0:
				logger.Infof("Orderer.EtcdRaft.Options.ElectionTick unset, setting to %v", genesisDefaults.Orderer.EtcdRaft.Options.ElectionTick)
				ord.EtcdRaft.Options.ElectionTick = genesisDefaults.Orderer.EtcdRaft.Options.ElectionTick

			case ord.EtcdRaft.Options.HeartbeatTick == 0:
				logger.Infof("Orderer.EtcdRaft.Options.HeartbeatTick unset, setting to %v", genesisDefaults.Orderer.EtcdRaft.Options.HeartbeatTick)
				ord.EtcdRaft.Options.HeartbeatTick = genesisDefaults.Orderer.EtcdRaft.Options.HeartbeatTick

			case ord.EtcdRaft.Options.MaxInflightBlocks == 0:
				logger.Infof("Orderer.EtcdRaft.Options.MaxInflightBlocks unset, setting to %v", genesisDefaults.Orderer.EtcdRaft.Options.MaxInflightBlocks)
				ord.EtcdRaft.Options.MaxInflightBlocks = genesisDefaults.Orderer.EtcdRaft.Options.MaxInflightBlocks

			case ord.EtcdRaft.Options.SnapshotIntervalSize == 0:
				logger.Infof("Orderer.EtcdRaft.Options.SnapshotIntervalSize unset, setting to %v", genesisDefaults.Orderer.EtcdRaft.Options.SnapshotIntervalSize)
				ord.EtcdRaft.Options.SnapshotIntervalSize = genesisDefaults.Orderer.EtcdRaft.Options.SnapshotIntervalSize

			case len(ord.EtcdRaft.Consenters) == 0:
				logger.Errorf("%s configuration did not specify any consenter", EtcdRaft)

			default:
				break second_loop
			}
		}

		if _, err := time.ParseDuration(ord.EtcdRaft.Options.TickInterval); err != nil {
			logger.Errorf("Etcdraft TickInterval (%s) must be in time duration format", ord.EtcdRaft.Options.TickInterval)
		}

		// validate the specified members for Options
		if ord.EtcdRaft.Options.ElectionTick <= ord.EtcdRaft.Options.HeartbeatTick {
			logger.Errorf("election tick must be greater than heartbeat tick")
		}

		for _, c := range ord.EtcdRaft.GetConsenters() {
			if c.Host == "" {
				logger.Errorf("consenter info in %s configuration did not specify host", EtcdRaft)
			}
			if c.Port == 0 {
				logger.Errorf("consenter info in %s configuration did not specify port", EtcdRaft)
			}
			if c.ClientTlsCert == nil {
				logger.Errorf("consenter info in %s configuration did not specify client TLS cert", EtcdRaft)
			}
			if c.ServerTlsCert == nil {
				logger.Errorf("consenter info in %s configuration did not specify server TLS cert", EtcdRaft)
			}
			clientCertPath := string(c.GetClientTlsCert())
			utils.TranslatePaths(configDir, &clientCertPath)
			c.ClientTlsCert = []byte(clientCertPath)
			serverCertPath := string(c.GetServerTlsCert())
			utils.TranslatePaths(configDir, &serverCertPath)
			c.ServerTlsCert = []byte(serverCertPath)
		}
	default:
		logger.Errorf("unknown orderer type: %s", ord.OrdererType)
	}
}
