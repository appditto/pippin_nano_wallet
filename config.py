from dotenv import load_dotenv
load_dotenv()

import argparse
from util.env import Env
from version import __version__

class Config(object):
    _instance = None

    def __init__(self):
        raise RuntimeError('Call instance() instead')

    @classmethod
    def instance(cls) -> 'Config':
        if cls._instance is None:
            cls._instance = cls.__new__(cls)
            parser = argparse.ArgumentParser(description=f"Pippin {'BANANO' if Env.banano() else 'Nano'} Wallet API v{__version__}")
            parser.add_argument('-l', '--log-file', type=str, help='Log file location', default='/tmp/pippin_wallet.log')
            parser.add_argument('-p', '--port', type=int, help='Port to listen on', default=11338)
            parser.add_argument('-u', '--node-url', type=str, help='URL of the node', default='[::1]')
            parser.add_argument('-np', '--node-port', type=int, help='Port of the node', default=7072 if Env.banano() else 7076)
            parser.add_argument('--debug', action='store_true', help='Runs in debug mode if specified', default=False)

            options, unknown = parser.parse_known_args()

            # Parse options
            cls.log_file = options.log_file
            cls.debug = options.debug
            cls.node_url = options.node_url
            cls.node_port = options.node_port
            cls.port = options.port
            cls.host = '127.0.0.1'
        return cls._instance