package wallet

import (
	"coins_wallet/payment"
	"context"
	"encoding/json"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

const (
	// API error codes.
	notPositivePaymentAmountErrCode = "NOT_POSITIVE_PAYMENT_AMOUNT"
	lowBalanceErrCode               = "LOW_BALANCE"
	paymentToSameAccountErrCode     = "PAYMENT_TO_SAME_ACCOUNT"
	fromAccountNotFoundErrCode      = "FROM_ACCOUNT_NOT_FOUND"
	toAccountNotFoundErrCode        = "TO_ACCOUNT_NOT_FOUND"
	differentCurrenciesErrCode      = "DIFFERENT_CURRENCIES"
	incorrectRequestErrCode         = "INCORRECT_REQUEST"
	internalErrorErrCode            = "INTERNAL_ERROR"
)

func MakeHandler(s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}
	sendPaymentHandler := kithttp.NewServer(
		makeSendPaymentEndpoint(s),
		decodeSendPaymentRequest,
		encodeResponse,
		opts...,
	)
	getAllPaymentsHandler := kithttp.NewServer(
		makeGetAllPaymentsEndpoint(s),
		decodeGetAllPaymentsRequest,
		encodeResponse,
		opts...,
	)
	getAllAccountsHandler := kithttp.NewServer(
		makeGetAllAccountsEndpoint(s),
		decodeGetAllAccountsRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/wallet/v1/payments", sendPaymentHandler).Methods("POST")
	r.Handle("/wallet/v1/payments", getAllPaymentsHandler).Methods("GET")
	r.Handle("/wallet/v1/accounts", getAllAccountsHandler).Methods("GET")

	return r
}

type decodingError struct {
	Details string
}

func (de *decodingError) Error() string {
	return de.Details
}

func decodeSendPaymentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	fromAccountId := r.PostFormValue("from_account")
	toAccountId := r.PostFormValue("to_account")
	if fromAccountId == "" || toAccountId == "" {
		return nil, &decodingError{"'from_account' and 'to_account' are required"}
	}

	amount, err := strconv.ParseFloat(r.PostFormValue("amount"), 64)
	if err != nil {
		return nil, &decodingError{"'amount' is required and must have a float format"}
	}

	return &sendPaymentRequest{FromAccountId: fromAccountId, ToAccountId: toAccountId, Amount: amount}, nil
}

func decodePaginationRequest(r *http.Request) (*paginationRequest, error) {
	var offset, limit int
	offsetText := r.FormValue("offset")
	if offsetText != "" {
		var err error
		offset, err = strconv.Atoi(offsetText)
		if err != nil {
			return nil, &decodingError{"'offset' must be an int"}
		}
	}
	limitText := r.FormValue("limit")
	if limitText != "" {
		var err error
		limit, err = strconv.Atoi(limitText)
		if err != nil {
			return nil, &decodingError{"'limit' must be an int"}
		}
	}
	return &paginationRequest{Offset: offset, Limit: limit}, nil
}

func decodeGetAllPaymentsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	decoded, err := decodePaginationRequest(r)
	if err != nil {
		return nil, err
	}
	return &getAllPaymentsRequest{paginationRequest: decoded}, nil
}

func decodeGetAllAccountsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	decoded, err := decodePaginationRequest(r)
	if err != nil {
		return nil, err
	}
	return &getAllAccountsRequest{paginationRequest: decoded}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var httpStatusCode int
	switch response.(type) {
	case *sendPaymentResponse:
		httpStatusCode = http.StatusCreated
	default:
		httpStatusCode = http.StatusOK
	}
	w.WriteHeader(httpStatusCode)
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var errorCode string
	var httpStatusCode int
	switch err.(type) {
	case *decodingError:
		errorCode, httpStatusCode = incorrectRequestErrCode, http.StatusBadRequest
	case *DifferentCurrenciesError:
		errorCode, httpStatusCode = differentCurrenciesErrCode, http.StatusConflict
	default:
		switch err {
		case PaymentToSameAccountErr:
			errorCode, httpStatusCode = paymentToSameAccountErrCode, http.StatusBadRequest
		case payment.LowBalanceErr:
			errorCode, httpStatusCode = lowBalanceErrCode, http.StatusBadRequest
		case NotPositivePaymentAmount:
			errorCode, httpStatusCode = notPositivePaymentAmountErrCode, http.StatusBadRequest
		case FromAccountNotFound:
			errorCode, httpStatusCode = fromAccountNotFoundErrCode, http.StatusNotFound
		case ToAccountNotFound:
			errorCode, httpStatusCode = toAccountNotFoundErrCode, http.StatusNotFound
		default:
			errorCode, httpStatusCode = internalErrorErrCode, http.StatusInternalServerError
		}
	}
	w.WriteHeader(httpStatusCode)
	errorMessage := strings.ToUpper(err.Error()[:1]) + err.Error()[1:] + "."
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    errorCode,
			"message": errorMessage,
		},
	})
}
