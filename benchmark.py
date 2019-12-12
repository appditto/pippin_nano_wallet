# A simple bench mark for pippin
# 1) Create new wallet
# 2) Create 21 accounts
# 3) Load 1 account with 20 banano
# 4) Send 1 banano to 20 accounts
# 5) Consolidate it back into 1 account
# 6) Return to sender
import argparse
import asyncio
import datetime

from aiohttp import web, TCPConnector, ClientSession, AsyncResolver

parser = argparse.ArgumentParser(description="Pippin Benchmark")
parser.add_argument('--node-url', type=str, help='Node or Pippin RPC location', required=True)
options = parser.parse_args()

loop = asyncio.new_event_loop()

class RPCClient(object):
    _instance = None

    def __init__(self):
        raise RuntimeError('Call instance() instead')

    @classmethod
    def instance(cls) -> 'RPCClient':
        if cls._instance is None:
            cls._instance = cls.__new__(cls)
            cls.node_url = options.node_url
            cls.connector = TCPConnector(family=0 ,resolver=AsyncResolver())
            cls.session = ClientSession(connector=cls.connector)
        return cls._instance


    @classmethod
    async def close(cls):
        if hasattr(cls, 'session') and cls.session is not None:
            await cls.session.close()
        if cls._instance is not None:
            cls._instance = None

    async def make_request(self, req_json: dict):
        async with self.session.post(self.node_url ,json=req_json, timeout=300) as resp:
            return await resp.json()

async def do_send(rpc: RPCClient, source, destination, wallet):
    result = await rpc.make_request({
                'action': 'send',
                'source': source,
                'destination': destination,
                'amount': '100000000000000000000000000000',
                'wallet': wallet
            })
    return (result['block'], destination)

async def run_test():
    rpc = RPCClient.instance()
    # Wallet create
    wallet = (await rpc.make_request({
        'action':'wallet_create'
    }))['wallet']
    # Accounts_create
    accounts = (await rpc.make_request({
        'action':'accounts_create',
        'count': 21,
        'wallet': wallet
    }))['accounts']
    # Wait for user to send 50 BANANO to first account
    confirmed = False
    while not confirmed:
        try:
            print(f"Send 20 BANANO to {accounts[0]}")
            g = input("Type 'done' after it's been sent:")
            if g != 'done':
                print(f"Please type 'done' when you've sent it")
            else:
                confirmed = True
        except KeyboardInterrupt:
            break
    start_time = datetime.datetime.utcnow()
    pending_hash = (await rpc.make_request({
        'action':'pending',
        'account':accounts[0],
        'count':1,
        'include_active': True
    }))['blocks'][0]
    blocks_info = await rpc.make_request({
        'action': 'block_info',
        'json_block': True,
        'hash': pending_hash
    })
    sender = blocks_info['contents']['account']
    # Receive block
    print("receiving")
    received = (await rpc.make_request({
        'action':'receive',
        'wallet': wallet,
        'block': pending_hash,
        'account': accounts[0]
    }))['block']
    send_hashes = {}
    # Send to 20 accounts
    print("Distributing to accounts")
    tasks = []
    for a in accounts:
        if a == accounts[0]:
            continue
        tasks.append(do_send(rpc, accounts[0], a, wallet))
    ret = await asyncio.gather(*tasks)
    for r in ret:
        send_hashes[r[1]] = r[0]
    print("Receiving on all accounts")
    # Receive all blocks and send them back
    tasks = []
    for a, h in send_hashes.items():
        tasks.append(rpc.make_request({
            'action': 'receive',
            'wallet': wallet,
            'block': h,
            'account': a
        }))
    await asyncio.gather(*tasks)
    print(f"Conolidating back to {accounts[0]}")
    tasks = []
    for a, h in send_hashes.items():
        tasks.append(rpc.make_request({
            'action': 'send',
            'source': a,
            'destination': accounts[0],
            'amount': '100000000000000000000000000000',
            'wallet': wallet
        }))
    await asyncio.gather(*tasks)
    await asyncio.sleep(5)
    # Receive all pendings on account[0]
    pendings = (await rpc.make_request({
        'action':'pending',
        'account':accounts[0],
        'count':50,
        'include_active': True
    }))['blocks']
    tasks = []
    for p in pendings:
        tasks.append(rpc.make_request({
            'action': 'receive',
            'wallet': wallet,
            'block': p,
            'account': accounts[0]
        }))
    await asyncio.gather(*tasks)
    # Return to sender
    print(f"Returning to {sender}")
    send = (await rpc.make_request({
        'action': 'send',
        'source': accounts[0],
        'destination': sender,
        'amount': '2000000000000000000000000000000',
        'wallet': wallet
    }))['block']
    final_time = (datetime.datetime.utcnow() - start_time).total_seconds() - 5
    print(f"Finished in {final_time} seconds")

if __name__ == "__main__":
    loop.run_until_complete(run_test())
    loop.run_until_complete(RPCClient.close())
    loop.close()