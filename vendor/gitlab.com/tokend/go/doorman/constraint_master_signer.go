package doorman

func (doorman *Doorman) MasterSigner(_ string) SignerConstraint {
	return doorman.SignerOf(doorman.MasterAddress)
}
