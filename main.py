# Install uvloop
try:
	import uvloop
	uvloop.install()
except ImportError:
	print("Couldn't install uvloop, falling back to the slower asyncio event loop")

import asyncio
import logging

from aiohttp import web, log
from db.redis import RedisDB
from db.tortoise_config import DBConfig
from config import Config
from logging.handlers import TimedRotatingFileHandler, WatchedFileHandler
from network.rpc_client import RPCClient
from network.work_client import WorkClient
from server.pippin_server import PippinServer
from util.env import Env

# Set and patch nanopy
import nanopy
nanopy.account_prefix = 'ban_' if Env.banano() else 'nano_'
if Env.banano():
    nanopy.standard_exponent = 29
    nanopy.work_difficulty = 'fffffe0000000000'

# Configuration
config = Config.instance()

# Setup logger
if config.debug:
    logging.basicConfig(level=logging.DEBUG)
else:
    root = logging.getLogger('aiohttp.server')
    logging.basicConfig(level=logging.INFO)
    handler = WatchedFileHandler(config.log_file)
    formatter = logging.Formatter("%(asctime)s;%(levelname)s;%(message)s", "%Y-%m-%d %H:%M:%S %z")
    handler.setFormatter(formatter)
    root.addHandler(handler)
    root.addHandler(TimedRotatingFileHandler(config.log_file, when="d", interval=1, backupCount=100))  

if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    try:
        # Initialize database first
        log.server_logger.info("Initializing database")
        loop.run_until_complete(DBConfig().init_db())
        # Setup server
        server = PippinServer(config.host, config.port)
        # Check is remote node is alive
        try:
            is_alive = loop.run_until_complete(RPCClient.instance().is_alive())
        except Exception:
            log.server_logger.exception("Couldn't do is_alive RPC call")
            is_alive = False
        finally:
            if not is_alive:
                log.server_logger.error(f"Error: Could not connect to remote node at {Config.instance().node_url}")
                exit(1)
        # Start server
        log.server_logger.info(f"Pippin server started at {config.host}:{config.port}")
        loop.run_until_complete(server.start())
        loop.run_forever()
    except Exception:
        log.server_logger.exception("Pippin exited with exception")
    except BaseException:
        pass
    finally:
        log.server_logger.info("Pipping is exiting")
        tasks = [
            RPCClient.close(),
            server.stop(),
            RedisDB.close(),
            WorkClient.close()
        ]
        loop.run_until_complete(asyncio.wait(tasks))
        loop.close()