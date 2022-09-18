package models

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"
	"golang.org/x/crypto/blake2b"
)

// StateBlock is a block from the nano protocol
type StateBlock struct {
	Type           string `json:"type" mapstructure:"type"`
	Hash           string `json:"hash" mapstructure:"hash"`
	Account        string `json:"account" mapstructure:"account"`
	Previous       string `json:"previous" mapstructure:"previous"`
	Representative string `json:"representative" mapstructure:"representative"`
	Balance        string `json:"balance" mapstructure:"balance"`
	Link           string `json:"link" mapstructure:"link"`
	LinkAsAccount  string `json:"link_as_account" mapstructure:"link_as_account"`
	Work           string `json:"work" mapstructure:"work"`
	Signature      string `json:"signature" mapstructure:"signature"`
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
	previous, err := hex.DecodeString(b.Previous)
	if err != nil {
		return err
	}
	h.Write(previous)
	pubkey, err = utils.AddressToPub(b.Representative)
	if err != nil {
		return err
	}
	h.Write(pubkey)
	// COnvert balance to big int
	balance, ok := big.NewInt(0).SetString(b.Balance, 10)
	if !ok {
		return errors.New("Invalid balance")
	}
	h.Write(balance.FillBytes(make([]byte, 16)))
	link, err := hex.DecodeString(b.Link)
	if err != nil {
		return err
	}
	h.Write(link)
	b.Hash = hex.EncodeToString(h.Sum(nil))
	return nil
}

func (b *StateBlock) Sign(privateKey ed25519.PrivateKey) error {
	if err := b.computeHash(); err != nil {
		return err
	}
	hash, err := hex.DecodeString(b.Hash)
	if err != nil {
		return err
	}
	sig := ed25519.Sign(privateKey, hash)
	b.Signature = hex.EncodeToString(sig)
	return nil
}
