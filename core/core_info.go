package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Info struct {
	CoreVersion                string `json:"build"`
	MasterAccountID            string `json:"master_account_id"`
	CommissionAccountID        string `json:"commission_account_id"`
	NetworkPassphrase          string `json:"network"`
	MasterExchangeName         string
	TxExpirationPeriod         int64 `json:"tx_expiration_period"`
	WithdrawalDetailsMaxLength int64 `json:"withdrawal_details_max_length"`
}

type InfoResponse struct {
	Info Info `json:"info"`
}

func (i *Info) validate() error {
	errorProvider := func(name string) error {
		return errors.New(fmt.Sprintf("%s must not be empty. Please check connection with stellar-core", name))
	}
	if i.NetworkPassphrase == "" {
		return errorProvider("NetworkPassphrase")
	}

	if i.MasterAccountID == "" {
		return errorProvider("MasterAccountID")
	}

	if i.CommissionAccountID == "" {
		return errorProvider("CommissionAccountID")
	}

	if i.TxExpirationPeriod <= 0 {
		return errorProvider("TxExpirationPeriod")
	}

	if i.WithdrawalDetailsMaxLength <= 0 {
		return errorProvider("WithdrawalDetailsMaxLength")
	}
	return nil
}

// Returns validated Info, if fails to get or info is invalid - returns error
func GetStellarCoreInfo(coreURl string) (Info, error) {
	if coreURl == "" {
		return Info{}, errors.New("Invalid Stellar Core URl")
	}

	resp, err := http.Get(fmt.Sprint(coreURl, "/info"))
	if err != nil {
		return Info{}, err
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Info{}, err
	}

	var response InfoResponse
	err = json.Unmarshal(contents, &response)
	if err != nil {
		return Info{}, err
	}

	err = response.Info.validate()
	return response.Info, err
}
