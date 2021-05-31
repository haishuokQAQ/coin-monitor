package bscscan

type GetContractSourceCodeResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Result  []*SmartContractDetail `json:"result"`
}
