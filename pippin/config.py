
import pathlib
from dotenv import load_dotenv
load_dotenv()
from pippin.util.utils import Utils
load_dotenv(dotenv_path=Utils.get_project_root().joinpath(pathlib.PurePath('.env')))

from aiohttp import log
from pippin.version import __version__
from pippin.util.utils import Utils
from pippin.util.validators import Validators

import secrets
import yaml

DEFAULT_BANANO_REPS = [
    'ban_1ka1ium4pfue3uxtntqsrib8mumxgazsjf58gidh1xeo5te3whsq8z476goo',
    'ban_1cake36ua5aqcq1c5i3dg7k8xtosw7r9r7qbbf5j15sk75csp9okesz87nfn',
    'ban_1fomoz167m7o38gw4rzt7hz67oq6itejpt4yocrfywujbpatd711cjew8gjj'
]

DEFAULT_NANO_REPS = [
    'nano_1x7biz69cem95oo7gxkrw6kzhfywq4x5dupw4z1bdzkb74dk9kpxwzjbdhhs',
    'nano_1thingspmippfngcrtk1ofd3uwftffnu4qu9xkauo9zkiuep6iknzci3jxa6',
    'nano_1natrium1o3z5519ifou7xii8crpxpk8y65qmkih8e8bpsjri651oza8imdd',
    'nano_3o7uzba8b9e1wqu5ziwpruteyrs3scyqr761x7ke6w1xctohxfh5du75qgaj'
]

class Config(object):
    _instance = None

    def __init__(self):
        raise RuntimeError('Call instance() instead')

    @classmethod
    def instance(cls) -> 'Config':
        if cls._instance is None:
            cls._instance = cls.__new__(cls)
            try:
                with open(f"{Utils.get_project_root().joinpath(pathlib.PurePath('config.yaml'))}", "r") as in_yaml:
                    cls.yaml = list(yaml.load_all(in_yaml, Loader=yaml.FullLoader))[0]
            except FileNotFoundError:
                cls.yaml = None
            # Parse options
            cls.banano = cls.get_yaml_property('wallet', 'banano', False)
            cls.log_file = cls.get_yaml_property('server', 'log_file', default='/tmp/pippin_wallet.log')
            cls.debug = cls.get_yaml_property('server', 'debug', default=False)
            cls.stdout = cls.get_yaml_property('server', 'log_to_stdout', default=False)
            cls.node_url = cls.get_yaml_property('server', 'node_rpc_url', default='http://[::1]:7072' if cls.banano else 'http://[::1]:7076')
            cls.node_ws_url = cls.get_yaml_property('server', 'node_ws_url', None)
            cls.port = cls.get_yaml_property('server', 'port', default=11338)
            cls.host = cls.get_yaml_property('server', 'host', default='127.0.0.1')
            cls.work_peers = cls.get_yaml_property('wallet', 'work_peers', [])
            cls.node_work_generate = cls.get_yaml_property('wallet', 'node_work_generate', False)
            cls.receive_minimum = cls.get_yaml_property('wallet', 'receive_minimum', 1000000000000000000000000000 if cls.banano else 1000000000000000000000000)
            cls.auto_receive_on_send = cls.get_yaml_property('wallet', 'auto_receive_on_send', True)
            cls.max_work_processes = cls.get_yaml_property('wallet', 'max_work_processes', 1)
            cls.max_sign_threads = cls.get_yaml_property('wallet', 'max_sign_threads', 1)
            if not cls.banano:
                cls.preconfigured_reps = cls.get_yaml_property('wallet', 'preconfigured_representatives_nano', default=None)
            else:
                cls.preconfigured_reps = cls.get_yaml_property('wallet', 'preconfigured_representatives_banano', default=None)
            # Enforce that all reps are valid
            if cls.preconfigured_reps is not None:
                cls.preconfigured_reps = set(cls.preconfigured_reps)
                for r in cls.preconfigured_reps:
                    if not Validators.is_valid_address(r):
                        log.server_logger.warn(f"{r} is not a valid representative!")
                        cls.preconfigured_reps.remove(r)
                if len(cls.preconfigured_reps) == 0:
                    cls.preconfigured_reps = None
            # Go to default if None
            if cls.preconfigured_reps is None:
                cls.preconfigured_reps = DEFAULT_BANANO_REPS if cls.banano else DEFAULT_NANO_REPS

        return cls._instance

    @classmethod
    def get_yaml_property(cls, category: str, subcategory: str, default):
        """Get a property from yaml config"""
        if cls.yaml is None:
            return default
        elif category in cls.yaml and cls.yaml[category] is not None and subcategory in cls.yaml[category]:
            return cls.yaml[category][subcategory]
        return default

    def get_random_rep(self) -> str:
        """Returns a random representative"""
        return secrets.choice(self.preconfigured_reps)
