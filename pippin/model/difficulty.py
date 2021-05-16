import nanolib
import pippin.config as config

class DifficultyModel(object):
  _instance = None

  def __init__(self):
    raise RuntimeError('Call instance() instead')

  @classmethod
  def instance(cls) -> 'DifficultyModel':
    if cls._instance is None:
      cls._instance = cls.__new__(cls)
      cls.receive_difficulty = 'fffffe0000000000' if config.Config.instance().banano else 'fffffe0000000000'
      cls.send_difficulty = 'fffffe0000000000' if config.Config.instance().banano else 'fffffff800000000'

    return cls._instance