package wallet

import (
	"context"
	"strings"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/database"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/stretchr/testify/assert"
)

var nanoWallet *NanoWallet

func init() {
	ctx := context.Background()
	dbconn, _ := database.GetSqlDbConn(true)
	client, _ := database.NewEntClient(dbconn)
	if err := client.Schema.Create(ctx); err != nil {
		panic(err)
	}
	nanoWallet = &NanoWallet{
		DB:  client,
		Ctx: ctx,
	}
}

func TestWalletCreate(t *testing.T) {
	// Predictable seed
	seed, _ := utils.GenerateSeed(strings.NewReader("9f729340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040b"))

	wallet, err := nanoWallet.WalletCreate(seed)
	assert.Nil(t, err)

	assert.Equal(t, false, wallet.Encrypted)
	assert.Equal(t, seed, wallet.Seed)
}
