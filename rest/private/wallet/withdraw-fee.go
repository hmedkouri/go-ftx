package wallet

import (
	"net/http"
)

type RequestForWithdrawFees struct {
	Coin    string  `json:"coin"`
	Size    float64 `json:"size"`
	Address string  `json:"address"`
	// Optionals
	Tag      string `json:"tag,omitempty"`
	Methods  string `json:"method,omitempty"`
}

type ResponseForWithdrawFees struct {
	Methods  string `json:"method"`
	Address string `json:"address"`	
	Fee  float64 `json:"fee"`
	Congested bool `json:"congested"`
}

func (req *RequestForWithdrawFees) Path() string {
	return "/wallet/withdrawal_fee"
}

func (req *RequestForWithdrawFees) Method() string {
	return http.MethodGet
}

func (req *RequestForWithdrawFees) Query() string {
	return ""
}

func (req *RequestForWithdrawFees) Payload() []byte {
	b, err := json.Marshal(req)
	if err != nil {
		return nil
	}
	return b
}
