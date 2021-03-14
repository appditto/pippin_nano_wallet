import nanopy
import pippin.config as config

class DifficultyModel(object):
  _instance = None

  def __init__(self):
    raise RuntimeError('Call instance() instead')

  @classmethod
  def instance(cls) -> 'DifficultyModel':
    if cls._instance is None:
      cls._instance = cls.__new__(cls)
      cls.base_receive_difficulty = 'fffffe0000000000' if config.Config.instance().banano else 'fffffe0000000000'
      cls.base_send_difficulty = 'fffffe0000000000' if config.Config.instance().banano else 'fffffff800000000'
      cls.receive_difficulty = cls.base_receive_difficulty
      cls.send_difficulty = cls.base_send_difficulty

    return cls._instance

  def adjusted_receive_difficulty(self, difficulty: str) -> str:
      """Ensure reasonable difficulty limits because the node is unreasonable sometimes"""
      if nanolib.work.derive_work_multiplier(difficulty, base_difficulty=self.base_receive_difficulty) > 8:
          return nanolib.work.derive_work_difficulty(8, base_difficulty=self.base_receive_difficulty)
      return difficulty

  def adjusted_send_difficulty(self, difficulty: str) -> str:
      """Ensure reasonable difficulty limits because the node is unreasonable sometimes"""
      if nanolib.work.derive_work_multiplier(difficulty, base_difficulty=self.base_send_difficulty) > 8:
          return nanolib.work.derive_work_difficulty(8, base_difficulty=self.base_send_difficulty)
      return difficulty