package resource

import "bullioncoin.githost.io/development/keychain/db2/keychain"

type JWK struct {
	KeyType   string `json:"kty"`
	Key       string `json:"k"`
	Algorithm string `json:"alg"`
}

func (jwk *JWK) Populate(key *keychain.Key) {
	jwk.KeyType = "oct"
	jwk.Key = key.Key
	jwk.Algorithm = "A256GCM"
}
