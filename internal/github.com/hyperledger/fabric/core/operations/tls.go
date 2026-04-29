/*
Copyright IBM Corp All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
/*
Notice: This file has been modified for Hyperledger Fabric SDK Go usage.
Please review third_party pinning scripts and patches for more details.
*/

package operations

import (
	"io/ioutil"

	"github.com/Duanraudon/cryptogm/tls"
	"github.com/Duanraudon/cryptogm/x509"
	"github.com/Duanraudon/fabric-sdk-go-gm/internal/github.com/hyperledger/fabric/sdkinternal/pkg/comm"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/core/cryptosuite"
)

type TLS struct {
	Enabled            bool
	CertFile           string
	KeyFile            string
	ClientCertRequired bool
	ClientCACertFiles  []string
}

func (t TLS) Config() (*tls.Config, error) {
	var tlsConfig *tls.Config
	var cfgCipherSuites []uint16

	if t.Enabled {
		cert, err := tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		for _, caPath := range t.ClientCACertFiles {
			caPem, err := ioutil.ReadFile(caPath)
			if err != nil {
				return nil, err
			}
			caCertPool.AppendCertsFromPEM(caPem)
		}

		var gmUsed bool
		if cryptosuite.SecAlgo == "SHA2" {
			cfgCipherSuites = comm.DefaultTLSCipherSuites
			gmUsed = false
		} else {
			cfgCipherSuites = comm.DefaultGMTLSCipherSuites
			gmUsed = true
		}

		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			CipherSuites: cfgCipherSuites,
			ClientCAs:    caCertPool,
			GMUsed:       gmUsed,
		}
		if t.ClientCertRequired {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		} else {
			tlsConfig.ClientAuth = tls.VerifyClientCertIfGiven
		}
	}

	return tlsConfig, nil
}
