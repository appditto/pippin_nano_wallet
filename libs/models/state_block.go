package models

import (
	"math/big"

	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"
	"golang.org/x/crypto/blake2b"
)

// StateBlock is a block from the nano protocol
type StateBlock struct {
	Type           string  `json:"type"`
	Hash           []byte  `json:"hash"`
	Account        string  `json:"account"`
	Previous       []byte  `json:"previous"`
	Representative string  `json:"representative"`
	Balance        big.Int `json:"balance"`
	Link           []byte  `json:"link"`
	Work           string  `json:"work"`
	Signature      []byte  `json:"signature"`
}

func (b *StateBlock) computeHash() error {
	h, err := blake2b.New256(nil)
	if err != nil {
		return err
	}
	h.Write(make([]byte, 31))
	h.Write([]byte{6})
	pubkey, err := utils.AddressToPub(b.Account)
	if err != nil {
		return err
	}
	h.Write(pubkey)
	h.Write(b.Previous)
	pubkey, err = utils.AddressToPub(b.Representative)
	if err != nil {
		return err
	}
	h.Write(pubkey)
	h.Write(b.Balance.FillBytes(make([]byte, 16)))
	h.Write(b.Link)
	b.Hash = h.Sum(nil)
	return nil
}

func (b *StateBlock) Sign(privateKey ed25519.PrivateKey) error {
	if err := b.computeHash(); err != nil {
		return err
	}
	b.Signature = ed25519.Sign(privateKey, b.Hash)
	return nil
}
