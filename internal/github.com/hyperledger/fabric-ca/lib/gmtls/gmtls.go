/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
/*
Notice: This file has been modified for Hyperledger Fabric SDK Go usage.
Please review third_party pinning scripts and patches for more details.
*/

package gmtls

import (
	"time"

	gmtls "github.com/Duanraudon/cryptogm/tls"
	x509GM "github.com/Duanraudon/cryptogm/x509"
	"github.com/Duanraudon/fabric-sdk-go-gm/internal/github.com/hyperledger/fabric-ca/sdkinternal/pkg/util"
	factory "github.com/Duanraudon/fabric-sdk-go-gm/internal/github.com/hyperledger/fabric-ca/sdkpatch/cryptosuitebridge"
	log "github.com/Duanraudon/fabric-sdk-go-gm/internal/github.com/hyperledger/fabric-ca/sdkpatch/logbridge"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/common/providers/core"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/core/cryptosuite"
	"github.com/pkg/errors"
)

// DefaultCipherSuites is a set of strong TLS cipher suites
var DefaultCipherSuites = []uint16{
	gmtls.GMTLS_SM2_WITH_SM4_SM3,
	gmtls.GMTLS_ECDHE_SM2_WITH_SM4_SM3,
}

// ClientTLSConfig defines the key material for a TLS client
type ClientTLSConfig struct {
	Enabled     bool     `skip:"true"`
	CertFiles   [][]byte `help:"A list of comma-separated PEM-encoded trusted certificate bytes"`
	Client      KeyCertFiles
	TlsCertPool *x509GM.CertPool
}

// KeyCertFiles defines the files need for client on TLS
type KeyCertFiles struct {
	KeyFile  []byte `help:"PEM-encoded key bytes when mutual authentication is enabled"`
	CertFile []byte `help:"PEM-encoded certificate bytes when mutual authenticate is enabled"`
}

// func LoadZHX509KeyPair(certPEMBlock, keyBytes []byte) (gmtls.Certificate, error) {
// 	fail := func(err error) (gmtls.Certificate, error) { return gmtls.Certificate{}, err }

// 	var cert gmtls.Certificate
// 	var skippedBlockTypes []string
// 	for {
// 		var certDERBlock *pem.Block
// 		certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
// 		if certDERBlock == nil {
// 			break
// 		}
// 		if certDERBlock.Type == "CERTIFICATE" {
// 			cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
// 		} else {
// 			skippedBlockTypes = append(skippedBlockTypes, certDERBlock.Type)
// 		}
// 	}

// 	if len(cert.Certificate) == 0 {
// 		if len(skippedBlockTypes) == 0 {
// 			return fail(errors.New("tls: failed to find any PEM data in certificate input"))
// 		}
// 		if len(skippedBlockTypes) == 1 && strings.HasSuffix(skippedBlockTypes[0], "PRIVATE KEY") {
// 			return fail(errors.New("tls: failed to find certificate PEM data in certificate input, but did find a private key; PEM inputs may have been switched"))
// 		}
// 		return fail(fmt.Errorf("tls: failed to find \"CERTIFICATE\" PEM block in certificate input after skipping PEM blocks of the following types: %v", skippedBlockTypes))
// 	}

// 	priAdapter, err := zhsm2.CreateSm2KeyAdapter(strings.TrimSpace(string(keyBytes)))
// 	if err != nil {
// 		return fail(fmt.Errorf("tls: failed to load zh sm2 key adapter, Got err: %s", err))
// 	}

// 	x509Cert, err := x509GM.ParseCertificate(cert.Certificate[0])
// 	if err != nil {
// 		return fail(err)
// 	}

// 	cert.PrivateKey = priAdapter

// 	//check pub and private match
// 	switch pub := x509Cert.PublicKey.(type) {
// 	case *rsa.PublicKey:
// 		priv, ok := cert.PrivateKey.(*rsa.PrivateKey)
// 		if !ok {
// 			return fail(errors.New("tls: private key type does not match public key type"))
// 		}
// 		if pub.N.Cmp(priv.N) != 0 {
// 			return fail(errors.New("tls: private key does not match public key"))
// 		}
// 	case *ecdsa.PublicKey:
// 		pub, _ = x509Cert.PublicKey.(*ecdsa.PublicKey)
// 		switch pub.Curve {
// 		case sm2.P256():
// 			zhPubKey := priAdapter.PublicKey()
// 			if pub.X.Cmp(zhPubKey.X) != 0 || pub.Y.Cmp(zhPubKey.Y) != 0 {
// 				return fail(errors.New("tls: zh sm2 private key does not match public key"))
// 			}
// 		default:
// 			return fail(errors.New("tls: private key type does not match public key type"))
// 		}
// 	default:
// 		return fail(errors.New("tls: unknown public key algorithm"))
// 	}

// 	return cert, nil
// }

// GetClientTLSConfig creates a gmtls.Config object from certs and roots
func GetClientTLSConfig(cfg *ClientTLSConfig, csp core.CryptoSuite) (*gmtls.Config, error) {
	var certs []gmtls.Certificate

	if csp == nil {
		csp = factory.GetDefault()
	}

	if cfg.Client.CertFile != nil {
		err := checkCertDates(cfg.Client.CertFile)
		if err != nil {
			return nil, err
		}

		//try zh-ce first
		var clientCert gmtls.Certificate
		// clientCert, err = LoadZHX509KeyPair(cfg.Client.CertFile, cfg.Client.KeyFile)
		// if err != nil {
		// 	clientCert, err = gmtls.X509KeyPair(cfg.Client.CertFile, cfg.Client.KeyFile)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// }
		// Can handle all of ECDAS, RSA and SM2.
		clientCert, err = gmtls.X509KeyPair(cfg.Client.CertFile, cfg.Client.KeyFile)
		if err != nil {
			return nil, err
		}

		certs = append(certs, clientCert)
	} else {
		log.Debug("Client TLS certificate and/or key file not provided")
	}
	rootCAPool := cfg.TlsCertPool

	if rootCAPool == nil {
		rootCAPool, err := x509GM.SystemCertPool()
		if err != nil {
			log.Debugf("Failed to load system cert pool, switching to empty cert pool ")
			rootCAPool = x509GM.NewCertPool()
		}

		if len(cfg.CertFiles) == 0 {
			return nil, errors.New("No trusted root certificates for TLS were provided")
		}

		for _, cacert := range cfg.CertFiles {
			ok := rootCAPool.AppendCertsFromPEM(cacert)
			if !ok {
				return nil, errors.New("Failed to process certificate")
			}
		}
	}
	var config *gmtls.Config
	if cryptosuite.SecAlgo == "SHA2" {
		config = &gmtls.Config{
			Certificates: certs,
			RootCAs:      rootCAPool,
		}
	} else {
		config = &gmtls.Config{
			GMUsed:       true,
			Certificates: certs,
			RootCAs:      rootCAPool,
		}
	}

	return config, nil
}

func checkCertDates(certPEM []byte) error {
	log.Debug("Check client TLS certificate for valid dates")

	cert, err := util.GetX509CertificateFromPEM(certPEM)
	if err != nil {
		return err
	}

	notAfter := cert.NotAfter
	currentTime := time.Now().UTC()

	if currentTime.After(notAfter) {
		return errors.New("Certificate provided has expired")
	}

	notBefore := cert.NotBefore
	if currentTime.Before(notBefore) {
		return errors.New("Certificate provided not valid until later date")
	}

	return nil
}
