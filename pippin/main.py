# Install uvloop
import sys
try:
    if sys.platform not in ('win32', 'cygwin', 'cli'):
        import uvloop
        uvloop.install()
except ImportError:
	print("Couldn't install uvloop, falling back to the slower asyncio event loop")

import argparse
import asyncio
import pathlib
import os
import shutil
import logging

from aiohttp import web, log
from pippin.db.redis import RedisDB
from tortoise import Tortoise
from pippin.db.tortoise_config import DBConfig
from pippin.config import Config
from logging.handlers import TimedRotatingFileHandler, WatchedFileHandler
from pippin.network.rpc_client import RPCClient
from pippin.network.work_client import WorkClient
from pippin.server.pippin_server import PippinServer
from pippin.util.nano_util import NanoUtil
from pippin.util.utils import Utils

parser = argparse.ArgumentParser(description="Pippin Server")
parser.add_argument('--generate-config', action='store_true', help='Generate sample configuration file and exit', default=False)
options = parser.parse_args()

# Create sample file if not exists
config_dir = Utils.get_project_root()
sample_file = config_dir.joinpath(pathlib.PurePath('sample.config.yaml'))
real_file = config_dir.joinpath(pathlib.PurePath('config.yaml'))
if (not os.path.isfile(sample_file) and not os.path.isfile(real_file)) or options.generate_config:
    ref_file = pathlib.Path(__file__).parent.joinpath(pathlib.PurePath('sample.config.yaml'))
    shutil.copyfile(ref_file, sample_file)
    print(f"Sample configuration created at: {sample_file}")
    if options.generate_config:
        exit(0)

# Configuration
config = Config.instance()

# Set and patch nanopy
import nanopy
nanopy.account_prefix = 'ban_' if config.banano else 'nano_'
if config.banano:
    nanopy.standard_exponent = 29
    nanopy.work_difficulty = 'fffffe0000000000'

# Setup logger
if config.debug:
    logging.basicConfig(level=logging.DEBUG)
elif config.stdout:
    logging.basicConfig(level=logging.INFO)
else:
    root = logging.getLogger('aiohttp.server')
    logging.basicConfig(level=logging.INFO)
    handler = WatchedFileHandler(config.log_file)
    formatter = logging.Formatter("%(asctime)s;%(levelname)s;%(message)s", "%Y-%m-%d %H:%M:%S %z")
    handler.setFormatter(formatter)
    root.addHandler(handler)
    root.addHandler(TimedRotatingFileHandler(config.log_file, when="d", interval=1, backupCount=100))  

def main():
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
        log.server_logger.info(f"Pippin server starting at {config.host}:{config.port}")
        tasks = [
            server.start()
        ]
        # Check if DPoW should be started
        dpow_client = WorkClient.instance().dpow_client
        if dpow_client is not None:
            loop.run_until_complete(dpow_client.setup())
            tasks.append(dpow_client.loop())
        loop.run_until_complete(asyncio.wait(tasks))
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
            WorkClient.close(),
            NanoUtil.close(),
            Tortoise.close_connections()
        ]
        loop.run_until_complete(asyncio.wait(tasks))
        loop.close()

if __name__ == "__main__":
    main()