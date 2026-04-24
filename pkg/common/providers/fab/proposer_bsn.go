package fab

import (
	//pb "github.com/Duanraudon/fabric-sdk-go-gm/third_party/github.com/hyperledger/fabric/protos/peer"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type RequestProposal struct {
	TransactionProposal *TransactionProposal
	SignProposal        *pb.SignedProposal
}
