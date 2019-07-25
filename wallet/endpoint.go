package wallet

import (
	"context"
	"github.com/georgysavva/generic_wallet/account"
	"github.com/georgysavva/generic_wallet/payment"
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
	Offset *int
	Limit  *int
}

type getAllPaymentsRequest struct {
	*paginationRequest
}

type getAllPaymentsResponse struct {
	Results     []*payment.Payment `json:"results"`
	TotalNumber int                `json:"total_number"`
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
		return &getAllPaymentsResponse{Results: payments, TotalNumber: totalNumber}, nil
	}
}

type getAllAccountsRequest struct {
	*paginationRequest
}

type getAllAccountsResponse struct {
	Results     []*account.Account `json:"results"`
	TotalNumber int                `json:"total_number"`
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
		return &getAllAccountsResponse{Results: accounts, TotalNumber: totalNumber}, nil
	}
}
