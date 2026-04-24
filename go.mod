// Copyright SecureKey Technologies Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//
// GM-enabled Fabric SDK Go
// Supports SM2/SM3/SM4 Chinese National Cryptographic Standard

module github.com/Duanraudon/fabric-sdk-go-gm

//https://github.com/Duanraudon/fabric-sdk-go-gm.git

go 1.20

require (
	github.com/Duanraudon/cryptogm v1.0.3
	github.com/Knetic/govaluate v3.0.0+incompatible
	github.com/cloudflare/cfssl v1.4.1
	github.com/go-kit/kit v0.8.0
	github.com/golang/mock v1.4.3
	github.com/golang/protobuf v1.3.3
	github.com/hyperledger/fabric-config v0.0.5
	github.com/hyperledger/fabric-lib-go v1.0.0
	github.com/hyperledger/fabric-protos-go v0.0.0-20200707132912-fee30f3ccd23
	github.com/miekg/pkcs11 v1.0.3
	github.com/mitchellh/mapstructure v1.3.2
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.1.0
	github.com/spf13/cast v1.3.1
	github.com/spf13/viper v1.1.1
	github.com/stretchr/testify v1.5.1
	golang.org/x/crypto v0.0.0-20220919173607-35f4265a4bc0
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2
	google.golang.org/grpc v1.29.1
	gopkg.in/yaml.v2 v2.3.0
)

require (
	github.com/VividCortex/gohistogram v1.0.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/go-logfmt/logfmt v0.4.0 // indirect
	github.com/google/certificate-transparency-go v1.0.21 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/prometheus/common v0.6.0 // indirect
	github.com/prometheus/procfs v0.0.3 // indirect
	github.com/spf13/afero v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/weppos/publicsuffix-go v0.5.0 // indirect
	github.com/zmap/zcrypto v0.0.0-20190729165852-9051775e6a2e // indirect
	github.com/zmap/zlint v0.0.0-20190806154020-fd021b4cfbeb // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55 // indirect
)

replace google.golang.org/grpc => ./google.golang.org/grpc
