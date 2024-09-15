package wallet

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/appditto/pippin_nano_wallet/libs/database"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent"
	entblock "github.com/appditto/pippin_nano_wallet/libs/database/ent/block"
	nanorpc "github.com/appditto/pippin_nano_wallet/libs/rpc"
	"github.com/appditto/pippin_nano_wallet/libs/rpc/models/requests"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"
	"github.com/appditto/pippin_nano_wallet/libs/wallet/models"
	"github.com/mitchellh/mapstructure"
)

var ErrBlockNotFound = errors.New("block not found")
var ErrInsufficientBalance = errors.New("insufficient balance")
var ErrSameRepresentative = errors.New("same representative")

// The core that creates and publishes send, receive, and change blocks
// See: https://docs.nano.org/protocol-design/blocks/

// Retrieve a block for an account or adhoc account, error if it doesn't exist
// Intended to retrieve a block with a particular send ID, to see if it's already been created
func (w *NanoWallet) GetBlockFromDatabase(wallet *ent.Wallet, address string, sendID string) (*ent.Block, error) {
	if wallet == nil {
		return nil, ErrInvalidWallet
	}

	// Determine if wallet is locked or not
	_, err := GetDecryptedKeyFromStorage(wallet, "seed")
	if err != nil {
		return nil, err
	}

	// Check if account exists
	acc, err := w.GetAccount(wallet, address)
	if err != nil {
		return nil, err
	}

	// Get block
	var block *ent.Block
	if acc != nil {
		block, err = w.DB.Block.Query().Where(entblock.AccountID(acc.ID), entblock.SendID(sendID)).First(w.Ctx)
		if ent.IsNotFound(err) {
			return nil, ErrBlockNotFound
		} else if err != nil {
			return nil, err
		}
	}

	return block, nil
}

// ** Low level block creations, not intended for use by the user **
func (w *NanoWallet) createReceiveBlock(wallet *ent.Wallet, receiver *ent.Account, hash string, precomputedWork *string, bpowKey *string) (*models.StateBlock, error) {
	if wallet == nil {
		return nil, ErrInvalidWallet
	} else if receiver == nil {
		return nil, ErrInvalidAccount
	}
	blockInfo, err := w.RpcClient.MakeBlockInfoRequest(hash)
	if err != nil {
		return nil, err
	} else if blockInfo == nil {
		return nil, ErrBlockNotFound
	}
	// Get account info
	isOpen := true
	accountInfo, err := w.RpcClient.MakeAccountInfoRequest(receiver.Address)
	if errors.Is(err, nanorpc.ErrAccountNotFound) {
		isOpen = false
	} else if err != nil {
		return nil, err
	}

	var workbase string
	if isOpen {
		workbase = accountInfo.Frontier
	} else {
		pub, err := utils.AddressToPub(receiver.Address, w.Config.Wallet.Banano)
		if err != nil {
			return nil, err
		}
		workbase = hex.EncodeToString(pub)
	}

	// Build other block fields
	var previous string
	if isOpen {
		previous = accountInfo.Frontier
	} else {
		previous = "0000000000000000000000000000000000000000000000000000000000000000"
	}

	var representative string

	if isOpen {
		representative = accountInfo.Representative
	} else {
		if wallet.Representative != nil {
			representative = *wallet.Representative
		} else {
			rep, err := w.Config.GetRandomRep()
			if err != nil {
				return nil, err
			}
			representative = rep
		}
	}

	var balance *big.Int
	receiveAmount, ok := big.NewInt(0).SetString(blockInfo.Amount, 10)
	if !ok {
		return nil, errors.New("Unable to parse balance")
	}
	if !isOpen {
		balance = receiveAmount
	} else {
		currentBalance, ok := big.NewInt(0).SetString(accountInfo.Balance, 10)
		if !ok {
			return nil, errors.New("Unable to parse balance")
		}
		balance = big.NewInt(0).Add(receiveAmount, currentBalance)
	}

	var work string
	if precomputedWork == nil {
		key := ""
		if bpowKey != nil {
			key = *bpowKey
		}
		work, err = w.WorkClient.WorkGenerateMeta(workbase, 1, true, false, key)
		if err != nil {
			return nil, err
		}
	} else {
		work = *precomputedWork
	}

	stateBlock := &models.StateBlock{
		Type:           "state",
		Account:        receiver.Address,
		Previous:       previous,
		Representative: representative,
		Balance:        balance.String(),
		Link:           hash,
		Work:           work,
		Banano:         w.Config.Wallet.Banano,
	}

	// Get the private key for this account
	var priv ed25519.PrivateKey
	if receiver.PrivateKey != nil {
		decoded, err := hex.DecodeString(*receiver.PrivateKey)
		if err != nil {
			return nil, err
		}
		priv = ed25519.PrivateKey(decoded)
	} else {
		sd, err := GetDecryptedKeyFromStorage(wallet, "seed")
		if err != nil {
			return nil, err
		}
		_, priv, _ = utils.KeypairFromSeed(sd, uint32(*receiver.AccountIndex))
	}

	// Sign the block
	err = stateBlock.Sign(priv)
	if err != nil {
		return nil, err
	}

	return stateBlock, nil
}

// Receive all without locking the wallet
func (w *NanoWallet) receiveAll(wallet *ent.Wallet, acc *ent.Account, bpowKey *string) (int, error) {
	if wallet == nil {
		return 0, ErrInvalidWallet
	} else if acc == nil {
		return 0, ErrInvalidAccount
	}
	receivedCount := 0
	// Get pending
	pending, err := w.RpcClient.MakeReceivableRequest(acc.Address, w.Config.Wallet.ReceiveMinimum)
	if err != nil {
		return receivedCount, err
	}
	if len(pending.Blocks) == 0 {
		return receivedCount, nil
	}

	// Create and publish blocks
	for hash := range pending.Blocks {
		sb, err := w.createReceiveBlock(wallet, acc, hash, nil, bpowKey)
		if err != nil {
			return receivedCount, err
		}

		// Publish block
		subtype := "receive"
		resp, err := w.RpcClient.MakeProcessRequest(requests.ProcessRequest{
			BaseRequest: requests.BaseRequest{
				Action: "process",
			},
			Subtype:   &subtype,
			JsonBlock: true,
			Block:     *sb,
		})
		if err != nil || !utils.Validate64HexHash(resp.Hash) {
			return receivedCount, err
		}
		receivedCount++
	}
	return receivedCount, nil
}

func (w *NanoWallet) createSendBlock(wallet *ent.Wallet, sender *ent.Account, amount string, destination string, precomputedWork *string, bpowKey *string) (*models.StateBlock, error) {
	if wallet == nil {
		return nil, ErrInvalidWallet
	} else if sender == nil {
		return nil, ErrInvalidAccount
	}

	sendAmount, ok := big.NewInt(0).SetString(amount, 10)
	if !ok {
		return nil, errors.New("Unable to parse send amount")
	}

	// Get account info
	accountInfo, err := w.RpcClient.MakeAccountInfoRequest(sender.Address)
	if errors.Is(err, nanorpc.ErrAccountNotFound) {
		if w.Config.Wallet.AutoReceiveOnSend == nil || !*w.Config.Wallet.AutoReceiveOnSend {
			return nil, ErrInsufficientBalance
		}
		// See if account has a pending balance to open the accountt
		bal, err := w.RpcClient.MakeAccountBalanceRequest(sender.Address)
		if err != nil {
			return nil, err
		}
		receivable, ok := big.NewInt(0).SetString(bal.Receivable, 10)
		if !ok {
			return nil, errors.New("Unable to parse receivable amount")
		}
		if receivable.Cmp(sendAmount) < 0 {
			return nil, ErrInsufficientBalance
		}
		receivedCount, err := w.receiveAll(wallet, sender, bpowKey)
		if err != nil {
			return nil, err
		}
		if receivedCount == 0 {
			return nil, ErrInsufficientBalance
		}
		// Re-get accountInfo
		accountInfo, err = w.RpcClient.MakeAccountInfoRequest(sender.Address)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// Convert balance to big int
	balanceBigInt, ok := big.NewInt(0).SetString(accountInfo.Balance, 10)
	if !ok {
		return nil, errors.New("Unable to parse balance")
	}

	// Check if balance is sufficient
	if sendAmount.Cmp(balanceBigInt) > 0 {
		if w.Config.Wallet.AutoReceiveOnSend == nil || !*w.Config.Wallet.AutoReceiveOnSend {
			return nil, ErrInsufficientBalance
		}
		// Automatically receive blocks to see if we can make up the difference
		receivedCount, _ := w.receiveAll(wallet, sender, bpowKey)
		if receivedCount > 0 {
			// Re-check balance
			accountInfo, err = w.RpcClient.MakeAccountInfoRequest(sender.Address)
			if err != nil {
				return nil, err
			}
			balanceBigInt, ok = big.NewInt(0).SetString(accountInfo.Balance, 10)
			if !ok {
				return nil, errors.New("Unable to parse balance")
			}
			if sendAmount.Cmp(balanceBigInt) > 0 {
				return nil, ErrInsufficientBalance
			}
		} else {
			return nil, ErrInsufficientBalance
		}
	}

	workbase := accountInfo.Frontier

	// Build other block fields
	previous := accountInfo.Frontier

	var representative string
	if wallet.Representative != nil {
		representative = *wallet.Representative
	} else {
		rep, err := w.Config.GetRandomRep()
		if err != nil {
			return nil, err
		}
		representative = rep
	}

	// Calculate new balance, subtracing sendAmount from balanceBigInt
	newBalance := balanceBigInt.Sub(balanceBigInt, sendAmount)

	var work string
	if precomputedWork == nil {
		key := ""
		if bpowKey != nil {
			key = *bpowKey
		}
		difficulty := 1
		if !w.Config.Wallet.Banano {
			difficulty = 64
		}
		work, err = w.WorkClient.WorkGenerateMeta(workbase, difficulty, true, false, key)
		if err != nil {
			return nil, err
		}
	} else {
		work = *precomputedWork
	}

	// Link is pubkey of destination
	link, err := utils.AddressToPub(destination, w.Config.Wallet.Banano)
	if err != nil {
		return nil, errors.New("Invalid destination address")
	}

	stateBlock := &models.StateBlock{
		Type:           "state",
		Account:        sender.Address,
		Previous:       previous,
		Representative: representative,
		Balance:        newBalance.String(),
		Link:           hex.EncodeToString(link),
		Work:           work,
		Banano:         w.Config.Wallet.Banano,
	}

	// Get the private key for this account
	var priv ed25519.PrivateKey
	if sender.PrivateKey != nil {
		decoded, err := hex.DecodeString(*sender.PrivateKey)
		if err != nil {
			return nil, err
		}
		priv = ed25519.PrivateKey(decoded)
	} else {
		sd, err := GetDecryptedKeyFromStorage(wallet, "seed")
		if err != nil {
			return nil, err
		}
		_, priv, _ = utils.KeypairFromSeed(sd, uint32(*sender.AccountIndex))
	}

	// Sign the block
	err = stateBlock.Sign(priv)
	if err != nil {
		return nil, err
	}

	return stateBlock, nil
}

func (w *NanoWallet) createChangeBlock(wallet *ent.Wallet, changer *ent.Account, representative string, precomputedWork *string, bpowKey *string, onlyIfDifferent bool) (*models.StateBlock, error) {
	if wallet == nil {
		return nil, ErrInvalidWallet
	} else if changer == nil {
		return nil, ErrInvalidAccount
	}

	// Get account info
	accountInfo, err := w.RpcClient.MakeAccountInfoRequest(changer.Address)
	if err != nil {
		return nil, err
	}

	if onlyIfDifferent && accountInfo.Representative == representative {
		return nil, ErrSameRepresentative
	}

	workbase := accountInfo.Frontier

	// Build other block fields
	previous := accountInfo.Frontier

	var work string
	if precomputedWork == nil {
		key := ""
		if bpowKey != nil {
			key = *bpowKey
		}
		difficulty := 1
		if !w.Config.Wallet.Banano {
			difficulty = 64
		}
		work, err = w.WorkClient.WorkGenerateMeta(workbase, difficulty, true, false, key)
		if err != nil {
			return nil, err
		}
	} else {
		work = *precomputedWork
	}

	stateBlock := &models.StateBlock{
		Type:           "state",
		Account:        changer.Address,
		Previous:       previous,
		Representative: representative,
		Balance:        accountInfo.Balance,
		Link:           "0000000000000000000000000000000000000000000000000000000000000000",
		Work:           work,
		Banano:         w.Config.Wallet.Banano,
	}

	// Get the private key for this account
	var priv ed25519.PrivateKey
	if changer.PrivateKey != nil {
		decoded, err := hex.DecodeString(*changer.PrivateKey)
		if err != nil {
			return nil, err
		}
		priv = ed25519.PrivateKey(decoded)
	} else {
		sd, err := GetDecryptedKeyFromStorage(wallet, "seed")
		if err != nil {
			return nil, err
		}
		_, priv, _ = utils.KeypairFromSeed(sd, uint32(*changer.AccountIndex))
	}

	// Sign the block
	err = stateBlock.Sign(priv)
	if err != nil {
		return nil, err
	}

	return stateBlock, nil
}

// The user facing APIs intended to be  for block creation/publishing
// They are done in a locked context

// Receive single block
func (w *NanoWallet) CreateAndPublishReceiveBlock(wallet *ent.Wallet, source string, hash string, work *string, bpowKey *string) (string, error) {
	if wallet == nil {
		return "", ErrInvalidWallet
	}

	acc, err := w.GetAccount(wallet, source)
	if err != nil {
		return "", err
	}

	// Obtain lock
	lock, err := database.GetRedisDB().Locker.Obtain(w.Ctx, fmt.Sprintf("acl:%s", acc.Address), time.Second*30, &database.LockRetryStrategy)
	if err != nil {
		return "", database.ErrLockNotObtained
	}
	defer lock.Release(w.Ctx)

	sb, err := w.createReceiveBlock(wallet, acc, hash, work, bpowKey)
	if err != nil {
		return "", err
	}

	// Publish block
	subtype := "receive"
	resp, err := w.RpcClient.MakeProcessRequest(requests.ProcessRequest{
		BaseRequest: requests.BaseRequest{
			Action: "process",
		},
		Subtype:   &subtype,
		JsonBlock: true,
		Block:     *sb,
	})
	if err != nil || !utils.Validate64HexHash(resp.Hash) {
		return "", err
	}
	return resp.Hash, nil
}

// Receive all blocks in all accounts on wallet, respecting receive minimum
func (w *NanoWallet) ReceiveAllBlocks(wallet *ent.Wallet, source string, bpowKey *string) (int, error) {
	if wallet == nil {
		return 0, ErrInvalidWallet
	}

	acc, err := w.GetAccount(wallet, source)
	if err != nil {
		return 0, err
	}

	// Obtain lock
	// Longer lock since this culd be long running
	lock, err := database.GetRedisDB().Locker.Obtain(w.Ctx, fmt.Sprintf("acl:%s", acc.Address), time.Second*300, &database.LockRetryStrategy)
	if err != nil {
		return 0, database.ErrLockNotObtained
	}
	defer lock.Release(w.Ctx)

	return w.receiveAll(wallet, acc, bpowKey)
}

func (w *NanoWallet) CreateAndPublishSendBlock(wallet *ent.Wallet, amount string, source string, destination string, id *string, work *string, bpowKey *string) (string, error) {
	if wallet == nil {
		return "", ErrInvalidWallet
	}
	acc, err := w.GetAccount(wallet, source)
	if err != nil {
		return "", err
	}

	// Obtain lock
	lock, err := database.GetRedisDB().Locker.Obtain(w.Ctx, fmt.Sprintf("acl:%s", acc.Address), time.Second*30, &database.LockRetryStrategy)
	if err != nil {
		return "", database.ErrLockNotObtained
	}
	defer lock.Release(w.Ctx)

	// This is our idempotent send test, we don't create a new send block if a send with this ID has already been created from this account
	if id != nil {
		block, err := w.GetBlockFromDatabase(wallet, source, *id)
		if !errors.Is(err, ErrBlockNotFound) && err != nil {
			return "", err
		} else if block != nil {
			// Now we can just republish...
			var sb models.StateBlock
			if err := mapstructure.Decode(block.Block, &sb); err != nil {
				return "", err
			}
			subtype := "send"
			// Call process with same block
			w.RpcClient.MakeProcessRequest(requests.ProcessRequest{
				BaseRequest: requests.BaseRequest{
					Action: "process",
				},
				Subtype:   &subtype,
				JsonBlock: true,
				Block:     sb,
			})
			return strings.ToUpper(sb.Hash), nil
		}
	}

	sb, err := w.createSendBlock(wallet, acc, amount, destination, work, bpowKey)
	if err != nil {
		return "", err
	}

	// Publish block
	subtype := "send"
	resp, err := w.RpcClient.MakeProcessRequest(requests.ProcessRequest{
		BaseRequest: requests.BaseRequest{
			Action: "process",
		},
		Subtype:   &subtype,
		JsonBlock: true,
		Block:     *sb,
	})
	if err != nil || !utils.Validate64HexHash(resp.Hash) {
		return "", err
	}

	// If the ID is set save it in database for indempotency
	if id != nil {
		var asInterface map[string]interface{}
		inrec, _ := json.Marshal(sb)
		json.Unmarshal(inrec, &asInterface)
		_, err := w.DB.Block.Create().SetAccount(acc).SetBlock(asInterface).SetBlockHash(resp.Hash).SetSubtype("send").SetSendID(*id).Save(w.Ctx)
		if err != nil {
			return "", err
		}
	}

	return resp.Hash, nil
}

func (w *NanoWallet) CreateAndPublishChangeBlock(wallet *ent.Wallet, address string, representative string, work *string, bpowKey *string, onlyIfDifferent bool) (string, error) {
	if wallet == nil {
		return "", ErrInvalidWallet
	}
	acc, err := w.GetAccount(wallet, address)
	if err != nil {
		return "", err
	}

	// Obtain lock
	lock, err := database.GetRedisDB().Locker.Obtain(w.Ctx, fmt.Sprintf("acl:%s", acc.Address), time.Second*30, &database.LockRetryStrategy)
	if err != nil {
		return "", database.ErrLockNotObtained
	}
	defer lock.Release(w.Ctx)

	sb, err := w.createChangeBlock(wallet, acc, representative, work, bpowKey, onlyIfDifferent)
	if err != nil {
		return "", err
	}

	// Publish block
	subtype := "change"
	resp, err := w.RpcClient.MakeProcessRequest(requests.ProcessRequest{
		BaseRequest: requests.BaseRequest{
			Action: "process",
		},
		Subtype:   &subtype,
		JsonBlock: true,
		Block:     *sb,
	})
	if err != nil || !utils.Validate64HexHash(resp.Hash) {
		return "", err
	}

	return resp.Hash, nil
}
