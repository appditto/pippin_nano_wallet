package wallet

import (
	"strings"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/database/ent"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"
	"github.com/stretchr/testify/assert"
)

func TestGetAccount(t *testing.T) {
	// Predictable seed
	seed, _ := utils.GenerateSeed(strings.NewReader("63bbf524e6bae9e2bf22c9967f8fd093d1ce906b5dc581ba2f3529156efbae92"))

	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	assert.Equal(t, false, wallet.Encrypted)
	assert.Equal(t, seed, wallet.Seed)

	// Create an account and an adhoc account
	account, err := MockWallet.AccountCreate(wallet, nil)
	assert.Nil(t, err)

	_, priv, _ := ed25519.GenerateKey(strings.NewReader("1f729340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040d"))
	adhocAcct, err := MockWallet.AdhocAccountCreate(wallet, priv)
	assert.Nil(t, err)

	// Retrieve account
	retrieved, err := MockWallet.GetAccount(wallet, account.Address)
	assert.Nil(t, err)
	assert.Equal(t, account.Address, retrieved.Address)

	// Retrieve adhoc account
	retrievedAdhoc, err := MockWallet.GetAccount(wallet, adhocAcct.Address)
	assert.Nil(t, err)
	assert.Equal(t, adhocAcct.Address, retrievedAdhoc.Address)

	// Retrieve unknown account
	_, err = MockWallet.GetAccount(wallet, "nano_1efa1gxbitary1urzix9h13nkzadtz71n3auyj7uztb8i4qbtipu8cxz61ee")
	assert.ErrorIs(t, ErrAccountNotFound, err)
}

func TestAccountCreate(t *testing.T) {
	// Predictable seed
	seed, _ := utils.GenerateSeed(strings.NewReader("5f729340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040c"))

	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	// Create a couple accounts and ensure they are sequential
	acct, err := MockWallet.AccountCreate(wallet, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, *acct.AccountIndex)
	assert.Equal(t, "nano_3tdqk8ghsdfapzhrag5978izd19minorfmxergefkecdsbyxaw6og4fejs89", acct.Address)

	acct, err = MockWallet.AccountCreate(wallet, nil)
	assert.Nil(t, err)
	assert.Equal(t, 2, *acct.AccountIndex)
	assert.Equal(t, "nano_1frwge7oebdn87jip7k3sa1uuyf4yxxjh8jg67i69r7smf7tddj1gr6yremf", acct.Address)

	// Test collision
	idx := 2
	acct, err = MockWallet.AccountCreate(wallet, &idx)
	assert.ErrorIs(t, ErrAccountExists, err)

	// Test at index
	idx = 3
	acct, err = MockWallet.AccountCreate(wallet, &idx)
	assert.Nil(t, err)
	assert.Equal(t, "69d5b476d8e892213a17511b977de91eb41adc9b8bf46e25a42da3d4de70427bd194c238a97aa0aa436cfb4507f2bd4c1e8e4adb5f91f0e56eb30ee4da723493", *acct.PrivateKey)
	assert.Equal(t, "nano_3nenrawckyo1ob3psyt71zsdtm1yjs7fpqwjy5kpxergwmf96f6mdenouyyk", acct.Address)

	// Test now that sequence is broken
	acct, err = MockWallet.AccountCreate(wallet, nil)
	assert.Nil(t, err)
	assert.Equal(t, 4, *acct.AccountIndex)
	assert.Equal(t, "nano_1w84hd687n5nywb1zfjat7ghe7czgo39tyeqbnbi67icr4fornhqgh3w5q6a", acct.Address)

	// Check that it fails if wallet is locked
	MockWallet.EncryptWallet(wallet, "password")
	_, err = MockWallet.AccountCreate(wallet, nil)
	assert.ErrorIs(t, ErrWalletLocked, err)
}

func TestAccountCreateBadInput(t *testing.T) {
	// Empty seed
	_, err := MockWallet.AccountCreate(nil, nil)
	assert.ErrorIs(t, ErrInvalidWallet, err)
}

func TestAdhocAccountCreate(t *testing.T) {
	// Predictable seed
	seed, _ := utils.GenerateSeed(strings.NewReader("aa729340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040d"))

	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	// Create a couple accounts and ensure they are sequential
	acct, err := MockWallet.AccountCreate(wallet, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, *acct.AccountIndex)
	assert.Equal(t, "nano_3suihcm3txrfcnipecixeu7kcdm8jisyb1osy14j1r5na1c17g7kkbbu143o", acct.Address)

	_, priv, _ := ed25519.GenerateKey(strings.NewReader("1f729340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040d"))
	adhocAcct, err := MockWallet.AdhocAccountCreate(wallet, priv)
	assert.Nil(t, err)
	assert.Equal(t, "31663732393334306530376565653639616261633034396332666464346133632c9c941001f4236f487aac98848320ce746f7b809af76e5e795835ef30022424", *adhocAcct.PrivateKey)
	assert.Equal(t, "nano_1d6wkia15x35fx69od6rik3k3mmnfxxr38qqfsh9kp3oxwr16b36eztouu3a", adhocAcct.Address)
	assert.Equal(t, wallet.ID, adhocAcct.WalletID)

	// Test duplicate returns same account
	adhocAcct, err = MockWallet.AdhocAccountCreate(wallet, priv)
	assert.Nil(t, err)
	assert.Equal(t, "31663732393334306530376565653639616261633034396332666464346133632c9c941001f4236f487aac98848320ce746f7b809af76e5e795835ef30022424", *adhocAcct.PrivateKey)
	assert.Equal(t, "nano_1d6wkia15x35fx69od6rik3k3mmnfxxr38qqfsh9kp3oxwr16b36eztouu3a", adhocAcct.Address)
	assert.Equal(t, wallet.ID, adhocAcct.WalletID)

	// Test we get non-adhoc account if we are trying to re-create it
	_, priv, _ = utils.KeypairFromSeed(seed, 1)
	acct, err = MockWallet.AdhocAccountCreate(wallet, priv)
	assert.Equal(t, "nano_3suihcm3txrfcnipecixeu7kcdm8jisyb1osy14j1r5na1c17g7kkbbu143o", acct.Address)
	assert.Equal(t, 1, *acct.AccountIndex)

	// Check that it fails if wallet is locked
	MockWallet.EncryptWallet(wallet, "password")
	_, err = MockWallet.AdhocAccountCreate(wallet, priv)
	assert.ErrorIs(t, ErrWalletLocked, err)
}

func TestAdhocAccountCreateBadInput(t *testing.T) {
	_, err := MockWallet.AdhocAccountCreate(nil, nil)
	assert.ErrorIs(t, ErrInvalidWallet, err)
	// Bad private key
	_, err = MockWallet.AdhocAccountCreate(&ent.Wallet{}, nil)
	assert.ErrorIs(t, ErrInvalidPrivKey, err)
}

func TestAccountsCreate(t *testing.T) {
	// Predictable seed
	seed, _ := utils.GenerateSeed(strings.NewReader("9f729340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040e"))

	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	// Create a couple accounts and ensure they are sequential
	accts, err := MockWallet.AccountsCreate(wallet, 10)
	assert.Nil(t, err)
	assert.Len(t, accts, 10)
	assert.Equal(t, 1, *accts[0].AccountIndex)
	assert.Equal(t, 2, *accts[1].AccountIndex)
	assert.Equal(t, 3, *accts[2].AccountIndex)
	assert.Equal(t, 4, *accts[3].AccountIndex)
	assert.Equal(t, 5, *accts[4].AccountIndex)
	assert.Equal(t, 6, *accts[5].AccountIndex)
	assert.Equal(t, 7, *accts[6].AccountIndex)
	assert.Equal(t, 8, *accts[7].AccountIndex)
	assert.Equal(t, 9, *accts[8].AccountIndex)
	assert.Equal(t, 10, *accts[9].AccountIndex)
	assert.Equal(t, "nano_3hntkbk1q6pn1n8481shemojcmtpxxpjbojm7h5h5p6jz53bahjuif6d8j4f", accts[0].Address)
	assert.Equal(t, "nano_1dh7j8bw1pi8xur1zste3c7oqc1ef3yttory4snhufb5haj1ctzxapgpqhzq", accts[1].Address)
	assert.Equal(t, "nano_1qhc3k9p1dxb45f9jrsz3afydqdxoo8y4deeu75q34pmskjgyafq67fd5k1r", accts[2].Address)
	assert.Equal(t, "nano_3cmpq338qbq6hhhijia6ack71kas64snza6r89zhspffzaimxnokaofpa1wt", accts[3].Address)
	assert.Equal(t, "nano_319uakdyoq48h3zswtwt8ei7fjdq4mue3dizzz3fbgnteqb5zqeycygteerd", accts[4].Address)
	assert.Equal(t, "nano_1xew8dmxca5fbo16hgg88p9sgfpamig1kdos5te68hkstj3kf4bi1cyb3qdt", accts[5].Address)
	assert.Equal(t, "nano_33iiz1hkggxirkspip8hef9uiiguhay5ykjdydnyuresxxt1mkkrfz9qeiak", accts[6].Address)
	assert.Equal(t, "nano_3gxjsq4b55fd6hmzroa7wgtmw7wy7nwjofspewzhdwdo5dck9sa85rnfujrn", accts[7].Address)
	assert.Equal(t, "nano_3xoyiewop6gk4g8jaynkf7bhn5baciwtdppkkoprsbafmimcmaiywz6pjqt5", accts[8].Address)
	assert.Equal(t, "nano_1pigk4bbfhpdg3u5b8sdhij8c78n66eap9478p4dhgamyqt7herco57z6ib7", accts[9].Address)

	// Lets break the sequence
	idx := 11
	_, err = MockWallet.AccountCreate(wallet, &idx)
	assert.Nil(t, err)
	idx++
	_, err = MockWallet.AccountCreate(wallet, &idx)
	assert.Nil(t, err)
	idx++
	_, err = MockWallet.AccountCreate(wallet, &idx)
	assert.Nil(t, err)

	// Create 3 more
	accts, err = MockWallet.AccountsCreate(wallet, 3)
	assert.Nil(t, err)
	assert.Len(t, accts, 3)
	assert.Equal(t, 14, *accts[0].AccountIndex)
	assert.Equal(t, 15, *accts[1].AccountIndex)
	assert.Equal(t, 16, *accts[2].AccountIndex)

	// Check that it fails if wallet is locked
	MockWallet.EncryptWallet(wallet, "password")
	_, err = MockWallet.AccountsCreate(wallet, 100)
	assert.ErrorIs(t, ErrWalletLocked, err)
}

func TestAccountsCreateBadInput(t *testing.T) {
	_, err := MockWallet.AccountsCreate(nil, 10)
	assert.ErrorIs(t, ErrInvalidWallet, err)
	_, err = MockWallet.AccountsCreate(&ent.Wallet{}, 0)
	assert.ErrorIs(t, ErrInvalidAccountCount, err)
}

func TestAccountsList(t *testing.T) {
	// Predictable seed
	seed, _ := utils.GenerateSeed(strings.NewReader("1111111111111111abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040e"))

	// +1 This by default creates an account
	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	// Create some accounts and adhoc accounts to make sure they come back in results
	// +5
	_, err = MockWallet.AccountsCreate(wallet, 5)
	assert.Nil(t, err)

	_, priv, _ := utils.KeypairFromSeed(seed, 100)

	// +1
	_, err = MockWallet.AdhocAccountCreate(wallet, priv)

	// List accounts
	accts, err := MockWallet.AccountsList(wallet, 100)
	assert.Nil(t, err)
	// We actually have 7 accounts because WalletCreate implicitly creates the first one
	// So this isn't a mistake
	assert.Len(t, accts, 7)

	assert.Contains(t, accts, "nano_1g7kfyzdqggtmoop8kwiwu4kh7camtaukpc6gkt8mybm4ezx457wjoxyhiwb")
	assert.Contains(t, accts, "nano_1ucbt3tm7y65tbz4pgxmnkce3ycgpiqcp8un7cybc3cxtapeangti3i9apbx")
	assert.Contains(t, accts, "nano_3pws13shsun9b5o1febnzj8bb1nqa99no9fe9t3hm5r6ehi1mznzp31ri11s")
	assert.Contains(t, accts, "nano_1og8yj39i88nc48qseiofnt3raxi7ttsomiz6n6nbo5qmfh954hynzzfhjk7")
	assert.Contains(t, accts, "nano_1xbhmnrt6x1ji4pf6u1m3fckxifxusbmpsp6nuyhhpnfj719ja8z4aubfmud")
	assert.Contains(t, accts, "nano_13687htr3rp5wxfyxpjjzctydd46rca6cpm377ce9eccctpy5ht3zb79b444")
	assert.Contains(t, accts, "nano_1bzpyc67m6hzhm8egshbnyseowohs11d7hkcw4ksz8guetsyegkx3r1ns6s4")

	// Check that it fails if wallet is locked
	MockWallet.EncryptWallet(wallet, "password")
	_, err = MockWallet.AccountsList(wallet, 100)
	assert.ErrorIs(t, ErrWalletLocked, err)
}

func TestAccountExists(t *testing.T) {
	// Create a test wallet
	seed, _ := utils.GenerateSeed(strings.NewReader("cd92b5cf58554077ccd450a09ccbefaea04229c7d9bc54abba06dcaf7ed01af9"))
	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	// Check that account exists if asle
	exists, err := MockWallet.AccountExists(wallet, "nano_33mhuqjxr166czm4y37xk7emfnt4zogxmqrhbfxyngrkbdchmpsk6qehhm3n")
	assert.Nil(t, err)
	assert.False(t, exists)

	// Create an adhoc account
	_, priv, _ := utils.KeypairFromSeed(seed, 100)
	_, err = MockWallet.AdhocAccountCreate(wallet, priv)
	assert.Nil(t, err)

	exists, err = MockWallet.AccountExists(wallet, "nano_1hpq679fqnsahkjz4d66nantwsjbkd1erjbycinbrmhcfnuketzhfaeoptu6")
	assert.Nil(t, err)
	assert.True(t, exists)

	// Create some accounts
	_, err = MockWallet.AccountCreate(wallet, nil)
	exists, err = MockWallet.AccountExists(wallet, "nano_1pidkij46sqyf7gan8fugj693z5ornpf449tikop83dwsuosy1o5164p1jry")
	assert.True(t, exists)
}
