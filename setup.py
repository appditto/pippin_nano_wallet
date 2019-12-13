# coding: utf8
import re
import sys
import pippin.version

from setuptools import find_packages, setup

if sys.version_info < (3, 7):
    raise RuntimeError("Pippin requires Python >= 3.7")

def requirements() -> list:
    return open("requirements.txt", "rt").read().splitlines()

setup(
    # Application name:
    name="pippin-wallet",
    # Version number:
    version=pippin.version.__version__,
    # Application author details:
    author="Appditto LLC",
    author_email="hello@appditto.com",
    # License
    license="MIT LIcense",
    # Packages
    packages=find_packages(include=["pippin*"]),
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