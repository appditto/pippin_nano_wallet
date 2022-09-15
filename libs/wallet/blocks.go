package wallet

import (
	"errors"
	"math/big"

	"github.com/appditto/pippin_nano_wallet/libs/database/ent"
	entblock "github.com/appditto/pippin_nano_wallet/libs/database/ent/block"
	"github.com/appditto/pippin_nano_wallet/libs/rpc/models/requests"
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
	acc, adhoc, err := w.GetAccount(wallet, address)
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
	} else {
		block, err = w.DB.Block.Query().Where(entblock.AdhocAccountID(adhoc.ID), entblock.SendID(sendID)).First(w.Ctx)
		if ent.IsNotFound(err) {
			return nil, ErrBlockNotFound
		} else if err != nil {
			return nil, err
		}
	}

	return block, nil
}

func (w *NanoWallet) CreateAndPublishSendBlock(wallet *ent.Wallet, amount big.Int, source string, destination string, id string, work string, bpowKy string) (string, error) {
	_, _, err := w.GetAccount(wallet, source)
	if err != nil {
		return "", err
	}

	// This is our indempotent send test, we don't create a new send block if a send with this ID has already been created from this account
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
