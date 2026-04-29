package gmtls

import (
	"time"

	"github.com/pkg/errors"
	gmx509 "github.com/tjfoc/gmsm/x509"
	tls "github.com/tjfoc/gmtls"
	"github.com/tjfoc/gmtls/gmcredentials"

	cryptox509 "github.com/Duanraudon/cryptogm/x509"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/common/providers/fab"
	"google.golang.org/grpc/credentials"
)

// Manager creates gmtls transport credentials for peer and orderer connections.
type Manager struct {
	config fab.EndpointConfig
}

// New creates a gmtls manager.
func New(config fab.EndpointConfig) *Manager {
	return &Manager{config: config}
}

// NewTransportCredentials builds gmtls transport credentials using the supplied server CA and name.
func (m *Manager) NewTransportCredentials(cert *cryptox509.Certificate, serverName string) (credentials.TransportCredentials, error) {
	certPool := gmx509.NewCertPool()
	if cert != nil && len(cert.Raw) > 0 {
		gmCert, err := gmx509.ParseCertificate(cert.Raw)
		if err != nil {
			return nil, err
		}
		certPool.AddCert(gmCert)
	}

	tlsConf := &tls.Config{
		RootCAs:      certPool,
		ServerName:   serverName,
		Certificates: []tls.Certificate{},
	}
	tlsConf.VerifyPeerCertificate = verifyPeerCertificate

	return gmcredentials.NewTLS(tlsConf), nil
}

func verifyPeerCertificate(rawCerts [][]byte, verifiedChains [][]*gmx509.Certificate) error {
	for _, certDER := range rawCerts {
		cert, err := gmx509.ParseCertificate(certDER)
		if err != nil {
			return err
		}
		if err := validateCertificateDates(cert); err != nil {
			return err
		}
	}

	for _, chain := range verifiedChains {
		for _, cert := range chain {
			if err := validateCertificateDates(cert); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateCertificateDates(cert *gmx509.Certificate) error {
	if cert == nil {
		return nil
	}
	now := time.Now().UTC()
	if now.Before(cert.NotBefore) {
		return errors.New("certificate provided is not valid until later date")
	}
	if now.After(cert.NotAfter) {
		return errors.New("certificate provided has expired")
	}
	return nil
}
