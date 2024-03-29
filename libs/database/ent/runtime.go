// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/appditto/pippin_nano_wallet/libs/database/ent/account"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent/block"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent/schema"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent/wallet"
	"github.com/google/uuid"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	accountFields := schema.Account{}.Fields()
	_ = accountFields
	// accountDescAddress is the schema descriptor for address field.
	accountDescAddress := accountFields[2].Descriptor()
	// account.AddressValidator is a validator for the "address" field. It is called by the builders before save.
	account.AddressValidator = accountDescAddress.Validators[0].(func(string) error)
	// accountDescPrivateKey is the schema descriptor for private_key field.
	accountDescPrivateKey := accountFields[4].Descriptor()
	// account.PrivateKeyValidator is a validator for the "private_key" field. It is called by the builders before save.
	account.PrivateKeyValidator = accountDescPrivateKey.Validators[0].(func(string) error)
	// accountDescWork is the schema descriptor for work field.
	accountDescWork := accountFields[5].Descriptor()
	// account.DefaultWork holds the default value on creation for the work field.
	account.DefaultWork = accountDescWork.Default.(bool)
	// accountDescCreatedAt is the schema descriptor for created_at field.
	accountDescCreatedAt := accountFields[6].Descriptor()
	// account.DefaultCreatedAt holds the default value on creation for the created_at field.
	account.DefaultCreatedAt = accountDescCreatedAt.Default.(func() time.Time)
	// accountDescID is the schema descriptor for id field.
	accountDescID := accountFields[0].Descriptor()
	// account.DefaultID holds the default value on creation for the id field.
	account.DefaultID = accountDescID.Default.(func() uuid.UUID)
	blockFields := schema.Block{}.Fields()
	_ = blockFields
	// blockDescBlockHash is the schema descriptor for block_hash field.
	blockDescBlockHash := blockFields[2].Descriptor()
	// block.BlockHashValidator is a validator for the "block_hash" field. It is called by the builders before save.
	block.BlockHashValidator = blockDescBlockHash.Validators[0].(func(string) error)
	// blockDescSendID is the schema descriptor for send_id field.
	blockDescSendID := blockFields[4].Descriptor()
	// block.SendIDValidator is a validator for the "send_id" field. It is called by the builders before save.
	block.SendIDValidator = blockDescSendID.Validators[0].(func(string) error)
	// blockDescSubtype is the schema descriptor for subtype field.
	blockDescSubtype := blockFields[5].Descriptor()
	// block.SubtypeValidator is a validator for the "subtype" field. It is called by the builders before save.
	block.SubtypeValidator = blockDescSubtype.Validators[0].(func(string) error)
	// blockDescCreatedAt is the schema descriptor for created_at field.
	blockDescCreatedAt := blockFields[6].Descriptor()
	// block.DefaultCreatedAt holds the default value on creation for the created_at field.
	block.DefaultCreatedAt = blockDescCreatedAt.Default.(func() time.Time)
	// blockDescID is the schema descriptor for id field.
	blockDescID := blockFields[0].Descriptor()
	// block.DefaultID holds the default value on creation for the id field.
	block.DefaultID = blockDescID.Default.(func() uuid.UUID)
	walletFields := schema.Wallet{}.Fields()
	_ = walletFields
	// walletDescSeed is the schema descriptor for seed field.
	walletDescSeed := walletFields[1].Descriptor()
	// wallet.SeedValidator is a validator for the "seed" field. It is called by the builders before save.
	wallet.SeedValidator = walletDescSeed.Validators[0].(func(string) error)
	// walletDescRepresentative is the schema descriptor for representative field.
	walletDescRepresentative := walletFields[2].Descriptor()
	// wallet.RepresentativeValidator is a validator for the "representative" field. It is called by the builders before save.
	wallet.RepresentativeValidator = walletDescRepresentative.Validators[0].(func(string) error)
	// walletDescEncrypted is the schema descriptor for encrypted field.
	walletDescEncrypted := walletFields[3].Descriptor()
	// wallet.DefaultEncrypted holds the default value on creation for the encrypted field.
	wallet.DefaultEncrypted = walletDescEncrypted.Default.(bool)
	// walletDescWork is the schema descriptor for work field.
	walletDescWork := walletFields[4].Descriptor()
	// wallet.DefaultWork holds the default value on creation for the work field.
	wallet.DefaultWork = walletDescWork.Default.(bool)
	// walletDescCreatedAt is the schema descriptor for created_at field.
	walletDescCreatedAt := walletFields[5].Descriptor()
	// wallet.DefaultCreatedAt holds the default value on creation for the created_at field.
	wallet.DefaultCreatedAt = walletDescCreatedAt.Default.(func() time.Time)
	// walletDescID is the schema descriptor for id field.
	walletDescID := walletFields[0].Descriptor()
	// wallet.DefaultID holds the default value on creation for the id field.
	wallet.DefaultID = walletDescID.Default.(func() uuid.UUID)
}
