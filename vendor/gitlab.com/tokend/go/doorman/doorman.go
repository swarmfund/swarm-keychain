package doorman

import (
	"net/http"

	"gitlab.com/tokend/go/signcontrol"
)

type SignerConstraint func(r *http.Request) error

type Doorman struct {
	// SkipSignatureCheck disable signature validation
	SkipSignatureCheck bool
	// PassAllChecks disable constraints validation completely, any request will succeed
	PassAllChecks bool
	// AccountQ used to get account details during constraint checks
	AccountQ AccountQ
	// MasterAddress master account address used for admin signature checks
	MasterAddress string
}

// Check ensures request passes at least one constraint
func Check(r *http.Request, constraints ...SignerConstraint) error {
	for _, constraint := range constraints {
		switch err := constraint(r); err {
		case nil:
			// request passed constraint check
			return nil
		case signcontrol.ErrNotAllowed:
			// check failed, let's try next one
			continue
		default:
			// probably runtime issue
			return err
		}
	}

	// request failed all checks
	return signcontrol.ErrNotAllowed
}
