# coding: utf8
import re
import sys
import pippin.version

from setuptools import find_packages, setup

if sys.version_info < (3, 6):
    raise RuntimeError("Pippin requires Python >= 3.6")

def requirements() -> list:
    try:
        ret = open("requirements.txt", "rt").read().splitlines()
        if sys.platform not in ('win32', 'cygwin', 'cli'):
            ret.append('uvloop>=0.14.0')
        return ret
    except FileNotFoundError:
        ret = [
            'tortoise-orm>=0.15.24,<0.16',
            'aiosqlite>=0.10.0',
            'asyncpg>=0.20.0',
            'pypika<0.37.0',
            'aiomysql>=0.0.20',
            'bitstring>=3.1.6',
            'aiodns>=2.0.0',
            'aioredis>=1.3.0',
            'aioredlock==0.3.0',
            'python-dotenv>=0.10.3',
            'python-rapidjson>=0.9.1',
            'nanopy>=20.0.0',
            'aiohttp>=3.6.2',
            'pyyaml>=5.3.1',
            'pycryptodome>=3.9.4',
            'aiounittest>=1.3.1',
            'websockets>=8.1',
            'setuptools'
        ]
        if sys.platform not in ('win32', 'cygwin', 'cli'):
            ret.append('uvloop>=0.14.0')
        return ret
setup(
    # Application name:
    name="pippin-wallet",
    # Version number:
    version=pippin.version.__version__,
    # Application author details:
    author="Appditto LLC",
    author_email="hello@appditto.com",
    # License
    license="MIT License",
    # Packages
    packages=find_packages(include=["pippin*"]),
    package_data={'pippin': ['*.yaml']},
    zip_safe=True,
    # Details
    url="https://github.com/appditto/pippin_nano_wallet",
    description="A production-ready, high-performance developer wallet for Nano and BANANO.",
    long_description=open("README.md", "r").read(),
    long_description_content_type="text/markdown",
    classifiers=[
        "License :: OSI Approved :: MIT License",
        "Development Status :: 5 - Production/Stable",
        "Intended Audience :: Developers",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.6",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Framework :: AsyncIO",
        "Topic :: Security :: Cryptography",
        "Operating System :: POSIX",
        "Operating System :: MacOS :: MacOS X",
    ],
    keywords=(
        "cryptocurrency wallet nano banano "
        "bitcoin api aiohttp "
        "async asyncio aio"
    ),
    # Dependent packages (distributions)
    install_requires=requirements(),
    entry_points={
        'console_scripts': [
            'pippin-server = pippin.main:main',
            'pippin-cli = pippin.pippin_cli:main'
        ]
    }
)
