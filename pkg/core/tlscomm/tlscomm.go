package tlscomm

import (
	cryptox509 "github.com/Duanraudon/cryptogm/x509"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/client/common/verifier"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/common/providers/fab"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/core/config/comm"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/core/cryptosuite"
	gmtlscomm "github.com/Duanraudon/fabric-sdk-go-gm/pkg/core/tlscomm/gmtls"
	"google.golang.org/grpc/credentials"
)

// NewTransportCredentials creates transport credentials that match the active crypto suite.
func NewTransportCredentials(cert *cryptox509.Certificate, serverName string, config fab.EndpointConfig) (credentials.TransportCredentials, error) {
	if cryptosuite.SecAlgo == "SM3" {
		return gmtlscomm.New(config).NewTransportCredentials(cert, serverName)
	}

	tlsConfig, err := comm.TLSConfig(cert, serverName, config)
	if err != nil {
		return nil, err
	}
	tlsConfig.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*cryptox509.Certificate) error {
		return verifier.VerifyPeerCertificate(rawCerts, verifiedChains)
	}

	return credentials.NewTLS(tlsConfig), nil
}
