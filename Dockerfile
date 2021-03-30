FROM python:3.9-slim-buster

RUN apt-get update \
 && apt-get install -y build-essential \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/app

COPY requirements.txt ./
RUN pip install --trusted-host pypi.org --no-cache-dir pippin-wallet==1.1.21

RUN mkdir PippinData
COPY docker.config.yaml PippinData/config.yaml

CMD ["pippin-server"]