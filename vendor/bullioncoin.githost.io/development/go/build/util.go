package build

import (
	"bullioncoin.githost.io/development/go/keypair"
	"bullioncoin.githost.io/development/go/support/errors"
	"bullioncoin.githost.io/development/go/xdr"
)

func setAccountId(addressOrSeed string, aid *xdr.AccountId) error {
	kp, err := keypair.Parse(addressOrSeed)
	if err != nil {
		return err
	}

	if aid == nil {
		return errors.New("aid is nil in setAccountId")
	}

	return aid.SetAddress(kp.Address())
}
