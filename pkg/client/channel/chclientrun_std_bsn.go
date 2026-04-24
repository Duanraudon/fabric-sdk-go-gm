package channel

import (
	"github.com/Duanraudon/fabric-sdk-go-gm/pkg/client/channel/invoke"
)

func callBsnExecute(cc *Client, request Request, options ...RequestOption) (Response, error) {
	return cc.InvokeHandler(invoke.NewBsnExecuteHandler(), request, options...)
}

func callBsnQuery(cc *Client, request Request, options ...RequestOption) (Response, error) {
	return cc.InvokeHandler(invoke.NewBsnQueryHandler(), request, options...)
}

func callBsnExecuteHasOpts(cc *Client, request Request, opts invoke.TxnHeaderOptsProvider, options ...RequestOption) (Response, error) {
	return cc.InvokeHandler(invoke.NewBsnExecuteHandlerHasOpts(opts), request, options...)
}

func callBsnQueryHasOpts(cc *Client, request Request, opts invoke.TxnHeaderOptsProvider, options ...RequestOption) (Response, error) {
	return cc.InvokeHandler(invoke.NewBsnQueryHandlerHasOpts(opts), request, options...)
}
