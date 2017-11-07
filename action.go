package keychain

import (
	"database/sql"
	"net/http"
	"net/url"

	"errors"

	"bullioncoin.githost.io/development/go/signcontrol"
	"bullioncoin.githost.io/development/go/xdr"
	"gitlab.com/distributed_lab/tokend/keychain/actions"
	"gitlab.com/distributed_lab/tokend/keychain/db2"
	"gitlab.com/distributed_lab/tokend/keychain/db2/core"
	"gitlab.com/distributed_lab/tokend/keychain/db2/keychain"
	"gitlab.com/distributed_lab/tokend/keychain/httpx"
	"gitlab.com/distributed_lab/tokend/keychain/log"
	"gitlab.com/distributed_lab/tokend/keychain/render/problem"
	"github.com/zenazn/goji/web"
	"gitlab.com/distributed_lab/logan"
)

// Action is the "base type" for all actions in horizon.  It provides
// structs that embed it with access to the App struct.
//
// Additionally, this type is a trigger for go-codegen and causes
// the file at Action.tmpl to be instantiated for each struct that
// embeds Action.
type Action struct {
	actions.Base
	App *App
	Log *log.Entry

	kq *keychain.Q
	cq core.QInterface
}

func (action *Action) GetAccountIdByBalance(balanceID string) (*string, error) {
	var balance core.Balance
	err := action.CoreQ().BalanceByID(&balance, balanceID)
	if err != nil {
		return nil, err
	}
	return &balance.AccountID, nil
}

func (action *Action) IsAllowed(ownersOfData ...string) {
	if action.Err != nil {
		return
	}

	if len(ownersOfData) == 0 {
		action.Err = errors.New("ownersOfData must not be empty")
		action.Log.WithError(action.Err)
		return
	}

	for _, ownerOfData := range ownersOfData {
		if action.Err != nil && action.Err.Error() != problem.NotAllowed.Error() {
			return
		}
		action.Err = nil
		action.isAllowed(ownerOfData)
		if action.Err == nil {
			return
		}
	}
}

func (action *Action) isAllowed(ownerOfData string) {
	//return if develop mode without signatures is used
	if action.App.config.SkipCheck {
		return
	}

	isSigner := action.IsAccountSigner(action.App.CoreInfo.MasterAccountID, action.Signer)
	if action.Err != nil {
		return
	}

	if isSigner != nil && *isSigner {
		action.IsAdmin = true
		return
	}

	// only master or master signers can access this data
	if ownerOfData == "" || ownerOfData == action.App.CoreInfo.MasterAccountID {
		action.Err = &problem.NotAllowed
		return
	}

	isSigner = action.IsAccountSigner(ownerOfData, action.Signer)
	if action.Err != nil {
		return
	}

	if ownerOfData == action.Signer && isSigner == nil {
		return
	}

	if isSigner != nil && *isSigner {
		return
	}

	action.Err = &problem.NotAllowed
}

// CoreQ provides access to queries that access the stellar core database.
func (action *Action) CoreQ() core.QInterface {
	if action.cq == nil {
		action.cq = &core.Q{Repo: action.App.CoreRepo(action.Ctx)}
	}
	return action.cq
}

// HistoryQ provides access to queries that access the history portion of
// horizon's database.
func (action *Action) KeychainQ() *keychain.Q {
	if action.kq == nil {
		action.kq = &keychain.Q{Repo: action.App.KeychainRepo(action.Ctx)}
	}

	return action.kq
}

// GetPageQuery is a helper that returns a new db.PageQuery struct initialized
// using the results from a call to GetPagingParams()
func (action *Action) GetPageQuery() db2.PageQuery {
	if action.Err != nil {
		return db2.PageQuery{}
	}

	r, err := db2.NewPageQuery(action.GetPagingParams())

	if err != nil {
		action.Err = err
	}

	return r
}

// Prepare sets the action's App field based upon the goji context
func (action *Action) Prepare(c web.C, w http.ResponseWriter, r *http.Request) {
	base := &action.Base
	base.Prepare(c, w, r)
	action.App = action.GojiCtx.Env["app"].(*App)

	base.SkipCheck = action.App.config.SkipCheck //pass config variable to base (since base can't read one)

	base.Signer = r.Header.Get(signcontrol.PublicKeyHeader)

	if action.Ctx != nil {
		action.Log = log.Ctx(action.Ctx)
	} else {
		action.Log = log.DefaultLogger
	}
}

// ValidateCursorAsDefault ensures that the cursor parameter is valid in the way
// it is normally used, i.e. it is either the string "now" or a string of
// numerals that can be parsed as an int64.
func (action *Action) ValidateCursorAsDefault() {
	if action.Err != nil {
		return
	}

	if action.GetString(actions.ParamCursor) == "now" {
		return
	}

	action.GetInt64(actions.ParamCursor)
}

// BaseURL returns the base url for this requestion, defined as a url containing
// the Host and Scheme portions of the request uri.
func (action *Action) BaseURL() *url.URL {
	return httpx.BaseURL(action.Ctx)
}

// IsAccountSigner load core account by `accountId` and checks to see if any of the signers is`signer`
func (action *Action) IsAccountSigner(accountId, signer string) *bool {
	logan.NewWithJSONFormatter()
	var account core.Account
	isSigner := new(bool)
	err := action.CoreQ().AccountByAddress(&account, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		action.Log.WithError(err).Error("Failed to load account")
		action.Err = &problem.ServerError
		*isSigner = false
		return isSigner
	}

	signers, err := action.GetSigners(&account)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load signers")
		action.Err = &problem.ServerError
		*isSigner = false
		return isSigner
	}

	for i := range signers {
		if signer == signers[i].Publickey {
			*isSigner = true
			return isSigner
		}
	}
	*isSigner = false
	return isSigner
}

func (action *Action) GetSigners(account *core.Account) ([]core.Signer, error) {
	// commission and sequence provider accounts are managed by master account signers
	if account.AccountType == int32(xdr.AccountTypeCommission) {
		var masterAccount core.Account
		err := action.CoreQ().AccountByAddress(&masterAccount, action.App.CoreInfo.MasterAccountID)
		if err != nil {
			action.Log.WithError(err).Error("Failed to get master account from db")
			return nil, err
		}

		return action.GetSigners(&masterAccount)
	}

	var signers []core.Signer
	err := action.CoreQ().SignersByAddress(&signers, account.Accountid)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get signers")
		return nil, err
	}

	// is master key allowed
	if account.Thresholds[0] <= 0 {
		return signers, nil
	}

	signers = append(signers, core.Signer{
		Accountid:  account.Accountid,
		Publickey:  account.Accountid,
		Weight:     int32(account.Thresholds[0]),
		SignerType: action.getMasterSignerType(),
		Identity:   0,
	})

	return signers, nil
}

func (action *Action) getMasterSignerType() int32 {
	result := int32(0)
	for i := range xdr.SignerTypeAll {
		result |= int32(xdr.SignerTypeAll[i])
	}
	return result
}
