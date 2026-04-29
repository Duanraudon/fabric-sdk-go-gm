/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"github.com/Duanraudon/cryptogm/tls"

	"github.com/Duanraudon/cryptogm/x509"

	"github.com/Duanraudon/fabric-sdk-go-gm/internal/github.com/hyperledger/fabric/sdkinternal/pkg/comm"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/common/providers/fab"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/core/cryptosuite"
	"github.com/pkg/errors"
)

// TLSConfig returns the appropriate config for TLS including the root CAs,
// certs for mutual TLS, and server host override. Works with certs loaded either from a path or embedded pem.
func TLSConfig(cert *x509.Certificate, serverName string, config fab.EndpointConfig) (*tls.Config, error) {

	if cert != nil {
		config.TLSCACertPool().Add(cert)
	}

	certPool, err := config.TLSCACertPool().Get()
	if err != nil {
		return nil, err
	}

	// 根据密码算法选择 CipherSuites 和 GMUsed 标志
	var cfgCipherSuites []uint16
	var gmUsed bool
	if cryptosuite.SecAlgo == "SHA2" {
		cfgCipherSuites = comm.DefaultTLSCipherSuites
		gmUsed = false
	} else {
		cfgCipherSuites = comm.DefaultGMTLSCipherSuites
		gmUsed = true
	}

	return &tls.Config{
		RootCAs:      certPool,
		Certificates: config.TLSClientCerts(),
		ServerName:   serverName,
		CipherSuites: cfgCipherSuites,
		GMUsed:       gmUsed,
	}, nil
}

// TLSCertHash is a utility method to calculate the SHA256 hash of the configured certificate (for usage in channel headers)
func TLSCertHash(config fab.EndpointConfig) ([]byte, error) {
	certs := config.TLSClientCerts()

	hashAlgo := cryptosuite.SecAlgo

	if len(certs) == 0 {
		return computeHash([]byte(""), hashAlgo)
	}

	cert := certs[0]
	if len(cert.Certificate) == 0 {
		return computeHash([]byte(""), hashAlgo)
	}

	return computeHash(cert.Certificate[0], hashAlgo)
}

// computeHash computes hash for given bytes using underlying cryptosuite default
func computeHash(msg []byte, hashAlgo string) ([]byte, error) {
	var h []byte
	var err error

	if hashAlgo == "SHA2" {
		h, err = cryptosuite.GetDefault().Hash(msg, cryptosuite.GetSHA256Opts())
	} else {
		h, err = cryptosuite.GetDefault().Hash(msg, cryptosuite.GetGMSM3Opts())
	}

	if err != nil {
		return nil, errors.WithMessage(err, "failed to compute tls cert hash")
	}
	return h, err
}
