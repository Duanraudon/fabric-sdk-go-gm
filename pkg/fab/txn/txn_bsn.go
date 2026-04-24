package txn

import (
	reqContext "context"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/common/providers/fab"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/context"
	"github.com/golang/protobuf/proto"
	//"github.com/Duanraudon/fabric-sdk-go-gm/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric-protos-go/common"
	//protos_utils "github.com/Duanraudon/fabric-sdk-go-gm/third_party/github.com/hyperledger/fabric/protos/utils"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
)

// Send send a transaction to the chain’s orderer service (one or more orderer endpoints) for consensus and committing to the ledger.
func BsnSend(reqCtx reqContext.Context, tx *fab.Transaction, orderers []fab.Orderer) (*fab.TransactionResponse, error) {
	//GatewayLog.Logs( "Send Start",)
	if len(orderers) == 0 {
		return nil, errors.New("orderers is nil")
	}

	//for _,o := range orderers{
	// //GatewayLog.Logs( "Send Orderer",o.URL())
	//}

	if tx == nil {
		return nil, errors.New("transaction is nil")
	}
	if tx.Proposal == nil || tx.Proposal.Proposal == nil {
		return nil, errors.New("proposal is nil")
	}

	// the original header
	hdr, err := GetHeader(tx.Proposal.Proposal.Header)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal proposal header failed")
	}
	//替换为Admin的Sign的
	ctx, ok := context.RequestClientContext(reqCtx)
	if !ok {
		return nil, errors.New("failed get client context from reqContext for signPayload")
	}

	creator, err := ctx.Serialize()

	if err != nil {
		return nil, errors.Wrap(err, "creator failed")
	}

	sign := &common.SignatureHeader{}

	err = proto.Unmarshal(hdr.SignatureHeader, sign)

	if err != nil {
		return nil, errors.Wrap(err, "sign failed")
	}
	sign.Creator = creator

	hdr.SignatureHeader, err = proto.Marshal(sign)

	if err != nil {
		return nil, errors.Wrap(err, "sign  Marshal failed")
	}

	// serialize the tx
	txBytes, err := GetBytesTransaction(tx.Transaction)
	if err != nil {
		return nil, err
	}

	// create the payload
	payload := common.Payload{Header: hdr, Data: txBytes}

	transactionResponse, err := BroadcastPayload(reqCtx, &payload, orderers)
	if err != nil {
		return nil, err
	}
	//GatewayLog.Logs( "Send End",)
	return transactionResponse, nil
}

// GetHeader Get Header from bytes
func GetHeader(bytes []byte) (*common.Header, error) {
	hdr := &common.Header{}
	err := proto.Unmarshal(bytes, hdr)
	return hdr, errors.Wrap(err, "error unmarshaling Header")
}

// GetBytesTransaction get the bytes of Transaction from the message

func GetBytesTransaction(tx *peer.Transaction) ([]byte, error) {
	bytes, err := proto.Marshal(tx)
	return bytes, errors.Wrap(err, "error unmarshaling Transaction")
}
