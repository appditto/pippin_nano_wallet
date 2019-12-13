import logging
import os
from tortoise import Tortoise
from pippin.util.utils import Utils
import pathlib

class DBConfig(object):
    def __init__(self, mock = False):
        self.logger = logging.getLogger()
        self.modules = {'db': ['pippin.db.models.wallet', 'pippin.db.models.account', 'pippin.db.models.adhoc_account', 'pippin.db.models.block']}
        self.mock = mock
        if self.mock:
            self.use_postgres = False
            self.use_mysql = False
            return
        # Postgres
        self.use_postgres = False
        self.postgres_db = os.getenv('POSTGRES_DB')
        self.postgres_user = os.getenv('POSTGRES_USER')
        self.postgres_password = os.getenv('POSTGRES_PASSWORD')
        self.postgres_host = os.getenv('POSTGRES_HOST', '127.0.0.1')
        self.postgres_port = os.getenv('POSTGRES_PORT', 5432)
        if self.postgres_db is not None and self.postgres_user is not None and self.postgres_password is not None:
            self.use_postgres = True
        elif self.postgres_db is not None or self.postgres_user is not None or self.postgres_password is not None:
            raise Exception("ERROR: Postgres is not properly configured. POSTGRES_DB, POSTGRES_USER, and POSTGRES_PASSWORD environment variables are all required.")
        # MySQL
        self.use_mysql = False
        if not self.use_postgres:
            self.mysql_db = os.getenv('MYSQL_DB')
            self.mysql_user = os.getenv('MYSQL_USER')
            self.mysql_password = os.getenv('MYSQL_PASSWORD')
            self.mysql_host = os.getenv('MYSQL_HOST', '127.0.0.1')
            self.mysql_port = os.getenv('MYSQL_PORT', 3306)
            if self.mysql_db is not None and self.mysql_user is not None and self.mysql_password is not None:
                self.use_mysql = True
            elif self.mysql_db is not None or self.mysql_user is not None or self.mysql_password is not None:
                raise Exception("ERROR: Postgres is not properly configured. MYSQL_DB, MYSQL_USER, and MYSQL_PASSWORD environment variables are all required.")

    async def init_db(self):
        if self.use_postgres:
            self.logger.info(f"Using PostgreSQL Database {self.postgres_db}")
            await Tortoise.init(
                db_url=f'postgres://{self.postgres_user}:{self.postgres_password}@{self.postgres_host}:{self.postgres_port}/{self.postgres_db}',
                modules=self.modules
            )
        elif self.use_mysql:
            self.logger.info(f"Using MySQL Database {self.mysql_db}")
            await Tortoise.init(
                db_url=f'mysql://{self.mysql_user}:{self.mysql_password}@{self.mysql_host}:{self.mysql_port}/{self.mysql_db}',
                modules=self.modules
            )
        else:
            self.logger.info(f"Using SQLite database pippin.db")
            dbpath = Utils.get_project_root().joinpath(pathlib.PurePath('pippin.db')) if not self.mock else Utils.get_project_root().joinpath(pathlib.PurePath('mock.db'))
            await Tortoise.init(
                db_url=f'sqlite://{dbpath}',
                modules=self.modules
            )
        # Create tables
        await Tortoise.generate_schemas(safe=True)