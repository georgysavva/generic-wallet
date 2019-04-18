package wallet

import (
	"context"
	"generic_wallet/account"
	"generic_wallet/payment"
	"github.com/go-kit/kit/endpoint"
)

type sendPaymentRequest struct {
	FromAccountId string
	ToAccountId   string
	Amount        float64
}

type sendPaymentResponse struct {
	Ok bool `json:"ok"`
}

func makeSendPaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*sendPaymentRequest)
		err := s.SendPayment(ctx, req.FromAccountId, req.ToAccountId, req.Amount)
		if err != nil {
			return nil, err
		}
		return &sendPaymentResponse{Ok: true}, nil
	}
}

type paginationRequest struct {
	Offset int
	Limit  int
}

type getAllPaymentsRequest struct {
	*paginationRequest
}

type getAllPaymentsResponse struct {
	Results []*payment.Payment `json:"results"`
	Total   int                `json:"total"`
}

func makeGetAllPaymentsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*getAllPaymentsRequest)
		payments, totalNumber, err := s.GetAllPayments(ctx, req.Offset, req.Limit)
		if err != nil {
			return nil, err
		}
		if payments == nil {
			payments = []*payment.Payment{}
		}
		return &getAllPaymentsResponse{Results: payments, Total: totalNumber}, nil
	}
}

type getAllAccountsRequest struct {
	*paginationRequest
}

type getAllAccountsResponse struct {
	Results []*account.Account `json:"results"`
	Total   int                `json:"total"`
}

func makeGetAllAccountsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*getAllAccountsRequest)
		accounts, totalNumber, err := s.GetAllAccounts(ctx, req.Offset, req.Limit)
		if err != nil {
			return nil, err
		}
		if accounts == nil {
			accounts = []*account.Account{}
		}
		return &getAllAccountsResponse{Results: accounts, Total: totalNumber}, nil
	}
}
