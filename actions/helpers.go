package actions

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"gitlab.com/tokend/go/strkey"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/keychain/db2"
	"gitlab.com/tokend/keychain/db2/core"
	"gitlab.com/tokend/keychain/render/problem"
	"gitlab.com/tokend/keychain/utils"
)

const (
	// ParamCursor is a query string param name
	ParamCursor = "cursor"
	// ParamOrder is a query string param name
	ParamOrder = "order"
	// ParamLimit is a query string param name
	ParamLimit = "limit"
)

// GetString retrieves a string from either the URLParams, form or query string.
// This method uses the priority (URLParams, Form, Query).
func (base *Base) GetString(name string) string {
	if base.Err != nil {
		return ""
	}

	fromURL, ok := base.GojiCtx.URLParams[name]

	if ok {
		return fromURL
	}

	if base.isJson {
		fromJson := base.JsonValue(name)
		if fromJson != "" {
			return fromJson
		}
	} else {
		fromForm := base.R.FormValue(name)

		if fromForm != "" {
			return fromForm
		}
	}

	return base.R.URL.Query().Get(name)
}

// GetNonEmptyString retrieves an string from the action parameter of the given name.
// Populates err if the value is an empty string
func (base *Base) GetNonEmptyString(name string) string {
	if base.Err != nil {
		return ""
	}

	value := base.GetString(name)
	if value == "" {
		base.SetInvalidField(name, errors.New("Must not be empty."))
	}

	return value
}

func (base *Base) GetBalanceIDAsString(name string) string {
	if base.Err != nil {
		return ""
	}

	rawValue := base.GetNonEmptyString(name)
	if base.Err != nil {
		return ""
	}

	_, err := strkey.Decode(strkey.VersionByteBalanceID, rawValue)
	if err != nil {
		base.SetInvalidField(name, err)
	}

	return rawValue
}

func (base *Base) GetRestrictedString(name string, minLength, maxLength int) string {
	rawValue := base.GetNonEmptyString(name)
	if base.Err != nil {
		return ""
	}
	length := len(rawValue)
	validLength := (length >= minLength) && (length <= maxLength)
	if !validLength {
		base.SetInvalidField(name, errors.New(fmt.Sprintf(" is not %d-%d characters", minLength, maxLength)))
		return ""
	}

	return rawValue
}

// GetNonEmptyString retrieves an string from the action parameter of the given name.
// Populates err if the value is an empty string
func (base *Base) GetStringWithFlag(name string, nonEmpty bool) string {
	if nonEmpty {
		return base.GetNonEmptyString(name)
	}

	return base.GetString(name)
}

// GetMobileNumber retrieves an string from the action parameter
// of the given name and split whitespaces
// Populates err if the value is an empty string
func (base *Base) GetStingWithoutWhitespaces(name string) string {
	str := base.GetNonEmptyString(name)
	if base.Err != nil {
		return ""
	}

	return strings.Replace(str, " ", "", -1)
}

// GetInt64 retrieves an int64 from the action parameter of the given name.
// Populates err if the value is not a valid int64
func (base *Base) GetInt64(name string) int64 {
	if base.Err != nil {
		return 0
	}

	asStr := base.GetString(name)

	if asStr == "" {
		return 0
	}

	asI64, err := strconv.ParseInt(asStr, 10, 64)

	if err != nil {
		base.SetInvalidField(name, err)
		return 0
	}

	return asI64
}

// GetInt32 retrieves an int32 from the action parameter of the given name.
// Populates err if the value is not a valid int32
func (base *Base) GetInt32(name string) int32 {
	if base.Err != nil {
		return 0
	}

	asStr := base.GetString(name)

	if asStr == "" {
		return 0
	}

	asI64, err := strconv.ParseInt(asStr, 10, 32)

	if err != nil {
		base.SetInvalidField(name, err)
		return 0
	}

	return int32(asI64)
}

func (base *Base) GetOptionalBool(name string) *bool {
	if base.Err != nil {
		return nil
	}

	asStr := base.GetString(name)
	if asStr == "" {
		return nil
	}

	result, err := strconv.ParseBool(asStr)
	if err != nil {
		base.SetInvalidField(name, err)
		return nil
	}

	return &result
}

// GetInt32 retrieves an int32 from the action parameter of the given name.
// Populates err if the value is not a valid int32
func (base *Base) GetBool(name string) bool {
	if base.Err != nil {
		return false
	}

	asStr := base.GetString(name)

	if asStr == "true" {
		return true
	} else {
		return false
	}
}

// GetUInt64 retrieves an uint64 from the action parameter of the given name.
// Populates err if the value is not a valid uint64
func (base *Base) GetUInt64(name string) uint64 {
	if base.Err != nil {
		return 0
	}

	asStr := base.GetString(name)

	if asStr == "" {
		return 0
	}

	asUI64, err := strconv.ParseUint(asStr, 10, 64)

	if err != nil {
		base.SetInvalidField(name, err)
		return 0
	}

	return asUI64
}

// GetPagingParams returns the cursor/order/limit triplet that is the
// standard way of communicating paging data to a horizon endpoint.
func (base *Base) GetPagingParams() (cursor string, order string, limit uint64) {
	if base.Err != nil {
		return
	}

	cursor = base.GetString(ParamCursor)
	order = base.GetString(ParamOrder)
	// TODO: add GetUint64 helpers
	limit = uint64(base.GetInt64(ParamLimit))

	if lei := base.R.Header.Get("Last-Event-ID"); lei != "" {
		cursor = lei
	}

	return
}

// GetPageQuery is a helper that returns a new db.PageQuery struct initialized
// using the results from a call to GetPagingParams()
func (base *Base) GetPageQuery() db2.PageQuery {
	if base.Err != nil {
		return db2.PageQuery{}
	}

	r, err := db2.NewPageQuery(base.GetPagingParams())

	if err != nil {
		base.Err = err
	}

	return r
}

// GetAddress retrieves a stellar address.  It confirms the value loaded is a
// valid stellar address, setting an invalid field error if it is not.
func (base *Base) GetAddress(name string) (result string) {
	if base.Err != nil {
		return
	}

	result = base.GetString(name)

	_, err := strkey.Decode(strkey.VersionByteAccountID, result)

	if err != nil {
		base.SetInvalidField(name, err)
	}

	return result
}

// GetAccountID retireves an xdr.AccountID by attempting to decode a stellar
// address at the provided name.
func (base *Base) GetAccountID(name string) (result xdr.AccountId) {
	raw, err := strkey.Decode(strkey.VersionByteAccountID, base.GetString(name))

	if base.Err != nil {
		return
	}

	if err != nil {
		base.SetInvalidField(name, err)
		return
	}

	var key xdr.Uint256
	copy(key[:], raw)

	result, err = xdr.NewAccountId(xdr.CryptoKeyTypeKeyTypeEd25519, key)
	if err != nil {
		base.SetInvalidField(name, err)
		return
	}

	return
}

// SetInvalidField establishes an error response triggered by an invalid
// input field from the user.
func (base *Base) SetInvalidField(name string, reason error) {
	br := problem.BadRequest

	br.Extras = map[string]interface{}{}
	br.Extras["invalid_field"] = name
	br.Extras["reason"] = reason.Error()

	base.Err = &br
}

// Path returns the current action's path, as determined by the http.Request of
// this action
func (base *Base) Path() string {
	return base.R.URL.Path
}

// ValidateBodyType sets an error on the action if the requests Content-Type
//  is not `application/x-www-form-urlencoded`
func (base *Base) ValidateBodyType() {
	c := base.R.Header.Get("Content-Type")

	if c == "" {
		return
	}

	mt, _, err := mime.ParseMediaType(c)

	if err != nil {
		base.Err = err
		return
	}

	switch {
	case mt == "application/x-www-form-urlencoded":
		return
	case mt == "multipart/form-data":
		return
	case mt == "application/json":
		base.isJson = true
		return
	default:
		base.Err = &problem.UnsupportedMediaType
	}
}

func (base *Base) GetCoreAccount(name string, coreQ *core.Q) *core.Account {
	if base.Err != nil {
		return nil
	}

	accountID := base.GetString(name)
	if accountID == "" {
		return nil
	}

	return base.GetCoreAccountByAddress(accountID, coreQ)
}

func (base *Base) GetCoreAccountByAddress(accountID string, coreQ *core.Q) *core.Account {
	if base.Err != nil {
		return nil
	}

	var account core.Account
	err := coreQ.AccountByAddress(&account, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			base.Err = &problem.NotFound
			return nil
		}
		base.Err = &problem.ServerError
		return nil
	}
	return &account
}

func (base *Base) UnmarshalBody(dest interface{}) {
	if !base.isJson {
		base.Err = &problem.UnsupportedMediaType
		return
	}
	decoder := json.NewDecoder(base.R.Body)
	err := decoder.Decode(&dest)
	if err != nil {
		base.Err = &problem.BadRequest
		return
	}
	base.ValidateToProblem(utils.ValidateStruct("", dest))
}

func (base *Base) ValidateToProblem(ok bool, result *utils.ValidateError) {
	if !ok {
		if result != nil {
			base.SetInvalidField(result.Name, result.Reason)
			return
		}
		base.Err = &problem.BadRequest
		return
	}
}

func (base *Base) ParseResponse(response *http.Response) (p *problem.P) {
	switch response.StatusCode {
	case http.StatusOK:
		{
			p = &problem.Success
			break
		}
	case http.StatusNotFound:
		{
			p = &problem.NotFound
			break
		}
	case http.StatusUnauthorized:
		{
			p = &problem.NotAllowed
			break
		}
	case http.StatusBadRequest:
		{
			p = &problem.BadRequest
			break
		}
	case http.StatusInternalServerError:
		{
			p = &problem.ServerError
			break
		}
	default:
		p = &problem.P{
			Type:   response.Status,
			Title:  response.Status,
			Status: response.StatusCode,
		}
	}
	return p
}
