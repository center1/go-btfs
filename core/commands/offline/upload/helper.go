package upload

import (
	"encoding/base64"
	logging "github.com/ipfs/go-log"
	ic "github.com/libp2p/go-libp2p-core/crypto"
)

var log = logging.Logger("upload")

func convertPubKeyFromString(pubKeyStr string) (ic.PubKey, error) {
	raw, err := base64.StdEncoding.DecodeString(pubKeyStr)
	if err != nil {
		return nil, err
	}
	return ic.UnmarshalPublicKey(raw)
}

func convertToPubKey(pubKeyStr string) (ic.PubKey, error) {
	pubKey, err := convertPubKeyFromString(pubKeyStr)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

