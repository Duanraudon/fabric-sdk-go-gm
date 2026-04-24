/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ext

import (
	"encoding/json"

	"github.com/Duanraudon/fabric-sdk-go-gm/internal/github.com/hyperledger/fabric/sdkinternal/configtxgen/encoder"
	localconfig "github.com/Duanraudon/fabric-sdk-go-gm/internal/github.com/hyperledger/fabric/sdkinternal/configtxgen/genesisconfig"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/fab/resource/genesisconfig"
	cb "github.com/hyperledger/fabric-protos-go/common"
)

// NewConsortiumOrgGroup ...
func NewConsortiumOrgGroup(org genesisconfig.Organization) (*cb.ConfigGroup, error) {
	b, err := json.Marshal(org)
	if err != nil {
		return nil, err
	}
	o := &localconfig.Organization{}
	err = json.Unmarshal(b, o)
	if err != nil {
		return nil, err
	}
	cfgG, err := encoder.NewConsortiumOrgGroup(o)
	if err != nil {
		return nil, err
	}
	return cfgG, nil
}
