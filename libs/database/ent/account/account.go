// Code generated by ent, DO NOT EDIT.

package account

import (
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the account type in the database.
	Label = "account"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldWalletID holds the string denoting the wallet_id field in the database.
	FieldWalletID = "wallet_id"
	// FieldAddress holds the string denoting the address field in the database.
	FieldAddress = "address"
	// FieldAccountIndex holds the string denoting the account_index field in the database.
	FieldAccountIndex = "account_index"
	// FieldPrivateKey holds the string denoting the private_key field in the database.
	FieldPrivateKey = "private_key"
	// FieldWork holds the string denoting the work field in the database.
	FieldWork = "work"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// EdgeWallet holds the string denoting the wallet edge name in mutations.
	EdgeWallet = "wallet"
	// EdgeBlocks holds the string denoting the blocks edge name in mutations.
	EdgeBlocks = "blocks"
	// Table holds the table name of the account in the database.
	Table = "accounts"
	// WalletTable is the table that holds the wallet relation/edge.
	WalletTable = "accounts"
	// WalletInverseTable is the table name for the Wallet entity.
	// It exists in this package in order to avoid circular dependency with the "wallet" package.
	WalletInverseTable = "wallets"
	// WalletColumn is the table column denoting the wallet relation/edge.
	WalletColumn = "wallet_id"
	// BlocksTable is the table that holds the blocks relation/edge.
	BlocksTable = "blocks"
	// BlocksInverseTable is the table name for the Block entity.
	// It exists in this package in order to avoid circular dependency with the "block" package.
	BlocksInverseTable = "blocks"
	// BlocksColumn is the table column denoting the blocks relation/edge.
	BlocksColumn = "account_id"
)

// Columns holds all SQL columns for account fields.
var Columns = []string{
	FieldID,
	FieldWalletID,
	FieldAddress,
	FieldAccountIndex,
	FieldPrivateKey,
	FieldWork,
	FieldCreatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// AddressValidator is a validator for the "address" field. It is called by the builders before save.
	AddressValidator func(string) error
	// PrivateKeyValidator is a validator for the "private_key" field. It is called by the builders before save.
	PrivateKeyValidator func(string) error
	// DefaultWork holds the default value on creation for the "work" field.
	DefaultWork bool
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)
