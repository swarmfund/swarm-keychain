package doorman

import (
	"net/http"
	"gitlab.com/tokend/go/signcontrol"
)

func (doorman *Doorman) SignerOf(address string) SignerConstraint {
	return func(r *http.Request) error {
		if doorman.PassAllChecks {
			return nil
		}

		signer, err := signcontrol.CheckSignature(r)
		if err != nil {
			return err
		}

		if signer == address {
			return nil
		}

		signers, err := doorman.AccountQ.Signers(address)
		if err != nil {
			return err
		}

		// TODO make it readable
		for _, accountSigner := range signers {
			if accountSigner.AccountID == signer && accountSigner.Weight > 0 {
				return nil
			}
		}
		return signcontrol.ErrNotAllowed
	}
}
