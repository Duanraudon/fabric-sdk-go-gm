package txn

import (
	reqContext "context"
	"github.com/pkg/errors"
	"sync"

	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/common/errors/multi"
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/common/providers/fab"
	//pb "github.com/Duanraudon/fabric-sdk-go-gm/third_party/github.com/hyperledger/fabric/protos/peer"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// SendProposal sends a TransactionProposal to ProposalProcessor.
func SendBsnProposal(reqCtx reqContext.Context, signedProposal *pb.SignedProposal, targets []fab.ProposalProcessor) ([]*fab.TransactionProposalResponse, error) {

	if len(targets) < 1 {
		return nil, errors.New("targets is required")
	}

	for _, p := range targets {
		if p == nil {
			return nil, errors.New("target is nil")
		}
	}

	targets = getTargetsWithoutDuplicates(targets)

	request := fab.ProcessProposalRequest{SignedProposal: signedProposal}

	var responseMtx sync.Mutex
	var transactionProposalResponses []*fab.TransactionProposalResponse
	var wg sync.WaitGroup
	errs := multi.Errors{}

	for _, p := range targets {
		wg.Add(1)
		go func(processor fab.ProposalProcessor) {
			defer wg.Done()

			// TODO: The RPC should be timed-out.
			//resp, err := processor.ProcessTransactionProposal(context.NewRequestOLD(ctx), request)
			resp, err := processor.ProcessTransactionProposal(reqCtx, request)
			if err != nil {
				logger.Debugf("Received error response from txn proposal processing: %s", err)
				responseMtx.Lock()
				errs = append(errs, err)
				responseMtx.Unlock()
				return
			}

			responseMtx.Lock()
			transactionProposalResponses = append(transactionProposalResponses, resp)
			responseMtx.Unlock()
		}(p)
	}
	wg.Wait()
	return transactionProposalResponses, errs.ToError()
}
