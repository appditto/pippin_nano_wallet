package wallet

import (
	"encoding/hex"
	"errors"
	"math/big"

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
func (w *NanoWallet) createReceiveBlock(wallet *ent.Wallet, receiver *ent.Account, hash string, precomputedWork *string) (*models.StateBlock, error) {
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
		pub, err := utils.AddressToPub(receiver.Address)
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
	curB, ok := big.NewInt(0).SetString(blockInfo.Amount, 10)
	if !ok {
		return nil, errors.New("Unable to parse balance")
	}
	if !isOpen {
		balance = curB
	} else {
		receiveAmount, ok := big.NewInt(0).SetString(blockInfo.Amount, 10)
		if !ok {
			return nil, errors.New("Unable to block info amount")
		}
		balance = big.NewInt(0).Add(receiveAmount, curB)
	}

	var work string
	if precomputedWork == nil {
		work, err = w.WorkClient.WorkGenerateMeta(workbase, 1, true, false, "")
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

func (w *NanoWallet) CreateAndPublishSendBlock(wallet *ent.Wallet, amount big.Int, source string, destination string, id string, work string, bpowKy string) (string, error) {
	_, err := w.GetAccount(wallet, source)
	if err != nil {
		return "", err
	}

	// This is our idempotent send test, we don't create a new send block if a send with this ID has already been created from this account
	if id != "" {
		block, err := w.GetBlockFromDatabase(wallet, source, id)
		if !ent.IsNotFound(err) && err != nil {
			return "", err
		} else if block != nil {
			// Now we can just republish...
			var sb models.StateBlock
			if err := mapstructure.Decode(block.Block, &sb); err != nil {
				return "", err
			}
			subtype := "send"
			resp, err := w.RpcClient.MakeProcessRequest(requests.ProcessRequest{
				BaseRequest: requests.BaseRequest{
					Action: "process",
				},
				Subtype:   &subtype,
				JsonBlock: true,
				Block:     sb,
			})
			if err != nil {
				return "", err
			}
			return resp.Hash, nil
		}
	}
	// # Get account info
	// account_balance = await RPCClient.instance().account_balance(self.account.address)
	// if account_balance is None:
	// 		return None

	// # Check balance
	// if amount > int(account_balance['balance']):
	// 		# Auto-receive blocks if they have it pending
	// 		if config.Config.instance().auto_receive_on_send and int(account_balance['balance']) + int(account_balance['pending']) >= amount:
	// 				await self._receive_all()
	// 				account_info = await RPCClient.instance().account_info(self.account.address)
	// 				if account_info is None:
	// 						return None
	// 				if amount > int(account_info['balance']):
	// 						raise InsufficientBalance(account_info['balance'])
	// 		else:
	// 				raise InsufficientBalance(account_balance['balance'])

	return "", nil
}

// async def send(self, amount: int, destination: str, id: str = None, work: str = None, bpow_key: str = None) -> dict:
// """Create a send block and return hash of published block
// 		amount is in RAW"""

// async with await (await RedisDB.instance().get_lock_manager()).lock(f"pippin:{self.account.address}") as lock:
// 		# See if block exists, if ID specified
// 		# If so just rebroadcast it and return the hash
// 		if id is not None:
// 				if not self.adhoc():
// 						block = await Block.filter(send_id=id, account=self.account).first()
// 				else:
// 						block = await Block.filter(send_id=id, adhoc_account=self.account).first()
// 				if block is not None:
// 						await RPCClient.instance().process(block.block)
// 						return {"block": block.block_hash.upper()}
// 		# Create block
// 		state_block = await self._send_block_create(amount, destination, id=id, work=work, bpow_key=bpow_key)
// 		# Publish block
// 		resp = await self.publish(state_block, subtype='send')
// 		# Cache if ID specified
// 		if resp is not None and 'block' in resp:
// 				# Cache block in database if it has id specified
// 				if id is not None:
// 						block = Block(
// 								account=self.account if not self.adhoc() else None,
// 								adhoc_account=self.account if self.adhoc() else None,
// 								block_hash=nanopy.block_hash(state_block),
// 								block=state_block,
// 								send_id=id,
// 								subtype='send'
// 						)
// 						await block.save()
// 		return resp
