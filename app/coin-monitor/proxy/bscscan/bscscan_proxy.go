package bscscan

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

var bscScanRestyClient *resty.Client

func InitEtherScanClient() {
	bscScanRestyClient = resty.New().SetHostURL("https://api.bscscan.com")
}

func GetContractSourceCodeByAddress(ctx context.Context, address string) (*GetContractSourceCodeResponse, error) {
	responseBody := &GetContractSourceCodeResponse{}
	resp, err := bscScanRestyClient.R().SetQueryParams(map[string]string{
		"apiKey":  "MAHVXMJXMQWV1V7FTHDUTTR53F4ZYSDUDQ",
		"module":  "contract",
		"action":  "getsourcecode",
		"address": address,
	}).
		SetResult(responseBody).
		Get("api")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Http Code %+v.Message %+v", resp.StatusCode(), string(resp.Body())))
	}
	return responseBody, nil
}

func GetTransactionsByAddress(ctx context.Context, address string) ([]*TransactionDetail, error) {
	responseBody := &struct {
		Status  string               `json:"status"`
		Message string               `json:"message"`
		Result  []*TransactionDetail `json:"result"`
	}{}
	/**
	  module:account
	  action:txlist
	  address:0x773355277126cbDCf8EB80702f6bc1A3Cb843Bbb
	  startblock:0
	  endblock:99999999
	  sort:desc
	  apikey:IIUZ1TNS9QJVHCTQJNYSHV1ZJ5ARXGSNUH

	*/
	resp, err := bscScanRestyClient.R().SetQueryParams(map[string]string{
		"sort":    "desc",
		"apiKey":  "IIUZ1TNS9QJVHCTQJNYSHV1ZJ5ARXGSNUH",
		"module":  "account",
		"action":  "txlist",
		"address": address,
	}).
		SetResult(responseBody).
		Get("api")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Http Code %+v.Message %+v", resp.StatusCode(), string(resp.Body())))
	}
	return responseBody.Result, nil
}
