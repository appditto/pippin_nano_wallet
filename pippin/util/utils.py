from pathlib import Path, PurePath
import os

class Utils(object):
    """Generic utilities"""
    @staticmethod
    def get_project_root():
        home_path = Path.home()
        pippin_path = home_path.joinpath(PurePath('PippinData'))
        if not os.path.exists(pippin_path):
            os.makedirs(pippin_path)
        return pippin_path