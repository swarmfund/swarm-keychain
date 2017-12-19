package coreinfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"gitlab.com/swarmfund/api/log"
)

// Connector is a structure with methods for getting core info.
type Connector struct {
	url    *url.URL
	client *http.Client
}

// Info represents a response of Horizon / endpoint
// and contains base core information.
type Info struct {
	HistoryLatestLedger  int64  `json:"history_latest_ledger"`
	HistoryElderLedger   int64  `json:"history_elder_ledger"`
	CoreLatestLedger     int64  `json:"core_latest_ledger"`
	CoreElderLedger      int64  `json:"core_elder_ledger"`
	NetworkPassphrase    string `json:"network_passphrase"`
	CommissionAccountID  string `json:"commission_account_id"`
	OperationalAccountID string `json:"operational_account_id"`
	MasterAccountID      string `json:"master_account_id"`
	MasterExchangeName   string `json:"master_exchange_name"`
}

var coreInfo *Info

// NewConnector is returns new instance of Connector.
func NewConnector(coreURL string) (*Connector, error) {
	u, err := url.Parse(coreURL)
	if err != nil {
		return nil, err
	}

	conn := &Connector{
		url:    u,
		client: &http.Client{},
	}
	_ = conn.Info()
	return conn, nil
}

// Info returns a cached core info.
func (c *Connector) Info() Info {
	if coreInfo != nil {
		return *coreInfo
	}

	info, err := c.getCoreInfo()
	if err != nil {
		panic(err)
	}

	go c.runUpdater()
	coreInfo = info
	return *coreInfo

}

// GetMasterAccountID is returns accountID of the Master account,
// method for implementing the data.CoreInfoI.
func (c *Connector) GetMasterAccountID() string {
	return c.Info().MasterAccountID
}

// runUpdater starts an infinite loop in which it updates coreInfo every hour.
func (c *Connector) runUpdater() {
	entry := log.WithField("service", "corer")
	var info *Info
	var err error

	for {
		time.Sleep(1 * time.Hour)

		info, err = c.getCoreInfo()
		if err != nil {
			entry.WithError(err).Error("unable to update core info")
			continue
		}

		coreInfo = info
		entry.Debug("core info updated")
	}
}

// getCoreInfo get date from the Core /info endpoint.
func (c *Connector) getCoreInfo() (*Info, error) {
	c.url.Path = "/info"
	resp, err := c.client.Get(c.url.String())
	if err != nil {
		return nil, err
	}

	var response Info
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	return &response, err
}
