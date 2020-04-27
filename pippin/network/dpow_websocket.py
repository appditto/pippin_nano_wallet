from aiohttp import log
import asyncio
import websockets
import rapidjson as json
import traceback
import nanopy
import nanolib

class ConnectionClosed(Exception):
    pass

class DpowClient(object):
    """Websocket client for DPoW & BoomPoW"""
    DPOW_SERVER = 'wss://dpow.nanocenter.org/service_ws/'
    BPOW_SERVER = 'wss://bpow.banano.cc/service_ws/'

    def __init__(self, dpow_user: str, dpow_key: str, work_futures: dict = {}, bpow: bool = False):
        """work_futures is a dict of id, futures to keep track of requests and responses"""
        self.bpow = bpow
        self.url = self.BPOW_SERVER if self.bpow else self.DPOW_SERVER
        self.ws: websockets.WebSocketClientProtocol = None
        self.stop = False
        self.dpow_user = dpow_user
        self.dpow_key = dpow_key
        self.work_futures = work_futures

    def adjust_difficulty(self, difficulty: str) -> str:
        """Ensure sane difficulty limits for DPoW/BPoW"""
        if nanolib.work.derive_work_multiplier(difficulty, base_difficulty=nanopy.work_difficulty) > 8:
            return nanolib.work.derive_work_difficulty(8, base_difficulty=nanopy.work_difficulty)
        return difficulty

    async def request_work(self, id: str, hash: str, difficulty: str):
        """Request work from DPoW/BPoW WS"""
        if self.stop or self.ws is None:
            return
        try:
            if self.ws.closed:
                raise ConnectionClosed()
            req = {
                "user": self.dpow_user,
                "api_key": self.dpow_key,
                "hash": hash,
                "id": id,
                "timeout": 15,
                "difficulty": self.adjust_difficulty(difficulty)
            }
            await self.ws.send(json.dumps(req))
        except Exception:
            raise ConnectionClosed()        

    async def setup(self):
        try:
            self.ws = await websockets.connect(self.url)
            log.server_logger.info(f"Connected to {'BoomPoW' if self.bpow else 'Distributed PoW'} Service")
        except Exception as e:
            log.server_logger.critical("DPOW WS: Error connecting to websocket server.")
            log.server_logger.error(traceback.format_exc())
            raise

    async def close(self):
        self.stop = True
        await self.ws.wait_closed()

    async def reconnect_forever(self):
        log.server_logger.warn("DPOW WS: Attempting websocket reconnection every 30 seconds...")
        while not self.stop:
            try:
                await self.setup()
                log.server_logger.warn("DPOW WS: Connected to websocket!")
                break
            except:
                log.server_logger.debug("DPOW WS: Websocket reconnection failed")
                await asyncio.sleep(30)

    async def loop(self):
        while not self.stop:
            try:
                rec = json.loads(await self.ws.recv())
                request_id = rec.get("id", None)
                if request_id is not None:
                    try:
                        result = self.work_futures[str(request_id)]
                        if not result.done():
                            result.set_result(rec)
                    except KeyError:
                        pass                    
            except KeyboardInterrupt:
                break
            except websockets.exceptions.ConnectionClosed as e:
                log.server_logger.error(f"DPOW WS: Connection closed to websocket. Code: {e.code} , reason: {e.reason}.")
                await self.reconnect_forever()
            except Exception as e:
                log.server_logger.critical(f"DPOW WS: Unknown exception while handling getting a websocket message:\n{traceback.format_exc()}")
                await self.reconnect_forever()
